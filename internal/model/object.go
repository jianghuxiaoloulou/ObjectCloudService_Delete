package model

import (
	"WowjoyProject/ObjectCloudService_Delete/global"
	"WowjoyProject/ObjectCloudService_Delete/pkg/general"
)

// 自动获取需要删除的数据
func GetDeleteData() {
	if global.RunStatus {
		global.Logger.Info("上次获取的数据没有消耗完，等待消耗完，再获取数据....")
		return
	}
	global.RunStatus = true

	global.Logger.Info("***开始获取需要删除的数据***")
	sql := `select ins.instance_key,ins.file_name,im.img_file_name,sl.ip,sl.s_virtual_dir
	from instance ins
	left join image im on im.instance_key = ins.instance_key
	left join study_location sl on sl.n_station_code = ins.location_code
	left join file_remote fr on ins.instance_key = fr.instance_key
	left join study s on ins.study_key = s.study_key 
	where fr.dcm_file_exist = 1 and (fr.dcm_file_exist_obs_local = 1 or fr.dcm_file_exist_obs_cloud = 1)
	and fr.dcm_update_time_retrieve >= ? and fr.dcm_update_time_retrieve <= ?`
	switch global.GeneralSetting.DelImgFlag {
	case "001":
		sql += ` and s.modality != "US"`
		sql += ` and s.modality != "ES"`
	case "010":
		sql += ` and s.modality = "US"`
	case "100":
		sql += ` and s.modality = "ES"`
	case "011":
		sql += ` and s.modality != "ES"`
	case "101":
		sql += ` and s.modality != "US"`
	case "110":
		sql += ` and s.modality in ("US","ES")`
	}
	global.Logger.Debug("查询数据的时间段是：", global.GeneralSetting.QueryStartTime, " :to: ", global.GeneralSetting.QueryEndTime)
	global.Logger.Info(sql)
	rows, err := global.ReadDBEngine.Query(sql, global.GeneralSetting.QueryStartTime, global.GeneralSetting.QueryEndTime)
	if err != nil {
		global.Logger.Fatal(err)
		global.RunStatus = false
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
	global.GeneralSetting.QueryStartTime = GetNextData(global.GeneralSetting.QueryStartTime)
	global.GeneralSetting.QueryEndTime = GetNextData(global.GeneralSetting.QueryEndTime)
	global.RunStatus = false
}

// 跟新删除文件标志
func UpdateDeleteStatus(key int64, filetype global.FileType, status bool) {
	switch filetype {
	case global.DCM:
		if status {
			global.Logger.Info("***DCM文件删除成功，更新状态***")
			sql := `update file_remote fr set fr.dcm_file_exist = 0 where fr.instance_key = ?;`
			global.WriteDBEngine.Exec(sql, key)
		} else {
			global.Logger.Info("***DCM文件删除失败，更新状态***")
			sql := `update file_remote fr set fr.dcm_file_exist = 2 where fr.instance_key = ?;`
			global.WriteDBEngine.Exec(sql, key)
		}
	case global.JPG:
		if status {
			global.Logger.Info("***JPG文件删除成功，更新状态***")
			sql := `update file_remote fr set fr.img_file_exist = 0 where fr.instance_key = ?;`
			global.WriteDBEngine.Exec(sql, key)
		} else {
			global.Logger.Info("***JPG文件删除失败，更新状态***")
			sql := `update file_remote fr set fr.img_file_exist = 2 where fr.instance_key = ?;`
			global.WriteDBEngine.Exec(sql, key)
		}
	}
}

// 获取下一天的时间段
func GetNextData(time string) string {
	sql := `select DATE_ADD(?,INTERVAL 1 DAY);`
	row := global.WriteDBEngine.QueryRow(sql, time)
	var temptime string
	err := row.Scan(&temptime)
	if err != nil {
		global.Logger.Error(err)
		return time
	}
	return temptime
}
