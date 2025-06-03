package app

import (
	"strings"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Fanteria/EANBaker/core"
)

type OptsPage struct {
	// file       openFileDialog
	csvComma   inputField
	textHeader inputField
	eanHeader  inputField
	pdfFile    inputField
	saveBtn    widget.Clickable
}

func (o *OptsPage) SetFromGenerator(generator *core.Generator) {
	if generator.CsvComma == 0 {
		o.csvComma.SetText("")
	} else {
		o.csvComma.SetText(string(generator.CsvComma))
	}
	o.textHeader.SetText(generator.TextHeader)
	o.eanHeader.SetText(generator.EanHeader)
	o.pdfFile.SetText(generator.PdfPath)
}

func (o *OptsPage) optsPage(
	gtx C,
	th *material.Theme,
	generator *core.Generator,
	message *Message,
) D {
	return layout.Center.Layout(gtx, func(gtx C) D {
		return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx C) D {
			return layout.Flex{
				Axis:    layout.Vertical,
				Spacing: layout.SpaceAround,
			}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return material.H4(th, "Options").Layout(gtx)
				}),
				layout.Rigid(inset(layout.Inset{Top: unit.Dp(15)}, o.csvComma.GetWidget(th))),
				layout.Rigid(inset(layout.Inset{Top: unit.Dp(15)}, o.textHeader.GetWidget(th))),
				layout.Rigid(inset(layout.Inset{Top: unit.Dp(15)}, o.eanHeader.GetWidget(th))),
				layout.Rigid(inset(layout.Inset{Top: unit.Dp(15)}, o.pdfFile.GetWidget(th))),
				layout.Rigid(inset(layout.Inset{Top: unit.Dp(20)}, func(gtx C) D {
					if o.saveBtn.Clicked(gtx) {
						commaStr := o.csvComma.GetText()
						commaStr = strings.TrimSpace(commaStr)
						comma, err := core.CommaFromString(commaStr)
						if err != nil {
							message.setError(err)
						}
						generator.CsvComma = comma

						generator.TextHeader = o.textHeader.GetText()
						generator.EanHeader = o.eanHeader.GetText()
						generator.PdfPath = o.pdfFile.GetText()
					}
					return material.Button(th, &o.saveBtn, "Save").Layout(gtx)
				})),
			)
		})
	})
}
