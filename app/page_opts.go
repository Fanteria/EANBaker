package app

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Fanteria/EANBaker/core"
)

type OptsPage struct {
	csvComma     *inputField
	textHeader   *inputField
	eanHeader    *inputField
	timesHeader  *inputField
	pdfFile      *inputField
	saveBtn      widget.Clickable
	timesEachEan *inputField
}

// Renders the options page layout with configuration input fields and save functionality.
// Handles validation and updating of generator settings including CSV separator, headers,
// PDF path, and barcode repetition count. Returns the dimensions of the rendered layout.
func (o *OptsPage) optsPage(
	th *material.Theme,
	generator *core.Generator,
	message *Message,
) []layout.FlexChild {
	return []layout.FlexChild{
		layout.Rigid(func(gtx C) D {
			return material.H4(th, "Options").Layout(gtx)
		}),
		layout.Rigid(inset(layout.Inset{Top: unit.Dp(15)}, o.timesEachEan.GetWidget(th))),
		layout.Rigid(inset(layout.Inset{Top: unit.Dp(15)}, o.csvComma.GetWidget(th))),
		layout.Rigid(inset(layout.Inset{Top: unit.Dp(15)}, o.textHeader.GetWidget(th))),
		layout.Rigid(inset(layout.Inset{Top: unit.Dp(15)}, o.eanHeader.GetWidget(th))),
		layout.Rigid(inset(layout.Inset{Top: unit.Dp(15)}, o.timesHeader.GetWidget(th))),
		layout.Rigid(inset(layout.Inset{Top: unit.Dp(15)}, o.pdfFile.GetWidget(th))),
	}
}
