package snippets

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Snippet struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Content string `json:"content,omitempty"`
}

type SnippetManager interface {
	ListSnippets() ([]Snippet, error)
	SaveSnippet(name string, content string) error
	GetSnippet(name string) (string, error)
	DeleteSnippet(name string) error
}

type defaultSnippetManager struct {
	mu      sync.RWMutex
	baseDir string
}

func NewManager() (SnippetManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	snippetsDir := filepath.Join(homeDir, ".goscratched", "snippets")
	if err := os.MkdirAll(snippetsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create snippets dir: %w", err)
	}

	return &defaultSnippetManager{
		baseDir: snippetsDir,
	}, nil
}

func (m *defaultSnippetManager) ListSnippets() ([]Snippet, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	entries, err := os.ReadDir(m.baseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read snippets dir: %w", err)
	}

	var list []Snippet
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go") {
			fullPath := filepath.Join(m.baseDir, entry.Name())
			list = append(list, Snippet{
				Name: entry.Name(),
				Path: fullPath,
			})
		}
	}
	return list, nil
}

func (m *defaultSnippetManager) SaveSnippet(name string, content string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !strings.HasSuffix(name, ".go") {
		name = name + ".go"
	}

	filePath := filepath.Join(m.baseDir, name)
	return os.WriteFile(filePath, []byte(content), 0644)
}

func (m *defaultSnippetManager) GetSnippet(name string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if !strings.HasSuffix(name, ".go") {
		name = name + ".go"
	}

	filePath := filepath.Join(m.baseDir, name)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read snippet %s: %w", name, err)
	}
	return string(data), nil
}

func (m *defaultSnippetManager) DeleteSnippet(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !strings.HasSuffix(name, ".go") {
		name = name + ".go"
	}

	filePath := filepath.Join(m.baseDir, name)
	return os.Remove(filePath)
}
