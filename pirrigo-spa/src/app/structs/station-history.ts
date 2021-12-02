export class StationHistory {
    ID: number
    StationID: number
    ScheduleID: number
    Duration: number
    StartTime: Date
}

export class StationHistoryResponse {
    history: StationHistory[]
}
