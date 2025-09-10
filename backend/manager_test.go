package bsd_testtool_test

import (
	bsd_testtool "bsd_testtool/backend"
	"testing"
)

func TestManagerInit(t *testing.T) {
	bsd_testtool.GlobalManager.Init("./Apps")
}

func TestManagerLoadConfig(t *testing.T) {
	bsd_testtool.GlobalManager.Init("./Apps")
	err := bsd_testtool.GlobalManager.LoadApp("111")
	t.Log(err)
	app := bsd_testtool.GlobalManager.GetNowApp()
	if app != nil {
		app.PrintConfig()
	}

}
