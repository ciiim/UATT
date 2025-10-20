package bsd_testtool

import (
	"bsd_testtool/backend/types"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

type IOAction struct {
	// ioModuleCtx IOModuleCtx
	TimeoutMs int        `json:"TimeoutMs"`
	Modules   []IOModule `json:"Modules"`
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

	// send构建计算时填充，发送时使用
	// receive接收后填充，校验时使用
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

func (r *ReceiveAction) GetContext() *IOModuleCtx {
	ctx := IOModuleCtx{
		moduleUIDMap: make(map[types.ActionUID]IOModule),
	}
	for i, m := range r.Modules {
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
	getRes := FmtGetVar(fill.UseVar, ctx)
	if getRes == nil {
		return 0, nil, errors.New("wrong var")
	}

	res = make([]byte, 0)

	switch v := getRes.(type) {
	case []byte:
		res = v
		length = len(v)
	case []int:
		for _, b := range v {
			res = append(res, byte(b))
		}
		length = len(v)
	case byte:
		res = append(res, v)
		length = 1
	case int:
		res = append(res, byte(v))
		length = 1
	default:
		err = errors.New("wrong var type")
	}
	return
}

func (fill *IOFillModule) check(ctx *ActionContext, input []byte) (equal bool, err error) {
	length, res, err := fill.fill(ctx)
	if err != nil {
		return false, err
	}

	if length != len(input) {
		return false, nil
	}

	// 每个字节进行判断，遇到b[] 是-1的就跳过检查
	for i, c := range res {
		if input[i] != c {
			return false, nil
		}
	}

	equal = true
	err = nil

	return
}

func (fixed *IOFixedModule) check(input []byte) (equal bool, err error) {
	if len(fixed.FixedContent) != len(input) {
		return false, nil
	}

	// 每个字节进行判断，遇到b[] 是-1的就跳过检查
	for i, c := range fixed.FixedContent {
		if c == -1 {
			continue
		}
		if int(input[i]) != c {
			fmt.Printf("check failed input: %d, c: %d\n", input[i], c)
			return false, nil
		}
	}

	equal = true
	err = nil

	return
}

func (c *IOCustomModule) check(input []byte) (equal bool, err error) {

	// 可变长度的接收跳过检查
	if c.ReceiveVarLengthModuleUID != -1 && c.ReceiveVarLengthModuleUID != 0 {
		return true, nil
	}

	if len(c.CustomContent) != len(input) {
		return false, nil
	}

	// 每个字节进行判断，遇到b[] 是-1的就跳过检查
	for i, c := range c.CustomContent {
		if c == -1 {
			continue
		}
		if int(input[i]) != c {
			return false, nil
		}
	}

	equal = true
	err = nil

	return
}

func (calc *IOCalcModule) check(ctx *IOModuleCtx, input []byte) (equal bool, err error) {

	length, calRes, err := calc.calc(ctx)
	if err != nil {
		return false, err
	}

	if length != len(input) {
		return false, nil
	}

	// 每个字节进行判断，遇到b[] 是-1的就跳过检查
	for i, c := range calRes {
		if input[i] != c {
			return false, nil
		}
	}
	equal = true
	err = nil

	return
}

func (calc *IOCalcModule) calc(ctx *IOModuleCtx) (length int, res []byte, err error) {

	// 先拿到计算范围内的数据
	subBytes := make([][]byte, len(calc.CalcInputModulesUID))

	for i, uid := range calc.CalcInputModulesUID {
		sm, has := ctx.moduleUIDMap[uid]
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

	sentLength, err := ctx.serial.Write(b)
	if err != nil {
		ctx.SetController(&EnginControllor{nextUID: StopUID})
		return err
	}

	if sentLength != len(b) {
		ctx.SetController(&EnginControllor{nextUID: StopUID})
		return fmt.Errorf("sent %d, expect %d", sentLength, len(b))
	}

	ctx.SetController(&defaultNextControl)

	return nil
}

func (r *ReceiveAction) doAction(ctx *ActionContext) error {

	modCtx := r.GetContext()

	modCtx.subBytes = make([][]byte, len(r.Modules))

	// 一个个模块读取，遇到Custom里指定UID的，就向前查找要读取的字节数
	for i, m := range r.Modules {
		recvLength, _ := m.getBasicInfo()
		if c, ok := m.(*IOCustomModule); ok {
			varLengthUID := c.ReceiveVarLengthModuleUID
			input := modCtx.subBytes[modCtx.moduleUIDMap[varLengthUID].GetIndex()]
			temp := make([]byte, 8)
			copy(temp[8-len(input):], input)
			recvLength = int(binary.BigEndian.Uint64(temp))
		}
		fmt.Printf("recvLength:%d\n", recvLength)
		recvBuffer := make([]byte, recvLength)

		totalLength := 0
		for totalLength != recvLength {
			rLength, err := ctx.serial.Read(recvBuffer[totalLength:])
			if err != nil {
				ctx.SetController(&EnginControllor{nextUID: StopUID})
				return err
			}
			fmt.Printf("read %d bytes, %v\n", rLength, recvBuffer)
			totalLength += rLength
		}

		modCtx.subBytes[i] = recvBuffer
	}

	fullBytes := bytes.Join(modCtx.subBytes, nil)

	ctx.LastSerialBuffer = fullBytes

	uid, err := CheckReceiveBytesArray(r, ctx, modCtx)
	if err != nil {
		ctx.SetController(&EnginControllor{nextUID: StopUID})
		return fmt.Errorf("err:%s, uid:%d", err, uid)
	}

	ctx.SetController(&defaultNextControl)

	return nil
}

type CalcFn func(b []byte) []byte

var CalcFnMap map[string]CalcFn = map[string]CalcFn{
	"Length2BytesLE": Length2BytesLE,
	"Length1BytesLE": Length1BytesLE,
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

func Length1BytesLE(b []byte) []byte {
	return []byte{byte(len(b))}
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
