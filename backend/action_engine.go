package bsd_testtool

type variable struct {
	varName string
	varType string
	v       any
}

var variableList []*variable

var labelList []string

type ActionContext struct {
	app *App

	nowIndex     int
	nowModuleUID int

	LastModuleName   string
	LastSerialBuffer []byte
	LastExecResult   error
}

type ActionEngine struct {
}

func (a *ActionEngine) Start() {

}

func (a *ActionEngine) Step() {

}

func (a *ActionEngine) Stop() {

}
