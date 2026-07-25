package ui

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/panda-coder/GoScratch/pkg/db"
	"github.com/panda-coder/GoScratch/pkg/snippets"
)

type SidebarNode struct {
	ID       string
	Label    string
	Kind     string
	Path     string
	ConnID   string
	Table    string
	Children []*SidebarNode
}

type Sidebar struct {
	app        *App
	dbMgr      db.DBManager
	snipMgr    snippets.SnippetManager
	treeWidget *widget.Tree

	mu    sync.RWMutex
	nodes map[string]*SidebarNode
	roots []string
}

func NewSidebar(a *App) (*Sidebar, error) {
	dbMgr, err := db.NewManager()
	if err != nil {
		return nil, fmt.Errorf("failed to init db manager: %w", err)
	}

	snipMgr, err := snippets.NewManager()
	if err != nil {
		return nil, fmt.Errorf("failed to init snippets manager: %w", err)
	}

	s := &Sidebar{
		app:     a,
		dbMgr:   dbMgr,
		snipMgr: snipMgr,
		nodes:   make(map[string]*SidebarNode),
	}

	s.setupTreeWidget()
	s.Refresh()

	return s, nil
}

func (s *Sidebar) setupTreeWidget() {
	s.treeWidget = widget.NewTree(
		func(id string) []string {
			s.mu.RLock()
			defer s.mu.RUnlock()
			if id == "" {
				return s.roots
			}
			node := s.nodes[id]
			if node == nil {
				return nil
			}
			var childIDs []string
			for _, child := range node.Children {
				childIDs = append(childIDs, child.ID)
			}
			return childIDs
		},
		func(id string) bool {
			s.mu.RLock()
			defer s.mu.RUnlock()
			if id == "" {
				return len(s.roots) > 0
			}
			node := s.nodes[id]
			return node != nil && len(node.Children) > 0
		},
		func(branch bool) fyne.CanvasObject {
			lbl := widget.NewLabel("")
			lbl.TextStyle = fyne.TextStyle{Monospace: false}
			return lbl
		},
		func(id string, branch bool, cell fyne.CanvasObject) {
			s.mu.RLock()
			node := s.nodes[id]
			s.mu.RUnlock()

			lbl := cell.(*widget.Label)
			if node == nil {
				lbl.SetText("")
				return
			}
			lbl.SetText(node.Label)
		},
	)

	s.treeWidget.OnSelected = func(id string) {
		s.mu.RLock()
		node := s.nodes[id]
		s.mu.RUnlock()

		if node == nil {
			return
		}

		switch node.Kind {
		case "SNIPPET":
			content, err := s.snipMgr.GetSnippet(node.Label)
			if err == nil {
				s.app.editor.SetText(content)
				s.app.statusLabel.SetText("Snippet carregado: " + node.Label)
			}
		case "TABLE":
			code := fmt.Sprintf(`// Consulta automática na tabela %s
db, err := sql.Open("sqlite", %q)
if err != nil {
	fmt.Println("Erro na conexão:", err)
	return
}
defer db.Close()

rows, err := db.Query("SELECT * FROM %s LIMIT 50")
if err != nil {
	fmt.Println("Erro na query:", err)
	return
}

Dump(rows)`, node.Table, node.Path, node.Table)
			s.app.editor.SetText(code)
			s.app.statusLabel.SetText("Gerado snippet para tabela: " + node.Table)
		}
	}
}

func (s *Sidebar) Refresh() {
	go func() {
		newNodes := make(map[string]*SidebarNode)
		var newRoots []string

		snippetsRoot := &SidebarNode{
			ID:    "root_snippets",
			Label: "📂 Meus Snippets",
			Kind:  "ROOT_SNIPPETS",
		}
		newNodes["root_snippets"] = snippetsRoot
		newRoots = append(newRoots, "root_snippets")

		snipList, _ := s.snipMgr.ListSnippets()
		for _, snip := range snipList {
			childID := "snip_" + snip.Name
			childNode := &SidebarNode{
				ID:    childID,
				Label: "📄 " + snip.Name,
				Kind:  "SNIPPET",
				Path:  snip.Path,
			}
			newNodes[childID] = childNode
			snippetsRoot.Children = append(snippetsRoot.Children, childNode)
		}

		dbRoot := &SidebarNode{
			ID:    "root_dbs",
			Label: "🗄️ Conexões de Banco",
			Kind:  "ROOT_DBS",
		}
		newNodes["root_dbs"] = dbRoot
		newRoots = append(newRoots, "root_dbs")

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		connList, _ := s.dbMgr.ListConnections()
		for _, connCfg := range connList {
			connNodeID := "db_" + connCfg.ID
			connNode := &SidebarNode{
				ID:     connNodeID,
				Label:  "🔌 " + connCfg.Name,
				Kind:   "DB",
				ConnID: connCfg.ID,
				Path:   connCfg.ConnStr,
			}
			newNodes[connNodeID] = connNode
			dbRoot.Children = append(dbRoot.Children, connNode)

			tables, err := s.dbMgr.GetTables(ctx, connCfg.ID)
			if err == nil {
				for _, tbl := range tables {
					tblNodeID := connNodeID + "_tbl_" + tbl.Name
					tblNode := &SidebarNode{
						ID:     tblNodeID,
						Label:  "📊 " + tbl.Name,
						Kind:   "TABLE",
						ConnID: connCfg.ID,
						Table:  tbl.Name,
						Path:   connCfg.ConnStr,
					}
					newNodes[tblNodeID] = tblNode
					connNode.Children = append(connNode.Children, tblNode)

					for _, col := range tbl.Columns {
						colNodeID := tblNodeID + "_col_" + col.Name
						pkTag := ""
						if col.IsPrimaryKey {
							pkTag = " (PK)"
						}
						colNode := &SidebarNode{
							ID:    colNodeID,
							Label: fmt.Sprintf("🔹 %s [%s]%s", col.Name, col.DataType, pkTag),
							Kind:  "COLUMN",
						}
						newNodes[colNodeID] = colNode
						tblNode.Children = append(tblNode.Children, colNode)
					}
				}
			}
		}

		fyne.Do(func() {
			s.mu.Lock()
			s.nodes = newNodes
			s.roots = newRoots
			s.mu.Unlock()
			s.treeWidget.Refresh()
		})
	}()
}

func (s *Sidebar) BuildContainer() fyne.CanvasObject {
	addConnBtn := widget.NewButtonWithIcon("Nova Conexão", theme.ContentAddIcon(), func() {
		s.showAddConnectionDialog()
	})
	addConnBtn.Importance = widget.LowImportance

	saveSnipBtn := widget.NewButtonWithIcon("Salvar Snippet", theme.DocumentSaveIcon(), func() {
		s.showSaveSnippetDialog()
	})
	saveSnipBtn.Importance = widget.LowImportance

	btnBar := container.NewHBox(addConnBtn, saveSnipBtn)

	header := container.NewVBox(
		widget.NewLabelWithStyle(" Explorer", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		btnBar,
		widget.NewSeparator(),
	)

	return container.NewBorder(header, nil, nil, nil, s.treeWidget)
}

func (s *Sidebar) showAddConnectionDialog() {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Ex: SQLite Dev DB")

	connStrEntry := widget.NewEntry()
	connStrEntry.SetPlaceHolder("Ex: ./meu_banco.db ou :memory:")

	form := dialog.NewForm("Adicionar Conexão SQL", "Salvar", "Cancelar", []*widget.FormItem{
		widget.NewFormItem("Nome da Conexão:", nameEntry),
		widget.NewFormItem("String de Conexão:", connStrEntry),
	}, func(ok bool) {
		if !ok || strings.TrimSpace(nameEntry.Text) == "" || strings.TrimSpace(connStrEntry.Text) == "" {
			return
		}
		id := "conn_" + fmt.Sprintf("%d", time.Now().UnixNano())
		cfg := db.ConnectionConfig{
			ID:      id,
			Name:    nameEntry.Text,
			Driver:  db.DriverSQLite,
			ConnStr: connStrEntry.Text,
		}
		if err := s.dbMgr.SaveConnection(cfg); err == nil {
			s.Refresh()
			s.app.statusLabel.SetText("Nova conexão salva: " + cfg.Name)
		}
	}, s.app.mainWindow)

	form.Resize(fyne.NewSize(450, 250))
	form.Show()
}

func (s *Sidebar) showSaveSnippetDialog() {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Ex: teste_query.go")

	form := dialog.NewForm("Salvar Snippet", "Salvar", "Cancelar", []*widget.FormItem{
		widget.NewFormItem("Nome do Arquivo:", nameEntry),
	}, func(ok bool) {
		if !ok || strings.TrimSpace(nameEntry.Text) == "" {
			return
		}
		fileName := nameEntry.Text
		if !strings.HasSuffix(fileName, ".go") {
			fileName += ".go"
		}
		content := s.app.editor.Text
		if err := s.snipMgr.SaveSnippet(fileName, content); err == nil {
			s.Refresh()
			s.app.statusLabel.SetText("Snippet salvo: " + fileName)
		}
	}, s.app.mainWindow)

	form.Resize(fyne.NewSize(400, 200))
	form.Show()
}
