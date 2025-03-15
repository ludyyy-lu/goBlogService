package global

import (
	"github.com/ludyyy-lu/goBlogService/pkg/logger"
	"github.com/ludyyy-lu/goBlogService/pkg/setting"
)

var (
	ServerSetting   *setting.ServerSettingS
	AppSetting      *setting.AppSettingS
	DatabaseSetting *setting.DatabaseSettingS
	Logger          *logger.Logger
)
