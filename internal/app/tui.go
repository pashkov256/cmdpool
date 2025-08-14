package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/pashkov256/cmdpool/internal/config"
	"github.com/pashkov256/cmdpool/internal/executor"
	"github.com/rivo/tview"
)

// TUI represents the terminal user interface
type TUI struct {
	app           *tview.Application
	executor      *executor.Executor
	config        *config.Config
	mainLayout    *tview.Flex
	commandPanels []*CommandPanel
	statusBar     *tview.TextView
	helpBar       *tview.TextView
	selectedPanel int
}

// CommandPanel represents a single command display panel
type CommandPanel struct {
	*tview.Box
	command  *executor.Command
	output   *tview.TextView
	status   *tview.TextView
	title    *tview.TextView
	expanded bool
}

// NewTUI creates a new TUI instance
func NewTUI() *TUI {
	tui := &TUI{
		app:           tview.NewApplication(),
		executor:      executor.NewExecutor(),
		commandPanels: make([]*CommandPanel, 0),
		selectedPanel: 0,
	}

	tui.setupUI()
	tui.setupKeyBindings()
	tui.setupUpdateLoop()

	return tui
}

// setupUI initializes the user interface
func (tui *TUI) setupUI() {
	// Create main layout
	tui.mainLayout = tview.NewFlex().SetDirection(tview.FlexRow)

	// Create command panels area
	panelsArea := tview.NewFlex().SetDirection(tview.FlexColumn)
	tui.mainLayout.AddItem(panelsArea, 0, 1, true)

	// Create status bar
	tui.statusBar = tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("cmdpool - Ready").
		SetTextColor(tcell.ColorYellow)

	// Create help bar
	tui.helpBar = tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("↑↓: Navigate | Enter: Expand | r: Restart | s: Stop | +: Add | q: Quit").
		SetTextColor(tcell.ColorGray)

	// Add status and help bars
	tui.mainLayout.AddItem(tui.statusBar, 1, 0, false)
	tui.mainLayout.AddItem(tui.helpBar, 1, 0, false)

	// Set root
	tui.app.SetRoot(tui.mainLayout, true)
}

// setupKeyBindings sets up keyboard shortcuts
func (ui *TUI) setupKeyBindings() {
	ui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyUp:
			ui.selectPreviousPanel()
			return nil
		case tcell.KeyDown:
			ui.selectNextPanel()
			return nil
		case tcell.KeyLeft:
			ui.selectPreviousPanel()
			return nil
		case tcell.KeyRight:
			ui.selectNextPanel()
			return nil
		case tcell.KeyEnter:
			ui.expandSelectedPanel()
			return nil
		case tcell.KeyRune:
			switch event.Rune() {
			case 'r':
				ui.restartSelectedCommand()
				return nil
			case 's':
				ui.stopSelectedCommand()
				return nil
			case '+':
				ui.addNewCommand()
				return nil
			case 'q':
				ui.quit()
				return nil
			}
		}
		return event
	})
}

// setupUpdateLoop starts the UI update loop
func (ui *TUI) setupUpdateLoop() {
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for range ticker.C {
			ui.app.QueueUpdateDraw(func() {
				ui.updateUI()
			})
		}
	}()
}

// updateUI updates the interface elements
func (ui *TUI) updateUI() {
	// Update status bar
	commands := ui.executor.GetCommands()
	running := 0
	done := 0
	failed := 0

	for _, cmd := range commands {
		switch cmd.Status {
		case executor.StatusRunning:
			running++
		case executor.StatusDone:
			done++
		case executor.StatusFailed:
			failed++
		}
	}

	statusText := fmt.Sprintf("cmdpool - Running: %d | Done: %d | Failed: %d", running, done, failed)
	ui.statusBar.SetText(statusText)

	// Update command panels
	for _, panel := range ui.commandPanels {
		panel.updateDisplay()
	}
}

// selectNextPanel selects the next panel
func (ui *TUI) selectNextPanel() {
	if len(ui.commandPanels) == 0 {
		return
	}

	ui.selectedPanel = (ui.selectedPanel + 1) % len(ui.commandPanels)
	ui.updatePanelSelection()
}

// selectPreviousPanel selects the previous panel
func (ui *TUI) selectPreviousPanel() {
	if len(ui.commandPanels) == 0 {
		return
	}

	ui.selectedPanel = (ui.selectedPanel - 1 + len(ui.commandPanels)) % len(ui.commandPanels)
	ui.updatePanelSelection()
}

// updatePanelSelection updates the visual selection
func (ui *TUI) updatePanelSelection() {
	for i, panel := range ui.commandPanels {
		if i == ui.selectedPanel {
			panel.SetBorderColor(tcell.ColorYellow)
			panel.SetBorder(true)
		} else {
			panel.SetBorderColor(tcell.ColorGray)
			panel.SetBorder(true)
		}
	}
}

// expandSelectedPanel expands the selected panel to full screen
func (ui *TUI) expandSelectedPanel() {
	if len(ui.commandPanels) == 0 || ui.selectedPanel >= len(ui.commandPanels) {
		return
	}

	panel := ui.commandPanels[ui.selectedPanel]
	panel.expand()
}

// restartSelectedCommand restarts the selected command
func (ui *TUI) restartSelectedCommand() {
	if len(ui.commandPanels) == 0 || ui.selectedPanel >= len(ui.commandPanels) {
		return
	}

	panel := ui.commandPanels[ui.selectedPanel]
	if err := ui.executor.RestartCommand(panel.command.ID); err != nil {
		// Show error in status bar
		ui.statusBar.SetText(fmt.Sprintf("Error restarting command: %v", err))
		ui.statusBar.SetTextColor(tcell.ColorRed)
	}
}

// stopSelectedCommand stops the selected command
func (ui *TUI) stopSelectedCommand() {
	if len(ui.commandPanels) == 0 || ui.selectedPanel >= len(ui.commandPanels) {
		return
	}

	panel := ui.commandPanels[ui.selectedPanel]
	if err := ui.executor.StopCommand(panel.command.ID); err != nil {
		// Show error in status bar
		ui.statusBar.SetText(fmt.Sprintf("Error stopping command: %v", err))
		ui.statusBar.SetTextColor(tcell.ColorRed)
	}
}

// addNewCommand adds a new command via input dialog
func (ui *TUI) addNewCommand() {
	// Create input dialog
	inputField := tview.NewInputField().
		SetLabel("Command: ").
		SetFieldWidth(50)

	modal := tview.NewModal().
		SetText("Enter command to execute:").
		AddButtons([]string{"Run", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonIndex == 0 { // Run button
				command := inputField.GetText()
				if command != "" {
					ui.executor.RunCommands([]string{command})
					ui.app.SetRoot(ui.mainLayout, true)
				}
			} else {
				ui.app.SetRoot(ui.mainLayout, true)
			}
		})

	// Show modal
	ui.app.SetRoot(modal, true)
}

// quit exits the application
func (ui *TUI) quit() {
	ui.executor.Stop()
	ui.app.Stop()
}

// AddCommand adds a new command panel
func (ui *TUI) AddCommand(command *executor.Command) {
	panel := NewCommandPanel(command)
	ui.commandPanels = append(ui.commandPanels, panel)

	// Add to panels area
	panelsArea := ui.mainLayout.GetItem(0).(*tview.Flex)
	panelsArea.AddItem(panel, 0, 1, false)

	// Update selection
	ui.updatePanelSelection()
}

// Run starts the TUI
func (ui *TUI) Run() error {
	return ui.app.Run()
}

// NewCommandPanel creates a new command panel
func NewCommandPanel(command *executor.Command) *CommandPanel {
	panel := &CommandPanel{
		Box:      tview.NewBox().SetBorder(true),
		command:  command,
		expanded: false,
	}

	// Create title
	panel.title = tview.NewTextView().
		SetText(command.Name).
		SetTextAlign(tview.AlignCenter).
		SetTextColor(tcell.ColorWhite)

	// Create status
	panel.status = tview.NewTextView().
		SetText(string(command.Status)).
		SetTextAlign(tview.AlignRight).
		SetTextColor(tcell.ColorYellow)

	// Create output
	panel.output = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetTextColor(tcell.ColorGreen)

	// Set up layout
	panel.SetBorder(true)
	panel.SetTitle(fmt.Sprintf(" %s ", command.Name))

	return panel
}

// updateDisplay updates the panel display
func (panel *CommandPanel) updateDisplay() {
	// Update status
	statusText := string(panel.command.Status)
	statusColor := tcell.ColorWhite

	switch panel.command.Status {
	case executor.StatusRunning:
		statusText = "🟢 Running"
		statusColor = tcell.ColorGreen
	case executor.StatusDone:
		statusText = "✅ Done"
		statusColor = tcell.ColorGreen
	case executor.StatusFailed:
		statusText = "🔴 Failed"
		statusColor = tcell.ColorRed
	case executor.StatusStopped:
		statusText = "⏹️ Stopped"
		statusColor = tcell.ColorYellow
	}

	panel.status.SetText(statusText)
	panel.status.SetTextColor(statusColor)

	// Update output - use public methods instead of accessing private fields
	output := strings.Join(panel.command.GetOutput(), "\n")
	panel.output.SetText(output)
}

// expand expands the panel to full screen
func (panel *CommandPanel) expand() {
	// Create full screen view
	fullScreen := tview.NewFlex().SetDirection(tview.FlexRow)

	// Header with title and status
	header := tview.NewFlex().SetDirection(tview.FlexColumn)
	header.AddItem(panel.title, 0, 1, false)
	header.AddItem(panel.status, 0, 1, false)

	// Output area
	outputArea := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetTextColor(tcell.ColorGreen)

	// Set output text
	output := strings.Join(panel.command.GetOutput(), "\n")
	outputArea.SetText(output)

	// Add components
	fullScreen.AddItem(header, 3, 0, false)
	fullScreen.AddItem(outputArea, 0, 1, true)

	// Add close button
	closeBtn := tview.NewButton("Close (ESC)").SetSelectedFunc(func() {
		// Return to main view
		// This would need to be implemented with proper navigation
	})

	fullScreen.AddItem(closeBtn, 1, 0, false)

	// Show full screen
	// This would need proper navigation implementation
}

// RunTUI starts the TUI application
func RunTUI() error {
	tui := NewTUI()
	return tui.Run()
}
