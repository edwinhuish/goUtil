// +build windows

package thuOS

import (
	fileFunc "github.com/Mengdch/goUtil/FileTools"
	"github.com/edwinhuish/win"
	"golang.org/x/sys/windows"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
	"unsafe"
)

const win11Ver = 22000
func GetProcessId(name string) (uint32, error) {
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return 0, err
	}
	defer windows.CloseHandle(snapshot)
	var procEntry windows.ProcessEntry32
	procEntry.Size = uint32(unsafe.Sizeof(procEntry))
	if err = windows.Process32First(snapshot, &procEntry); err != nil {
		return 0, err
	}
	for {
		if windows.UTF16ToString(procEntry.ExeFile[:]) == name {
			return procEntry.ProcessID, nil
		}
		err = windows.Process32Next(snapshot, &procEntry)
		if err != nil {
			return 0, err
		}
	}
}

func UserHomeDir() string {
	home := winHome()
	if len(home) == 0 {
		return ""
	}
	return filepath.Join(home, "AppData", "Local")
}

func winHome() string {
	env := "USERPROFILE"
	var home string
	if home = os.Getenv(env); home == "" {
		return ""
	}
	return home
}
func GetFocus() uint64 {
	hForeWnd := win.GetForegroundWindow()
	dwForeID := win.GetWindowThreadProcessId(hForeWnd, nil)
	dwCurID := win.GetCurrentThreadId()
	win.AttachThreadInput(int32(dwForeID), int32(dwCurID), true)
	wnd := win.GetFocus()
	win.AttachThreadInput(int32(dwForeID), int32(dwCurID), false)
	return uint64(wnd)
}
func SetTop(hWnd uint64) bool {
	w := win.HWND(hWnd)
	if !win.IsWindow(w) {
		return false
	}
	win.ShowWindow(w, win.SW_SHOWNOACTIVATE)
	i := 4
	hForeWnd := win.GetForegroundWindow()
	if uint64(hForeWnd) == hWnd {
		return true
	}
	oldFw := hForeWnd
	var dwForeID, dwCurID uint32
	for i > 0 {
		dwForeID = win.GetWindowThreadProcessId(hForeWnd, nil)
		dwCurID = win.GetWindowThreadProcessId(w, nil)
		if dwCurID > 0 && dwForeID > 0 && dwForeID != dwCurID {
			win.AttachThreadInput(int32(dwForeID), int32(dwCurID), true)
		}
		win.SetForegroundWindow(w)
		win.SwitchToThisWindow(w, false)
		win.SetWindowPos(w, win.HWND_TOP, 0, 0, 0, 0, win.SWP_NOSIZE|win.SWP_NOMOVE)
		win.SetActiveWindow(w)
		win.SetFocus(w)
		if dwCurID > 0 && dwForeID > 0 && dwForeID != dwCurID {
			win.AttachThreadInput(int32(dwForeID), int32(dwCurID), false)
		}
		hForeWnd = win.GetForegroundWindow()
		if uint64(hForeWnd) == hWnd {
			win.SetActiveWindow(w)
			win.SetFocus(w)
			break
		}
		fmt.Println(oldFw, hForeWnd, hWnd)
		time.Sleep(10 * time.Millisecond)
		i--
	}
	return true
}
func GetOSVer() uint32 {
	ver := windows.RtlGetVersion()
	if ver != nil {
		return ver.MajorVersion
	}
	return 0
}
func GetOSName() (string, error) {
	ver := windows.RtlGetVersion()
	if ver != nil {
		switch ver.MajorVersion {
		case 0, 1, 2, 3, 4:
			return "Windows NT", nil
		case 5:
			switch ver.MinorVersion {
			case 0:
				return "Windows 2000", nil
			case 1:
				return "Windows XP", nil
			case 2:
				return "Windows Server 2003", nil
			}
		case 6:
			switch ver.MinorVersion {
			case 0:
				if ver.ProductType != 1 {
					return "Windows Server 2008", nil
				} else {
					return "Windows Vista", nil
				}
			case 1:
				if ver.ProductType != 1 {
					return "Windows Server 2008 R2", nil
				} else {
					return "Windows 7", nil
				}
			case 2:
				if ver.ProductType != 1 {
					return "Windows Server 2012", nil
				} else {
					return "Windows 8", nil
				}
			case 3:
				if ver.ProductType != 1 {
					return "Windows Server 2012 R2", nil
				} else {
					return "Windows 8.1", nil
				}
			}
		case 10:
			if ver.BuildNumber >= win11Ver {
				return "Windows 11", nil
			}
			if ver.ProductType != 1 {
				return "Windows Server 2016", nil
			} else {
				return "Windows 10", nil
			}
		}
		return "windows" + strconv.FormatInt(int64(ver.MajorVersion), 10) + "." + strconv.FormatInt(int64(ver.MinorVersion), 10), nil
	}
	return "", nil
}

const (
	PROCESS_CREATE_PROCESS            = 0x0080
	PROCESS_CREATE_THREAD             = 0x0002
	PROCESS_DUP_HANDLE                = 0x0040
	PROCESS_QUERY_INFORMATION         = 0x0400
	PROCESS_QUERY_LIMITED_INFORMATION = 0x1000
	PROCESS_SET_INFORMATION           = 0x0200
	PROCESS_SET_QUOTA                 = 0x0100
	PROCESS_SUSPEND_RESUME            = 0x0800
	PROCESS_TERMINATE                 = 0x0001
	PROCESS_VM_OPERATION              = 0x0008
	PROCESS_VM_READ                   = 0x0010
	PROCESS_VM_WRITE                  = 0x0020

	PROCESS_ALL_ACCESS = (PROCESS_CREATE_PROCESS | PROCESS_CREATE_THREAD | PROCESS_DUP_HANDLE | PROCESS_QUERY_INFORMATION | PROCESS_QUERY_LIMITED_INFORMATION | PROCESS_SET_INFORMATION | PROCESS_SET_QUOTA | PROCESS_SUSPEND_RESUME | PROCESS_TERMINATE | PROCESS_VM_OPERATION | PROCESS_VM_WRITE | PROCESS_VM_READ)

	GENERIC_WRITE         = 0x40000000
	FILE_SHARE_WRITE      = 0x00000002
	CREATE_ALWAYS         = 0x2
	FILE_ATTRIBUTE_NORMAL = 0x80

	MiniDumpNormal                         = 0x00000000
	MiniDumpWithDataSegs                   = 0x00000001
	MiniDumpWithFullMemory                 = 0x00000002
	MiniDumpWithHandleData                 = 0x00000004
	MiniDumpFilterMemory                   = 0x00000008
	MiniDumpScanMemory                     = 0x00000010
	MiniDumpWithUnloadedModules            = 0x00000020
	MiniDumpWithIndirectlyReferencedMemory = 0x00000040
	MiniDumpFilterModulePaths              = 0x00000080
	MiniDumpWithProcessThreadData          = 0x00000100
	MiniDumpWithPrivateReadWriteMemory     = 0x00000200
	MiniDumpWithoutOptionalData            = 0x00000400
	MiniDumpWithFullMemoryInfo             = 0x00000800
	MiniDumpWithThreadInfo                 = 0x00001000
	MiniDumpWithCodeSegs                   = 0x00002000
	MiniDumpWithoutAuxiliaryState          = 0x00004000
	MiniDumpWithFullAuxiliaryState         = 0x00008000
	MiniDumpWithPrivateWriteCopyMemory     = 0x00010000
	MiniDumpIgnoreInaccessibleMemory       = 0x00020000
	MiniDumpWithTokenInformation           = 0x00040000
	MiniDumpWithModuleHeaders              = 0x00080000
	MiniDumpFilterTriage                   = 0x00100000
	MiniDumpWithAvxXStateContext           = 0x00200000
	MiniDumpWithIptTrace                   = 0x00400000
	MiniDumpValidTypeFlags                 = 0x007fffff
)

var pid uint32
var miniDumpWriteDump *windows.LazyProc

func init() {
	pid = windows.GetCurrentProcessId()
	dbghelp := windows.NewLazySystemDLL("Dbghelp.dll")
	miniDumpWriteDump = dbghelp.NewProc("MiniDumpWriteDump")
	win.SetUnhandledExceptionFilter(onException)
}

type exceptionInfo struct {
	ThreadId          uint32
	ExceptionPointers uintptr
	ClientPointers    uint32
}

func logRecord(value, error string) {
	if log != nil {
		log(value, error)
	}
}
func onException(param uintptr) uintptr {
	go logRecord("", "onException")
	if miniDumpWriteDump == nil {
		logRecord("onException", "miniDumpWriteDump nil")
		return 0
	}
	pHandle := windows.CurrentProcess()
	homePath := UserHomeDir()
	dir := filepath.Join(homePath, "one2much")
	err := fileFunc.MakeDir(dir)
	if err != nil {
		logRecord("fileFunc.MakeDir"+dir, err.Error())
		return 0
	}
	name := filepath.Join(dir, "one2much"+strconv.FormatInt(time.Now().Unix(), 32)+".dmp")
	fromString, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		logRecord("onException.UTF16PtrFromString:"+name, err.Error())
		return 0
	}
	var sa windows.SecurityAttributes
	fHandle, err := windows.CreateFile(fromString, GENERIC_WRITE, FILE_SHARE_WRITE, &sa, CREATE_ALWAYS, FILE_ATTRIBUTE_NORMAL, 0)
	if err != nil {
		logRecord("onException.CreateFile:"+name, err.Error())
		return 0
	}
	var info exceptionInfo
	info.ExceptionPointers = param
	info.ThreadId = windows.GetCurrentThreadId()
	success, _, err := miniDumpWriteDump.Call(uintptr(pHandle), uintptr(pid), uintptr(fHandle), MiniDumpNormal, uintptr(unsafe.Pointer(&info)), 0, 0)
	windows.CloseHandle(fHandle)
	if success != 1 {
		logRecord("onException.miniDumpWriteDump", err.Error())
		return 0
	}
	if up != nil {
		up(name)
	}
	return 0
}
func IsWin7() bool {
	ver := windows.RtlGetVersion()
	if ver != nil && ver.MajorVersion < 7 && ver.MinorVersion < 2 { // win8
		return true
	}
	return false
}
func IsWin11() bool {
	ver := windows.RtlGetVersion()
	if ver != nil && ver.MajorVersion == 10 && ver.BuildNumber >= win11Ver { // win8
		return true
	}
	return false
}
