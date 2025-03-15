package setting

import "github.com/spf13/viper"

//对读取配置的行为进行封装，以便应用程序的使用

type Setting struct{
	vp *viper.Viper
}
//初始化本项目的基础属性
func NewSetting() (*Setting, error){
	vp := viper.New()
	vp.SetConfigName("config")
	vp.AddConfigPath("configs/")
	vp.SetConfigType("yaml")
	err := vp.ReadInConfig()
	if err != nil {
		return nil,err
	}
	return &Setting{vp},nil
}