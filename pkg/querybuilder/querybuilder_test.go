package querybuilder

import "testing"

func TestSuggestKeywords(t *testing.T) {
	b := NewSimpleBuilder()
	res := b.Suggest("SE")
	if len(res) == 0 {
		t.Fatalf("expected suggestions for 'SE'")
	}
	if res[0].Value != "SELECT" {
		t.Fatalf("expected first suggestion to be SELECT, got %s", res[0].Value)
	}
}
