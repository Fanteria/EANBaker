package app

import (
	"fmt"
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

// Sets the message to an error state if the provided error is not nil.
func (m *Message) setError(err error) {
	if err != nil {
		m.message = err.Error()
		m.messageType = Error
	}
}

const NAME string = "EANBaker"

// Starts the GUI application in a separate goroutine.
// Creates a new window and runs the UI event loop.
// Exits the program when the window is closed.
// Returns an error if the GUI fails to start.
func RunGui(logger *core.MultiLogger) error {
	go func() {
		window := new(app.Window)
		err := runUI(window, logger)
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

// Inset wraps a widget with the specified inset padding.
// Returns a new widget function that applies the inset layout
// to the provided widget.
func inset(inset layout.Inset, widget layout.Widget) layout.Widget {
	return func(gtx C) D {
		return inset.Layout(gtx, widget)
	}
}

// Handles the main UI event loop for the application window.
// Manages page switching between main and options pages,
// handles button clicks, and renders the UI based on current state.
// Processes window events until destruction.
func runUI(w *app.Window, log *core.MultiLogger) error {
	var message Message

	messageBtn := widget.Clickable{}
	optionsBtn := widget.Clickable{}
	openOptions := false

	var ops op.Ops
	th := material.NewTheme()

	generator, err := core.LoadGenerator("./." + NAME + ".json", log.Logger)
	if err != nil {
		generator = &core.Generator{TimesEachEAN: 1}
	}

	mainPage := MainPage{
		file:       NewOpenFileDialog("Choose file"),
		textHeader: NewInputField(generator.TextHeader, "Text", "Text column header"),
		eanHeader:  NewInputField(generator.EanHeader, "EAN", "EAN column header"),
		timesHeader:  NewInputField(generator.TimesHeader, "Times", "Times column header (print once if empty)"),
	}

	optsPage := OptsPage{
		csvComma:     NewInputField(string(generator.CsvComma), "Csv sep", "Csv column separator"),
		textHeader:   NewInputField(generator.TextHeader, "Text", "Text column header"),
		eanHeader:    NewInputField(generator.EanHeader, "EAN", "EAN column header"),
		timesHeader:  NewInputField(generator.TimesHeader, "Times", "Times column header (print once if empty)"),
		pdfFile:      NewInputField(generator.PdfPath, "Pdf path", "Static path to generated pdf."),
		timesEachEan: NewInputField(fmt.Sprint(generator.TimesEachEAN), "Times each EAN", "Number of times each EAN code will be printed in the output PDF."),
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
							return mainPage.mainPage(gtx, th, generator, &message, log.Logger)
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
