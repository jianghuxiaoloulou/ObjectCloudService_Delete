package object

import (
	"WowjoyProject/ObjectCloudService_Delete/global"
	"WowjoyProject/ObjectCloudService_Delete/internal/model"
	"WowjoyProject/ObjectCloudService_Delete/pkg/general"
)

type ObjectData struct {
	Key   int64           // 目标key
	File  string          // 文件路径
	Type  global.FileType // 文件类型
	Count int             // 文件执行次数
}

// 封装对象相关操作

func NewObject(data global.ObjectData) *ObjectData {
	return &ObjectData{
		Key:   data.InstanceKey,
		File:  data.FilePath,
		Type:  data.Type,
		Count: data.Count,
	}
}

// 删除文件
func (obj *ObjectData) DeleteFile() {
	if general.DeleteFile(obj.File) {
		global.Logger.Info(obj.Key, " :文件删除成功，更新标志")
		model.UpdateDeleteStatus(obj.Key, obj.Type, true)
	} else {
		global.Logger.Info(obj.Key, " :文件删除失败")
		if !ReDo(obj) {
			global.Logger.Info("数据补偿失败", obj.Key)
			global.Logger.Info(obj.Key, " :文件删除失败，更新标志")
			// 上传失败更新数据库
			model.UpdateDeleteStatus(obj.Key, obj.Type, false)
		}
	}
}

// 补偿操作
func ReDo(obj *ObjectData) bool {
	global.Logger.Info("开始补偿操作：", obj.Key)
	if obj.Count < global.GeneralSetting.Count {
		obj.Count += 1
		data := global.ObjectData{
			InstanceKey: obj.Key,
			FilePath:    obj.File,
			Type:        obj.Type,
			Count:       obj.Count,
		}
		global.ObjectDataChan <- data
		return true
	}
	return false
}
