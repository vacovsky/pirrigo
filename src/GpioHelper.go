package main

import (
	"fmt"
	"time"

	"github.com/stianeikeland/go-rpio"
)

func Example() {
	pin := rpio.Pin(10)

	pin.Output() // Output mode
	pin.High()   // Set pin High
	pin.Low()    // Set pin Low
	pin.Toggle() // Toggle pin (Low -> High -> Low)

	pin.Input()       // Input mode
	res := pin.Read() // Read state from pin (High / Low)
	fmt.Println(res)
	pin.Mode(rpio.Output) // Alternative syntax
	pin.Write(rpio.High)  // Alternative syntax
	pin.PullUp()
	pin.PullDown()
	pin.PullOff()

	pin.Pull(rpio.PullUp)
	rpio.Close()

}

func GpioActivator(gpio int, state bool, seconds int) {
	if SETTINGS.SimulateGpioActivity {
		GpioSimulation(gpio, state, seconds)
	} else {
		GpioActivate(gpio, state, seconds)
	}
}

func GpioSimulation(gpio int, state bool, seconds int) {
	fmt.Println("GPIO Sinulation starting:",
		"\nTime:", time.Now(), "\nGPIO:", gpio,
		"\nDesired State", state,
		"\nDuration (seconds):", seconds)
	fmt.Println("Active!", time.Now())
	for seconds > 0 {
		fmt.Println(seconds)

		time.Sleep(time.Duration(seconds) * time.Second)
		seconds -= 1
	}
	fmt.Println("Deactivated!", time.Now())

}

//TODO

func GpioActivate(gpio int, state bool, seconds int) {
	pin := rpio.Pin(10)
	defer rpio.Close()
	fmt.Println(pin)
}
