package bsd_testtool

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
)

type IOSubModule interface {
	GetUID() ModuleUID
	SetIndex(index int)
	GetIndex() int
	getBasicInfo() (int, []byte)
}

type IOSubModuleBase struct {
	index        int
	ModuleTypeID int `json:"ModuleTypeID"`
	ModuleUID    int `json:"ModuleUID"`
}

type IOSubModuleFill struct {
	IOSubModuleBase
	FillLength int    `json:"FillLength"`
	UseVar     string `json:"UseVar"`
}

type IOSubModuleFixed struct {
	IOSubModuleBase
	FixedContent []int `json:"FixedContent"`
}

type IOSubModuleCalc struct {
	IOSubModuleBase
	Mode                string `json:"Mode"`
	CalcFunc            string `json:"CalcFunc"`
	CalcTiming          string `json:"CalcTiming"`
	PlaceholderBytes    []int  `json:"PlaceholderBytes"`
	CalcInputModulesUID []int  `json:"CalcInputModulesUID"`
}

type IOSubModuleCustom struct {
	IOSubModuleBase
	CustomContent []int `json:"CustomContent"`
}

type IOModuleFeatureField struct {
	TimeoutMs  int
	SubModules []IOSubModule
}

func (i *IOModuleFeatureField) UnmarshalJSON(b []byte) error {
	type Temp struct {
		TimeoutMs  int               `json:"TimeoutMs"`
		SubModules []json.RawMessage `json:"SubModules"`
	}

	var t Temp

	if err := json.Unmarshal(b, &t); err != nil {
		return err
	}

	i.TimeoutMs = t.TimeoutMs

	for _, raw := range t.SubModules {
		var base IOSubModuleBase
		if err := json.Unmarshal(raw, &base); err != nil {
			return err
		}
		switch SubModuleTypeID(base.ModuleTypeID) {
		case FillSubMT:
			var fill IOSubModuleFill
			if err := json.Unmarshal(raw, &fill); err != nil {
				return err
			}
			i.SubModules = append(i.SubModules, &fill)
		case FixedSubMT:
			var fixed IOSubModuleFixed
			if err := json.Unmarshal(raw, &fixed); err != nil {
				return err
			}
			i.SubModules = append(i.SubModules, &fixed)
		case CalcSubMT:
			var calc IOSubModuleCalc
			if err := json.Unmarshal(raw, &calc); err != nil {
				return err
			}
			i.SubModules = append(i.SubModules, &calc)
		case CustomSubMT:
			var custom IOSubModuleCustom
			if err := json.Unmarshal(raw, &custom); err != nil {
				return err
			}
			i.SubModules = append(i.SubModules, &custom)
		default:
			return errors.New("unsupport sub module")
		}
	}

	return nil
}

func (i *IOSubModuleBase) GetUID() ModuleUID {
	return ModuleUID(i.ModuleUID)
}

func (i *IOSubModuleBase) SetIndex(index int) {
	i.index = index
}

func (i *IOSubModuleBase) GetIndex() int {
	return i.index
}

type IOSubModuleCtx struct {
	// UID map, 用来给Calc的模块提供计算来源
	subUIDMap map[ModuleUID]IOSubModule

	// Fill的模块默认放在Now里面，不依赖前后子模块。只依赖ActionEngine的上下文
	calcNowArr []IOSubModule

	// 计算时机在总长度计算完毕，Now的计算完成后，拼接各个模块的结果前
	calcPostArr []IOSubModule

	subBytes [][]byte
}

func doSend(ctx *ActionContext, m *ModuleBase) error {
	return nil
}

func doReceive(ctx *ActionContext, m *ModuleBase) error {
	return nil
}

func (i *IOModuleFeatureField) GetContext() *IOSubModuleCtx {
	ctx := IOSubModuleCtx{
		subUIDMap: make(map[ModuleUID]IOSubModule),
	}
	for i, m := range i.SubModules {
		ctx.subUIDMap[m.GetUID()] = m

		m.SetIndex(i)

		switch t := m.(type) {
		case *IOSubModuleFill:
			ctx.calcNowArr = append(ctx.calcNowArr, m)
		case *IOSubModuleCalc:
			if t.CalcTiming == "Post" {
				ctx.calcPostArr = append(ctx.calcPostArr, m)
			} else {
				ctx.calcNowArr = append(ctx.calcNowArr, m)
			}
		}
	}

	return &ctx
}

func (fixed *IOSubModuleFixed) getBasicInfo() (length int, res []byte) {
	length = len(fixed.FixedContent)
	res = make([]byte, length)
	for i, b := range fixed.FixedContent {
		res[i] = byte(b & 0xFF)
	}
	return
}

func (fill *IOSubModuleFill) getBasicInfo() (length int, res []byte) {
	re := regexp.MustCompile(`\{(\d+|\d:\S+)\}`)
	fmt.Printf("re.FindStringIndex(fill.UseVar): %v\n", re.FindStringIndex(fill.UseVar))
	return
}

// 只返回占位符的长度
func (calc *IOSubModuleCalc) getBasicInfo() (length int, res []byte) {
	length = len(calc.PlaceholderBytes)
	return
}

func (custom *IOSubModuleCustom) getBasicInfo() (length int, res []byte) {
	length = len(custom.CustomContent)
	res = make([]byte, length)
	for i, b := range custom.CustomContent {
		res[i] = byte(b)
	}
	return
}

func (fill *IOSubModuleFill) fill(ctx *ActionContext) (length int, res []byte, err error) {
	return
}

func (calc *IOSubModuleCalc) check(ctx *IOSubModuleCtx, b []byte) (equal bool, err error) {
	return
}

func (calc *IOSubModuleCalc) calc(ctx *IOSubModuleCtx) (length int, res []byte, err error) {

	// 先拿到计算范围内的数据
	subBytes := make([][]byte, len(calc.CalcInputModulesUID))

	for i, uid := range calc.CalcInputModulesUID {
		sm, has := ctx.subUIDMap[ModuleUID(uid)]
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
