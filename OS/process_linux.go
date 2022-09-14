// +build linux

package thuOS

import (
	"os"
)

func UserHomeDir() string {
	return os.Getenv("HOME")
}
func GetSystemUserDir() []FolderShow {
	names := []FolderShow{{Path: "Desktop", Name: "桌面"}, {Path: "Documents", Name: "文档"}, {Path: "Downloads", Name: "下载"}}
	return names
}
func GetLogicalDrives() ([]FolderShow, error) {
	var drives []FolderShow
	return drives, nil
}
func ExecFile(file string) error {
	return nil
}
func DeleteReboot(filePath string) (err error) {
	return nil
}
func GetFreeBytes(filePath string) (free uint64, suc bool) {
	return
}
func GetFocus() uint64                         { return 0 }
func SetTop(hWnd uint64)                       {}
func GetProcessId(name string) (uint32, error) { return 0, nil }
