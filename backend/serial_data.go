package bsd_testtool

import (
	"bsd_testtool/backend/types"
	"bytes"
	"errors"
	"fmt"
)

var ErrCheckFailed = errors.New("check failed")

func BuildSendBytesArray(s *SendAction, actionCtx *ActionContext) ([]byte, error) {
	totalLength := 0

	ctx := s.GetContext()

	ctx.subBytes = make([][]byte, len(s.Modules))

	var fullBytes []byte

	var length int
	var resBytes []byte
	var err error

	// 依次解析SubModule
	for i, sm := range s.Modules {
		switch sm.(type) {
		// Fill的先不算总长度，读出了数据之后再计算进去，不然要读两次数据，浪费空间
		case *IOFillModule:
			break
		default:
			length, resBytes = sm.getBasicInfo()

			totalLength += length
			// 其中可能会有一些resBytes是nil
			ctx.subBytes[i] = resBytes
		}

	}

	// fmt.Printf("total length:%d\n", totalLength)

	// 执行now的计算模块
	for _, sm := range ctx.calcNowArr {
		switch t := sm.(type) {
		case *IOFillModule:
			_, resBytes, err = t.fill(actionCtx)
			totalLength += len(resBytes)
		case *IOCalcModule:
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
		case *IOFillModule:
			_, resBytes, err = t.fill(actionCtx)
			totalLength += len(resBytes)
		case *IOCalcModule:
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

	fmt.Printf("send bytes %v\n", fullBytes)

	return fullBytes, nil
}

// 如果检查不过，返回第一个检查不过的模块UID
func CheckReceiveBytesArray(r *ReceiveAction, actionCtx *ActionContext, modCtx *IOModuleCtx) (types.ActionUID, error) {

	// 检查模块计算结果
	for _, sm := range modCtx.moduleUIDMap {
		switch t := sm.(type) {
		case *IOFillModule:
			e, err := t.check(actionCtx, modCtx.subBytes[t.GetIndex()])
			if err != nil {
				return t.GetUID(), err
			}
			if !e {
				return t.GetUID(), ErrCheckFailed
			}
		case *IOFixedModule:
			e, err := t.check(modCtx.subBytes[t.GetIndex()])
			if err != nil {
				return t.GetUID(), err
			}
			if !e {
				return t.GetUID(), ErrCheckFailed
			}
		case *IOCustomModule:
			e, err := t.check(modCtx.subBytes[t.GetIndex()])
			if err != nil {
				return t.GetUID(), err
			}
			if !e {
				return t.GetUID(), ErrCheckFailed
			}
		case *IOCalcModule:
			e, err := t.check(modCtx, modCtx.subBytes[t.GetIndex()])
			if err != nil {
				return t.GetUID(), err
			}
			if !e {
				return t.GetUID(), ErrCheckFailed
			}
		default:
			return -1, fmt.Errorf("unsupport sub module type, UID: %d", sm.GetUID())
		}
	}

	return -1, nil
}
