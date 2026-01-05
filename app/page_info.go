package app

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Fanteria/EANBaker/core"
	"github.com/Fanteria/EANBaker/values"
)

type InfoPage struct {
	saveBtn widget.Clickable
}

func (i *InfoPage) infoPage(th *material.Theme, message *Message, log *core.MultiLogger) []layout.FlexChild {
	return []layout.FlexChild{
		layout.Rigid(func(gtx C) D {
			return material.H4(th, "Info").Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return material.Label(th, 16, "Version: "+values.Version).Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return material.Label(th, 16, "Commit: "+values.Commit).Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return material.Label(th, 16, "Build date: "+values.Date).Layout(gtx)
		}),
		layout.Rigid(inset(layout.Inset{Top: unit.Dp(20)}, func(gtx C) D {
			if i.saveBtn.Clicked(gtx) {
				path := "EANBaker.log"
				err := log.SaveToFile(path)
				if err != nil {
					message.setError(err)
				} else {
					message.setInfo(fmt.Sprintf("Log file '%s' saved", path))
				}
			}
			return material.Button(th, &i.saveBtn, "Save log file").Layout(gtx)
		})),
	}
}
