package app

import (
	"io"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/explorer"
)

type openFileDialog struct {
	btnText  string
	fileBtn  widget.Clickable
	result   chan openFileResult
	data     *string
	filename string
}
type openFileResult struct {
	data     string
	filename string
	err      error
}

// Creates a new file dialog with the specified button text.
// Returns an openFileDialog instance configured with
// the provided button text, an initialized clickable
// widget, and a result channel.
func NewOpenFileDialog(btnText string) openFileDialog {
	return openFileDialog{
		btnText: btnText,
		fileBtn: widget.Clickable{},
		result:  make(chan openFileResult),
		data:    nil,
	}
}

// Opens a new file picker dialog window in a separate goroutine.
// Creates a new app window with an explorer widget to allow
// file selection. Sends the result (file content and name or error)
// through the result channel.
func (o *openFileDialog) openNewDialogWindow() {
	o.data = nil
	go func() {
		window := new(app.Window)
		picker := explorer.NewExplorer(window)
		file, err := picker.ChooseFile()
		if err != nil {
			o.result <- openFileResult{err: err}
		}
		// If dialog is closed, just do nothing.
		if file == nil {
			return
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			o.result <- openFileResult{err: err}
		}
		if f, ok := file.(*os.File); ok {
			o.result <- openFileResult{data: string(data), filename: f.Name()}
		} else {
			o.result <- openFileResult{data: string(data)}
		}
	}()
}

// Checks for results from the file picker dialog without blocking.
// If a result is available, it updates the message with success
// or error information and stores the file data and filename
// in the dialog instance.
func (o *openFileDialog) checkResult(msg *Message) {
	select {
	case res := <-o.result:
		if res.err != nil {
			msg.message = res.err.Error()
			msg.messageType = Error
		} else {
			msg.message = "File loaded"
			msg.messageType = Info
		}
		o.data = &res.data
		o.filename = res.filename
	default:
	}
}

// Returns a layout widget for the file picker button.
// The widget displays a button that triggers the file
// dialog when clicked. It also checks for dialog results
// and updates the provided message accordingly.
func (o *openFileDialog) GetWidget(th *material.Theme, msg *Message) layout.Widget {
	o.checkResult(msg)
	return func(gtx C) D {
		btnText := o.btnText
		if o.GetFileName() != "" {
			btnText += ": " + o.GetFileName()
		}
		button := material.Button(th, &o.fileBtn, btnText)
		if o.fileBtn.Clicked(gtx) {
			o.openNewDialogWindow()
		}
		return button.Layout(gtx)
	}
}

// Returns a pointer to the loaded file content as a string.
// Returns nil if no file has been successfully loaded.
func (o *openFileDialog) GetFileContent() *string {
	return o.data
}

// Returns an empty string if no file has been loaded
// or if the file source doesn't provide a name.
func (o *openFileDialog) GetFileName() string {
	return o.filename
}

func (o *openFileDialog) Reset() {
	o.data = nil
	o.filename = ""
}
