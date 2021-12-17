package main

import (
	"WowjoyProject/ObjectCloudService_Delete/global"
	"WowjoyProject/ObjectCloudService_Delete/internal/model"
	"WowjoyProject/ObjectCloudService_Delete/pkg/general"
	"WowjoyProject/ObjectCloudService_Delete/pkg/object"
	"WowjoyProject/ObjectCloudService_Delete/pkg/workpattern"

	"github.com/robfig/cron"
)

// @title 本地存储文件删除服务
// @version 1.0.0.1
// @description 文件删除服务
// @termsOfService https://github.com/jianghuxiaoloulou/OBbjectCloudService.git
func main() {
	global.Logger.Info("*******开始运行文件删除服务********")

	global.ObjectDataChan = make(chan global.ObjectData)

	// 注册工作池，传入任务
	// 参数1 初始化worker(工人)设置最大线程数
	wokerPool := workpattern.NewWorkerPool(global.GeneralSetting.MaxThreads)
	// 有任务就去做，没有就阻塞，任务做不过来也阻塞
	wokerPool.Run()
	// 处理任务
	go func() {
		for {
			select {
			case data := <-global.ObjectDataChan:
				sc := &Dosomething{key: data}
				wokerPool.JobQueue <- sc
			}
		}
	}()

	run()
}

type Dosomething struct {
	key global.ObjectData
}

func (d *Dosomething) Do() {
	global.Logger.Info("正在处理的数据是：", d.key)
	//处理封装对象操作
	obj := object.NewObject(d.key)
	obj.DeleteFile()
}

func run() {
	// 方式一：
	// for {
	// 	time.Sleep(time.Second * 10)
	// 	model.AutoUploadObjectData()
	// }

	// 方式二：获取任务(定时任务)
	MyCron := cron.New()
	MyCron.AddFunc(global.GeneralSetting.CronSpec, func() {
		global.Logger.Info("开始执行定时任务")
		work()
	})

	MyCron.Start()

	defer MyCron.Stop()

	select {}
}

func work() {
	global.Logger.Info("***获取磁盘空间***")
	size := int(general.GetDiskSize(global.GeneralSetting.TargetPath))
	global.Logger.Info("***目标磁盘剩余空间：", size, "GB")
	if size > global.GeneralSetting.TargetValue {
		global.Logger.Info("***目标磁盘剩余空间未达到设定值，程序暂时不执行删除任务***")
		return
	} else {
		global.Logger.Info("***目标磁盘剩余空间小于设定值，执行删除任务***")
		model.GetDeleteData()
	}
}
