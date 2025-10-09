package bsd_testtool

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const (
	DefaultPath   = "./"
	DefaultFolder = "Apps"
)

type Manager struct {
	ctx context.Context

	// App文件存放文件夹
	appFolder     string
	appFileNames  []string
	nowApp        *App
	runningEngine *ActionEngine
}

var ErrNotFoundApp error = errors.New("could not found app")
var ErrAppExist error = errors.New("app exist")

var GlobalManager Manager

func (m *Manager) Startup(ctx context.Context) {
	m.ctx = ctx
}

func (m *Manager) Init(appFolder string) error {
	if appFolder == "" {
		appFolder = "./Apps"
	}
	_, err := os.Stat(appFolder)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			if err = os.Mkdir(appFolder, os.ModeDir); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	folder, err := os.Open(appFolder)
	if err != nil {
		return err
	}
	defer folder.Close()

	entry, err := folder.ReadDir(-1)
	if err != nil {
		return err
	}

	for _, e := range entry {
		if e.IsDir() {
			continue
		}
		if strings.Contains(e.Name(), ".json") {
			m.appFileNames = append(m.appFileNames, e.Name())
		}
	}

	m.appFolder = appFolder

	return nil
}

func (m *Manager) GetAllAppName() []string {
	fmt.Printf("m.appFileNames: %v\n", m.appFileNames)
	return m.appFileNames
}

func (m *Manager) GetNowApp() *App {
	return m.nowApp
}

func (m *Manager) Start() error {
	if m.nowApp == nil {
		return ErrNotFoundApp
	}

	if err := GlobalSerial.OpenSerial(); err != nil {
		return err
	}

	engine := NewActionEngine(m.nowApp)

	if err := engine.PreCompile(); err != nil {
		return err
	}

	m.runningEngine = engine

	engine.StartAsync()

	return nil
}

func (m *Manager) Stop() error {

	if m.runningEngine == nil {
		return nil
	}

	m.runningEngine.Stop()

	m.runningEngine = nil

	return GlobalSerial.CloseSerial()
}

func (m *Manager) GetAppSettings() (AppConfigSettings, error) {
	if m.nowApp == nil {
		return AppConfigSettings{}, ErrNotFoundApp
	}

	var settings AppConfigSettings

	settings.AppName = m.nowApp.appFileName
	settings.LogEnable = m.nowApp.logEnable
	settings.LogExportEnable = m.nowApp.logExportEnable
	settings.LogExportLoaction = m.nowApp.logExportLoaction
	settings.SerialConfig = struct {
		BaudRate int    "json:\"BaudRate\""
		DataBits int    "json:\"DataBits\""
		Parity   string "json:\"Parity\""
		StopBits int    "json:\"StopBits\""
	}(*m.nowApp.serialConfig)

	return settings, nil
}

func (m *Manager) GetActionList() ([]ConfigActionBaseJson, error) {
	if m.nowApp == nil {
		return nil, ErrNotFoundApp
	}

	return m.nowApp.GetActionList()
}

func (m *Manager) CreateApp(basic AppConfigSettings) error {

	appName := basic.AppName + ".json"

	var loadAppName string
	for _, a := range m.appFileNames {
		if a == appName {
			loadAppName = a
			break
		}
	}
	if loadAppName != "" {
		return ErrAppExist
	}

	basicJson, err := json.Marshal(&basic)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(m.appFolder, appName))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(basicJson)
	if err != nil {
		return err
	}

	m.appFileNames = append(m.appFileNames, appName)

	return nil
}

func (m *Manager) DeleteApp(appFileName string) error {
	var loadAppName string
	for _, a := range m.appFileNames {
		if a == appFileName {
			loadAppName = a
			break
		}
	}
	if loadAppName == "" {
		return ErrNotFoundApp
	}

	if m.nowApp != nil && m.nowApp.appFileName == appFileName {
		m.nowApp = nil
	}

	if err := os.Remove(filepath.Join(m.appFolder, appFileName)); err != nil {
		return err
	}

	m.appFileNames = slices.DeleteFunc(m.appFileNames, func(fileName string) bool {
		return appFileName == fileName
	})

	return nil
}

func (m *Manager) LoadApp(appName string) error {
	var loadAppName string
	for _, a := range m.appFileNames {
		if a == appName {
			loadAppName = a
		}
	}
	if loadAppName == "" {
		return ErrNotFoundApp
	}
	appConfigFile, err := os.Open(filepath.Join(m.appFolder, loadAppName))
	if err != nil {
		return err
	}
	defer appConfigFile.Close()

	appConfigStat, err := appConfigFile.Stat()
	if err != nil {
		return err
	}

	appConfigJson := make([]byte, appConfigStat.Size())

	readBytes, err := appConfigFile.Read(appConfigJson)
	if err != nil {
		return err
	}
	if readBytes != int(appConfigStat.Size()) {
		return errors.New("read app config error")
	}

	appConfig, err := ParseAppConfig(appConfigJson)
	if err != nil {
		return err
	}

	m.nowApp = NewApp(loadAppName, appConfig)

	return nil
}

func (m *Manager) SyncAppSettings(settings AppConfigSettings) error {
	if m.nowApp == nil {
		return ErrNotFoundApp
	}

	nowApp := m.nowApp

	nowApp.logEnable = settings.LogEnable
	nowApp.logExportEnable = settings.LogExportEnable
	nowApp.logExportLoaction = settings.LogExportLoaction
	nowApp.serialConfig = (*SerialConfig)(&settings.SerialConfig)

	return nil
}

// 同步前端Actions到后端App的Action链表里
func (m *Manager) SyncActions(actions []ConfigActionBaseJson) error {
	if m.nowApp == nil {
		return ErrNotFoundApp
	}

	newActions := make([]ConfigActionBase, len(actions))

	for i, a := range actions {
		newActions[i] = a.ToBase()
	}

	return m.nowApp.FullUpdateActions(newActions)
}

func (m *Manager) SaveApp() error {

	if m.nowApp == nil {
		return ErrNotFoundApp
	}

	//重新组装AppConfig
	config := AppConfig{
		AppName: m.nowApp.AppName,
		SerialConfig: struct {
			BaudRate int    "json:\"BaudRate\""
			DataBits int    "json:\"DataBits\""
			Parity   string "json:\"Parity\""
			StopBits int    "json:\"StopBits\""
		}(*m.nowApp.serialConfig),
		LogEnable:         m.nowApp.logEnable,
		LogExportEnable:   m.nowApp.logExportEnable,
		LogExportLoaction: m.nowApp.logExportLoaction,
		Actions:           make([]ConfigActionBase, 0),
	}

	// 组装Actions
	nowAction := m.nowApp.firstAction.next

	for nowAction != m.nowApp.firstAction {
		configActionBase := ConfigActionBase{
			ActionUID:        nowAction.actionUID,
			ActionTypeID:     nowAction.actionTypeID,
			ActionType:       nowAction.actionTypeStr,
			Name:             nowAction.name,
			TypeFeatureField: nowAction.action,
		}
		config.Actions = append(config.Actions, configActionBase)
		nowAction = nowAction.next
	}

	configJson, err := json.Marshal(&config)
	if err != nil {
		return err
	}

	//覆盖写入
	configFile, err := os.Create(filepath.Join(m.appFolder, m.nowApp.appFileName))
	if err != nil {
		return err
	}
	defer configFile.Close()

	if _, err = configFile.Write(configJson); err != nil {
		return err
	}

	return nil
}

func (m *Manager) GetAllSerial() ([]string, error) {
	return GlobalSerial.GetAllPort()
}

func (m *Manager) SelectSerialCom(com string) {
	GlobalSerial.SelectPort(com)
}
