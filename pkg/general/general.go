package general

import (
	"os"
	"path/filepath"
	"syscall"
	"unsafe"
)

// 基础函数
func GetFilePath(file, ip, virpath string) (path string) {
	path += "\\\\"
	path += ip
	path += "\\"
	path += virpath
	path += "\\"
	path += file
	return
}

func GetFileSize(filename string) int64 {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return 0
	}
	return fileInfo.Size()
}

// 检查文件路径
func CheckPath(path string) {
	dir, _ := filepath.Split(path)
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(dir, os.ModePerm)
		}
	}
}

func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func DeleteFile(dest string) int {
	if !Exist(dest) {
		// 文件不存在
		return 0
	} else {
		err := os.Remove(dest)
		if err != nil {
			// 文件删除失败
			return 2
		}
		// 删除成功
		return 1
	}
}

// 返回磁盘大小 GB
func GetDiskSize(disk string) (size int64) {
	kernel := syscall.MustLoadDLL("kernel32.dll")
	proc := kernel.MustFindProc("GetDiskFreeSpaceExW")
	lpFreeBytesAvailable := int64(0)
	proc.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(disk))),
		uintptr(unsafe.Pointer(&lpFreeBytesAvailable)))
	size = lpFreeBytesAvailable / 1024 / 1024 / 1024.0 //GB
	return size
}
