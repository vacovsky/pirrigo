export class StationSchedule {
    ID: number
    StartDate: Date
    EndDate: Date
    Sunday: boolean
    Monday: boolean
    Tuesday: boolean
    Wednesday: boolean
    Thursday: boolean
    Friday: boolean
    Saturday: boolean
    StationID: number
    StartTime: number
    Duration: number
    Repeating: boolean
}

export class StationScheduleResponse {
    schedule: StationSchedule[]
}
