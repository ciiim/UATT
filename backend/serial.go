package bsd_testtool

import (
	"errors"

	"go.bug.st/serial"
)

type Serial struct {
	port     serial.Port
	mode     serial.Mode
	portName string
}

var ErrSerialNotOpen error = errors.New("serial not open")

var GlobalSerial Serial

func (s *Serial) GetAllPort() ([]string, error) {
	return serial.GetPortsList()
}

func (s *Serial) GetMode() *serial.Mode {
	return &s.mode
}

func (s *Serial) SetModeS(m *serial.Mode) {
	s.mode.BaudRate = m.BaudRate
	s.mode.DataBits = m.DataBits
	s.mode.Parity = m.Parity
	s.mode.StopBits = m.StopBits
}

func (s *Serial) SetMode(baudRate int, dataBits int, parity serial.Parity, stopBits serial.StopBits) {
	s.mode.BaudRate = baudRate
	s.mode.DataBits = dataBits
	s.mode.Parity = parity
	s.mode.StopBits = stopBits
}

func (s *Serial) SelectPort(port string) {
	s.portName = port
}

func (s *Serial) OpenSerial() error {
	p, err := serial.Open(s.portName, &s.mode)
	if err != nil {
		return err
	}
	s.port = p
	return nil
}

func (s *Serial) CloseSerial() error {
	if s.port == nil {
		return ErrSerialNotOpen
	}
	if err := s.port.Close(); err != nil {
		return err
	}
	s.port = nil
	return nil
}

func (s *Serial) Write(buffer []byte) (int, error) {
	if s.port == nil {
		return 0, ErrSerialNotOpen
	}
	return s.port.Write(buffer)
}

func (s *Serial) Read(buffer []byte) (int, error) {
	if s.port == nil {
		return 0, ErrSerialNotOpen
	}
	return s.port.Read(buffer)
}
