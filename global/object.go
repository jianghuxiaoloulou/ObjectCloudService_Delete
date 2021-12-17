package global

// 文件类型
type FileType int

const (
	DCM FileType = iota // DCM 文件
	JPG                 // JPG 文件
)

type ObjectData struct {
	InstanceKey int64    // instance_key 目标key
	FilePath    string   // 文件路径
	Type        FileType // 文件类型
	Count       int      // 文件执行次数
}

var (
	ObjectDataChan chan ObjectData
)
