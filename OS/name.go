package thuOS

import (
	"os"
	"os/user"
	"runtime"
)

var log func(string, string)
var up func(string)

func SetLog(lf func(string, string), u func(string)) {
	log = lf
	up = u
}

func GetName() (deviceName, deviceId string, err error) {
	// 获取设备名
	u, err := user.Current()
	if err != nil {
		return
	}
	name, _ := os.Hostname()
	deviceName = runtime.GOOS + "(" + runtime.GOARCH + ")" + name + "\\" + u.Username
	deviceId = u.Gid
	return
}
