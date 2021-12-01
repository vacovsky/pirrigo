export class Station {
    ID: number
    GPIO: number
    Notes: string
    Description: string
}

export class StationResponse {
    stations: Station[]
}

