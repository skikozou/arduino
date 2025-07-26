package main

import (
	"fmt"
	"io"

	"github.com/sirupsen/logrus"
	"github.com/tarm/serial"
)

func Init() {
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:            true,
		DisableLevelTruncation: true,
		PadLevelText:           true,
	})
}

func main() {
	Init()
	ui := NewLogUI()

	go func() {
		if err := ui.Run(); err != nil {
			panic(err)
		}
	}()

	go func() {
		ui.Log("Loading COM ports...")

		results, err := getPort()
		index := -1

		for {
			for _, r := range *results {
				ui.Log(fmt.Sprintf(" %s - %s", r.DeviceID, r.Name))
			}
			port := ui.RequestInput("select COM port")

			for i, v := range *results {
				if v.DeviceID == port {
					index = i
					break
				}
			}

			if index != -1 {
				break
			}

			ui.Log("invaid input")
		}

		c := &serial.Config{
			Name: (*results)[index].DeviceID,
			Baud: 9600,
		}
		s, err := serial.OpenPort(c)
		if err != nil {
			logrus.Fatal(err)
		}
		defer s.Close()

		go Writer(ui, s)

		ms := make([]byte, 128)
		buf := make([]byte, 128)
		for {
			n, err := s.Read(buf)
			if err != nil && err != io.EOF {
				logrus.Fatal(err)
			}
			if n > 0 {
				ms = append(ms, buf[:n]...)
				if string(buf[n-1]) == "\n" {
					ui.Log(string(ms))
					ms = nil
				}
			}
		}
	}()

	select {}
}

func Writer(ui *LogUI, s *serial.Port) {
	for {
		ms := ui.RequestInput("message")
		s.Write([]byte(ms + "\n"))
	}
}
