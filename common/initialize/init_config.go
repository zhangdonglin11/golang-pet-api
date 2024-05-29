package initialize

import (
	"fmt"
	"golang-pet-api/common/global"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

func InitConfig() {
	yamlFile, err := ioutil.ReadFile("./settings.yaml")
	//错误就down机
	if err != nil {
		fmt.Println("配置文件读取失败！")
		panic(err)
	}
	//绑定值
	err = yaml.Unmarshal(yamlFile, &global.Config)
	if err != nil {
		fmt.Println("配置文件绑定失败！")
		panic(err)
	}
	fmt.Println("配置文件绑定成功！")
}
