package parser

import (
	"bsd_testtool/backend/context"
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type StringFmt struct{}

func (s *StringFmt) Sprintf(format string, ctx *context.Context) string {
	re := regexp.MustCompile(`\{(\d+|\d:\S+)\}`)
	return re.ReplaceAllStringFunc(format, func(s string) string {
		if ctx == nil {
			return "[ctx nil]"
		}
		idx := re.FindStringSubmatch(s)[1]
		arrIndexBegin, arrIndexEnd := -1, -1
		if strings.Contains(idx, ":") {
			arrFmtStrings := strings.Split(idx, ":")
			idx = arrFmtStrings[0]
			arrIndexStrings := strings.Split(arrFmtStrings[1], ",")
			if len(arrIndexStrings) == 2 {
				arrIndexEnd, _ = strconv.Atoi(arrIndexStrings[1])
			}
			arrIndexBegin, _ = strconv.Atoi(arrIndexStrings[0])
		switch idx {
		case "0":
			return ctx.LastModuleName
		case "1":
			if arrIndexBegin > len(ctx.LastSerialBuffer)-1 || arrIndexEnd > len(ctx.LastSerialBuffer)-1 {
				return "[index out of range]"
			}
			if arrIndexBegin != -1 && arrIndexEnd == -1 {
				return hex.EncodeToString(ctx.LastSerialBuffer[arrIndexBegin : arrIndexBegin+1])
			} else if arrIndexBegin != -1 {
				return hex.EncodeToString(ctx.LastSerialBuffer[arrIndexBegin:arrIndexEnd])
			} else {
				return hex.EncodeToString(ctx.LastSerialBuffer)
			}
		case "2":
			return ctx.LastExecResult.Error()
		}
		return "[nil]"
	})
}
