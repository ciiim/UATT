package bsd_testtool_test

import (
	bsd_testtool "bsd_testtool/backend"
	"testing"
)

func TestSerial(t *testing.T) {
	portList, _ := bsd_testtool.GlobalSerial.GetAllPort()
	t.Logf("port %v, len:%d", portList, len(portList))

}
