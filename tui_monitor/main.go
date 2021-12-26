package main

import (
	"fmt"
	"log"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	p2 := widgets.NewParagraph()
	p2.SetRect(0, 4, 40, 9)
	g2 := widgets.NewGauge()
	g2.SetRect(0, 0, 40, 4)

	g := widgets.NewGauge()
	g.SetRect(41, 0, 81, 4)
	p := widgets.NewParagraph()
	p.SetRect(41, 4, 81, 9)

	lc := widgets.NewPlot()
	lc.SetRect(0, 9, 40, 20)
	lc2 := widgets.NewPlot()
	lc2.SetRect(41, 9, 81, 20)

	g.Title = "Irrigation System: Idle"
	g.Percent = 0
	g.BarColor = ui.ColorRed
	g.BorderStyle.Fg = ui.ColorWhite
	g.TitleStyle.Fg = ui.ColorCyan

	p2.Title = "..."
	p2.Text = ``
	p2.TextStyle.Fg = ui.ColorWhite
	p2.BorderStyle.Fg = ui.ColorCyan

	g2.Title = "3-D Printer: Idle"
	g2.Percent = 100
	g2.BarColor = ui.ColorRed
	g2.BorderStyle.Fg = ui.ColorWhite
	g2.TitleStyle.Fg = ui.ColorCyan

	p.Title = "..."
	// p.Text = `Press C to cancel`
	p.TextStyle.Fg = ui.ColorWhite
	p.BorderStyle.Fg = ui.ColorCyan

	lc.Title = "Temperature & Humidity: Blacktent"
	lc.Data = make([][]float64, 1)
	lc.Data[0] = getDataForChart(BT_HOST)
	lc.AxesColor = ui.ColorGreen
	lc.LineColors[0] = ui.ColorYellow
	lc.Marker = widgets.MarkerBraille

	lc2.Title = "Temperature & Humidity: Greenhouse"
	lc2.Data = make([][]float64, 1)
	lc2.Data[0] = getDataForChart(GH_HOST)
	lc2.AxesColor = ui.ColorWhite
	lc2.LineColors[0] = ui.ColorMagenta
	lc2.Marker = widgets.MarkerBraille

	draw := func() {
		runStatus := getPirriRunStatus()
		if !runStatus.IsIdle {
			var remaining int64
			remaining, g.Percent = calcTimeDiff(runStatus)
			g.Title = fmt.Sprintf(`Zone %d is Active (%d)`, runStatus.StationID, g.Percent)
			if runStatus.IsManual {
				p.TextStyle.Fg = ui.ColorRed
				p.Title = fmt.Sprintf(`Unscheduled Zone Run Details`)
			} else {
				p.TextStyle.Fg = ui.ColorGreen
				p.Title = fmt.Sprintf(`Scheduled Run %d`, runStatus.ScheduleID)
			}
			p.Text = fmt.Sprintf(`Scheduled: %t
			Remaining: %d min / %d min (%d%%)
			`,
				!runStatus.IsManual,
				remaining/60-3,
				runStatus.Duration/60,
				g.Percent,
			)
		} else {
			g.Title = fmt.Sprintf(`Irrigation System: Idle`)
			p.Title = fmt.Sprintf(`...`)
			p.Text = fmt.Sprintf(``)
		}

		lc.Data[0] = getDataForChart(BT_HOST)
		lc2.Data[0] = getDataForChart(GH_HOST)

		printData := getPrintStatus()
		if printData.Progress.Completion == 100 || printData.Progress.Completion == 0 {
			g2.Title = fmt.Sprintf(`3-D Printer: Idle`)
			p2.Title = fmt.Sprintf(`...`)
			p2.Text = fmt.Sprintf(``)
		} else {
			g2.Title = fmt.Sprintf(`3-D Printer: %d%%`,
				int(printData.Progress.Completion),
				// int(100*(printData.Progress.PrintTimeLeft/printData.Progress.PrintTime)),
			)
			p2.Title = fmt.Sprintf(printData.Job.File.Name)
			g2.Percent = int(printData.Progress.Completion)
			p2.Text = fmt.Sprintf(`State: %s
		Remaining: %d of %d min
		Filament Len/Vol:  %d, %d
		`,
				printData.State,
				printData.Progress.PrintTimeLeft/60, printData.Progress.PrintTime/60,
				int(printData.Job.Filament.Tool0.Length), int(printData.Job.Filament.Tool0.Volume),
			)
		}

		lc.Title = fmt.Sprintf("Blacktent - Temp: %f", lc.Data[0][len(lc.Data)-1])
		lc2.Title = fmt.Sprintf("Greenhouse - Temp: %f", lc2.Data[0][len(lc2.Data)-1])

		ui.Render(g, lc, lc2, p, g2, p2)
	}

	draw()
	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second * 20).C

	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "c", "<MouseLeft>":
				cancelStationRun()
				draw()
			}
		case <-ticker:
			go draw()
		}
	}
}
