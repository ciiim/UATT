package bsd_testtool

type TestCaseFile struct {
	Name       string `toml:"name" mapstructure:"name"`
	ReqContent string `toml:"req_content" mapstructure:"req_content"`
	AckContent string `toml:"ack_content" mapstructure:"ack_content"`
}

type TestItemFile struct {
	Name       string         `toml:"name" mapstructure:"name"`
	ReqPattern string         `toml:"req_pattern" mapstructure:"req_pattern"`
	ExpectAck  string         `toml:"expect_ack" mapstructure:"expect_ack"`
	TestCases  []TestCaseFile `toml:"test_cases" mapstructure:"test_cases"`
	Timeout    int            `toml:"timeout" mapstructure:"timeout"`
}

type ManagerConfig struct {
	Baud        int            `toml:"baud" mapstructure:"baud"`
	VerifyType  int            `toml:"verify_type" mapstructure:"verify_type"`
	CommandHead string         `toml:"command_head" mapstructure:"command_head"`
	TestItems   []TestItemFile `toml:"test_items" mapstructure:"test_items"`
}

/*----*/
/*----*/

type config struct {
	Baud        int
	VerifyType  int
	CommandHead []byte
}

// case自定义字节
type testCustomIndex struct {
	Index       int
	CommandByte byte
}

type testCase struct {
	Name             string
	ReqContent       []testCustomIndex
	ExpectAckContent []testCustomIndex
}

type testItem struct {
	Name string
	// 指令模板
	ReqPattern []byte
	// 期望应答
	ExpectAck []byte

	// 可编辑的索引
	EditableReqIndex []int
	EditableAckIndex []int
	TestCases        []testCase

	Timeout int
}
