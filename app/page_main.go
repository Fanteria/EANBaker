package app

import (
	"fmt"
	"strings"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Fanteria/EANBaker/core"
)

type MainPage struct {
	file       openFileDialog
	textHeader inputField
	eanHeader  inputField
	submitBtn  widget.Clickable
}

func (m *MainPage) SetFromGenerator(generator *core.Generator) {
	m.textHeader.SetText(generator.TextHeader)
	m.eanHeader.SetText(generator.EanHeader)
}

func (m *MainPage) mainPage(
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
					return material.H4(th, "EANBaker").Layout(gtx)
				}),
				layout.Rigid(inset(layout.Inset{Top: unit.Dp(15)}, m.file.GetWidget(th, message))),
				layout.Rigid(inset(layout.Inset{Top: unit.Dp(15)}, m.textHeader.GetWidget(th))),
				layout.Rigid(inset(layout.Inset{Top: unit.Dp(15)}, m.eanHeader.GetWidget(th))),
				layout.Rigid(inset(layout.Inset{Top: unit.Dp(20)}, func(gtx C) D {
					if m.submitBtn.Clicked(gtx) {
						if m.file.GetFileContent() == nil {
							message.message = "Input file must be set."
							message.messageType = Error
						} else {
							err := func() error {

								table, err := core.TableFromCsv(
									strings.NewReader(*m.file.GetFileContent()),
									rune(generator.CsvComma))
								if err != nil {
									table, err = core.TableFromExcel(
										strings.NewReader(*m.file.GetFileContent()),
										0)
									if err != nil {
										return err
									}
								}

								records, err := core.RecordsFromTable(
									table,
									m.textHeader.GetText(),
									m.eanHeader.GetText())
								if err != nil {
									return err
								}

								pdfFile := "./" + NAME + ".pdf"
								if m.file.GetFileName() != "" {
									pdfFile = core.GeneratePdfPath(m.file.GetFileName())
								}
								if generator.PdfPath != "" {
									pdfFile = generator.PdfPath
								}

								pdf := core.NewPdf()
								pdf.AddPages(records)
								err = pdf.Save(pdfFile)
								if err != nil {
									return err
								} else {
									message.message = fmt.Sprintf("File %s saved.", pdfFile)
									message.messageType = Info
								}
								generator.TextHeader = m.textHeader.GetText()
								generator.EanHeader = m.eanHeader.GetText()
								generator.Save("./." + NAME + ".json")
								setHidden("./." + NAME + ".json")
								return nil
							}()

							if err != nil {
								message.setError(err)
							}
						}
					}
					return material.Button(th, &m.submitBtn, "Submit").Layout(gtx)
				})),
			)
		})
	})
}
