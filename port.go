package main

import (
	"github.com/StackExchange/wmi"
)

type SerialPortInfo struct {
	Status      string
	Name        string
	DeviceID    string
	Description string
	Caption     string
}

func getPort() (*[]SerialPortInfo, error) {
	var ports []SerialPortInfo
	query := "SELECT * FROM Win32_SerialPort"

	err := wmi.Query(query, &ports)
	if err != nil {
		return nil, err
	}
	return &ports, nil
}
