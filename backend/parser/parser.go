package parser

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

// commandParse2Bytes 可阅读字符转成byte字节
// e.g. "aa 55" -> 0xaa, 0x55
func commandParse2Bytes(str string) ([]byte, error) {
	str = strings.ReplaceAll(str, " ", "")
	return hex.DecodeString(str)
}

// parsePlaceHolder 解析模板字符串，将XX转换为testCustomIndex并替换成00
// e.g. "aa 55 XX XX 00" -> "aa 55 00 00 00" and []int{2, 3}
func parsePlaceHolder(patternStr string) (string, []int, error) {
	var (
		replaceStr string
		indexes    []int
	)
	// 以空格分割字符串
	patterns := strings.Split(patternStr, " ")
	for i, pattern := range patterns {
		if pattern == "XX" {
			replaceStr += "00 "
			indexes = append(indexes, i)
		} else {
			replaceStr += pattern + " "
		}
	}
	return strings.TrimSpace(replaceStr), indexes, nil
}

func (m *Manager) parse(mc *ManagerConfig) error {
	b, err := commandParse2Bytes(mc.CommandHead)
	if err != nil {
		return errors.New("[ERROR] error parsing command head")
	}

	m.config.CommandHead = b
	m.config.Baud = mc.Baud
	m.config.VerifyType = mc.VerifyType

	// 解析item
	for _, tomlItem := range mc.TestItems {

		// 解析请求模板字符串
		parsedReqPatternString, editableReqIndex, err := parsePlaceHolder(tomlItem.ReqPattern)
		if err != nil {
			return errors.New("[ERROR] error parsing pattern")
		}
		parsedReqPattern, err := commandParse2Bytes(parsedReqPatternString)
		if err != nil {
			return errors.New("[ERROR] error parsing pattern")
		}

		// 解析期望应答模板字符串
		parsedAckPatternString, editableAckIndex, err := parsePlaceHolder(tomlItem.ExpectAck)
		if err != nil {
			return errors.New("[ERROR] error parsing expect ack")
		}
		parsedAckPattern, err := commandParse2Bytes(parsedAckPatternString)
		if err != nil {
			return errors.New("[ERROR] error parsing expect ack")
		}

		fmt.Println("req index:", editableReqIndex)
		fmt.Println("ack index:", editableAckIndex)

		item := &testItem{
			Name:             tomlItem.Name,
			ReqPattern:       parsedReqPattern,
			ExpectAck:        parsedAckPattern,
			EditableReqIndex: editableReqIndex,
			EditableAckIndex: editableAckIndex,
		}

		if err := item.parseCase(&tomlItem); err != nil {
			return err
		}

		m.itemMap[item.Name] = item
		m.itemArray = append(m.itemArray, item)
	}
	return nil
}

// parseCase 解析CASE
func (i *testItem) parseCase(tomlItem *TestItemFile) error {
	for _, tomlCase := range tomlItem.TestCases {
		var testCaseTemp testCase

		if err := testCaseTemp.initCase(tomlCase.Name, i.EditableReqIndex, i.EditableAckIndex, tomlCase.ReqContent, tomlCase.AckContent); err != nil {
			return err
		}

		i.TestCases = append(i.TestCases, testCaseTemp)
		fmt.Println("case:", testCaseTemp)
	}
	return nil
}

func (c *testCase) initCase(name string, reqIndex, ackIndex []int, reqContent, ackContent string) error {
	c.Name = name

	reqBytes, err := commandParse2Bytes(reqContent)
	if err != nil {
		return errors.New("[ERROR] error parsing req content")
	}

	ackBytes, err := commandParse2Bytes(ackContent)
	if err != nil {
		return errors.New("[ERROR] error parsing ack content")
	}

	for i, index := range reqIndex {
		c.ReqContent = append(c.ReqContent, testCustomIndex{
			Index:       index,
			CommandByte: reqBytes[i],
		})
	}
	for i, index := range ackIndex {
		c.ExpectAckContent = append(c.ExpectAckContent, testCustomIndex{
			Index:       index,
			CommandByte: ackBytes[i],
		})
	}
	return nil
}
