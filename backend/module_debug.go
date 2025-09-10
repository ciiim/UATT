package bsd_testtool

import (
	"encoding/json"
	"errors"
)

type PrintModuleFeatureField struct {
	PrintString string `json:"PrintFmt"`
}

func unmarshalDebugModule(moduleTypeID ModuleTypeID, b []byte) (any, error) {
	switch moduleTypeID {
	case PrintMT:
		var p PrintModuleFeatureField
		if err := json.Unmarshal(b, &p); err != nil {
			return nil, err
		}
		return &p, nil
	default:
		return nil, errors.New("unsupport module")
	}
}
