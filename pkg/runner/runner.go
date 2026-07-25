package runner

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/panda-coder/GoScratch/pkg/dumper"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

type Mode string

const (
	ModeYaegi Mode = "YAEGI"
	ModeGoRun Mode = "GORUN"
)

const (
	DumpStartTag = "<<<GOSCRATCH_DUMP_START>>>"
	DumpEndTag   = "<<<GOSCRATCH_DUMP_END>>>"
)

type ExecutionResult struct {
	Stdout   string
	Stderr   string
	DumpData []*dumper.DumpNode
	ModeUsed Mode
	Err      error
}

type Runner interface {
	Execute(ctx context.Context, code string) (*ExecutionResult, error)
}

type defaultRunner struct {
	dumper        dumper.Dumper
	disableYaegi  bool
	allowFallback bool
}

type Option func(*defaultRunner)

func WithDisableYaegi(disable bool) Option {
	return func(r *defaultRunner) {
		r.disableYaegi = disable
	}
}

func WithAllowFallback(allow bool) Option {
	return func(r *defaultRunner) {
		r.allowFallback = allow
	}
}

func New(opts ...Option) Runner {
	r := &defaultRunner{
		dumper:        dumper.New(),
		allowFallback: true,
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (r *defaultRunner) Execute(ctx context.Context, code string) (*ExecutionResult, error) {
	wrappedCode := WrapCode(code)

	if !r.disableYaegi {
		res, err := r.executeYaegi(ctx, wrappedCode)
		if err == nil && res.Err == nil {
			return res, nil
		}

		if !r.allowFallback {
			return res, err
		}
	}

	return r.executeGoRun(ctx, wrappedCode)
}

func WrapCode(code string) string {
	trimmed := strings.TrimSpace(code)
	if strings.HasPrefix(trimmed, "package ") {
		if !strings.Contains(code, "goscratched") {
			lines := strings.Split(code, "\n")
			var newLines []string
			imported := false
			for _, line := range lines {
				newLines = append(newLines, line)
				if !imported && strings.HasPrefix(strings.TrimSpace(line), "package ") {
					newLines = append(newLines, `import . "goscratched/goscratched"`)
					imported = true
				}
			}
			return strings.Join(newLines, "\n")
		}
		return code
	}

	return fmt.Sprintf(`package main

import (
	"fmt"
	. "goscratched/goscratched"
)

func main() {
	_ = fmt.Sprintf
%s
}
`, indent(code, "\t"))
}

func indent(code, prefix string) string {
	lines := strings.Split(code, "\n")
	for i, line := range lines {
		if strings.TrimSpace(line) != "" {
			lines[i] = prefix + line
		}
	}
	return strings.Join(lines, "\n")
}

func (r *defaultRunner) executeYaegi(ctx context.Context, code string) (*ExecutionResult, error) {
	var stdoutBuf, stderrBuf bytes.Buffer

	i := interp.New(interp.Options{
		Stdout: &stdoutBuf,
		Stderr: &stderrBuf,
	})

	if err := i.Use(stdlib.Symbols); err != nil {
		return nil, fmt.Errorf("failed to load stdlib symbols in yaegi: %w", err)
	}

	var dumpNodes []*dumper.DumpNode
	var mu sync.Mutex

	dumpFunc := func(v any) {
		node := r.dumper.Inspect(v)
		mu.Lock()
		dumpNodes = append(dumpNodes, node)
		mu.Unlock()
	}

	symbols := map[string]map[string]reflect.Value{
		"goscratched/goscratched/goscratched": {
			"Dump": reflect.ValueOf(dumpFunc),
		},
		"goscratched/goscratched": {
			"Dump": reflect.ValueOf(dumpFunc),
		},
	}
	if err := i.Use(symbols); err != nil {
		return nil, fmt.Errorf("failed to register Dump symbol in yaegi: %w", err)
	}

	_, evalErr := i.EvalWithContext(ctx, code)

	res := &ExecutionResult{
		Stdout:   stdoutBuf.String(),
		Stderr:   stderrBuf.String(),
		DumpData: dumpNodes,
		ModeUsed: ModeYaegi,
		Err:      evalErr,
	}

	return res, nil
}

func (r *defaultRunner) executeGoRun(ctx context.Context, code string) (*ExecutionResult, error) {
	tempDir, err := os.MkdirTemp("", "goscratched-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir for go run: %w", err)
	}
	defer os.RemoveAll(tempDir)

	pkgDir := filepath.Join(tempDir, "goscratched")
	if err := os.MkdirAll(pkgDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create package dir: %w", err)
	}

	helperCode := fmt.Sprintf(`package goscratched

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type DumpNode struct {
	Name     string      %[1]sjson:"name"%[1]s
	Type     string      %[1]sjson:"type"%[1]s
	Kind     string      %[1]sjson:"kind"%[1]s
	ValueStr string      %[1]sjson:"value_str"%[1]s
	Children []*DumpNode %[1]sjson:"children,omitempty"%[1]s
}

func Dump(v any) {
	node := inspectValue("", reflect.ValueOf(v), make(map[uintptr]bool))
	data, _ := json.Marshal(node)
	fmt.Printf("%[2]s\n%%s\n%[3]s\n", string(data))
}

func inspectValue(name string, val reflect.Value, visited map[uintptr]bool) *DumpNode {
	if !val.IsValid() {
		return &DumpNode{Name: name, Type: "nil", Kind: "NIL", ValueStr: "<nil>"}
	}
	typ := val.Type()
	typeStr := typ.String()

	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		node := &DumpNode{Name: name, Type: typeStr, Kind: "SLICE", ValueStr: fmt.Sprintf("len=%%d", val.Len())}
		for i := 0; i < val.Len(); i++ {
			node.Children = append(node.Children, inspectValue(fmt.Sprintf("%%d", i), val.Index(i), visited))
		}
		return node
	case reflect.Map:
		node := &DumpNode{Name: name, Type: typeStr, Kind: "MAP", ValueStr: fmt.Sprintf("len=%%d", val.Len())}
		for _, k := range val.MapKeys() {
			node.Children = append(node.Children, inspectValue(fmt.Sprintf("%%v", k.Interface()), val.MapIndex(k), visited))
		}
		return node
	case reflect.Struct:
		node := &DumpNode{Name: name, Type: typeStr, Kind: "STRUCT"}
		for i := 0; i < val.NumField(); i++ {
			f := val.Field(i)
			if f.CanInterface() {
				node.Children = append(node.Children, inspectValue(typ.Field(i).Name, f, visited))
			}
		}
		return node
	default:
		return &DumpNode{Name: name, Type: typeStr, Kind: "PRIMITIVE", ValueStr: fmt.Sprintf("%%v", val.Interface())}
	}
}
`, "`", DumpStartTag, DumpEndTag)

	if err := os.WriteFile(filepath.Join(pkgDir, "goscratched.go"), []byte(helperCode), 0644); err != nil {
		return nil, fmt.Errorf("failed to write helper code: %w", err)
	}

	goModCode := "module temp\n\ngo 1.22\n"
	if err := os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte(goModCode), 0644); err != nil {
		return nil, fmt.Errorf("failed to write temp go.mod: %w", err)
	}

	mainPath := filepath.Join(tempDir, "main.go")
	adjustedCode := strings.ReplaceAll(code, `"goscratched/goscratched"`, `"temp/goscratched"`)
	if err := os.WriteFile(mainPath, []byte(adjustedCode), 0644); err != nil {
		return nil, fmt.Errorf("failed to write temp main.go: %w", err)
	}

	cmd := exec.CommandContext(ctx, "go", "run", ".")
	cmd.Dir = tempDir

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err = cmd.Run()

	stdoutRaw := stdoutBuf.String()
	stdoutClean, dumpNodes := parseDumpData(stdoutRaw)

	res := &ExecutionResult{
		Stdout:   stdoutClean,
		Stderr:   stderrBuf.String(),
		DumpData: dumpNodes,
		ModeUsed: ModeGoRun,
		Err:      err,
	}

	return res, nil
}

func parseDumpData(rawStdout string) (string, []*dumper.DumpNode) {
	lines := strings.Split(rawStdout, "\n")
	var cleanLines []string
	var nodes []*dumper.DumpNode

	inDump := false
	var dumpJSONBuilder strings.Builder

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == DumpStartTag {
			inDump = true
			dumpJSONBuilder.Reset()
			continue
		}
		if trimmed == DumpEndTag {
			inDump = false
			var node dumper.DumpNode
			if err := json.Unmarshal([]byte(dumpJSONBuilder.String()), &node); err == nil {
				nodes = append(nodes, &node)
			}
			continue
		}

		if inDump {
			dumpJSONBuilder.WriteString(line)
		} else {
			cleanLines = append(cleanLines, line)
		}
	}

	return strings.Join(cleanLines, "\n"), nodes
}

func TimeoutContext(parent context.Context, duration time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, duration)
}
