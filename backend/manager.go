package bsd_testtool

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const (
	DefaultPath   = "./"
	DefaultFolder = "Apps"
)

type Manager struct {

	// App文件存放文件夹
	appFolder string
	apps      []*App
	nowApp    *App
}

var ErrNotFoundApp error = errors.New("could not found app")

var GlobalManager Manager

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

		if strings.Contains(e.Name(), ".app") {
			app, err := getBasicApp(filepath.Join(appFolder, e.Name()))
			if err != nil {
				fmt.Printf("get app[%s] failed: %v\n", e.Name(), err)
				continue
			}
			m.apps = append(m.apps, app)
		}
	}

	return nil
}

func (m *Manager) GetAllAppName() []string {
	appNames := make([]string, len(m.apps))
	for _, app := range m.apps {
		appNames = append(appNames, app.AppName)
	}
	return appNames
}

func (m *Manager) GetNowApp() *App {
	return m.nowApp
}

func (m *Manager) LoadApp(appName string) error {
	var app *App
	for _, a := range m.apps {
		if a.AppName == appName {
			app = a
		}
	}
	if app == nil {
		return ErrNotFoundApp
	}

	m.nowApp = app

	appConfigFile, err := os.Open(m.nowApp.appFileLocation)
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

	m.nowApp.config = appConfig

	return nil
}

func getBasicApp(appFile string) (*App, error) {
	file, err := os.Open(appFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	fileBuffer := make([]byte, fileStat.Size())

	readBytes, err := file.Read(fileBuffer)
	if err != nil || int64(readBytes) != fileStat.Size() {
		return nil, err
	}

	var app App

	if err = json.Unmarshal(fileBuffer, &app); err != nil {
		return nil, err
	}

	app.appFileLocation = appFile

	// fmt.Printf("app.AppName: %v\n", app.AppName)

	return &app, nil
}
