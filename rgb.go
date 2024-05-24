package main

import (
    "github.com/jacobsa/go-serial/serial"
)

func CtrlRGB(val string) int {
    options := serial.OpenOptions {
        PortName: "/dev/ttyUSB0",
        BaudRate: 115200,
        DataBits: 8,
        StopBits: 1,
        MinimumReadSize: 8,
    }

    port, err := serial.Open(options)
    if err != nil {
        return -1
    }
    defer port.Close()

    data := []byte(val)
    _, err = port.Write(data)
    if err != nil {
        return -2
    }

    return 0
}
