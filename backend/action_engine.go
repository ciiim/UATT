package bsd_testtool

import "errors"

type ControlType ModuleTypeID

type Index int

type variable struct {
	varName string
	varType string
	v       any
}

type controlFlow struct {
	controlType ControlType

	TrueIndex  Index // IF 成立执行的下一个 index
	FalseIndex Index // IF 不成立执行的下一个 index
	EndIndex   Index // FOR/IF 最终要跳到的 BLOCKEND
	LoopStart  Index // FOR 循环开始的 index
}

type ActionContext struct {
	actionModuleList []ModuleBase

	// 变量表
	variableMap map[string]*variable

	// UID表，用于引用的快速获取
	moduleUIDMap map[ModuleUID]*ModuleBase

	// 控制流表，key为oneIndex，用于if else for的流程跳转
	controlFlowMap map[Index]*controlFlow

	// 标签表，用于goto跳转
	labelMap map[string]Index

	nowIndex     Index // 等价于PC指针
	nowModuleUID int

	LastModuleName   string
	LastSerialBuffer []byte
	LastExecResult   error
}

type ActionEngine struct {
	app *App

	ctx *ActionContext
}

func (a *ActionEngine) PreCompile() error {
	if a.app == nil {
		return ErrNotFoundApp
	}

	a.ctx = &ActionContext{}

	// 拷贝执行代码副本
	a.ctx.actionModuleList = a.app.config.Actions

	// 遍历代码
	for i := 0; i < len(a.ctx.actionModuleList); i++ {
		// 编码Index
		a.ctx.actionModuleList[i].Index = i

		//UID表
		a.ctx.moduleUIDMap[ModuleUID(a.ctx.actionModuleList[i].ModuleUID)] = &a.ctx.actionModuleList[i]

		// 变量表
		if ModuleTypeID(a.ctx.actionModuleList[i].ModuleTypeID) == DeclareMT {
			tmp, is := a.ctx.actionModuleList[i].TypeFeatureField.(DeclareModuleFeatureField)
			if !is {
				return errors.New("wrong module type and type id")
			}
			v := &variable{
				varName: tmp.VarName,
				varType: tmp.VarType,
			}
			switch tmp.VarType {
			case "number":
				v.v = tmp.VarNumberValue
			case "string":
				v.v = tmp.VarStringValue
			case "bytesarray":
				v.v = tmp.VarByteArrayValue
			}
			a.ctx.variableMap[tmp.VarName] = v
		}

		// label表
		if ModuleTypeID(a.ctx.actionModuleList[i].ModuleTypeID) == GotoLabelMT {
			tmp, is := a.ctx.actionModuleList[i].TypeFeatureField.(LabelModuleFeatureField)
			if !is {
				return errors.New("wrong module type and type id")
			}
			a.ctx.labelMap[tmp.LabelName] = Index(i)
		}

		// 控制流表
	}

	return nil
}

func (a *ActionEngine) Start() {
	a.ctx.nowIndex = 0
}

func (a *ActionEngine) doModule() {

}

func (a *ActionEngine) next() {

}

func (a *ActionEngine) Step() {

}

func (a *ActionEngine) Stop() {

}
