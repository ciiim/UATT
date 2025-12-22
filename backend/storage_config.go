package bsd_testtool

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const (
	RESERVED_SPACE_SIZE = 1024 * 1024 * 1 // 1MB

	ViewerExeName = "uart_testtool_viewer.exe"
)

var ReservedSpace [RESERVED_SPACE_SIZE]byte = [RESERVED_SPACE_SIZE]byte{'=', 'U', 'S', 'E', 'R', 'S', 'T', 'A', '='}

type StorageConfig struct {
	Canvas *CanvasConfig `json:"Canvas"`

	AppList []*AppConfig `json:"AppList"`
}

func GetStorageConfig() *StorageConfig {

	var cfg StorageConfig

	jsonStart := ReservedSpace[9:]

	if err := json.Unmarshal(jsonStart, &cfg); err != nil {
		fmt.Printf("GetStorageConfig Failed err: %s", err.Error())
		return nil
	}

	return &cfg
}

func WriteStorageConfig(canvas *CanvasConfig, useAppList []*AppConfig) error {
	var cfg StorageConfig = StorageConfig{
		Canvas:  canvas,
		AppList: useAppList,
	}
	return injectDataIntoExe(ViewerExeName, strings.Join([]string{canvas.CanvasFileName, "configured_viewer.exe"}, "-"), &cfg)
}

func injectDataIntoExe(templatePath, outputPath string, config *StorageConfig) error {
	// 1. 读取模板exe
	exeData, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("读取模板exe失败: %v", err)
	}

	// 2. 序列化用户数据
	configData, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %v", err)
	}

	// 3. 查找数据标记位置
	startMarker := []byte{'=', 'U', 'S', 'E', 'R', 'S', 'T', 'A', '='}

	startPos := bytes.Index(exeData, startMarker)
	if startPos == -1 {
		return fmt.Errorf("找不到起始标记")
	}

	endPos := startPos + RESERVED_SPACE_SIZE

	// 4. 替换标记之间的数据
	newExeData := make([]byte, 0, len(exeData))
	newExeData = append(newExeData, exeData[:startPos+len(startMarker)]...)
	newExeData = append(newExeData, configData...)

	// 填充剩余空间以保持exe大小一致
	remainingSpace := endPos - startPos - len(startMarker) - len(configData)
	if remainingSpace > 0 {
		newExeData = append(newExeData, make([]byte, remainingSpace)...)
	}

	newExeData = append(newExeData, exeData[endPos:]...)

	// 5. 写入新exe
	return os.WriteFile(outputPath, newExeData, 0755)
}
