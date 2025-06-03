package app

import (
	"image/color"
	"log"
	"os"

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

func (m *Message) setError(err error) {
	if err != nil {
		m.message = err.Error()
		m.messageType = Error
	}
}

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

func inset(inset layout.Inset, widget layout.Widget) layout.Widget {
	return func(gtx C) D {
		return inset.Layout(gtx, widget)
	}
}

func runUI(w *app.Window) error {
	var message Message

	messageBtn := widget.Clickable{}
	optionsBtn := widget.Clickable{}
	openOptions := false

	var ops op.Ops
	th := material.NewTheme()

	generator, err := core.LoadGenerator("./." + NAME + ".json")
	if err != nil {
		generator = &core.Generator{}
	}

	mainPage := MainPage{
		file:       NewOpenFileDialog("Choose file"),
		textHeader: NewInputField(generator.TextHeader, "Text", "Text column header"),
		eanHeader:  NewInputField(generator.EanHeader, "EAN", "EAN column header"),
	}

	optsPage := OptsPage{
		csvComma:   NewInputField(string(generator.CsvComma), "Csv sep", "Csv column separator"),
		textHeader: NewInputField(generator.TextHeader, "Text", "Text column header"),
		eanHeader:  NewInputField(generator.EanHeader, "EAN", "EAN column header"),
		pdfFile:    NewInputField(generator.PdfPath, "Pdf path", "Static path to generated pdf."),
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
						if openOptions {
							return optsPage.optsPage(gtx, th, generator, &message)
						} else {
							return mainPage.mainPage(gtx, th, generator, &message)
						}
					}),
					layout.Stacked(func(gtx C) D {
						return layout.Inset{Top: unit.Dp(20), Left: unit.Dp(20)}.Layout(gtx, func(gtx C) D {
							if message.message == "" {
								return layout.Dimensions{}
							}
							if messageBtn.Clicked(gtx) {
								message.message = ""
								return layout.Dimensions{}
							}
							button := material.Button(th, &messageBtn, message.message+" ×")
							switch message.messageType {
							case Error:
								button.Background = color.NRGBA{R: 176, G: 0, B: 32, A: 255}
							case Info:
								button.Background = color.NRGBA{R: 76, G: 175, B: 80, A: 255}
							}
							return button.Layout(gtx)
						})
					}),
					layout.Expanded(func(gtx C) D {
						if optionsBtn.Clicked(gtx) {
							openOptions = !openOptions
							if openOptions {
								optsPage.SetFromGenerator(generator)
							} else {
								mainPage.SetFromGenerator(generator)
							}
						}
						return layout.NE.Layout(gtx, func(gtx C) D {
							return layout.Inset{Top: unit.Dp(20), Right: unit.Dp(20)}.Layout(gtx, func(gtx C) D {
								var buttonIcon string
								if openOptions {
									buttonIcon = "×"
								} else {
									buttonIcon = "☰"
								}
								button := material.Button(th, &optionsBtn, buttonIcon)
								return button.Layout(gtx)
							})
						})
					}),
				)

			})

			e.Frame(gtx.Ops)
		}
	}
}
