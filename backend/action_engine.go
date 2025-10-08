package bsd_testtool

import (
	"bsd_testtool/backend/types"
	"errors"
	"fmt"
	"io"
	"time"
)

const (
	StopUID  types.ActionUID = -100
	DummyUID types.ActionUID = -1
)

var ErrInvalidAction = errors.New("invalid action")
var ErrInvalidLabel = errors.New("invalid label")
var ErrInvalidExpression = errors.New("invalid expression")
var ErrNoController = errors.New("no controller")

type ControlType = types.ActionTypeID

type VarType = string

const (
	VarNumber      VarType = "number"
	VarNumberArray VarType = "array"
)

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

	expressionAst *AstNode

	TrueUID      types.ActionUID // IF 成立执行的下一个 UID
	FalseUID     types.ActionUID // IF 不成立执行的下一个 UID
	EndUID       types.ActionUID // FOR/IF 最终要跳到的 BLOCKEND
	LoopStartUID types.ActionUID // FOR 循环开始的 UID
}

type EnginControllor struct {
	nextUID types.ActionUID
}

var defaultNextControl = EnginControllor{nextUID: -1}

type ActionContext struct {
	actionFirst *Action

	// 变量表
	variableMap map[string]*variable

	// UID表，用于引用的快速获取
	actionUIDMap map[types.ActionUID]*Action

	// 控制流表，key为{if, for}action的UID，用于if else for的流程跳转
	controlFlowMap map[types.ActionUID]*controlFlow

	// 标签表，用于goto跳转
	labelMap map[string]types.ActionUID

	nowAction *Action

	nowActionStatus ActionStatus

	log io.WriteCloser

	stepCh chan int

	controller *EnginControllor

	serial *Serial

	LastActionName   string
	LastSerialBuffer []byte
	LastExecResult   error
}

func (ctx *ActionContext) DefaultNextAction() {
	ctx.controller = &defaultNextControl
}

func (ctx *ActionContext) SetController(con *EnginControllor) {
	ctx.controller = con
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
		log:              a.app.log,
		actionFirst:      a.app.firstAction,
		actionUIDMap:     a.app.actionMap,
		variableMap:      make(map[string]*variable),
		controlFlowMap:   make(map[types.ActionUID]*controlFlow),
		labelMap:         make(map[string]types.ActionUID),
		stepCh:           make(chan int, 2),
		nowAction:        nil,
		LastActionName:   "NIL",
		LastSerialBuffer: nil,
		LastExecResult:   nil,
		nowActionStatus:  Idle,
		serial:           &GlobalSerial,
	}

	type tempControlInfo struct {
		ct types.ActionTypeID
		i  types.ActionUID
	}
	controlFlowStack := make([]tempControlInfo, 0)

	nowAction := a.ctx.actionFirst.next
	// 遍历代码
	for nowAction != a.ctx.actionFirst {

		t := types.ActionTypeID(nowAction.actionTypeID)

		// 变量表
		if t == types.DeclareAT {
			tmp, is := nowAction.action.(*DeclareAction)
			if !is {
				return ErrInvalidAction
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
		if t == types.GotoLabelAT {
			tmp, is := nowAction.action.(*LabelAction)
			if !is {
				return ErrInvalidAction
			}
			a.ctx.labelMap[tmp.LabelName] = nowAction.actionUID
		}

		// 控制流表
		switch t {
		case types.IfAT:

			ifField, has := nowAction.action.(*IfAction)
			if !has {
				return ErrInvalidAction
			}

			parser := NewParser(ifField.Condition)
			if parser == nil {
				return ErrInvalidExpression
			}

			// 加进map里面
			a.ctx.controlFlowMap[nowAction.actionUID] = &controlFlow{
				controlType:   t,
				expressionAst: parser.GetAST(),
				TrueUID:       nowAction.actionUID + 1,
			}

			// 压栈
			controlFlowStack = append(controlFlowStack, tempControlInfo{
				ct: t,
				i:  nowAction.actionUID,
			})

		case types.ElseAT:
			cf, has := a.ctx.controlFlowMap[controlFlowStack[len(controlFlowStack)-1].i]
			if !has {
				return errors.New("else without matching if")
			}
			cf.FalseUID = nowAction.actionUID + 1

			// 临时存IF
			ifcf := controlFlowStack[len(controlFlowStack)-1]

			// 出栈IF
			controlFlowStack = controlFlowStack[:len(controlFlowStack)-1]

			// 入栈ELSE, 但是index还是if的index，不然endblock会找不到对应的控制流程结构体
			controlFlowStack = append(controlFlowStack, tempControlInfo{
				ct: t,
				i:  ifcf.i,
			})
		case types.ForLabelAT:
			// 压栈
			controlFlowStack = append(controlFlowStack, tempControlInfo{
				ct: t,
				i:  nowAction.actionUID,
			})
		case types.EndBlockAT:

			top := controlFlowStack[len(controlFlowStack)-1]

			if top.ct == types.IfAT || top.ct == types.ElseAT {
				cf, has := a.ctx.controlFlowMap[controlFlowStack[len(controlFlowStack)-1].i]
				// 忽略多余的EndBlock，不算做错误
				if !has {
					continue
				}
				cf.EndUID = nowAction.actionUID + 1
			} else if top.ct == types.ForLabelAT {
				// 如果是for类型，还得往map里插入一个可以通过endblock的index查找的元素
				// 加进map里面
				a.ctx.controlFlowMap[nowAction.actionUID] = &controlFlow{
					controlType: t,

					// 这里不+1，因为要回到for上做检查
					LoopStartUID: top.i,

					EndUID: nowAction.actionUID + 1,
				}
			}

			//出栈
			controlFlowStack = controlFlowStack[:len(controlFlowStack)-1]
		}

		nowAction = nowAction.next
	}

	if len(controlFlowStack) != 0 {
		return fmt.Errorf("if or for statment not closed, at index:%d", controlFlowStack[0].i)
	}

	return nil
}

func (a *ActionEngine) GetStopReason() error {
	return a.ctx.LastExecResult
}

func (a *ActionEngine) StartSync() error {
	a.ctx.nowAction = a.ctx.actionFirst.next
	a.ctx.nowActionStatus = Running
	return a.innerStart()
}

func (a *ActionEngine) StartAsync() {
	a.ctx.nowAction = a.ctx.actionFirst.next
	a.ctx.nowActionStatus = Running
	go a.innerStart()
}

func (a *ActionEngine) innerStart() error {
	for a.ctx.nowAction.actionUID != DummyUID && a.ctx.nowActionStatus != WaitingStop {
		action := a.ctx.nowAction

		fmt.Printf("now action name %s, UID:%d\n", action.name, action.actionUID)
		time.Sleep(time.Millisecond * 500)

		if action.breakPoint {
			a.ctx.nowActionStatus = WaitingStep
			_, ok := <-a.ctx.stepCh
			if !ok {
				break
			}

		}

		a.ctx.LastExecResult = a.doAction(action.action)

		a.ctx.LastActionName = action.name

		if err := a.control(); err != nil {
			a.ctx.nowActionStatus = ErrorStopped
			return err
		}
	}

	if a.ctx.LastExecResult != nil {
		a.ctx.nowActionStatus = ErrorStopped
	}
	a.cleanCtx()
	return a.ctx.LastExecResult
}

func (a *ActionEngine) doAction(action IAction) error {
	return action.doAction(a.ctx)
}

func (a *ActionEngine) control() error {
	if a.ctx.controller == nil {
		return ErrNoController
	}

	con := a.ctx.controller
	a.ctx.controller = nil

	if con.nextUID == StopUID {
		a.ctx.nowAction = a.ctx.actionFirst
	} else if con.nextUID == DummyUID {
		a.ctx.nowAction = a.ctx.nowAction.next
	} else {
		a.ctx.nowAction = a.ctx.actionUIDMap[con.nextUID]
	}
	return nil
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
