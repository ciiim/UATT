package bsd_testtool

import (
	"bsd_testtool/backend/types"
	"bytes"
	"fmt"
	"regexp"
)

type IOAction struct {
	ioModuleCtx IOModuleCtx
	TimeoutMs   int        `json:"TimeoutMs"`
	Modules     []IOModule `json:"Modules"`
}

type SendAction struct {
	IOAction
}

type ReceiveAction struct {
	IOAction
}

type IOFillModule struct {
	IOModuleConfigFill
}

type IOCalcModule struct {
	IOModuleConfigCalc
}

type IOFixedModule struct {
	IOModuleConfigFixed
}

type IOCustomModule struct {
	IOModuleConfigCustom
}

type IOModule interface {
	GetUID() types.ModuleUID
	SetIndex(idx int)
	GetIndex() int
	getBasicInfo() (int, []byte)
}

type IOModuleCtx struct {
	// UID map, 用来给Calc的模块提供计算来源
	moduleUIDMap map[types.ModuleUID]IOModule

	// Fill的模块默认放在Now里面，不依赖前后子模块。只依赖ActionEngine的上下文
	calcNowArr []IOModule

	// 计算时机在总长度计算完毕，Now的计算完成后，拼接各个模块的结果前
	calcPostArr []IOModule

	subBytes [][]byte
}

func (f *IOFillModule) GetUID() types.ActionUID {
	return f.ModuleUID
}

func (f *IOFillModule) SetIndex(idx int) {
	f.Index = idx
}

func (f *IOFillModule) GetIndex() int {
	return f.Index
}

func (fill *IOFillModule) getBasicInfo() (length int, res []byte) {
	re := regexp.MustCompile(`\{(\d+|\d:\S+)\}`)
	fmt.Printf("re.FindStringIndex(fill.UseVar): %v\n", re.FindStringIndex(fill.UseVar))
	return
}

func (f *IOFixedModule) GetUID() types.ActionUID {
	return f.ModuleUID
}

func (f *IOFixedModule) SetIndex(idx int) {
	f.Index = idx
}

func (f *IOFixedModule) GetIndex() int {
	return f.Index
}

func (fixed *IOFixedModule) getBasicInfo() (length int, res []byte) {
	length = len(fixed.FixedContent)
	res = make([]byte, length)
	for i, b := range fixed.FixedContent {
		res[i] = byte(b & 0xFF)
	}
	return
}

func (calc *IOCalcModule) GetUID() types.ActionUID {
	return calc.ModuleUID
}

func (calc *IOCalcModule) SetIndex(idx int) {
	calc.Index = idx
}

func (calc *IOCalcModule) GetIndex() int {
	return calc.Index
}

// 只返回占位符的长度
func (calc *IOCalcModule) getBasicInfo() (length int, res []byte) {
	length = len(calc.PlaceholderBytes)
	return
}

func (custom *IOCustomModule) GetUID() types.ActionUID {
	return custom.ModuleUID
}

func (custom *IOCustomModule) SetIndex(idx int) {
	custom.Index = idx
}

func (custom *IOCustomModule) GetIndex() int {
	return custom.Index
}

func (custom *IOCustomModule) getBasicInfo() (length int, res []byte) {
	length = len(custom.CustomContent)
	res = make([]byte, length)
	for i, b := range custom.CustomContent {
		res[i] = byte(b)
	}
	return
}

func (s *SendAction) GetContext() *IOModuleCtx {
	ctx := IOModuleCtx{
		moduleUIDMap: make(map[types.ActionUID]IOModule),
	}
	for i, m := range s.Modules {
		ctx.moduleUIDMap[m.GetUID()] = m

		m.SetIndex(i)

		switch t := m.(type) {
		case *IOFillModule:
			ctx.calcNowArr = append(ctx.calcNowArr, m)
		case *IOCalcModule:
			if t.CalcTiming == "Post" {
				ctx.calcPostArr = append(ctx.calcPostArr, m)
			} else {
				ctx.calcNowArr = append(ctx.calcNowArr, m)
			}
		}
	}

	return &ctx
}

func (fill *IOFillModule) fill(ctx *ActionContext) (length int, res []byte, err error) {

	return
}

func (calc *IOCalcModule) check(ctx *IOModuleCtx, b []byte) (equal bool, err error) {
	return
}

func (calc *IOCalcModule) calc(ctx *IOModuleCtx) (length int, res []byte, err error) {

	// 先拿到计算范围内的数据
	subBytes := make([][]byte, len(calc.CalcInputModulesUID))

	for i, uid := range calc.CalcInputModulesUID {
		sm, has := ctx.moduleUIDMap[types.ActionUID(uid)]
		if !has {
			err = fmt.Errorf("cannot find module uid %d", uid)
			return
		}
		subBytes[i] = ctx.subBytes[sm.GetIndex()]
	}

	fullBytes := bytes.Join(subBytes, nil)

	// fmt.Printf("calc %s, fullbytes %v\n", calc.CalcFunc, fullBytes)

	// 丢进计算函数里
	calcFn := GetCalcFn(calc.CalcFunc)
	if calcFn == nil {
		return 0, nil, fmt.Errorf("no [%s] calc function", calc.CalcFunc)
	}

	res = calcFn(fullBytes)
	length = len(res)
	err = nil

	return
}

func (s *SendAction) doAction(ctx *ActionContext) error {

	b, err := BuildSendBytesArray(s, ctx)
	if err != nil {
		ctx.SetController(&EnginControllor{nextUID: StopUID})
		return err
	}

	_, err = ctx.serial.Write(b)
	if err != nil {
		ctx.SetController(&EnginControllor{nextUID: StopUID})
		return err
	}

	ctx.SetController(&defaultNextControl)

	return nil
}

func (r *ReceiveAction) doAction(ctx *ActionContext) error {

	ctx.SetController(&defaultNextControl)

	return nil
}

type CalcFn func(b []byte) []byte

var CalcFnMap map[string]CalcFn = map[string]CalcFn{
	"Length2BytesLE": Length2BytesLE,
	"Xor1Bytes":      Xor1Bytes,
	"Sum1Bytes":      Sum1Bytes,
}

func GetCalcFn(fnName string) CalcFn {
	return CalcFnMap[fnName]
}

func Length2BytesLE(b []byte) []byte {
	length := len(b)
	return []byte{
		byte((length >> 8) & 0xFF),
		byte(length & 0xFF),
	}
}

func Xor1Bytes(b []byte) []byte {
	res := byte(0)
	for _, tmp := range b {
		res ^= tmp
	}
	return []byte{res}
}

func Sum1Bytes(b []byte) []byte {
	res := byte(0)
	for _, tmp := range b {
		res += tmp
	}

	return []byte{res}
}
