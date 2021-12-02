export class Station {
    ID: number
    GPIO: number
    Notes: string
    Description: string
    Enabled: boolean
    FriendlyName: string
    LastRun: Date
    NextRun: Date
    Status: StationStatus
}

export class StationResponse {
    stations: Station[]
}

export class StationStatus {
    IsIdle: boolean
    IsManual: boolean
    StartTime: Date
    Duration: number
    ScheduleID: number
    StationID: number
    Cancel: boolean
}