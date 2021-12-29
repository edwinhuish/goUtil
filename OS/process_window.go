// +build windows

package thuOS

import (
	"golang.org/x/sys/windows"
	"os"
	"path/filepath"
	"unsafe"
)

func GetProcessEntry(name string) (*windows.ProcessEntry32, error) {
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, err
	}
	defer windows.CloseHandle(snapshot)
	var procEntry windows.ProcessEntry32
	procEntry.Size = uint32(unsafe.Sizeof(procEntry))
	if err = windows.Process32First(snapshot, &procEntry); err != nil {
		return nil, err
	}
	for {
		if windows.UTF16ToString(procEntry.ExeFile[:]) == name {
			return &procEntry, nil
		}
		err = windows.Process32Next(snapshot, &procEntry)
		if err != nil {
			return nil, err
		}
	}
}

func UserHomeDir() string {
	// TODO 其他系统适配
	home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	if home == "" {
		home = os.Getenv("USERPROFILE")
	}
	return filepath.Join(home, "AppData", "Local")
}
