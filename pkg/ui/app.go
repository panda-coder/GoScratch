package ui

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/panda-coder/GoScratch/pkg/dumper"
	"github.com/panda-coder/GoScratch/pkg/runner"
)

type App struct {
	fyneApp     fyne.App
	mainWindow  fyne.Window
	runner      runner.Runner
	sidebar     *Sidebar
	editor      *widget.Entry
	stdoutEntry *widget.Entry
	stderrEntry *widget.Entry
	statusLabel *widget.Label
	runBtn      *widget.Button
	stopBtn     *widget.Button
	dumpTab     *container.TabItem
	appTabs     *container.AppTabs

	mu            sync.Mutex
	currentCancel context.CancelFunc
	treeStore     *treeStore
	treeWidget    *widget.Tree
}

type treeStore struct {
	nodes map[string]*dumper.DumpNode
	roots []string
}

func NewApp() *App {
	a := &App{
		fyneApp: app.NewWithID("com.goscratched.app"),
		runner:  runner.New(),
		treeStore: &treeStore{
			nodes: make(map[string]*dumper.DumpNode),
		},
	}
	return a
}

func (a *App) Run() {
	a.mainWindow = a.fyneApp.NewWindow("GoScratch - REPL & Scratchpad Instantâneo com Banco de Dados")
	a.mainWindow.Resize(fyne.NewSize(1150, 700))

	a.setupUI()
	a.setupShortcuts()

	a.mainWindow.ShowAndRun()
}

func (a *App) setupUI() {
	// 1. Sidebar (Left)
	sidebar, err := NewSidebar(a)
	if err == nil {
		a.sidebar = sidebar
	}

	// 2. Editor Panel (Center)
	a.editor = widget.NewMultiLineEntry()
	a.editor.TextStyle = fyne.TextStyle{Monospace: true}
	a.editor.SetPlaceHolder("// Digite seu código Go aqui...\n// Exemplo:\nuser := struct {\n    Name string\n    Age  int\n}{\n    Name: \"Alice\",\n    Age:  30,\n}\nDump(user)\nfmt.Println(\"Execução concluída!\")")
	a.editor.SetText(`user := struct {
	Name string
	Age  int
}{
	Name: "Alice",
	Age:  30,
}

Dump(user)
fmt.Println("Executando GoScratch!")`)

	editorContainer := container.NewBorder(
		widget.NewLabelWithStyle(" 📝 Editor de Código", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		nil, nil, nil,
		a.editor,
	)

	// 3. Output Panel (Right)
	a.stdoutEntry = widget.NewMultiLineEntry()
	a.stdoutEntry.TextStyle = fyne.TextStyle{Monospace: true}
	a.stdoutEntry.Disable()

	a.stderrEntry = widget.NewMultiLineEntry()
	a.stderrEntry.TextStyle = fyne.TextStyle{Monospace: true}
	a.stderrEntry.Disable()

	consoleContainer := container.NewVSplit(
		container.NewBorder(widget.NewLabelWithStyle("Saída Standard (stdout)", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), nil, nil, nil, a.stdoutEntry),
		container.NewBorder(widget.NewLabelWithStyle("Erros (stderr)", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), nil, nil, nil, a.stderrEntry),
	)
	consoleContainer.SetOffset(0.7)

	consoleTab := container.NewTabItem("Console", consoleContainer)

	// Dump Tree Widget
	a.treeWidget = widget.NewTree(
		func(id string) []string {
			if id == "" {
				return a.treeStore.roots
			}
			node := a.treeStore.nodes[id]
			if node == nil {
				return nil
			}
			var childrenIDs []string
			for i := range node.Children {
				childrenIDs = append(childrenIDs, id+"."+strconv.Itoa(i))
			}
			return childrenIDs
		},
		func(id string) bool {
			if id == "" {
				return len(a.treeStore.roots) > 0
			}
			node := a.treeStore.nodes[id]
			return node != nil && len(node.Children) > 0
		},
		func(branch bool) fyne.CanvasObject {
			lbl := widget.NewLabel("")
			lbl.TextStyle = fyne.TextStyle{Monospace: true}
			return lbl
		},
		func(id string, branch bool, cell fyne.CanvasObject) {
			lbl := cell.(*widget.Label)
			node := a.treeStore.nodes[id]
			if node == nil {
				lbl.SetText("")
				return
			}
			lbl.SetText(formatNodeText(node))
		},
	)

	dumpContainer := container.NewBorder(
		widget.NewLabelWithStyle(" Inspeção de Dados (Dump)", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		nil, nil, nil,
		a.treeWidget,
	)
	a.dumpTab = container.NewTabItem("Inspeção (Dump)", dumpContainer)

	a.appTabs = container.NewAppTabs(consoleTab, a.dumpTab)

	// 4. Toolbar & Status Bar
	a.runBtn = widget.NewButtonWithIcon("Executar", theme.MediaPlayIcon(), func() {
		a.executeCode()
	})
	a.runBtn.Importance = widget.HighImportance

	a.stopBtn = widget.NewButtonWithIcon("Cancelar", theme.MediaStopIcon(), func() {
		a.cancelExecution()
	})
	a.stopBtn.Disable()

	a.statusLabel = widget.NewLabel("Pronto | Atalho: Ctrl+R ou Ctrl+Enter")

	topToolbar := container.NewHBox(
		a.runBtn,
		a.stopBtn,
		widget.NewSeparator(),
		a.statusLabel,
	)

	editorAndTabsSplit := container.NewHSplit(editorContainer, a.appTabs)
	editorAndTabsSplit.SetOffset(0.5)

	var mainLayout fyne.CanvasObject
	if a.sidebar != nil {
		sidebarSplit := container.NewHSplit(a.sidebar.BuildContainer(), editorAndTabsSplit)
		sidebarSplit.SetOffset(0.24)
		mainLayout = sidebarSplit
	} else {
		mainLayout = editorAndTabsSplit
	}

	mainContent := container.NewBorder(topToolbar, nil, nil, nil, mainLayout)
	a.mainWindow.SetContent(mainContent)
}

func (a *App) setupShortcuts() {
	shortcutRunR := &desktop.CustomShortcut{KeyName: fyne.KeyR, Modifier: fyne.KeyModifierControl}
	shortcutRunEnter := &desktop.CustomShortcut{KeyName: fyne.KeyReturn, Modifier: fyne.KeyModifierControl}

	a.mainWindow.Canvas().AddShortcut(shortcutRunR, func(shortcut fyne.Shortcut) {
		a.executeCode()
	})
	a.mainWindow.Canvas().AddShortcut(shortcutRunEnter, func(shortcut fyne.Shortcut) {
		a.executeCode()
	})
}

func (a *App) executeCode() {
	a.mu.Lock()
	if a.currentCancel != nil {
		a.mu.Unlock()
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	a.currentCancel = cancel
	a.runBtn.Disable()
	a.stopBtn.Enable()
	a.statusLabel.SetText("Executando...")
	a.mu.Unlock()

	code := a.editor.Text

	go func() {
		startTime := time.Now()
		res, err := a.runner.Execute(ctx, code)
		elapsed := time.Since(startTime)

		fyne.Do(func() {
			a.mu.Lock()
			a.currentCancel = nil
			a.runBtn.Enable()
			a.stopBtn.Disable()
			a.mu.Unlock()

			if err != nil && res == nil {
				a.statusLabel.SetText(fmt.Sprintf("Erro de Inicialização (%v)", elapsed.Truncate(time.Millisecond)))
				a.stderrEntry.SetText(err.Error())
				return
			}

			a.stdoutEntry.SetText(res.Stdout)
			a.stderrEntry.SetText(res.Stderr)
			if res.Err != nil {
				a.stderrEntry.SetText(a.stderrEntry.Text + "\n[Erro de Execução]: " + res.Err.Error())
			}

			a.updateDumpTree(res.DumpData)

			statusText := fmt.Sprintf("Concluído em %v | Engine: %s", elapsed.Truncate(time.Millisecond), res.ModeUsed)
			if res.Err != nil {
				statusText = fmt.Sprintf("Erro em %v | Engine: %s", elapsed.Truncate(time.Millisecond), res.ModeUsed)
			}
			a.statusLabel.SetText(statusText)

			if len(res.DumpData) > 0 {
				a.appTabs.Select(a.dumpTab)
			}
		})
	}()
}

func (a *App) cancelExecution() {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.currentCancel != nil {
		a.currentCancel()
		a.statusLabel.SetText("Execução Cancelada pelo Usuário.")
	}
}

func (a *App) updateDumpTree(nodes []*dumper.DumpNode) {
	store := &treeStore{
		nodes: make(map[string]*dumper.DumpNode),
	}

	for i, node := range nodes {
		id := strconv.Itoa(i)
		store.roots = append(store.roots, id)
		buildTreeNodes(id, node, store.nodes)
	}

	a.treeStore = store
	a.treeWidget.Refresh()
}

func buildTreeNodes(id string, node *dumper.DumpNode, nodes map[string]*dumper.DumpNode) {
	nodes[id] = node
	for i, child := range node.Children {
		childID := id + "." + strconv.Itoa(i)
		buildTreeNodes(childID, child, nodes)
	}
}

func formatNodeText(node *dumper.DumpNode) string {
	var parts []string
	if node.Name != "" {
		parts = append(parts, node.Name+":")
	}
	if node.Type != "" {
		parts = append(parts, "("+node.Type+")")
	}
	if node.ValueStr != "" {
		parts = append(parts, "= "+node.ValueStr)
	}
	return strings.Join(parts, " ")
}
