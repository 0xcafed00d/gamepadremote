package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/simulatedsimian/joystick"
	"github.com/tarm/serial"
)

type Config struct {
	Help          bool
	SerialDevice  string
	SerialSpeed   int
	NetHost       string
	NetPort       int
	ConsoleOutput bool
	RateMS        int
	JoystickIdx   int
}

var config Config

func init() {
	flag.BoolVar(&config.Help, "h", false, "display help")
	flag.BoolVar(&config.ConsoleOutput, "c", false, "output to console")
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

func exitOnError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func openComms(config Config) io.ReadWriteCloser {
	if config.ConsoleOutput {
		return os.Stdout
	}

	if len(config.NetHost) > 0 && config.NetPort != 0 {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", config.NetHost, config.NetPort))
		exitOnError(err)
		return conn
	}

	if len(config.SerialDevice) > 0 && config.SerialSpeed != 0 {
		serialcfg := serial.Config{Name: config.SerialDevice, Baud: config.SerialSpeed}
		port, err := serial.OpenPort(&serialcfg)
		exitOnError(err)
		return port
	}

	fmt.Fprintln(os.Stderr, "comms port incorrectly specified")
	flag.Usage()
	os.Exit(1)

	return nil
}

func openJoystick(config Config) joystick.Joystick {
	js, err := joystick.Open(config.JoystickIdx)
	exitOnError(err)

	return js
}

func checksum(s string) (ch byte) {

	for i := 0; i < len(s); i++ {
		ch += s[i]
	}
	ch = ^ch + 1
	return
}

func main() {
	flag.Parse()

	if len(flag.Args()) != 0 || config.Help {
		flag.Usage()
		os.Exit(1)
	}

	js := openJoystick(config)
	comms := openComms(config)
	defer comms.Close()

	ticker := time.NewTicker(time.Duration(config.RateMS) * time.Millisecond)

	if !config.ConsoleOutput {
		go func(r io.Reader) {
			io.Copy(os.Stdout, r)
		}(comms)
	}

	for {
		<-ticker.C
		state, err := js.Read()
		exitOnError(err)
		packet := fmt.Sprintf("!J%04x|%04x|%04x|%04x|%04x",
			uint16(state.Buttons),
			uint16(state.AxisData[0]),
			uint16(state.AxisData[1]),
			uint16(state.AxisData[3]),
			uint16(state.AxisData[4]))

		_, err = fmt.Fprintf(comms, "%s#%02x\n", packet, checksum(packet))
		exitOnError(err)
	}
}
