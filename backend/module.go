package bsd_testtool

// Module Type
const (
	subContainerMT = 0xFF
	printMT        = 0xFE

	gotoLabelMT = 0x00
	gotoMT      = 0x01

	forLabelMT = 0x02

	whileLabelMT = 0x03

	sendMT    = 0x04
	receiveMT = 0x05

	ifMT   = 0x06
	elseMT = 0x07

	exitMT = 0x08
)

type IModule interface {
}
