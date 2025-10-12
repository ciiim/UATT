package bsd_testtool

import (
	"bsd_testtool/backend/types"
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
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

type EngineMode int

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

const (
	Continuous EngineMode = iota

	Step
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

	engineMode EngineMode

	log []io.WriteCloser

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

type wailsLog struct {
	ctx context.Context
}

func (w *wailsLog) Write(b []byte) (int, error) {
	runtime.EventsEmit(w.ctx, "runtime-log", fmt.Sprintf("[INFO] %s\n", b))
	return len(b), nil
}

func (w *wailsLog) Close() error {
	return nil
}

func getWailsRuntimeLog(ctx context.Context) io.WriteCloser {
	return &wailsLog{
		ctx: ctx,
	}
}

type ActionEngine struct {
	app *App

	wailsCtx context.Context

	ctx *ActionContext

	stepCh chan struct{}

	stoppedCh chan struct{}

	stopCB func()
}

func NewActionEngine(app *App, ctx context.Context, stopCB func()) *ActionEngine {
	return &ActionEngine{
		app:       app,
		stepCh:    make(chan struct{}, 2),
		stoppedCh: make(chan struct{}),
		wailsCtx:  ctx,
		stopCB:    stopCB,
	}
}

func (a *ActionEngine) PreCompile() error {
	if a.app == nil {
		return ErrNotFoundApp
	}

	a.ctx = &ActionContext{
		log:              a.app.logs,
		actionFirst:      a.app.firstAction,
		actionUIDMap:     a.app.actionMap,
		variableMap:      make(map[string]*variable),
		controlFlowMap:   make(map[types.ActionUID]*controlFlow),
		labelMap:         make(map[string]types.ActionUID),
		nowAction:        nil,
		LastActionName:   "NIL",
		LastSerialBuffer: nil,
		LastExecResult:   nil,
		nowActionStatus:  Idle,
		serial:           &GlobalSerial,
	}

	a.ctx.log = append(a.ctx.log, getWailsRuntimeLog(a.wailsCtx))

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
			case VarNumber:
				v.v = tmp.VarNumberValue
			case VarNumberArray:
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
	a.ctx.engineMode = Continuous

	return a.innerStart()
}

func (a *ActionEngine) StartAsync() {
	a.ctx.nowAction = a.ctx.actionFirst.next
	a.ctx.nowActionStatus = Running
	a.ctx.engineMode = Continuous
	go a.innerStart()
}

func (a *ActionEngine) StepAsyncStart() {
	a.ctx.nowAction = a.ctx.actionFirst.next
	a.ctx.nowActionStatus = Running
	a.ctx.engineMode = Step
	go a.innerStart()
}

type ActionReportStatus = int

const (
	Ready ActionReportStatus = iota
	Now
	Done
)

type ActionReport struct {
	ActionName string          `json:"ActionName"`
	ActionUID  types.ActionUID `json:"ActionUID"`
	Result     string          `json:"Result"`
}

func (a *ActionEngine) innerStart() error {

	if a.wailsCtx != nil {
		runtime.EventsEmit(a.wailsCtx, "begin-action")
	}
	for a.ctx.nowAction.actionUID != DummyUID && a.ctx.nowActionStatus != WaitingStop {
		action := a.ctx.nowAction

		fmt.Printf("now action name %s, UID:%d\n", action.name, action.actionUID)
		// time.Sleep(time.Millisecond * 500)

		if a.ctx.engineMode == Continuous {
			if action.breakPoint {
				a.ctx.nowActionStatus = WaitingStep
				_, ok := <-a.stepCh
				if !ok {
					break
				}
			}
		} else if a.ctx.engineMode == Step {
			// 单步模式下，汇报即将执行的Action信息 用于引导UI显示
			if a.wailsCtx != nil {
				runtime.EventsEmit(a.wailsCtx, "ready-action", ActionReport{
					ActionName: action.name,
					ActionUID:  action.actionUID,
					Result:     "",
				})
			}
			a.ctx.nowActionStatus = WaitingStep
			_, ok := <-a.stepCh
			if !ok {
				break
			}
		}

		if a.wailsCtx != nil {
			runtime.EventsEmit(a.wailsCtx, "now-action", ActionReport{
				ActionName: action.name,
				ActionUID:  action.actionUID,
				Result:     "",
			})
		} else {
			println("no wails ctx")
		}

		a.ctx.LastExecResult = a.doAction(action.action)

		a.ctx.LastActionName = action.name

		// 汇报执行结果
		if a.wailsCtx != nil {
			runtime.EventsEmit(a.wailsCtx, "done-action", ActionReport{
				ActionName: action.name,
				ActionUID:  action.actionUID,
				Result: func() string {
					if a.ctx.LastExecResult != nil {
						return a.ctx.LastExecResult.Error()
					} else {
						return "success"
					}
				}(),
			})
		} else {
			println("no wails ctx")
		}

		if err := a.control(); err != nil {
			a.ctx.nowActionStatus = ErrorStopped
			return err
		}
	}

	if a.ctx.LastExecResult != nil {
		a.ctx.nowActionStatus = ErrorStopped
	}
	defer func() {
		select {
		case a.stoppedCh <- struct{}{}:
			a.cleanCtx()
			return
		case <-time.After(time.Millisecond * 50):
			a.cleanCtx()
			return
		}
	}()
	a.ctx.nowActionStatus = Stopped
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

	switch con.nextUID {
	case StopUID:
		a.ctx.nowAction = a.ctx.actionFirst
	case DummyUID:
		a.ctx.nowAction = a.ctx.nowAction.next
	default:
		a.ctx.nowAction = a.ctx.actionUIDMap[con.nextUID]
	}
	return nil
}

func (a *ActionEngine) Step() error {
	if a.ctx.nowActionStatus != WaitingStep {
		return nil
	}
	a.stepCh <- struct{}{}
	return nil
}

func (a *ActionEngine) Stop() {

	if a.ctx.engineMode == Step && a.ctx.nowActionStatus == WaitingStep {
		a.stepCh <- struct{}{}
	}
	println("now status", a.ctx.nowActionStatus)
	if a.ctx.nowActionStatus != Stopped {
		a.ctx.nowActionStatus = WaitingStop
		println("waiting stop")
		<-a.stoppedCh
	}
}

func (a *ActionEngine) cleanCtx() {
	close(a.stepCh)
	close(a.stoppedCh)
	runtime.EventsEmit(a.wailsCtx, "stopped")
	if a.stopCB != nil {
		a.stopCB()
	}
}
