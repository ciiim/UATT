package bsd_testtool

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
