export class Gpio {
    ID: number
    GPIO: number
    Notes: string
    Common: boolean
}

export class GpioResponse {
    gpios: Gpio[]
}
