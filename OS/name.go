package thuOS

import (
	"os/user"
	"runtime"
)

func GetName() (deviceName, deviceId string, err error) {
	// 获取设备名
	u, err := user.Current()
	if err != nil {
		return
	}

	deviceName = runtime.GOOS + "(" + runtime.GOARCH + ")" + u.Username
	deviceId = u.Gid
	return
}
