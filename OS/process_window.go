package thuOS

import (
	"golang.org/x/sys/windows"
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

