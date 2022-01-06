package model

import (
	"WowjoyProject/ObjectCloudService_Delete/global"
	"WowjoyProject/ObjectCloudService_Delete/pkg/general"
)

// 自动获取需要删除的数据
func GetDeleteData() {
	global.Logger.Info("***开始获取需要删除的数据***")
	sql := `select ins.instance_key,ins.file_name,im.img_file_name,sl.ip,sl.s_virtual_dir
	from instance ins
	left join image im on im.instance_key = ins.instance_key
	left join study_location sl on sl.n_station_code = ins.location_code
	left join file_remote fr on ins.instance_key = fr.instance_key
	where fr.dcm_file_exist = 1 and (fr.dcm_file_exist_obs_local = 1 or fr.dcm_file_exist_obs_cloud = 1)
	order by fr.dcm_update_time_retrieve asc limit ?;`
	// global.Logger.Debug(sql)
	rows, err := global.DBEngine.Query(sql, global.GeneralSetting.MaxTasks)
	if err != nil {
		global.Logger.Fatal(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		key := KeyData{}
		_ = rows.Scan(&key.instance_key, &key.dcmfile, &key.jpgfile, &key.ip, &key.virpath)
		if key.jpgfile.String != "" {
			file_path := general.GetFilePath(key.jpgfile.String, key.ip.String, key.virpath.String)
			global.Logger.Info("需要处理的文件名：", file_path)
			data := global.ObjectData{
				InstanceKey: key.instance_key.Int64,
				FilePath:    file_path,
				Type:        global.JPG,
				Count:       1,
			}
			global.ObjectDataChan <- data
		} else {
			global.Logger.Error(key.instance_key.Int64, ": JPG文件不存在")
			UpdateDeleteStatus(key.instance_key.Int64, global.JPG, true)
		}
		if key.dcmfile.String != "" {
			file_path := general.GetFilePath(key.dcmfile.String, key.ip.String, key.virpath.String)
			global.Logger.Info("需要处理的文件名：", file_path)
			data := global.ObjectData{
				InstanceKey: key.instance_key.Int64,
				FilePath:    file_path,
				Type:        global.DCM,
				Count:       1,
			}
			global.ObjectDataChan <- data
		} else {
			global.Logger.Error(key.instance_key.Int64, ": DCM文件不存在")
			UpdateDeleteStatus(key.instance_key.Int64, global.DCM, true)
		}
	}
}

// 跟新删除文件标志
func UpdateDeleteStatus(key int64, filetype global.FileType, status bool) {
	switch filetype {
	case global.DCM:
		if status {
			global.Logger.Info("***DCM文件删除成功，更新状态***")
			sql := `update file_remote fr set fr.dcm_file_exist = 0 where fr.instance_key = ?;`
			global.DBEngine.Exec(sql, key)
		} else {
			global.Logger.Info("***DCM文件删除失败，更新状态***")
			sql := `update file_remote fr set fr.dcm_file_exist = 2 where fr.instance_key = ?;`
			global.DBEngine.Exec(sql, key)
		}
	case global.JPG:
		if status {
			global.Logger.Info("***JPG文件删除成功，更新状态***")
			sql := `update file_remote fr set fr.img_file_exist = 0 where fr.instance_key = ?;`
			global.DBEngine.Exec(sql, key)
		} else {
			global.Logger.Info("***JPG文件删除失败，更新状态***")
			sql := `update file_remote fr set fr.img_file_exist = 2 where fr.instance_key = ?;`
			global.DBEngine.Exec(sql, key)
		}
	}
}
