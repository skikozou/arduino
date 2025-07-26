package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type LogUI struct {
	app        *tview.Application
	logView    *tview.TextView
	inputField *tview.InputField
	flex       *tview.Flex

	// 入力制御用
	inputEnabled bool
	inputPrompt  string
	inputChan    chan string
	mu           sync.Mutex
}

func NewLogUI() *LogUI {
	ui := &LogUI{
		app:          tview.NewApplication(),
		inputEnabled: false,
		inputChan:    make(chan string, 1),
	}

	ui.setupUI()
	return ui
}

func (ui *LogUI) setupUI() {
	// 上部ログ表示用のTextView
	ui.logView = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetWordWrap(true).
		SetWrap(true).
		SetChangedFunc(func() {
			ui.app.Draw()
		})
	ui.logView.SetBorder(true)

	// 下部入力用のInputField
	ui.inputField = tview.NewInputField().
		SetLabel("> ").
		SetFieldWidth(0).
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetFieldTextColor(tcell.ColorWhite).
		SetLabelColor(tcell.ColorWhite)
	ui.inputField.SetBorder(true)
	ui.inputField.SetBackgroundColor(tcell.ColorBlack)

	// 初期状態では入力を無効化
	ui.inputField.SetDisabled(true)

	// 入力完了時の処理
	ui.inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			ui.handleInput()
		}
	})

	// レイアウト作成
	ui.flex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(ui.logView, 0, 1, false).
		AddItem(ui.inputField, 3, 1, true)
}

func (ui *LogUI) handleInput() {
	ui.mu.Lock()
	defer ui.mu.Unlock()

	if !ui.inputEnabled {
		return
	}

	text := ui.inputField.GetText()
	if strings.TrimSpace(text) != "" {
		// 入力をログに表示（TextViewは任意のgoroutineから安全に書き込み可能）
		fmt.Fprintf(ui.logView, "[white]%s: %s[white]\n", ui.inputPrompt, text)
		ui.logView.ScrollToEnd()

		// 入力欄をクリア
		ui.inputField.SetText("")

		// 入力を無効化
		ui.inputEnabled = false
		ui.inputField.SetDisabled(true)
		ui.inputField.SetLabel("> ")
		ui.app.SetFocus(ui.logView)

		// チャンネルに送信
		go func() {
			ui.inputChan <- text
		}()
	}
}

func (ui *LogUI) enableInput(prompt string) {
	ui.mu.Lock()
	defer ui.mu.Unlock()

	ui.inputEnabled = true
	ui.inputPrompt = prompt

	ui.inputField.SetDisabled(false)
	ui.inputField.SetLabel("> ")
	ui.app.SetFocus(ui.inputField)
	ui.app.Draw() // Draw()は任意のgoroutineから安全に呼び出し可能
}

func (ui *LogUI) disableInput() {
	ui.mu.Lock()
	defer ui.mu.Unlock()

	ui.inputEnabled = false
	ui.inputField.SetDisabled(true)
	ui.inputField.SetLabel("> ")
	ui.app.SetFocus(ui.logView)
	ui.app.Draw()
}

// ログにメッセージを追加（TextViewは任意のgoroutineから安全）
func (ui *LogUI) Log(message string) {
	// 改行文字を適切に処理
	cleanMessage := strings.ReplaceAll(message, "\r\n", "\n")
	cleanMessage = strings.ReplaceAll(cleanMessage, "\r", "\n")
	fmt.Fprintf(ui.logView, "[yellow]%s[white]\n", cleanMessage)
	ui.logView.ScrollToEnd()
}

// シリアルポートからの生データ用のログメソッド
func (ui *LogUI) LogRaw(data []byte) {
	// バイナリデータを安全に表示
	output := ""
	for _, b := range data {
		if b >= 32 && b <= 126 { // 印刷可能文字
			output += string(b)
		} else if b == '\n' {
			output += "\n"
		} else if b == '\r' {
			// CRは無視（Windows改行対応）
			continue
		} else {
			output += fmt.Sprintf("\\x%02x", b)
		}
	}

	if strings.TrimSpace(output) != "" {
		fmt.Fprintf(ui.logView, "[cyan]RX: %s[white]\n", output)
		ui.logView.ScrollToEnd()
	}
}

// 入力を求める（ブロッキング）
func (ui *LogUI) RequestInput(prompt string) string {
	ui.Log(fmt.Sprintf("[green]%s", prompt))
	ui.enableInput(prompt)

	// 入力を待機
	input := <-ui.inputChan
	return input
}

// UIを開始
func (ui *LogUI) Run() error {
	return ui.app.SetRoot(ui.flex, true).EnableMouse(true).Run()
}

// UIを停止
func (ui *LogUI) Stop() {
	ui.app.Stop()
}
