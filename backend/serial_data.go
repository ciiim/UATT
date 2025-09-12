package bsd_testtool

import (
	"bytes"
	"fmt"
)

func BuildSendBytesArray(f *IOModuleFeatureField, actionCtx *ActionContext) ([]byte, error) {
	totalLength := 0

	ctx := f.GetContext()

	ctx.subBytes = make([][]byte, len(f.SubModules))

	var fullBytes []byte

	var length int
	var resBytes []byte
	var err error

	// 依次解析SubModule
	for i, sm := range f.SubModules {
		length, resBytes = sm.getBasicInfo()

		totalLength += length
		// 其中可能会有一些resBytes是nil
		ctx.subBytes[i] = resBytes
	}

	// fmt.Printf("total length:%d\n", totalLength)

	// 执行now的计算模块
	for _, sm := range ctx.calcNowArr {
		switch t := sm.(type) {
		case *IOSubModuleFill:
			_, resBytes, err = t.fill(actionCtx)
		case *IOSubModuleCalc:
			_, resBytes, err = t.calc(ctx)
		default:
			return nil, fmt.Errorf("unsupport sub module type, UID: %d", sm.GetUID())
		}
		if err != nil {
			return nil, err
		}

		ctx.subBytes[sm.GetIndex()] = resBytes
	}

	// 执行post的计算模块
	for _, sm := range ctx.calcPostArr {
		switch t := sm.(type) {
		case *IOSubModuleFill:
			_, resBytes, err = t.fill(actionCtx)
		case *IOSubModuleCalc:
			_, resBytes, err = t.calc(ctx)
		default:
			return nil, fmt.Errorf("unsupport sub module type, UID: %d", sm.GetUID())
		}
		if err != nil {
			return nil, err
		}

		ctx.subBytes[sm.GetIndex()] = resBytes
	}

	// 组装成完整数组
	fullBytes = bytes.Join(ctx.subBytes, nil)

	if len(fullBytes) != totalLength {
		return nil, fmt.Errorf("wrong length %d, expect %d bytes", len(fullBytes), totalLength)
	}

	return fullBytes, nil
}

// 如果检查不过，返回检查不过的模块UID
func CheckReceiveBytesArray(f *IOModuleFeatureField, actionCtx *ActionContext, checkBytes []byte) (ModuleUID, error) {

	// 依次检查非计算类模块同时拆分待检查的切片

	// 检查Now模块计算结果

	// 检查Post模块计算结果

	// bytes里面的-1是代表不检查这一项

	return 0, nil
}
