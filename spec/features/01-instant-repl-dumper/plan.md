# Desenho Técnico: Instant REPL & Rich Data Dump

## Arquitetura do Pacote `pkg/dumper`

```go
package dumper

type NodeKind string

const (
	KindPrimitive NodeKind = "PRIMITIVE"
	KindStruct    NodeKind = "STRUCT"
	KindMap       NodeKind = "MAP"
	KindSlice     NodeKind = "SLICE"
	KindPointer   NodeKind = "POINTER"
	KindInterface NodeKind = "INTERFACE"
	KindCircular  NodeKind = "CIRCULAR"
	KindNil       NodeKind = "NIL"
)

type DumpNode struct {
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	Kind     NodeKind    `json:"kind"`
	ValueStr string      `json:"value_str"`
	Children []*DumpNode `json:"children,omitempty"`
}

type Dumper interface {
	Inspect(v any) *DumpNode
}
```
