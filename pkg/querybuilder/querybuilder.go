package querybuilder

// Auto-generated basic implementation based on spec/features/03-query-builder/plan.md

type SuggestionKind string

const (
	KindKeyword SuggestionKind = "KEYWORD"
	KindTable   SuggestionKind = "TABLE"
	KindColumn  SuggestionKind = "COLUMN"
	KindClause  SuggestionKind = "CLAUSE"
)

type Suggestion struct {
	Label string         `json:"label"`
	Value string         `json:"value"`
	Kind  SuggestionKind `json:"kind"`
}

type Builder interface {
	Suggest(input string) []Suggestion
	Build(parts []string) string
}

// simpleBuilder is a minimal in-memory implementation used for tests and basic functionality.
type simpleBuilder struct {
	keywords []Suggestion
}

func NewSimpleBuilder() Builder {
	kb := []Suggestion{
		{Label: "SELECT", Value: "SELECT", Kind: KindKeyword},
		{Label: "FROM", Value: "FROM", Kind: KindKeyword},
		{Label: "WHERE", Value: "WHERE", Kind: KindKeyword},
	}
	return &simpleBuilder{keywords: kb}
}

func (s *simpleBuilder) Suggest(input string) []Suggestion {
	var out []Suggestion
	for _, k := range s.keywords {
		if len(input) == 0 || len(k.Value) >= len(input) && k.Value[:len(input)] == input {
			out = append(out, k)
		}
	}
	return out
}

func (s *simpleBuilder) Build(parts []string) string {
	return "" // noop for now
}
