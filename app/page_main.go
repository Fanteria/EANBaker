package app

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Fanteria/EANBaker/core"
)

type MainPage struct {
	file        openFileDialog
	textHeader  inputField
	eanHeader   inputField
	timesHeader inputField
	submitBtn   widget.Clickable
}

// Updates the main page input fields with values from the generator.
// Synchronizes the text header and EAN header fields with the generator's current settings.
func (m *MainPage) SetFromGenerator(generator *core.Generator) {
	m.textHeader.SetText(generator.TextHeader)
	m.eanHeader.SetText(generator.EanHeader)
	m.timesHeader.SetText(generator.EanHeader)
}

// Renders the main page layout with file selection, input fields, and submit functionality.
// Handles file processing (CSV/Excel), record extraction, PDF generation, and configuration saving.
// Returns the dimensions of the rendered layout.
func (m *MainPage) mainPage(
	gtx C,
	th *material.Theme,
	generator *core.Generator,
	message *Message,
	log *slog.Logger,
) D {
	return layout.Center.Layout(gtx, func(gtx C) D {
		return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx C) D {
			return layout.Flex{
				Axis:    layout.Vertical,
				Spacing: layout.SpaceAround,
			}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return material.H4(th, "EANBaker").Layout(gtx)
				}),
				layout.Rigid(inset(layout.Inset{Top: unit.Dp(15)}, m.file.GetWidget(th, message))),
				layout.Rigid(inset(layout.Inset{Top: unit.Dp(15)}, m.textHeader.GetWidget(th))),
				layout.Rigid(inset(layout.Inset{Top: unit.Dp(15)}, m.eanHeader.GetWidget(th))),
				layout.Rigid(inset(layout.Inset{Top: unit.Dp(15)}, m.timesHeader.GetWidget(th))),
				layout.Rigid(inset(layout.Inset{Top: unit.Dp(20)}, func(gtx C) D {
					if m.submitBtn.Clicked(gtx) {
						// Run button clicked function, if return error set it.
						message.setError(func() error {
							if m.file.GetFileContent() == nil {
								return errors.New("Input file must be set.")
							}
							// Set generator values
							generator.PdfPath = "./" + NAME + ".pdf"
							if m.file.GetFileName() != "" {
								generator.PdfPath = core.GeneratePdfPath(m.file.GetFileName())
							}
							generator.TextHeader = m.textHeader.GetText()
							generator.EanHeader = m.eanHeader.GetText()
							generator.TimesHeader = m.timesHeader.GetText()
							// generator.TimesEachEAN
							err := generator.Generate(
								m.file.GetFileName(),
								strings.NewReader(*m.file.GetFileContent()),
								log)
							if err != nil {
								return err
							}
							message.message = fmt.Sprintf("File %s saved.", generator.PdfPath)
							message.messageType = Info
							setHidden("./." + NAME + ".json")
							return nil
						}())
					}
					return material.Button(th, &m.submitBtn, "Submit").Layout(gtx)
				})),
			)
		})
	})
}
