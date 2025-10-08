package bsd_testtool_test

import (
	bsd_testtool "bsd_testtool/backend"
	"fmt"
	"testing"
)

func TestManagerInit(t *testing.T) {
	bsd_testtool.GlobalManager.Init("./Apps")
}

func TestSaveApp(t *testing.T) {
	bsd_testtool.GlobalManager.Init("./Apps")
	if err := bsd_testtool.GlobalManager.LoadApp("app-111.json"); err != nil {
		t.Error(err)
	}
	if err := bsd_testtool.GlobalManager.SaveApp(); err != nil {
		t.Error(err)
	}
}

func TestManagerLoadConfig(t *testing.T) {
	bsd_testtool.GlobalManager.Init("./Apps")
	err := bsd_testtool.GlobalManager.LoadApp("app-111.json")
	t.Log(err)
	app := bsd_testtool.GlobalManager.GetNowApp()

	ae := bsd_testtool.NewActionEngine(app)

	fmt.Printf("ae.PreCompile(): %v\n", ae.PreCompile())

	if err := ae.StartSync(); err != nil {
		t.Error(err)
	}

	fmt.Printf("ae.GetStopReason(): %v\n", ae.GetStopReason())

}
