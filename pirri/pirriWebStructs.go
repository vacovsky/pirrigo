package pirri

//ManualStationRun describes the data required to trigger a "manual" Station activation from the web front end.
type ManualStationRun struct {
	StationID int
	Duration  int
}
