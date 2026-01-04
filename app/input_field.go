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
	message    *Message
	setValue   func(value string) error
	getValue   func() string
}

// Creates a new input field widget with the specified configuration.
// Initializes the editor with the provided text (if not empty), sets single-line mode,
// and configures the field name and suggestion text.
func NewInputField(name string, suggestion string, message *Message, setValue func(value string) error, getValue func() string) inputField {
	ret := inputField{
		editor:     widget.Editor{SingleLine: true},
		name:       name,
		suggestion: suggestion,
		message:    message,
		setValue:   setValue,
		getValue:   getValue,
	}
	ret.Update()
	return ret
}

// Expect to set valid value
func (i *inputField) Update() {
	text := i.getValue()
	i.editor.SetText(text)
}

// Returns a layout widget for the input field.
// Creates a horizontal layout with a label (field name) and the text editor.
func (i *inputField) GetWidget(th *material.Theme) layout.Widget {
	return func(gtx C) D {
		return layout.UniformInset(unit.Dp(0)).Layout(gtx, func(gtx C) D {
			err := i.setValue(i.editor.Text())
			if err != nil {
				i.message.setError(err)
				i.editor.SetText(i.getValue())
			}
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
