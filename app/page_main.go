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
	textHeader  *inputField
	eanHeader   *inputField
	timesHeader *inputField
	pdfFile     *inputField
	submitBtn   widget.Clickable
}

// Renders the main page layout with file selection, input fields, and submit functionality.
// Handles file processing (CSV/Excel), record extraction, PDF generation, and configuration saving.
// Returns the dimensions of the rendered layout.
func (m *MainPage) mainPage(
	th *material.Theme,
	generator *core.Generator,
	message *Message,
	log *slog.Logger,
) []layout.FlexChild {
	pdfFileName := m.file.GetFileName()
	if pdfFileName != "" && generator.PdfPath == "" {
		generator.PdfPath = core.GeneratePdfPath(m.file.GetFileName())
		m.pdfFile.Update()
	}
	return []layout.FlexChild{
		layout.Rigid(func(gtx C) D {
			return material.H4(th, "EANBaker").Layout(gtx)
		}),
		layout.Rigid(inset(layout.Inset{Top: unit.Dp(15)}, m.file.GetWidget(th, message))),
		layout.Rigid(inset(layout.Inset{Top: unit.Dp(15)}, m.textHeader.GetWidget(th))),
		layout.Rigid(inset(layout.Inset{Top: unit.Dp(15)}, m.eanHeader.GetWidget(th))),
		layout.Rigid(inset(layout.Inset{Top: unit.Dp(15)}, m.timesHeader.GetWidget(th))),
		layout.Rigid(inset(layout.Inset{Top: unit.Dp(15)}, m.pdfFile.GetWidget(th))),
		layout.Rigid(inset(layout.Inset{Top: unit.Dp(20)}, func(gtx C) D {
			if m.submitBtn.Clicked(gtx) {
				// Run button clicked function, if return error set it.
				message.setError(func() error {
					log.Info("Try to generate", "generator", generator)
					if m.file.GetFileContent() == nil {
						return errors.New("Input file must be set.")
					}
					// Set generator values
					generator.CsvPath = m.file.GetFileName()

					if generator.PdfPath == "" {
						generator.PdfPath = "./" + NAME + ".pdf"
					}

					err := generator.Generate(
						m.file.GetFileName(),
						strings.NewReader(*m.file.GetFileContent()),
						log)
					if err != nil {
						return err
					}
					message.message = fmt.Sprintf("File %s saved.", generator.PdfPath)
					message.messageType = Info
					log.Info("File generated", "generator", generator)
					setHidden("./." + NAME + ".json")
					m.file.Reset()
					generator.PdfPath = ""
					m.pdfFile.Update()
					return nil
				}())
			}
			return material.Button(th, &m.submitBtn, "Submit").Layout(gtx)
		})),
	}
}
