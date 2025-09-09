package bsd_testtool

type VerifyFn func(command []byte) []byte

type IVerifyFn interface {
	Verify(command []byte) []byte
	Len() int
}

type Verifier struct {
	now   int
	funcs map[int]IVerifyFn
}

var verifier Verifier

func init() {
	verifier = Verifier{
		funcs: make(map[int]IVerifyFn),
	}
	verifier.funcs[0] = BytesAdd{length: 1}
}

func (v *Verifier) Verify(command []byte) []byte {
	return v.funcs[v.now].Verify(command)
}

func (v *Verifier) Len() int {
	return v.funcs[v.now].Len()
}

// BytesAdd 按字节和校验
type BytesAdd struct {
	length int
}

func (b BytesAdd) Verify(command []byte) []byte {
	checksum := byte(0)
	for i := range command {
		checksum += command[i]
	}
	return []byte{checksum}
}

func (b BytesAdd) Len() int {
	return b.length
}
