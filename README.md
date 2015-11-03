# gamepadremote
Game Pad Controller. Reads xbox 360 Joypad &amp; sends data over serial port, network or to console 

## Installation:
```bash
$ go get github.com/simulatedsimian/gamepadremote
```

## Usage:
```gamepadremote [options]
  -b int      serial baudrate (default 9600)
  -c	        output to console
  -d string   serial device name
  -h	        display help
  -j int    	gamepad index (default 0)
  -n string  	network host name
  -p int    	network port
  -r int    	sample rate in ms (default 100)
```
