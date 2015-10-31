package main

import (
	"flag"
	"fmt"
	"github.com/simulatedsimian/joystick"
	"io"
	"os"
	"time"
)

type Config struct {
	Help         bool
	SerialDevice string
	SerialSpeed  int
	NetHost      string
	NetPort      int
	RateMS       int
	JoystickIdx  int
}

var config Config

func init() {
	flag.BoolVar(&config.Help, "h", false, "display help")
	flag.StringVar(&config.SerialDevice, "d", "", "serial device name")
	flag.IntVar(&config.SerialSpeed, "b", 9600, "serial baudrate")
	flag.IntVar(&config.RateMS, "r", 100, "sample rate in ms")
	flag.IntVar(&config.JoystickIdx, "j", 0, "gamepad index (default 0)")
	flag.StringVar(&config.NetHost, "n", "", "network host name")
	flag.IntVar(&config.NetPort, "p", 0, "network port")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: gamepadremote [options]")
		flag.PrintDefaults()
	}
}

func openComms(config Config) io.ReadWriteCloser {
	if len(config.NetHost) > 0 && config.NetPort != 0 {

	}

	if len(config.SerialDevice) > 0 && config.SerialSpeed != 0 {

	}

	return nil
}

func openJoystick(config Config) joystick.Joystick {
	js, err := joystick.Open(config.JoystickIdx)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return js
}

func main() {
	flag.Parse()

	if len(flag.Args()) != 0 || config.Help {
		flag.Usage()
		os.Exit(1)
	}

	js := openJoystick(config)
	comms := openComms(config)
	//defer comms.Close()

	ticker := time.NewTicker(time.Duration(config.RateMS) * time.Millisecond)

	for {
		<-ticker.C
		fmt.Println("tick")

	}

	js = js
	comms = comms
}
