package main

type ManualStationRun struct {
	StationID int
	Duration  int
}

type StatsChart struct {
	ReportType int
	Labels     []string
	Series     []string
	Data       [][]int
}
