package test

import (
	"main/common"
	v1 "main/routes/v1"
	"testing"
)

func init() {
	common.DbConnections = common.NewDbConnection()
	common.GlobalConfig = common.NewConfig()
	if err := common.GlobalConfig.Load("yaml", "./env.yml"); err != nil {
		panic("配置读取失败")
	}
}
func TestContractCheckCososPass(t *testing.T) {
	start, _ := common.ParseTime("2006-01-02", "2020-08-16")
	end, _ := common.ParseTime("2006-01-02", "2020-09-16")
	if err := v1.CheckContractCocos(start, end, 7); err != nil {
		t.Fatal("检测失败")
	}
}
func TestContractCheckCososFaile1(t *testing.T) {
	start, _ := common.ParseTime("2006-01-02", "2020-01-05")
	end, _ := common.ParseTime("2006-01-02", "2020-02-05")
	if err := v1.CheckContractCocos(start, end, 7); err == nil {
		t.Fatal("检测失败")
	}
}

func TestContractCheckCososFaile2(t *testing.T) {
	start, _ := common.ParseTime("2006-01-02", "2020-01-01")
	end, _ := common.ParseTime("2006-01-02", "2020-02-05")
	if err := v1.CheckContractCocos(start, end, 7); err == nil {
		t.Fatal("检测失败")
	}
}
func TestContractCheckCososFaile3(t *testing.T) {
	start, _ := common.ParseTime("2006-01-02", "2020-01-01")
	end, _ := common.ParseTime("2006-01-02", "2020-01-08")
	if err := v1.CheckContractCocos(start, end, 7); err == nil {
		t.Fatal("检测失败")
	}
}
