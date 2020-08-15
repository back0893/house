package test

import (
	"main/common"
	"main/model"
	"testing"
)

func TestDb(t *testing.T) {
	common.DbConnections = common.NewDbConnection()
	common.GlobalConfig = common.NewConfig()
	if err := common.GlobalConfig.Load("yaml", "./env.yml"); err != nil {
		panic("配置读取失败")
	}
	house := common.DbConnections.Get("house")
	if house == nil {
		t.Error("nil")
		return
	}
	ml := model.House{}
	if house.Find(&ml).Error != nil {
		t.Error("nil")
		return
	}
	t.Log(ml.ID)
}

func TestCancel(t *testing.T) {
	common.DbConnections = common.NewDbConnection()
	common.GlobalConfig = common.NewConfig()
	if err := common.GlobalConfig.Load("yaml", "./env.yml"); err != nil {
		panic("配置读取失败")
	}
	house := common.DbConnections.Get("house")
	if house == nil {
		t.Error("nil")
		return
	}
	if house.Model(&model.Contract{ID: 3}).Update("cancel", 1).Error != nil {
		t.Fatal("错误")
		return
	}
}
