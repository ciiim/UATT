package bsd_testtool

type command struct {
	req       []byte
	expectAck []byte
	timeout   int
}
