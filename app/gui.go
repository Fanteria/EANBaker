package app

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"strings"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Fanteria/EANBaker/core"
)

type MessageType int

const (
	Info MessageType = iota
	Error
)

type Message struct {
	message     string
	messageType MessageType
}

// Globals
var message Message

const NAME string = "EANBaker"

func RunGui() error {
	go func() {
		window := new(app.Window)
		err := runUI(window)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
	return nil
}

type (
	C = layout.Context
	D = layout.Dimensions
)

type MainPage struct {
	file       openFileDialog
	textHeader inputField
	eanHeader  inputField
	submitBtn  widget.Clickable
}

func inset(inset layout.Inset, widget layout.Widget) layout.Widget {
	return func(gtx C) D {
		return inset.Layout(gtx, widget)
	}
}

func (m *MainPage) mainPage(gtx C, th *material.Theme, generator *core.Generator) D {
	return layout.Center.Layout(gtx, func(gtx C) D {
		return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx C) D {
			return layout.Flex{
				Axis:    layout.Vertical,
				Spacing: layout.SpaceAround,
			}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return material.H4(th, "EANBaker").Layout(gtx)
				}),
				layout.Rigid(inset(layout.Inset{Top: unit.Dp(15)}, m.file.GetWidget(th, &message))),
				layout.Rigid(inset(layout.Inset{Top: unit.Dp(15)}, m.textHeader.GetWidget(th))),
				layout.Rigid(inset(layout.Inset{Top: unit.Dp(15)}, m.eanHeader.GetWidget(th))),
				layout.Rigid(inset(layout.Inset{Top: unit.Dp(20)}, func(gtx C) D {
					if m.submitBtn.Clicked(gtx) {
						if m.file.GetFileContent() == nil {
							message.message = "Input file must be set."
							message.messageType = Error
						} else {
							records, err := core.RecordsFromCsv(
								strings.NewReader(*m.file.GetFileContent()),
								m.textHeader.GetText(),
								m.eanHeader.GetText())

							pdfFile := "./" + NAME + ".pdf"
							if m.file.GetFileName() != "" {
								pdfFile = core.GeneratePdfPath(m.file.GetFileName())
							}
							if generator.PdfPath != "" {
								pdfFile = generator.PdfPath
							}
							if err != nil {
								message.message = err.Error()
								message.messageType = Error
							} else {
								pdf := core.NewPdf()
								pdf.AddPages(records)
								err = pdf.Save(pdfFile)
								if err != nil {
									message.message = err.Error()
									message.messageType = Error
								} else {
									message.message = fmt.Sprintf("File %s saved.", pdfFile)
									message.messageType = Info
								}
								generator.TextHeader = m.textHeader.GetText()
								generator.EanHeader = m.eanHeader.GetText()
								generator.Save("./." + NAME + ".json")
								setHidden("./." + NAME + ".json")
							}
						}
					}
					return material.Button(th, &m.submitBtn, "Submit").Layout(gtx)
				})),
			)
		})
	})
}

func runUI(w *app.Window) error {
	messageBtn := widget.Clickable{}

	var ops op.Ops
	th := material.NewTheme()

	generator, err := core.LoadGenerator("./." + NAME + ".json")
	if err != nil {
		generator = &core.Generator{}
	}
	textHeader := ""
	eanHeader := ""
	if err == nil {
		textHeader = generator.TextHeader
		eanHeader = generator.EanHeader
	}

	mainPage := MainPage{
		file:       NewOpenFileDialog("Choose file"),
		textHeader: NewInputField(textHeader, "Text", "Text column header"),
		eanHeader:  NewInputField(eanHeader, "EAN", "EAN column header"),
	}

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// Layout the UI
			layout.Center.Layout(gtx, func(gtx C) D {

				return layout.Stack{}.Layout(gtx,

					layout.Expanded(func(gtx C) D {
						return layout.Dimensions{Size: gtx.Constraints.Max}

					}),
					layout.Expanded(func(gtx C) D {
						return mainPage.mainPage(gtx, th, generator)
					}),
					layout.Stacked(func(gtx C) D {
						if message.message == "" {
							return layout.Dimensions{}
						}
						if messageBtn.Clicked(gtx) {
							message.message = ""
							return layout.Dimensions{}
						}
						return layout.Inset{Top: unit.Dp(20), Left: unit.Dp(20)}.Layout(gtx, func(gtx C) D {
							button := material.Button(th, &messageBtn, message.message)
							switch message.messageType {
							case Error:
								button.Background = color.NRGBA{R: 176, G: 0, B: 32, A: 255}
							case Info:
								button.Background = color.NRGBA{R: 76, G: 175, B: 80, A: 255}
							}
							return button.Layout(gtx)
						})
					}),
				)

			})

			e.Frame(gtx.Ops)
		}
	}
}
