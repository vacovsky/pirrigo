# Pirrigo Fix Plan

## Status: COMPLETE (10/10)

### DONE

- [x] **Bug 1: Schedule edit dialog overwrites dates** — Removed date override in `eventClicked`.
- [x] **Bug 2: Analytics page — dead charts** — Backend chart #3 (`statsActivityPerStationByDOW`) now returns real data. Frontend `loadChartByID` uses the `chart` param instead of hardcoded `4`.
- [x] **Bug 3: Analytics — hardcoded 7-day window** — All stats endpoints accept `?days=N` query param. Default 7.
- [x] **Bug 4: Weather — dead Weather Underground API** — Replaced with Open-Meteo (free, no API key). Added `Latitude`/`Longitude` to Settings.Weather struct.
- [x] **Bug 5: Station edit dialog — dead code** — Removed empty if/else, duplicate `setStationGPIO`, dead `ngAfterViewInit`, unused `formatSliderLabel`.
- [x] **Bug 6: `postStationRun` fire-and-forget** — Now returns `Observable<StationResponse>`, caller subscribes with error handler.
- [x] **Bug 7: Dead `RabbitReactions.go` and `sdj/`** — Deleted.
- [x] **Bug 8: Dead `templates/` directory** — Deleted.
- [x] **Bug 9: `ioutil.ReadAll` deprecated** — Replaced with `io.ReadAll` in weather rewrite.
- [x] **Bug 10: API docs in source** — Stripped ~50 lines of commented API response examples from `apiclient.service.ts`.

## Tests

- `pirri/pirriWebStats_test.go` — 6 table-driven tests for `parseDaysParam` (default, explicit, zero, negative, garbage, one) + smoke test for chart #3.
- `weather/pirriWeather_test.go` — no-coords short-circuit test (returns Error status when lat/lon missing).
- `logging/pirriLogging.go` — fixed panic when `PIRRIGO_LOG_LOCATION` is unset (defaults to stderr).
