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

    // dayIsActiveHash: any;

    // public populatDayIsActiveHash(): void {
    //     this.dayIsActiveHash = {
    //         "Sunday": this.Sunday,
    //         "Monday": this.Monday,
    //         "Tuesday": this.Tuesday,
    //         "Wednesday": this.Wednesday,
    //         "Thursday": this.Thursday,
    //         "Friday": this.Friday,
    //         "Saturday": this.Saturday,
    //     }
    // }
}

export class StationScheduleResponse {
    stationSchedules: StationSchedule[]
}
