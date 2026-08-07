package results

// Row represents a single row of values keyed by column name.
type Row map[string]any

// ResultSet holds tabular results: ordered columns and rows.
type ResultSet struct {
	Columns []string `json:"columns"`
	Rows    []Row    `json:"rows"`
}

// Visualizer defines the minimal contract for a results visualizer.
type Visualizer interface {
	Load(result ResultSet) error
	SelectedRow() Row
}

// SimpleVisualizer is a small in-memory implementation used by tests.
type SimpleVisualizer struct {
	result   ResultSet
	selected int
}

// NewSimpleVisualizer constructs a SimpleVisualizer.
func NewSimpleVisualizer() *SimpleVisualizer {
	return &SimpleVisualizer{selected: -1}
}

// Load stores the provided ResultSet and resets selection.
func (s *SimpleVisualizer) Load(r ResultSet) error {
	s.result = r
	s.selected = -1
	return nil
}

// SelectedRow returns the currently selected row or nil when none selected.
func (s *SimpleVisualizer) SelectedRow() Row {
	if s.selected < 0 || s.selected >= len(s.result.Rows) {
		return nil
	}
	return s.result.Rows[s.selected]
}

// Select chooses a row index (out-of-range still allowed; SelectedRow will return nil).
func (s *SimpleVisualizer) Select(i int) {
	s.selected = i
}
