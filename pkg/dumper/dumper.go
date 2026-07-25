package dumper

import (
	"database/sql"
	"fmt"
	"reflect"
	"sort"
	"strconv"
)

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

type defaultDumper struct{}

func New() Dumper {
	return &defaultDumper{}
}

func (d *defaultDumper) Inspect(v any) *DumpNode {
	visited := make(map[uintptr]bool)
	return inspectValue("", reflect.ValueOf(v), visited)
}

var (
	errorType    = reflect.TypeOf((*error)(nil)).Elem()
	stringerType = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
)

func inspectValue(name string, val reflect.Value, visited map[uintptr]bool) *DumpNode {
	if !val.IsValid() {
		return &DumpNode{
			Name:     name,
			Type:     "nil",
			Kind:     KindNil,
			ValueStr: "<nil>",
		}
	}

	// Direct handling of *sql.Rows
	if val.CanInterface() {
		if rows, ok := val.Interface().(*sql.Rows); ok && rows != nil {
			return inspectSQLRows(name, rows)
		}
	}

	typ := val.Type()
	typeStr := typ.String()

	// Special handling for error interface implementations
	if val.CanInterface() && typ.Implements(errorType) {
		if (val.Kind() == reflect.Pointer || val.Kind() == reflect.Interface) && val.IsNil() {
			return &DumpNode{
				Name:     name,
				Type:     typeStr,
				Kind:     KindNil,
				ValueStr: "<nil>",
			}
		}
		errVal := val.Interface().(error)
		errStr := "<nil>"
		if errVal != nil {
			errStr = errVal.Error()
		}
		return &DumpNode{
			Name:     name,
			Type:     typeStr,
			Kind:     KindPrimitive,
			ValueStr: errStr,
		}
	}

	switch val.Kind() {
	case reflect.Interface:
		if val.IsNil() {
			return &DumpNode{
				Name:     name,
				Type:     typeStr,
				Kind:     KindNil,
				ValueStr: "<nil>",
			}
		}
		child := inspectValue(name, val.Elem(), visited)
		return child

	case reflect.Pointer:
		if val.IsNil() {
			return &DumpNode{
				Name:     name,
				Type:     typeStr,
				Kind:     KindNil,
				ValueStr: "<nil>",
			}
		}
		ptr := val.Pointer()
		if visited[ptr] {
			return &DumpNode{
				Name:     name,
				Type:     typeStr,
				Kind:     KindCircular,
				ValueStr: "[Circular Reference]",
			}
		}
		visited[ptr] = true
		defer func() { delete(visited, ptr) }()

		elemNode := inspectValue("*", val.Elem(), visited)
		return &DumpNode{
			Name:     name,
			Type:     typeStr,
			Kind:     KindPointer,
			ValueStr: fmt.Sprintf("%p", val.Interface()),
			Children: []*DumpNode{elemNode},
		}

	case reflect.Struct:
		node := &DumpNode{
			Name: name,
			Type: typeStr,
			Kind: KindStruct,
		}

		if val.CanInterface() && typ.Implements(stringerType) {
			node.ValueStr = val.Interface().(fmt.Stringer).String()
		}

		numFields := val.NumField()
		for i := 0; i < numFields; i++ {
			fieldTyp := typ.Field(i)
			fieldVal := val.Field(i)

			fieldName := fieldTyp.Name
			var child *DumpNode
			if fieldVal.CanInterface() {
				child = inspectValue(fieldName, fieldVal, visited)
			} else {
				child = &DumpNode{
					Name:     fieldName,
					Type:     fieldTyp.Type.String(),
					Kind:     KindPrimitive,
					ValueStr: "<unexported>",
				}
			}
			node.Children = append(node.Children, child)
		}
		return node

	case reflect.Map:
		if val.IsNil() {
			return &DumpNode{
				Name:     name,
				Type:     typeStr,
				Kind:     KindNil,
				ValueStr: "<nil>",
			}
		}
		ptr := val.Pointer()
		if visited[ptr] {
			return &DumpNode{
				Name:     name,
				Type:     typeStr,
				Kind:     KindCircular,
				ValueStr: "[Circular Reference]",
			}
		}
		visited[ptr] = true
		defer func() { delete(visited, ptr) }()

		node := &DumpNode{
			Name:     name,
			Type:     typeStr,
			Kind:     KindMap,
			ValueStr: fmt.Sprintf("len=%d", val.Len()),
		}

		keys := val.MapKeys()
		sortMapKeys(keys)

		for _, k := range keys {
			v := val.MapIndex(k)
			keyStr := fmt.Sprintf("%v", k.Interface())
			child := inspectValue(keyStr, v, visited)
			node.Children = append(node.Children, child)
		}
		return node

	case reflect.Slice, reflect.Array:
		if val.Kind() == reflect.Slice && val.IsNil() {
			return &DumpNode{
				Name:     name,
				Type:     typeStr,
				Kind:     KindNil,
				ValueStr: "<nil>",
			}
		}

		if val.Kind() == reflect.Slice {
			ptr := val.Pointer()
			if ptr != 0 {
				if visited[ptr] {
					return &DumpNode{
						Name:     name,
						Type:     typeStr,
						Kind:     KindCircular,
						ValueStr: "[Circular Reference]",
					}
				}
				visited[ptr] = true
				defer func() { delete(visited, ptr) }()
			}
		}

		node := &DumpNode{
			Name:     name,
			Type:     typeStr,
			Kind:     KindSlice,
			ValueStr: fmt.Sprintf("len=%d", val.Len()),
		}

		for i := 0; i < val.Len(); i++ {
			elemVal := val.Index(i)
			child := inspectValue(strconv.Itoa(i), elemVal, visited)
			node.Children = append(node.Children, child)
		}
		return node

	default:
		return &DumpNode{
			Name:     name,
			Type:     typeStr,
			Kind:     KindPrimitive,
			ValueStr: fmt.Sprintf("%v", val.Interface()),
		}
	}
}

func inspectSQLRows(name string, rows *sql.Rows) *DumpNode {
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return &DumpNode{
			Name:     name,
			Type:     "*sql.Rows",
			Kind:     KindNil,
			ValueStr: fmt.Sprintf("Error reading columns: %v", err),
		}
	}

	node := &DumpNode{
		Name: name,
		Type: "*sql.Rows",
		Kind: KindSlice,
	}

	rowCount := 0
	for rows.Next() {
		columns := make([]any, len(cols))
		columnPointers := make([]any, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			continue
		}

		rowNode := &DumpNode{
			Name: strconv.Itoa(rowCount),
			Type: "Row",
			Kind: KindStruct,
		}

		for i, colName := range cols {
			val := columns[i]
			if b, ok := val.([]byte); ok {
				val = string(b)
			}
			child := inspectValue(colName, reflect.ValueOf(val), make(map[uintptr]bool))
			rowNode.Children = append(rowNode.Children, child)
		}

		node.Children = append(node.Children, rowNode)
		rowCount++
	}
	node.ValueStr = fmt.Sprintf("rows=%d, cols=%d", rowCount, len(cols))

	return node
}

func sortMapKeys(keys []reflect.Value) {
	sort.Slice(keys, func(i, j int) bool {
		return fmt.Sprintf("%v", keys[i].Interface()) < fmt.Sprintf("%v", keys[j].Interface())
	})
}
