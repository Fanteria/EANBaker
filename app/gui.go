package app

import (
	"errors"
	"fmt"
	"image/color"
	"log"
	"os"
	"strconv"
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

const (
	PageMain = iota
	PageOptions
	PageInfo
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

func (m *Message) setInfo(s string) {
	m.message = s
	m.messageType = Info
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
	infoBtn := widget.Clickable{}
	actPage := PageMain

	var ops op.Ops
	th := material.NewTheme()

	generator, err := core.LoadGenerator("./."+NAME+".json", log.Logger)
	if err != nil {
		generator = &core.Generator{TimesEachEAN: 1}
	}

	textHeader := NewInputField("Text", "Text column header", &message, func(v string) error {
		generator.TextHeader = v
		return nil
	}, func() string { return generator.TextHeader })

	eanHeader := NewInputField("EAN", "EAN column header", &message, func(v string) error {
		generator.EanHeader = v
		return nil
	}, func() string { return generator.EanHeader })

	timesHeader := NewInputField("Times", "Times column header (print once if empty)", &message, func(v string) error {
		generator.TimesHeader = v
		return nil
	}, func() string { return generator.TimesHeader })

	csvComma := NewInputField("Csv sep", "Csv column separator", &message, func(v string) error {
		if v == "" {
			generator.CsvComma = ','
			return nil
		}
		comma, err := core.CommaFromString(strings.TrimSpace(v))
		if err != nil {
			return err
		}
		generator.CsvComma = comma
		return nil
	}, func() string { return string(generator.CsvComma) })

	pdfFile := NewInputField("Pdf path", "Static path to generated pdf.", &message, func(v string) error {
		generator.PdfPath = v
		return nil
	}, func() string { return generator.PdfPath })

	timesEachEan := NewInputField("Times each EAN", "Number of times each EAN code will be printed in the output PDF.", &message, func(v string) error {
		if v == "" {
			generator.TimesEachEAN = 1
			return nil
		}
		timesEachEan, err := strconv.ParseUint(strings.TrimSpace(v), 10, 0)
		if err != nil {
			return fmt.Errorf("Times each EAN must be positive integer not '%s'.", v)
		} else if timesEachEan <= 0 {
			return errors.New("Times each EAN must be positive integer not zero.")
		} else {
			generator.TimesEachEAN = uint(timesEachEan)
		}
		return nil
	}, func() string { return fmt.Sprint(generator.TimesEachEAN) })

	mainPage := MainPage{
		file:        NewOpenFileDialog("Choose file"),
		textHeader:  &textHeader,
		eanHeader:   &eanHeader,
		pdfFile:     &pdfFile,
		timesHeader: &timesHeader,
	}

	optsPage := OptsPage{
		csvComma:     &csvComma,
		textHeader:   &textHeader,
		eanHeader:    &eanHeader,
		timesHeader:  &timesHeader,
		pdfFile:      &pdfFile,
		timesEachEan: &timesEachEan,
	}

	infoPage := InfoPage {}

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
						var childs []layout.FlexChild
						switch actPage {
						case PageMain:
							childs = mainPage.mainPage(th, generator, &message, log.Logger)
						case PageOptions:
							childs = optsPage.optsPage(th)
						case PageInfo:
							childs = infoPage.infoPage(th, &message, log)
						}
						return layout.Center.Layout(gtx, func(gtx C) D {
							return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx C) D {
								return layout.Flex{
									Axis:    layout.Vertical,
									Spacing: layout.SpaceAround,
								}.Layout(gtx, childs...)
							})
						})
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
							if actPage == PageOptions {
								actPage = PageMain
							} else {
								actPage = PageOptions
							}
						}
						if infoBtn.Clicked(gtx) {
							if actPage == PageInfo {
								actPage = PageMain
							} else {
								actPage = PageInfo
							}
						}
						return layout.NE.Layout(gtx, func(gtx C) D {
							return layout.Inset{Top: unit.Dp(20), Right: unit.Dp(20)}.Layout(gtx, func(gtx C) D {
								return layout.Flex{
									Axis:    layout.Horizontal,
									Spacing: layout.SpaceBetween,
								}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										var buttonIcon string
										if actPage == PageInfo {
											buttonIcon = "×"
										} else {
											buttonIcon = "ℹ"
										}
										button := material.Button(th, &infoBtn, buttonIcon)
										width := gtx.Dp(unit.Dp(40))
										gtx.Constraints.Min.X = width
										gtx.Constraints.Max.X = width
										return button.Layout(gtx)
									}),
									layout.Rigid(func(gtx C) D {
										return layout.Inset{Left: unit.Dp(10)}.Layout(gtx, func(gtx C) D {
											var buttonIcon string
											if actPage == PageOptions {
												buttonIcon = "×"
											} else {
												buttonIcon = "☰"
											}
											button := material.Button(th, &optionsBtn, buttonIcon)
											width := gtx.Dp(unit.Dp(40))
											gtx.Constraints.Min.X = width
											gtx.Constraints.Max.X = width
											return button.Layout(gtx)
										})
									}),
								)
							})
						})
					}),
				)

			})

			e.Frame(gtx.Ops)
		}
	}
}
