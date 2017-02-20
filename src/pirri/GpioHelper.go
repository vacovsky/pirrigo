package main

import (
	"fmt"
	"time"

	"github.com/stianeikeland/go-rpio"
)

func gpioActivator(gpio int, state bool, seconds int) {
	if SETTINGS.SimulateGpioActivity {
		gpioSimulation(gpio, state, seconds)
	} else {
		gpioActivate(gpio, state, seconds)
	}
}

func gpioSimulation(gpio int, state bool, seconds int) {
	fmt.Println("GPIO Simulation starting:",
		"\nTime:", time.Now(),
		"\nGPIO:", gpio,
		"\nDesired State:", state,
		"\nDuration (seconds):", seconds)
	fmt.Println("Active!", time.Now())
	for seconds > 0 {
		time.Sleep(time.Second)
		seconds -= 1
	}
	fmt.Println("Deactivated!", time.Now())
}

func gpioActivate(gpio int, state bool, seconds int) {
	pin := rpio.Pin(gpio)
	defer rpio.Close()

	// activate gpio
	if state {
		pin.High()
	} else {
		pin.Low()
	}

	// start countdown
	for seconds > 0 {
		time.Sleep(time.Duration(1) * time.Second)
		seconds -= 1
	}

	if !state {
		pin.High()
	} else {
		pin.Low()
	}
}

/*func Example() {
//	pin := rpio.Pin(10)

//	pin.Output() // Output mode
//	pin.High()   // Set pin High
//	pin.Low()    // Set pin Low
//	pin.Toggle() // Toggle pin (Low -> High -> Low)

//	pin.Input()       // Input mode
//	res := pin.Read() // Read state from pin (High / Low)
//	fmt.Println(res)
//	pin.Mode(rpio.Output) // Alternative syntax
//	pin.Write(rpio.High)  // Alternative syntax
//	pin.PullUp()
//	pin.PullDown()
//	pin.PullOff()

//	pin.Pull(rpio.PullUp)
//	rpio.Close()
}*/
