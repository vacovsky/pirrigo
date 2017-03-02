package main

//GpioPin - describes a Raspberry Pi GPIO pin
type GpioPin struct {
	ID     int `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	GPIO   int `gorm:"not null;unique"`
	Notes  string
	Common bool `sql:"DEFAULT:false" gorm:"not null"`
}

func getGpioByPin() {
	gpio := db.Where("GPIO = ?", 4).Find(&GpioPin{}).Order("GPIO DESC")
	jsonifySQLResults(gpio)
}

// func (g *GpioPin) gpioActivator(gpio int, state bool, seconds int) {
//  if SETTINGS.SimulateGpioActivity {
//      gpioSimulation(gpio, state, seconds)
//  } else {
//      gpioActivate(gpio, state, seconds)
//  }
// }

// func (g *GpioPin) gpioSimulation(gpio int, state bool, seconds int) {
//  fmt.Println("GPIO Simulation starting:",
//      "\nTime:", time.Now(),
//      "\nGPIO:", gpio,
//      "\nDesired State:", state,
//      "\nDuration (seconds):", seconds)
//  fmt.Println("Active!", time.Now())
//  for seconds > 0 {
//      time.Sleep(time.Second)
//      seconds--
//  }
//  fmt.Println("Deactivated!", time.Now())
// }

// func (g *GpioPin) gpioActivate(gpio int, state bool, seconds int) {
//  pin := rpio.Pin(gpio)
//  defer rpio.Close()

//  // activate gpio
//  if state {
//      pin.High()
//  } else {
//      pin.Low()
//  }

//  // start countdown
//  for seconds > 0 {
//      time.Sleep(time.Duration(1) * time.Second)
//      seconds--
//  }

//  if !state {
//      pin.High()
//  } else {
//      pin.Low()
//  }
// }
