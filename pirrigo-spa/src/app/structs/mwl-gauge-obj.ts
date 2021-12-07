export class MwlGaugeObj {
    Value: number
    DialStartAngle: number
    DialEndAngle: number
    Animated: boolean
    AnimationDuration: number // seconds
    Max: number
    Min: number
    DialRadius: number
    Label: any
    ShowValue: number // Whether to show the value at the center of the gauge
    ValueClass: string // The CSS class of the gauge's text
    ValueDialClass: string // The CSS class of the gauge's fill (value dial)
    GaugeClass: string
    DialClass: string
    Color: any
    // GaugeCreated: EventEmitter<type> //Created Called when the gauge is created
}
