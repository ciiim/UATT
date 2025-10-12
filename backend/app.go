package bsd_testtool

import (
	"bsd_testtool/backend/types"
	"encoding/json"
	"errors"
	"io"
	"math/rand/v2"
	"os"
	"time"
)

var ErrOutOfIndex = errors.New("out of index")
var ErrNotFoundInsertAction = errors.New("not found insert action")

type SerialConfig struct {
	BaudRate int    `json:"BaudRate"`
	DataBits int    `json:"DataBits"`
	Parity   string `json:"Parity"`
	StopBits int    `json:"StopBits"`
}

type Action struct {
	prev *Action
	next *Action

	actionUID types.ActionUID

	actionTypeStr string

	actionTypeID types.ActionTypeID

	name string

	breakPoint bool

	// 执行interface
	action IAction
}

type App struct {
	AppName string

	logs []io.WriteCloser

	// 文件名
	appFileName string

	uidRand *rand.Rand

	serialConfig *SerialConfig

	logEnable         bool
	logExportEnable   bool
	logExportLoaction string

	firstAction *Action
	actionMap   map[types.ActionUID]*Action
}

func NewApp(appFileName string, config *AppConfig) *App {
	var app App

	app.uidRand = rand.New(rand.NewPCG(uint64(time.Now().Second()), uint64(time.Now().Nanosecond())))

	app.AppName = config.AppName
	app.appFileName = appFileName

	app.logEnable = config.LogEnable
	app.logExportEnable = config.LogExportEnable
	app.logExportLoaction = config.LogExportLoaction

	app.serialConfig = (*SerialConfig)(&config.SerialConfig)

	if app.logEnable {
		if app.logExportEnable {
			// TODO: 新建日志文件
		} else {
			app.logs = []io.WriteCloser{os.Stdout}
		}
	}

	app.actionMap = make(map[types.ActionUID]*Action)

	// 方便遍历
	dummy := new(Action)
	dummy.actionUID = DummyUID
	app.firstAction = dummy
	dummy.prev = dummy
	dummy.next = dummy

	// 把config的actions处理成链表
	for _, action := range config.Actions {
		a := new(Action)

		a.prev = app.firstAction.prev
		a.next = app.firstAction

		app.firstAction.prev.next = a
		app.firstAction.prev = a

		a.breakPoint = action.BreakPoint
		a.actionUID = types.ActionUID(action.ActionUID)
		a.actionTypeID = types.ActionTypeID(action.ActionTypeID)
		a.actionTypeStr = action.ActionType
		a.name = action.Name

		a.action = action.TypeFeatureField.(IConfig).ToAction()

		// map
		app.actionMap[types.ActionUID(action.ActionUID)] = a
	}

	return &app
}

func (a *App) FullUpdateActions(actions []ConfigActionBase) error {
	clear(a.actionMap)
	a.firstAction = nil

	// 方便遍历
	dummy := new(Action)
	dummy.actionUID = DummyUID
	a.firstAction = dummy
	dummy.prev = dummy
	dummy.next = dummy

	// 把config的actions处理成链表
	for _, action := range actions {
		na := new(Action)

		na.prev = a.firstAction.prev
		na.next = a.firstAction

		a.firstAction.prev.next = na
		a.firstAction.prev = na

		na.breakPoint = action.BreakPoint
		na.actionUID = types.ActionUID(action.ActionUID)
		na.actionTypeID = types.ActionTypeID(action.ActionTypeID)
		na.actionTypeStr = action.ActionType
		na.name = action.Name

		na.action = action.TypeFeatureField.(IConfig).ToAction()

		// map
		a.actionMap[types.ActionUID(action.ActionUID)] = na
	}

	return nil
}

func (a *App) GetAction(uid types.ActionUID) *Action {
	return a.actionMap[uid]
}

func (a *App) UpdateAction(uid types.ActionUID, action *Action) error {
	return nil
}

func (a *App) AddAction(insertAfterUID types.ActionUID, actionType types.ActionTypeID, actionTypeStr string, actionName string, action IAction) error {
	actionUID := a.uidRand.Int64()

	insertAction, has := a.actionMap[insertAfterUID]
	if !has {
		return ErrInvalidAction
	}

	newAction := &Action{
		prev:         insertAction,
		next:         insertAction.next,
		actionUID:    int(actionUID),
		actionTypeID: actionType,
		name:         actionName,
		breakPoint:   false,
		action:       action,
	}
	insertAction.next.prev = newAction
	insertAction.next = newAction

	a.actionMap[newAction.actionUID] = newAction

	return nil
}

func (a *App) RemoveAction(uid types.ActionUID) error {
	return nil
}

func (a *App) SwapAction(aUID types.ActionUID, bUID types.ActionUID) error {

	return nil

}

func (a *App) SetBreakPoint(uid types.ActionUID, enable bool) error {
	return nil
}

func (a *App) GetActionList() ([]ConfigActionBaseJson, error) {
	list := make([]ConfigActionBaseJson, 0)

	nowAction := a.firstAction.next

	for nowAction != a.firstAction {
		featField, err := json.Marshal(nowAction.action)
		if err != nil {
			return nil, err
		}
		list = append(list, ConfigActionBaseJson{
			ActionUID:        nowAction.actionUID,
			ActionType:       nowAction.actionTypeStr,
			ActionTypeID:     nowAction.actionTypeID,
			Name:             nowAction.name,
			BreakPoint:       nowAction.breakPoint,
			TypeFeatureField: featField,
		})
		nowAction = nowAction.next
	}

	return list, nil
}
