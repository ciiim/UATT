package bsd_testtool

func BuildSendBytesArray(f *IOModuleFeatureField, actionCtx *ActionContext) ([]byte, error) {
	totalLength := 0

	ctx := f.GetContext()

	subBytes := make([][]byte, len(f.SubModules))

	var length int
	var resBytes []byte

	// 依次解析SubModule
	for i, sm := range f.SubModules {
		switch t := sm.(type) {
		case *IOSubModuleFill:
			length, resBytes = t.get()
		case *IOSubModuleFixed:
			length, resBytes = t.get()
		case *IOSubModuleCalc:
			length, resBytes = t.get()
		case *IOSubModuleCustom:
			length, resBytes = t.get()
		}
		totalLength += length
		// 其中可能会有一些resBytes是nil
		subBytes[i] = resBytes
	}

}
