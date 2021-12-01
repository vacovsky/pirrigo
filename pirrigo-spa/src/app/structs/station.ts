export class Station {
    ID: number
    GPIO: number
    Notes: string
    Description: string
    Enabled: boolean
    FriendlyName: string
    LastRun: Date
    NextRun: Date
}

export class StationResponse {
    stations: Station[]
}

