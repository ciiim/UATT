package bsd_testtool

import (
	"errors"
	"fmt"
	"io"
)

type ControlType = ModuleTypeID

type Index int

type variable struct {
	varName string
	varType string
	v       any
}

type ActionStatus int

const (
	Idle ActionStatus = iota

	// 运行中
	Running

	// 已停止
	Stopped

	// 等待单步执行
	WaitingStep

	// 等待停止,中间状态
	WaitingStop

	// 因错误停止
	ErrorStopped
)

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

	nowIndex Index // 等价于PC指针

	nowActionStatus ActionStatus

	log io.WriteCloser

	stepCh chan int

	LastModuleName   string
	LastSerialBuffer []byte
	LastExecResult   error
}

type ActionEngine struct {
	app *App

	ctx *ActionContext
}

func NewActionEngine(app *App) *ActionEngine {
	return &ActionEngine{
		app: app,
	}
}

func (a *ActionEngine) PreCompile() error {
	if a.app == nil {
		return ErrNotFoundApp
	}

	a.ctx = &ActionContext{
		actionModuleList: make([]ModuleBase, 0),
		variableMap:      make(map[string]*variable),
		moduleUIDMap:     make(map[ModuleUID]*ModuleBase),
		controlFlowMap:   make(map[Index]*controlFlow),
		labelMap:         make(map[string]Index),
		stepCh:           make(chan int, 2),
		nowIndex:         -1,
		LastModuleName:   "NIL",
		LastSerialBuffer: nil,
		LastExecResult:   nil,
		nowActionStatus:  Idle,
	}

	// 拷贝执行代码副本
	a.ctx.actionModuleList = a.app.config.Actions

	type tempControlInfo struct {
		ct ModuleTypeID
		i  Index
	}
	controlFlowStack := make([]tempControlInfo, 0)

	// 遍历代码
	for i := 0; i < len(a.ctx.actionModuleList); i++ {
		// 编码Index
		a.ctx.actionModuleList[i].Index = i

		//UID表
		a.ctx.moduleUIDMap[ModuleUID(a.ctx.actionModuleList[i].ModuleUID)] = &a.ctx.actionModuleList[i]

		t := ModuleTypeID(a.ctx.actionModuleList[i].ModuleTypeID)

		// 变量表
		if t == DeclareMT {
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
		if t == GotoLabelMT {
			tmp, is := a.ctx.actionModuleList[i].TypeFeatureField.(LabelModuleFeatureField)
			if !is {
				return errors.New("wrong module type and type id")
			}
			a.ctx.labelMap[tmp.LabelName] = Index(i)
		}

		// 控制流表
		switch t {
		case IfMT:

			// 加进map里面
			a.ctx.controlFlowMap[Index(i)] = &controlFlow{
				controlType: t,
				TrueIndex:   Index(i + 1),
			}

			// 压栈
			controlFlowStack = append(controlFlowStack, tempControlInfo{
				ct: t,
				i:  Index(i),
			})

		case ElseMT:
			cf, has := a.ctx.controlFlowMap[controlFlowStack[len(controlFlowStack)-1].i]
			if !has {
				return errors.New("else without matching if")
			}
			cf.FalseIndex = Index(i + 1)

			// 临时存IF
			ifcf := controlFlowStack[len(controlFlowStack)-1]

			// 出栈IF
			controlFlowStack = controlFlowStack[:len(controlFlowStack)-1]

			// 入栈ELSE, 但是index还是if的index，不然endblock会找不到对应的控制流程结构体
			controlFlowStack = append(controlFlowStack, tempControlInfo{
				ct: t,
				i:  ifcf.i,
			})
		case ForLabelMT:
			// 压栈
			controlFlowStack = append(controlFlowStack, tempControlInfo{
				ct: t,
				i:  Index(i),
			})
		case EndBlockMT:

			top := controlFlowStack[len(controlFlowStack)-1]

			if top.ct == IfMT || top.ct == ElseMT {
				cf, has := a.ctx.controlFlowMap[controlFlowStack[len(controlFlowStack)-1].i]
				// 忽略多余的EndBlock，不算做错误
				if !has {
					continue
				}
				cf.EndIndex = Index(i + 1)
			} else if top.ct == ForLabelMT {
				// 如果是for类型，还得往map里插入一个可以通过endblock的index查找的元素
				// 加进map里面
				a.ctx.controlFlowMap[Index(i)] = &controlFlow{
					controlType: t,

					// 这里不+1，因为要回到for上做检查
					LoopStart: top.i,

					EndIndex: Index(i + 1),
				}
			}

			//出栈
			controlFlowStack = controlFlowStack[:len(controlFlowStack)-1]
		}
	}

	if len(controlFlowStack) != 0 {
		return fmt.Errorf("if or for statment not closed, at index:%d", controlFlowStack[0].i)
	}

	return nil
}

func (a *ActionEngine) StartSync() {
	a.ctx.nowIndex = 0
	a.ctx.nowActionStatus = Running
	a.innerStart()
}

func (a *ActionEngine) StartAsync() {
	a.ctx.nowIndex = 0
	a.ctx.nowActionStatus = Running
	go a.innerStart()
}

func (a *ActionEngine) innerStart() error {
	for a.ctx.nowIndex != Index(len(a.ctx.actionModuleList)-1) && a.ctx.nowActionStatus != WaitingStop {
		m := a.ctx.actionModuleList[a.ctx.nowIndex]

		if m.BreakPoint {
			a.ctx.nowActionStatus = WaitingStep
			_, ok := <-a.ctx.stepCh
			if !ok {
				break
			}

		}

		a.ctx.LastExecResult = a.doModule(&m)

		a.ctx.LastModuleName = m.Name

		a.next()
	}

	if a.ctx.LastExecResult != nil {
		a.ctx.nowActionStatus = ErrorStopped
	}
	a.cleanCtx()
	return a.ctx.LastExecResult
}

func (a *ActionEngine) doModule(m *ModuleBase) error {
	doFunc, has := ModuleFuncMap[ModuleTypeID(m.ModuleTypeID)]
	if !has {
		return fmt.Errorf("no found module type id %d's do function", m.ModuleTypeID)
	}
	return doFunc(a.ctx, m)
}

func (a *ActionEngine) next() {
	a.ctx.nowIndex++
}

func (a *ActionEngine) Step() {
	if a.ctx.nowActionStatus != WaitingStep {
		return
	}
	a.ctx.stepCh <- 0
}

func (a *ActionEngine) Stop() {
	a.cleanCtx()
}

func (a *ActionEngine) cleanCtx() {
	close(a.ctx.stepCh)
}
