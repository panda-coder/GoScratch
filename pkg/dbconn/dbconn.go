package dbconn

// Minimal, spec-compatible placeholder implementation for Database Connection Manager.

type ConnInfo struct {
	Name string `json:"name"`
	DSN  string `json:"dsn"`
}

type Manager struct {
	conns map[string]ConnInfo
}

func NewManager() *Manager {
	return &Manager{conns: make(map[string]ConnInfo)}
}

func (m *Manager) Add(name, dsn string) error {
	m.conns[name] = ConnInfo{Name: name, DSN: dsn}
	return nil
}

func (m *Manager) Get(name string) (ConnInfo, bool) {
	c, ok := m.conns[name]
	return c, ok
}

func (m *Manager) List() []ConnInfo {
	out := make([]ConnInfo, 0, len(m.conns))
	for _, v := range m.conns {
		out = append(out, v)
	}
	return out
}

func (m *Manager) Remove(name string) {
	delete(m.conns, name)
}
