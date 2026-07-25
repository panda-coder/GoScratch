package runner_test

import (
	"context"
	"testing"
	"time"

	"github.com/panda-coder/GoScratch/pkg/runner"
	"github.com/stretchr/testify/suite"
)

type RunnerTestSuite struct {
	suite.Suite
	runner runner.Runner
}

func (s *RunnerTestSuite) SetupTest() {
	s.runner = runner.New()
}

func (s *RunnerTestSuite) TestWrapCode() {
	codeNoPkg := `x := 10; fmt.Println(x)`
	wrapped := runner.WrapCode(codeNoPkg)
	s.Contains(wrapped, "package main")
	s.Contains(wrapped, "func main()")
	s.Contains(wrapped, "x := 10")

	codeWithPkg := `package main
func main() {
    fmt.Println("hello")
}`
	wrappedWithPkg := runner.WrapCode(codeWithPkg)
	s.Contains(wrappedWithPkg, "package main")
	s.Contains(wrappedWithPkg, "fmt.Println(\"hello\")")
}

func (s *RunnerTestSuite) TestExecuteYaegiSimple() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	code := `fmt.Println("hello goscratched")`
	res, err := s.runner.Execute(ctx, code)
	s.NoError(err)
	s.Require().NotNil(res)
	s.Nil(res.Err)
	s.Equal(runner.ModeYaegi, res.ModeUsed)
	s.Contains(res.Stdout, "hello goscratched")
}

func (s *RunnerTestSuite) TestExecuteYaegiDump() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	code := `
a := 42
Dump(a)
`
	res, err := s.runner.Execute(ctx, code)
	s.NoError(err)
	s.Require().NotNil(res)
	s.Nil(res.Err)
	s.Len(res.DumpData, 1)
	s.Equal("42", res.DumpData[0].ValueStr)
}

func (s *RunnerTestSuite) TestExecuteGoRunFallback() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	goRunRunner := runner.New(runner.WithDisableYaegi(true))

	code := `
a := 100
Dump(a)
fmt.Println("fallback output")
`
	res, err := goRunRunner.Execute(ctx, code)
	s.NoError(err)
	s.Require().NotNil(res)
	s.Nil(res.Err)
	s.Equal(runner.ModeGoRun, res.ModeUsed)
	s.Contains(res.Stdout, "fallback output")
	s.Len(res.DumpData, 1)
	s.Equal("100", res.DumpData[0].ValueStr)
}

func (s *RunnerTestSuite) TestExecuteTimeout() {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	code := `for {}`
	res, _ := s.runner.Execute(ctx, code)
	s.Require().NotNil(res)
	s.Error(res.Err)
}

func TestRunnerSuite(t *testing.T) {
	suite.Run(t, new(RunnerTestSuite))
}
