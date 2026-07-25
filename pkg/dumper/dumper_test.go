package dumper_test

import (
	"database/sql"
	"errors"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/panda-coder/GoScratch/pkg/dumper"
	"github.com/stretchr/testify/suite"
)

type DumperTestSuite struct {
	suite.Suite
	dumper dumper.Dumper
}

func (s *DumperTestSuite) SetupTest() {
	s.dumper = dumper.New()
}

func (s *DumperTestSuite) TestInspectPrimitive() {
	node := s.dumper.Inspect(42)
	s.Require().NotNil(node)
	s.Equal("int", node.Type)
	s.Equal(dumper.KindPrimitive, node.Kind)
	s.Equal("42", node.ValueStr)

	strNode := s.dumper.Inspect("hello goscratched")
	s.Equal("string", strNode.Type)
	s.Equal(dumper.KindPrimitive, strNode.Kind)
	s.Equal("hello goscratched", strNode.ValueStr)
}

func (s *DumperTestSuite) TestInspectNil() {
	node := s.dumper.Inspect(nil)
	s.Require().NotNil(node)
	s.Equal(dumper.KindNil, node.Kind)
	s.Equal("<nil>", node.ValueStr)
}

func (s *DumperTestSuite) TestInspectErrorInterface() {
	err := errors.New("something went wrong")
	node := s.dumper.Inspect(err)
	s.Require().NotNil(node)
	s.Equal(dumper.KindPrimitive, node.Kind)
	s.Equal("something went wrong", node.ValueStr)
}

type User struct {
	Name   string
	Age    int
	secret string
}

func (s *DumperTestSuite) TestInspectStruct() {
	u := User{Name: "Alice", Age: 30, secret: "hidden"}
	node := s.dumper.Inspect(u)

	s.Require().NotNil(node)
	s.Equal(dumper.KindStruct, node.Kind)
	s.Len(node.Children, 3)

	s.Equal("Name", node.Children[0].Name)
	s.Equal("Alice", node.Children[0].ValueStr)

	s.Equal("Age", node.Children[1].Name)
	s.Equal("30", node.Children[1].ValueStr)

	s.Equal("secret", node.Children[2].Name)
	s.Equal("<unexported>", node.Children[2].ValueStr)
}

func (s *DumperTestSuite) TestInspectSliceAndMap() {
	numbers := []int{10, 20}
	sliceNode := s.dumper.Inspect(numbers)

	s.Require().NotNil(sliceNode)
	s.Equal(dumper.KindSlice, sliceNode.Kind)
	s.Equal("len=2", sliceNode.ValueStr)
	s.Len(sliceNode.Children, 2)
	s.Equal("0", sliceNode.Children[0].Name)
	s.Equal("10", sliceNode.Children[0].ValueStr)

	m := map[string]int{"a": 1, "b": 2}
	mapNode := s.dumper.Inspect(m)
	s.Require().NotNil(mapNode)
	s.Equal(dumper.KindMap, mapNode.Kind)
	s.Len(mapNode.Children, 2)
	s.Equal("a", mapNode.Children[0].Name)
	s.Equal("1", mapNode.Children[0].ValueStr)
}

type Node struct {
	Value int
	Next  *Node
}

func (s *DumperTestSuite) TestInspectCircularReference() {
	n1 := &Node{Value: 1}
	n2 := &Node{Value: 2, Next: n1}
	n1.Next = n2

	node := s.dumper.Inspect(n1)
	s.Require().NotNil(node)
	s.Equal(dumper.KindPointer, node.Kind)
	s.Len(node.Children, 1)

	structNode := node.Children[0]
	s.Equal(dumper.KindStruct, structNode.Kind)
	s.Len(structNode.Children, 2)

	nextPtrNode := structNode.Children[1]
	s.Equal(dumper.KindPointer, nextPtrNode.Kind)

	n2Struct := nextPtrNode.Children[0]
	n2NextPtr := n2Struct.Children[1]
	s.Equal(dumper.KindCircular, n2NextPtr.Kind)
	s.Equal("[Circular Reference]", n2NextPtr.ValueStr)
}

func (s *DumperTestSuite) TestInspectSQLRows() {
	db, err := sql.Open("sqlite", ":memory:")
	s.Require().NoError(err)
	defer db.Close()

	_, err = db.Exec("CREATE TABLE users (id INT, name TEXT)")
	s.Require().NoError(err)
	_, err = db.Exec("INSERT INTO users VALUES (1, 'Bob'), (2, 'Charlie')")
	s.Require().NoError(err)

	rows, err := db.Query("SELECT id, name FROM users ORDER BY id ASC")
	s.Require().NoError(err)

	node := s.dumper.Inspect(rows)
	s.Require().NotNil(node)
	s.Equal("*sql.Rows", node.Type)
	s.Equal(dumper.KindSlice, node.Kind)
	s.Equal("rows=2, cols=2", node.ValueStr)
	s.Len(node.Children, 2)

	row0 := node.Children[0]
	s.Equal("Row", row0.Type)
	s.Len(row0.Children, 2)
	s.Equal("id", row0.Children[0].Name)
	s.Equal("1", row0.Children[0].ValueStr)
	s.Equal("name", row0.Children[1].Name)
	s.Equal("Bob", row0.Children[1].ValueStr)
}

func TestDumperSuite(t *testing.T) {
	suite.Run(t, new(DumperTestSuite))
}
