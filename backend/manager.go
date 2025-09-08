package bsd_testtool

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	DefaultPath = "."
	DefaultName = "default"
)

type Manager struct {
	viper     *viper.Viper
	config    config
	itemMap   map[string]*testItem
	itemArray []*testItem
}

var CManager Manager

func init() {

	CManager.itemMap = make(map[string]*testItem)
	CManager.itemArray = make([]*testItem, 0)

	v := viper.New()
	v.SetConfigType("toml")
	v.SetConfigName(DefaultName)
	v.AddConfigPath(DefaultPath)
	CManager.viper = v
	if err := v.ReadInConfig(); err != nil {
		fmt.Println("[ERROR] missing config file")
		panic(err)
	}
	mc := ManagerConfig{}
	if err := v.Unmarshal(&mc); err != nil {
		fmt.Println("[ERROR] missing necessary config field")
		panic(err)
	}

	if err := CManager.parse(&mc); err != nil {
		fmt.Println("[ERROR] parse config failed")
		panic(err)
	}
}
