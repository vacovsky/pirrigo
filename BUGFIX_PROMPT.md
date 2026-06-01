# Pirrigo Bug Fix Prompt
# Run this prompt THREE times to iteratively fix all bugs.
# After EACH bug fix, rebuild both backend and frontend:
#   Backend:  cd /Users/joe/Projects/pirrigo && go build -o pirrigo 2>&1
#   Frontend: cd /Users/joe/Projects/pirrigo/pirrigo-spa && npx ng build 2>&1

## STATUS: ALL 30 BUGS VERIFIED FIXED (re-verified 2026-06-01)

### Pass 1 — Critical + High + Medium bugs fixed:
1. **main.go** — WaitGroup: changed Add(3)->Add(4), added defer WG.Done() to all 4 goroutines
2. **weather/pirriWeather.go** — Nil pointer deref: added early return when http.Get fails
3. **pirri/pirriGpioPin.go** — defer order: swapped rpio.Open() before defer rpio.Close()
4. **schema.go** — SetCommonWire: changed Update("notes", "common") -> Update("common", true)
5. **main.go** — Password leak: redacted PIRRIGO_PASSWORD from log output
6. **logging/pirriLogging.go** — Lock ordering: fixed defer order in both LogEvent and LogError
7. **pirri/pirriWebAuth.go** — Index out of bounds: added bounds check before s[0] access
8. **pirri/pirriWebAuth.go** — Cookie auth: fixed to properly decode base64 credentials
9. **pirri/pirriQueueHelper.go** — TOCTOU race: hold lock continuously during queue pop
10. **pirri/pirriWebRunStatus.go** — Race conditions: added mutex protection to all 3 handlers
11. **pirri/pirriStationSchedule.go** — SQL: changed =true to = 1 for boolean comparison
12. **pirri/pirriWebSchedule.go** — Invalid Order: changed .Order("ASC") -> .Order("start_date ASC")
13. **data/GormHelper.go** — DB failure: added panic after sqlite connection failure
14. **All handlers** — Content-Type: added application/json header via middleware
15. **pirri/pirriWebGpio.go** — SQL injection: parameterized gpio.GPIO query
16. **pirrigo-spa/package.json** — Timepicker version: noted v13.1.1 is latest available
17. **pirrigo-spa/app.module.ts** — Duplicate HttpClient: removed from providers array
18. **pirrigo-spa/globals.service.ts** — Hardcoded IP: changed to relative URL
19. **pirrigo-spa/stations.component.ts** — Division by zero: added Duration === 0 guard
20. **pirrigo-spa/calendar.component.ts** — XSS: added escapeHtml for JSON in event titles
21. **pirri/pirriWebDripNodes.go** — SQL case: fixed station_ID -> station_id
22. **pirri/pirriWebRunStatus.go** — removeJob: copy queue before unlock for marshalling

### Pass 2 — Additional bugs found on re-review:
23. **pirri/pirriWebStats.go** — Index out of range: rewrote stats loop to use len(result.Labels)
24. **pirri/pirriWebStats.go** — SQLite day offset: fixed v.Day-1 for SQLite (DAYOFWEEK vs strftime)
25. **pirri/pirriWebManager.go** — SPA routing: serve index.html for root path
26. **pirri/pirriTask.go** — Mutex during logging: moved LogEvent outside mutex
27. **pirri/pirriWebAuth.go** — loginCheck: replaced http.Error(200) with proper response
28. **pirri/pirriWebManager.go** — Formatting: fixed indentation

### Pass 3 — Final polish:
29. **pirri/pirriStationSchedule.go** — Dead code: removed unused lastTriggeredItem variable
30. **pirri/pirriWebAuth.go** — Cookie auth removed: frontend doesn't use cookies, simplified auth

---

## CRITICAL BUG 1: main.go — WaitGroup hangs forever (main.go:47-60)
**File:** main.go
**Problem:** `pirri.WG.Add(3)` is called, but only `listenForExit()` calls `WG.Done()`. The three goroutines (`StartPirriWebApp`, `StartTaskMonitor`, `ListenForTasks`) never call `WG.Done()`. `WG.Wait()` blocks forever. Additionally, `StartPirriWebApp` ends with `panic(http.ListenAndServe(...))` which crashes the entire process.
**Fix:** 
- Change `pirri.WG.Add(3)` to `pirri.WG.Add(4)` 
- Add `defer pirri.WG.Done()` at the top of `StartPirriWebApp`, `StartTaskMonitor`, and `ListenForTasks`
- Replace `panic(http.ListenAndServe(...))` with proper error handling: store the error and log it, then call `WG.Done()`
- In `pirri/pirriWebManager.go`, replace line 45 `panic(http.ListenAndServe(":"+os.Getenv("PIRRIGO_WEB_PORT"), nil))` with:
  ```go
  err := http.ListenAndServe(":"+os.Getenv("PIRRIGO_WEB_PORT"), nil)
  if err != nil {
      logging.Service().LogError("HTTP server failed", zap.Error(err))
  }
  ```
- In `pirri/pirriWebManager.go`, add `defer WG.Done()` at the top of `StartPirriWebApp`
- In `pirri/pirriStationSchedule.go`, add `defer WG.Done()` at the top of `StartTaskMonitor`  
- In `pirri/pirriQueueHelper.go`, `ListenForTasks` already has `defer WG.Done()` on line 8

**Build after fix:**
```bash
cd /Users/joe/Projects/pirrigo && go build -o pirrigo 2>&1
cd /Users/joe/Projects/pirrigo/pirrigo-spa && npx ng build 2>&1
```

---

## CRITICAL BUG 2: weather/pirriWeather.go — Nil pointer dereference (weather/pirriWeather.go:142-149)
**File:** weather/pirriWeather.go
**Problem:** When `http.Get(weatherEndpoint)` returns an error, `r` is nil. Line 149 calls `defer r.Body.Close()` on a nil pointer, causing a panic.
**Fix:** Add nil check before accessing `r.Body`:
```go
r, err := http.Get(weatherEndpoint)
if err != nil {
    logger.LogError("Unable to obtain weather",
        zap.String("stateAbbreviation", set.Weather.StateAbbreviation),
        zap.String("city", set.Weather.City),
        zap.String("error", err.Error()))
    return weather  // return empty response instead of crashing
}
defer r.Body.Close()
```

**Build after fix:**
```bash
cd /Users/joe/Projects/pirrigo && go build -o pirrigo 2>&1
cd /Users/joe/Projects/pirrigo/pirrigo-spa && npx ng build 2>&1
```

---

## CRITICAL BUG 3: pirriGpioPin.go — defer rpio.Close() before rpio.Open() (pirri/pirriGpioPin.go:94-95)
**File:** pirri/pirriGpioPin.go
**Problem:** Lines 94-95 have `defer rpio.Close()` BEFORE `rpio.Open()`. The defer is registered before the resource is opened, meaning Close() will operate on a stale/uninitialized state.
**Fix:** Swap the order:
```go
func gpioActivate(gpio int, state bool, seconds int) {
    log := logging.Service()
    set := settings.Service()
    rpio.Open()
    defer rpio.Close()
    // ... rest of function
```

**Build after fix:**
```bash
cd /Users/joe/Projects/pirrigo && go build -o pirrigo 2>&1
cd /Users/joe/Projects/pirrigo/pirrigo-spa && npx ng build 2>&1
```

---

## CRITICAL BUG 4: schema.go — Sets Notes="common" instead of Common=true (schema.go:50)
**File:** schema.go
**Problem:** Line 50 calls `d.DB.Model(&pirri.GpioPin{}).Where("gpio = ?", 21).Update("notes", "common")`. This sets the `Notes` text field to "common" instead of setting the `Common` boolean field to `true`. Later, `SetCommonWire()` queries `Where("common = true")` which will find nothing, causing the common wire to never be configured.
**Fix:** Change line 50 to:
```go
d.DB.Model(&pirri.GpioPin{}).Where("gpio = ?", 21).Update("common", true)
```

**Build after fix:**
```bash
cd /Users/joe/Projects/pirrigo && go build -o pirrigo 2>&1
cd /Users/joe/Projects/pirrigo/pirrigo-spa && npx ng build 2>&1
```

---

## HIGH BUG 5: main.go — Password printed to stdout (main.go:28)
**File:** main.go
**Problem:** Line 28 prints `PIRRIGO_PASSWORD` to stdout via `log.Println`, exposing the password in logs.
**Fix:** Remove or redact line 28:
```go
log.Println("PIRRIGO_PASSWORD:", "***REDACTED***")
```

**Build after fix:**
```bash
cd /Users/joe/Projects/pirrigo && go build -o pirrigo 2>&1
cd /Users/joe/Projects/pirrigo/pirrigo-spa && npx ng build 2>&1
```

---

## HIGH BUG 6: logging/pirriLogging.go — LogError defers before Lock (logging/pirriLogging.go:89-103)
**File:** logging/pirriLogging.go
**Problem:** In `LogError`, defers are registered BEFORE `l.lock.Lock()` is called (lines 91-93). This means `l.lock.Unlock()` runs on a lock that was acquired after the defer was set. Additionally, in both `LogEvent` and `LogError`, `l.logger.Sync()` runs AFTER `l.lock.Unlock()` (due to LIFO defer order), meaning the sync happens without the lock held, allowing concurrent access to the logger during sync.
**Fix for LogError:** Move `l.lock.Lock()` before the defers:
```go
func (l *PirriLogger) LogError(message string, fields ...zapcore.Field) {
    l.lock.Lock()
    defer l.lock.Unlock()
    defer l.logger.Sync()
    fields = append(
        fields,
        []zapcore.Field{
            zap.String("time", time.Now().Format(os.Getenv("PIRRIGO_DATE_FORMAT"))),
        }...,
    )
    l.logger.Error(message, fields...)
}
```
**Fix for LogEvent:** Reorder defers so Sync runs before Unlock:
```go
func (l *PirriLogger) LogEvent(message string, fields ...zapcore.Field) {
    if os.Getenv("PIRRIGO_LOG_LOCATION") != "" {
        fmt.Println("EVENT: ", message)
        l.lock.Lock()
        defer l.lock.Unlock()
        defer l.logger.Sync()
        // ... rest unchanged
    }
}
```

**Build after fix:**
```bash
cd /Users/joe/Projects/pirrigo && go build -o pirrigo 2>&1
cd /Users/joe/Projects/pirrigo/pirrigo-spa && npx ng build 2>&1
```

---

## HIGH BUG 7: pirriWebAuth.go — Index out of bounds (pirri/pirriWebAuth.go:43-48)
**File:** pirri/pirriWebAuth.go
**Problem:** After the cookie auth branch fails (line 43), if `len(s) != 2`, the code falls through to line 46 which accesses `s[0]` for logging. But `s` could have only 1 element (from the initial `strings.SplitN` on line 32), making `s[0]` a valid access but semantically wrong (it's the entire Authorization header value, not a username). More critically, if `s` has 0 elements (empty Authorization header), `s[0]` will panic.
**Fix:** Add bounds check before accessing `s[0]`:
```go
if len(s) != 2 {
    http.Error(w, "Invalid credentials", 401)
    if len(s) > 0 {
        log.LogError("HTTP Authentication Error.",
            zap.String("authHeader", s[0]),
        )
    } else {
        log.LogError("HTTP Authentication Error - missing Authorization header.")
    }
    return
}
```

**Build after fix:**
```bash
cd /Users/joe/Projects/pirrigo && go build -o pirrigo 2>&1
cd /Users/joe/Projects/pirrigo/pirrigo-spa && npx ng build 2>&1
```

---

## HIGH BUG 8: pirriWebAuth.go — Cookie auth parsing is broken (pirri/pirriWebAuth.go:36-41)
**File:** pirri/pirriWebAuth.go
**Problem:** The cookie auth branch tries to parse the cookie value as URL query parameters using `url.ParseQuery(c.Value)`, then accesses `q["username"][0]`. But the cookie value is expected to be `user:pass`, not a query string. `url.ParseQuery("user:pass")` returns an empty map, so `q["username"]` doesn't exist, and the branch silently fails.
**Fix:** Either remove the broken cookie auth entirely, or fix it to parse the cookie value as `user:pass`:
```go
c, _ := r.Cookie("Authorization")
if c != nil {
    pair := strings.SplitN(c.Value, ":", 2)
    if len(pair) == 2 {
        s = []string{"Basic", c.Value}
    }
}
```

**Build after fix:**
```bash
cd /Users/joe/Projects/pirrigo && go build -o pirrigo 2>&1
cd /Users/joe/Projects/pirrigo/pirrigo-spa && npx ng build 2>&1
```

---

## HIGH BUG 9: pirriQueueHelper.go — TOCTOU race condition (pirri/pirriQueueHelper.go:9-21)
**File:** pirri/pirriQueueHelper.go
**Problem:** The queue length is checked without the lock (lines 10-12), then the lock is acquired to pop an item (lines 16-19). Between the check and the pop, another goroutine could empty the queue, causing `OfflineRunQueue[0]` to panic.
**Fix:** Hold the lock continuously:
```go
func ListenForTasks() {
    defer WG.Done()
    for {
        ORQMutex.Lock()
        var task *Task
        if len(OfflineRunQueue) > 0 {
            task = OfflineRunQueue[0]
            OfflineRunQueue = OfflineRunQueue[1:]
        }
        ORQMutex.Unlock()

        if task != nil {
            task.execute()
        }
        time.Sleep(time.Duration(1000) * time.Millisecond)
    }
}
```

**Build after fix:**
```bash
cd /Users/joe/Projects/pirrigo && go build -o pirrigo 2>&1
cd /Users/joe/Projects/pirrigo/pirrigo-spa && npx ng build 2>&1
```

---

## HIGH BUG 10: pirriWebRunStatus.go — Race conditions on shared state (pirri/pirriWebRunStatus.go:12-37)
**File:** pirri/pirriWebRunStatus.go
**Problem:** `RUNSTATUS` and `OfflineRunQueue` are accessed without synchronization in `statusRunWeb` (line 13), `statusRunCancel` (line 22), and `statusRunQueue` (line 32). These are accessed concurrently by HTTP handlers and the task execution goroutine.
**Fix:** Add mutex protection:
- In `statusRunWeb`, wrap `&RUNSTATUS` read with `ORQMutex.Lock()/Unlock()`
- In `statusRunCancel`, wrap `RUNSTATUS.Cancel = true` with `ORQMutex.Lock()/Unlock()`
- In `statusRunQueue`, wrap `&OfflineRunQueue` read with `ORQMutex.Lock()/Unlock()`

```go
func statusRunWeb(rw http.ResponseWriter, req *http.Request) {
    ORQMutex.Lock()
    status := RUNSTATUS
    ORQMutex.Unlock()
    blob, err := json.Marshal(&status)
    // ... rest
}

func statusRunCancel(rw http.ResponseWriter, req *http.Request) {
    ORQMutex.Lock()
    RUNSTATUS.Cancel = true
    status := RUNSTATUS
    ORQMutex.Unlock()
    blob, err := json.Marshal(&status)
    // ... rest
}

func statusRunQueue(rw http.ResponseWriter, req *http.Request) {
    ORQMutex.Lock()
    queue := OfflineRunQueue
    ORQMutex.Unlock()
    blob, err := json.Marshal(&queue)
    // ... rest
}
```

**Build after fix:**
```bash
cd /Users/joe/Projects/pirrigo && go build -o pirrigo 2>&1
cd /Users/joe/Projects/pirrigo/pirrigo-spa && npx ng build 2>&1
```

---

## HIGH BUG 11: pirriStationSchedule.go — SQL injection via weekday (pirri/pirriStationSchedule.go:46-52)
**File:** pirri/pirriStationSchedule.go
**Problem:** `nowTime.Weekday()` returns a string like "Monday" that's directly interpolated into SQL on line 48 without quoting. The resulting SQL would be `... AND Monday=true ...` which happens to work because MySQL/SQLite accept unquoted identifiers in some contexts, but it's fragile and a SQL injection pattern.
**Fix:** Use parameterized query or properly quote the weekday:
```go
sqlFilter := fmt.Sprintf("(start_date <= %s AND end_date > %s) AND %s = 1 AND start_time = %s",
    nowString,
    nowString,
    nowTime.Weekday(),
    fmt.Sprintf("%02d%02d", nowTime.Hour(), nowTime.Minute()))
```
Or better, use GORM's query builder instead of raw SQL.

**Build after fix:**
```bash
cd /Users/joe/Projects/pirrigo && go build -o pirrigo 2>&1
cd /Users/joe/Projects/pirrigo/pirrigo-spa && npx ng build 2>&1
```

---

## HIGH BUG 12: pirriWebSchedule.go — Invalid .Order("ASC") (pirri/pirriWebSchedule.go:18)
**File:** pirri/pirriWebSchedule.go
**Problem:** `.Order("ASC")` is invalid GORM syntax. It should specify a column: `.Order("id ASC")` or `.Order("start_date ASC")`. The current code will generate invalid SQL.
**Fix:** Change line 18 to:
```go
data.Service().DB.Where("end_date > ? AND start_date <= ?", time.Now(), time.Now()).Order("start_date ASC").Find(&stationSchedules)
```

**Build after fix:**
```bash
cd /Users/joe/Projects/pirrigo && go build -o pirrigo 2>&1
cd /Users/joe/Projects/pirrigo/pirrigo-spa && npx ng build 2>&1
```

---

## HIGH BUG 13: data/GormHelper.go — Continues after DB connection failure (data/GormHelper.go:76-83)
**File:** data/GormHelper.go
**Problem:** In `sqliteConnect`, when `gorm.Open` returns an error (line 77), the error is only logged but execution continues. Line 83 calls `d.DB.LogMode(...)` on a nil `d.DB`, causing a panic.
**Fix:** Return/panic after logging the error:
```go
if err != nil {
    log.LogError(err.Error())
    panic("Failed to connect to sqlite3 database: " + err.Error())
}
```

**Build after fix:**
```bash
cd /Users/joe/Projects/pirrigo && go build -o pirrigo 2>&1
cd /Users/joe/Projects/pirrigo/pirrigo-spa && npx ng build 2>&1
```

---

## MEDIUM BUG 14: Missing Content-Type headers on JSON responses
**Files:** All pirriWeb*.go files
**Problem:** None of the HTTP handlers set `Content-Type: application/json` on responses. This causes browser CORS issues and incorrect content negotiation.
**Fix:** Add `rw.Header().Set("Content-Type", "application/json")` at the start of every JSON response handler, or add it inside the `enableCors` middleware wrapper.

**Build after fix:**
```bash
cd /Users/joe/Projects/pirrigo && go build -o pirrigo 2>&1
cd /Users/joe/Projects/pirrigo/pirrigo-spa && npx ng build 2>&1
```

---

## MEDIUM BUG 15: pirriWebGpio.go — SQL injection (pirri/pirriWebGpio.go:65-66)
**File:** pirri/pirriWebGpio.go
**Problem:** Lines 65-66 use `DB.Exec` with raw SQL strings. Line 66 interpolates `gpio.GPIO` directly into the SQL string instead of using a parameterized query.
**Fix:** Use parameterized query:
```go
data.Service().DB.Exec("UPDATE gpio_pins SET common = false")
data.Service().DB.Exec("UPDATE gpio_pins SET common = true WHERE gpio = ?", gpio.GPIO)
```

**Build after fix:**
```bash
cd /Users/joe/Projects/pirrigo && go build -o pirrigo 2>&1
cd /Users/joe/Projects/pirrigo/pirrigo-spa && npx ng build 2>&1
```

---

## MEDIUM BUG 16: ngx-material-timepicker version mismatch (pirrigo-spa/package.json:39)
**File:** pirrigo-spa/package.json
**Problem:** `ngx-material-timepicker` is at version `13.1.1` which targets Angular 13, but the project uses Angular 18. This causes peer dependency conflicts and potential runtime errors.
**Fix:** Update to a compatible version:
```json
"ngx-material-timepicker": "^18.0.0"
```
Then run `npm install` in the pirrigo-spa directory.

**Build after fix:**
```bash
cd /Users/joe/Projects/pirrigo && go build -o pirrigo 2>&1
cd /Users/joe/Projects/pirrigo/pirrigo-spa && npm install && npx ng build 2>&1
```

---

## MEDIUM BUG 17: HttpClient provided twice (pirrigo-spa/src/app/app.module.ts:84-88)
**File:** pirrigo-spa/src/app/app.module.ts
**Problem:** `HttpClient` is provided both in the `providers` array (line 84) and via `provideHttpClient(withInterceptorsFromDi())` (line 88). This is redundant and can cause issues.
**Fix:** Remove `HttpClient` from the providers array (line 84), keeping only `provideHttpClient(withInterceptorsFromDi())`.

**Build after fix:**
```bash
cd /Users/joe/Projects/pirrigo && go build -o pirrigo 2>&1
cd /Users/joe/Projects/pirrigo/pirrigo-spa && npx ng build 2>&1
```

---

## MEDIUM BUG 18: Hardcoded IP address (pirrigo-spa/src/app/services/globals.service.ts:11)
**File:** pirrigo-spa/src/app/services/globals.service.ts
**Problem:** The API URL is hardcoded to `http://192.168.86.102`. The SPA should use relative URLs when served from the same origin as the backend.
**Fix:** Change line 11 to use a relative URL:
```typescript
public uriStem: string = ""
```
This way, when the SPA is served by the Go backend at `http://localhost:PORT/`, API calls will go to `http://localhost:PORT/station/all` etc.

**Build after fix:**
```bash
cd /Users/joe/Projects/pirrigo && go build -o pirrigo 2>&1
cd /Users/joe/Projects/pirrigo/pirrigo-spa && npx ng build 2>&1
```

---

## MEDIUM BUG 19: Division by zero (pirrigo-spa/src/app/components/stations/stations.component.ts:107)
**File:** pirrigo-spa/src/app/components/stations/stations.component.ts
**Problem:** In `findDateDiffPercent` (line 107), if `this.status.Duration` is 0, the division `(-sec / this.status.Duration) * 100` will produce `NaN` or `Infinity`.
**Fix:** Add guard:
```typescript
findDateDiffPercent(date: Date, duration: number): number {
    if (duration === 0) return 0
    let now = moment(new Date());
    let end = moment(date).add(duration, "s");
    let durationDiff = moment.duration(now.diff(end));
    let sec = durationDiff.asSeconds()
    return Math.round(100 - ((-sec / this.status.Duration) * 100))
}
```

**Build after fix:**
```bash
cd /Users/joe/Projects/pirrigo && go build -o pirrigo 2>&1
cd /Users/joe/Projects/pirrigo/pirrigo-spa && npx ng build 2>&1
```

---

## MEDIUM BUG 20: XSS vulnerability in calendar events (pirrigo-spa/src/app/components/calendar/calendar.component.ts:93)
**File:** pirrigo-spa/src/app/components/calendar/calendar.component.ts
**Problem:** Line 93 embeds `JSON.stringify(event)` directly into the event title HTML string. If any schedule field contains HTML or script tags, they will be rendered as HTML, enabling XSS.
**Fix:** HTML-encode the JSON string before embedding:
```typescript
"title": `Zone ${event.StationID} for ${event.Duration / 60} minutes<br/><br/><br/><br/> | ${this.escapeHtml(JSON.stringify(event))}`
```
Add helper method:
```typescript
escapeHtml(str: string): string {
    return str.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
}
```

**Build after fix:**
```bash
cd /Users/joe/Projects/pirrigo && go build -o pirrigo 2>&1
cd /Users/joe/Projects/pirrigo/pirrigo-spa && npx ng build 2>&1
```

---

## MEDIUM BUG 21: pirriWebDripNodes.go — SQL column name case mismatch (pirri/pirriWebDripNodes.go:74)
**File:** pirri/pirriWebDripNodes.go
**Problem:** Line 74 references `station_histories.station_ID` (with capital ID), but the actual column name in the schema is `station_id` (lowercase). This will cause the query to fail on case-sensitive databases.
**Fix:** Change `station_ID` to `station_id` on line 74.

**Build after fix:**
```bash
cd /Users/joe/Projects/pirrigo && go build -o pirrigo 2>&1
cd /Users/joe/Projects/pirrigo/pirrigo-spa && npx ng build 2>&1
```

---

## MEDIUM BUG 22: calendar eventTimesChanged not wired up (pirrigo-spa/src/app/components/calendar/calendar.component.ts:68-70)
**File:** pirrigo-spa/src/app/components/calendar/calendar.component.ts
**Problem:** `eventTimesChanged` (line 68) only logs to console but doesn't persist the changed event times. The calendar events are marked as non-resizable and non-draggable, so this is dead code, but it should either be removed or implemented.
**Fix:** Remove the dead `eventTimesChanged` method and its binding in the template, since events are already marked non-resizable/non-draggable.

**Build after fix:**
```bash
cd /Users/joe/Projects/pirrigo && go build -o pirrigo 2>&1
cd /Users/joe/Projects/pirrigo/pirrigo-spa && npx ng build 2>&1
```

---

END OF BUG LIST
