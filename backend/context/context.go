package context

type Context struct {
	LastModuleName   string
	LastSerialBuffer []byte
	LastExecResult   error
}

var GlobalCtx Context
