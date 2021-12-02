export class StationLogs {
    level: string
    message: string
    application: string
    interval: number
    time: string
}


export class StationLogsResponse {
    logs: StationLogs[]
}
