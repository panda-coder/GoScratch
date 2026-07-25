package snippets_test

import (
	"testing"

	"github.com/panda-coder/GoScratch/pkg/snippets"
	"github.com/stretchr/testify/suite"
)

type SnippetsTestSuite struct {
	suite.Suite
	mgr snippets.SnippetManager
}

func (s *SnippetsTestSuite) SetupTest() {
	mgr, err := snippets.NewManager()
	s.Require().NoError(err)
	s.mgr = mgr
}

func (s *SnippetsTestSuite) TestSaveAndGetSnippet() {
	snippetName := "test_demo.go"
	content := "println(\"hello snippets\")"

	err := s.mgr.SaveSnippet(snippetName, content)
	s.NoError(err)

	readContent, err := s.mgr.GetSnippet(snippetName)
	s.NoError(err)
	s.Equal(content, readContent)

	list, err := s.mgr.ListSnippets()
	s.NoError(err)
	s.NotEmpty(list)

	// Clean up
	s.Require().NoError(s.mgr.DeleteSnippet(snippetName))
}

func TestSnippetsSuite(t *testing.T) {
	suite.Run(t, new(SnippetsTestSuite))
}
