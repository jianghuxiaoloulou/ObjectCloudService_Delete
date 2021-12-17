package global

import (
	"WowjoyProject/ObjectCloudService_Delete/pkg/logger"
	"WowjoyProject/ObjectCloudService_Delete/pkg/setting"
)

var (
	ServerSetting   *setting.ServerSettingS
	GeneralSetting  *setting.GeneralSettingS
	DatabaseSetting *setting.DatabaseSettingS
	Logger          *logger.Logger
)
