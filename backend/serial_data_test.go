package bsd_testtool_test

import (
	bsd_testtool "bsd_testtool/backend"
	"testing"
)

var test bsd_testtool.SendAction = bsd_testtool.SendAction{
	IOAction: bsd_testtool.IOAction{
		TimeoutMs: 1000,
		Modules: []bsd_testtool.IOModule{
			&bsd_testtool.IOFixedModule{
				bsd_testtool.IOModuleConfigFixed{
					IOModuleConfigBase: bsd_testtool.IOModuleConfigBase{
						ModuleTypeID: 11,
						ModuleUID:    100,
					},
					FixedContent: []int{0xAA, 0x55},
				},
			},
			&bsd_testtool.IOCalcModule{
				bsd_testtool.IOModuleConfigCalc{
					IOModuleConfigBase: bsd_testtool.IOModuleConfigBase{
						ModuleTypeID: 12,
						ModuleUID:    101,
					},
					Mode:                "Calc",
					PlaceholderBytes:    []int{0x00, 0x00},
					CalcTiming:          "Now",
					CalcFunc:            "Length2BytesLE",
					CalcInputModulesUID: []int{102},
				},
			},
			&bsd_testtool.IOCustomModule{
				bsd_testtool.IOModuleConfigCustom{
					IOModuleConfigBase: bsd_testtool.IOModuleConfigBase{
						ModuleTypeID: 11,
						ModuleUID:    102,
					},
					CustomContent: []int{0x00, 0x01, 0x05, 0x11, 0x23},
				},
			},
			&bsd_testtool.IOCalcModule{
				bsd_testtool.IOModuleConfigCalc{
					IOModuleConfigBase: bsd_testtool.IOModuleConfigBase{
						ModuleTypeID: 12,
						ModuleUID:    103,
					},
					Mode:                "Calc",
					PlaceholderBytes:    []int{0x00},
					CalcTiming:          "Post",
					CalcFunc:            "Xor1Bytes",
					CalcInputModulesUID: []int{101, 102},
				},
			},
		},
	},
}

func TestBuildBytes(t *testing.T) {
	res, err := bsd_testtool.BuildSendBytesArray(&test, nil)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("res: %v\n", res)
}

func BenchmarkBuildBytes(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := bsd_testtool.BuildSendBytesArray(&test, nil)
		if err != nil {
			b.Error(err)
			return
		}
		// b.Logf("res: %v\n", res)
	}
}
