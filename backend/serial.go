package bsd_testtool

import (
	"errors"

	"github.com/tarm/serial"
)

type SerialManager struct {
	commandHead []byte
	port        *serial.Port
	verifier    *Verifier
}

func NewSerialManager(verifier *Verifier, head []byte) *SerialManager {
	return &SerialManager{
		commandHead: head,
		verifier:    verifier,
	}
}

func (s *SerialManager) Open(name string, baud int) error {
	port, err := serial.OpenPort(&serial.Config{
		Name: name,
		Baud: baud,
	})
	if err != nil {
		return err
	}
	s.port = port
	return nil
}

func (s *SerialManager) WriteWithVerify(body []byte) (int, error) {
	// add head
	body = append(s.commandHead, body...)
	// add verify
	verifyBytes := s.verifier.Verify(body)
	body = append(body, verifyBytes...)
	return s.port.Write(body)
}

func (s *SerialManager) ReadWithVerify(body []byte) (int, error) {
	_, err := s.port.Read(body)
	if err != nil {
		return 0, err
	}

	// verify head
	if string(body[:len(s.commandHead)]) != string(s.commandHead) {
		return 0, errors.New("head verify failed")
	}

	// verify body
	dataLen := len(body)
	verifyLen := s.verifier.Len()
	recvVerifyBytes := body[dataLen-verifyLen:]
	body = body[:dataLen-verifyLen]
	verifyBytes := s.verifier.Verify(body)
	if string(verifyBytes) != string(recvVerifyBytes) {
		return 0, errors.New("verify failed")
	}

	return dataLen, nil
}

func (s *SerialManager) Close() error {
	return s.port.Close()
}
