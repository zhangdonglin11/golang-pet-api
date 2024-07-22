package initialize

import (
	"go.uber.org/zap"
	"golang-pet-api/global"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

func InitConfig() {
	yamlFile, err := ioutil.ReadFile("./settings.yaml")
	//错误就down机
	if err != nil {
		zap.S().Fatalf("读取配置文件失败： %s", err.Error())
		panic(err)
	}
	//绑定值
	err = yaml.Unmarshal(yamlFile, &global.Config)
	if err != nil {
		zap.S().Fatalf("绑定配置失败： %s", err.Error())
		panic(err)
	}
	zap.S().Infof("配置信息: &v", global.Config)
}
