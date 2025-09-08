package app

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type inputField struct {
	editor     widget.Editor
	name       string
	suggestion string
}

// Creates a new input field widget with the specified configuration.
// Initializes the editor with the provided text (if not empty), sets single-line mode,
// and configures the field name and suggestion text.
func NewInputField(text string, name string, suggestion string) inputField {
	ret := inputField{
		editor:     widget.Editor{},
		name:       name,
		suggestion: suggestion,
	}

	if text != "" {
		ret.editor.SetText(text)
	}

	ret.editor.SingleLine = true
	return ret
}

// Returns the current text content of the input field.
func (i *inputField) GetText() string {
	return i.editor.Text()
}

// Updates the text content of the input field.
func (i *inputField) SetText(text string) {
	i.editor.SetText(text)
}

// Returns a layout widget for the input field.
// Creates a horizontal layout with a label (field name) and the text editor.
func (i *inputField) GetWidget(th *material.Theme) layout.Widget {
	return func(gtx C) D {
		return layout.UniformInset(unit.Dp(0)).Layout(gtx, func(gtx C) D {
			return layout.Flex{
				Axis:    layout.Horizontal,
				Spacing: layout.SpaceAround,
			}.Layout(gtx,
				layout.Rigid(inset(layout.Inset{Right: unit.Dp(5)}, func(gtx C) D {
					return material.Label(th, 16, i.name+":").Layout(gtx)
				})),
				layout.Rigid(func(gtx C) D {
					editor := material.Editor(th, &i.editor, i.suggestion)
					return editor.Layout(gtx)
				}),
			)
		})
	}
}
