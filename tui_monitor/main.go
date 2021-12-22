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

	g := widgets.NewGauge()
	g.Title = "No Active Zones"
	g.Percent = 100
	g.SetRect(0, 0, 40, 5)
	g.BarColor = ui.ColorRed
	g.BorderStyle.Fg = ui.ColorWhite
	g.TitleStyle.Fg = ui.ColorCyan

	p := widgets.NewParagraph()
	p.Title = "Text Box"
	p.Text = `Press C to cancel`
	p.SetRect(41, 0, 81, 5)
	p.TextStyle.Fg = ui.ColorWhite
	p.BorderStyle.Fg = ui.ColorCyan

	lc := widgets.NewPlot()
	lc.Title = "Temperature & Humidity: Blacktent"
	lc.Data = make([][]float64, 1)
	lc.Data[0] = getDataForChart(BT_HOST)
	lc.SetRect(0, 5, 40, 15)
	lc.AxesColor = ui.ColorGreen
	lc.LineColors[0] = ui.ColorYellow
	lc.Marker = widgets.MarkerBraille

	lc2 := widgets.NewPlot()
	lc2.Title = "Temperature & Humidity: Greenhouse"
	lc2.Data = make([][]float64, 1)
	lc2.Data[0] = getDataForChart(GH_HOST)
	lc2.SetRect(41, 5, 81, 15)
	lc2.AxesColor = ui.ColorWhite
	lc2.LineColors[0] = ui.ColorMagenta
	lc2.Marker = widgets.MarkerBraille

	draw := func() {
		runStatus := getPirriRunStatus()
		if !runStatus.IsIdle {
			var remaining int64
			remaining, g.Percent = calcTimeDiff(runStatus)
			g.Title = fmt.Sprintf(`Zone %d is Active`, runStatus.StationID)
			if runStatus.IsManual {
				p.TextStyle.Fg = ui.ColorRed
				p.Title = fmt.Sprintf(`Unscheduled Zone Run Details`)
			} else {
				p.TextStyle.Fg = ui.ColorGreen
				p.Title = fmt.Sprintf(`Scheduled Run %d`, runStatus.ScheduleID)
			}
			p.Text = fmt.Sprintf(`Scheduled: %t
			Remaining: %d min / %d min
			`,
				!runStatus.IsManual,
				remaining,
				runStatus.Duration/60,
			)
		} else {
			g.Title = fmt.Sprintf(`No Active Irrigation Zones`)
			p.Title = fmt.Sprintf(`...`)
			p.Text = fmt.Sprintf(``)
		}
		lc.Data[0] = getDataForChart(BT_HOST)
		lc2.Data[0] = getDataForChart(GH_HOST)

		lc.Title = fmt.Sprintf("Blacktent - Temp: %f", lc.Data[0][len(lc.Data)-1])
		lc2.Title = fmt.Sprintf("Greenhouse - Temp: %f", lc2.Data[0][len(lc2.Data)-1])
		ui.Render(g, lc, lc2, p)
	}

	tickerCount := 1
	draw()
	tickerCount++
	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second * 15).C

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
			tickerCount++
		}
	}
}
