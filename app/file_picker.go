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

func NewOpenFileDialog(btnText string) openFileDialog {
	return openFileDialog{
		btnText: btnText,
		fileBtn: widget.Clickable{},
		result:  make(chan openFileResult),
		data:    nil,
	}
}

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

func (o *openFileDialog) GetWidget(th *material.Theme, msg *Message) layout.Widget {
	o.checkResult(msg)
	return func(gtx C) D {
		button := material.Button(th, &o.fileBtn, o.btnText)
		if o.fileBtn.Clicked(gtx) {
			o.openNewDialogWindow()
		}
		return button.Layout(gtx)
	}
}

func (o *openFileDialog) GetFileContent() *string {
	return o.data
}

func (o *openFileDialog) GetFileName() string {
	return o.filename
}
