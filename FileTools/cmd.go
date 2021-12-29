package fileFunc

import (
	"bytes"
	"os/exec"
	"runtime"
)

func ExecCommand(param ... string) (errMsg, errStd, errOut, errCode string) {
	var cmd *exec.Cmd
	var name string
	var params []string
	switch runtime.GOOS {
	case "linux", "darwin":
		name = "bash"
		params = append(params, "-c")
	case "windows":
		name = "cmd.exe"
		params = append(params, "/c")
	}
	params = append(params, param...)
	cmd = exec.Command(name, params...)
	return execCmd(cmd)
}
func execCmd(cmd *exec.Cmd) (errMsg, errStd, errOut, errCode string) {
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		errCode = "CmdExecErr"
		errMsg = err.Error()
		errStd = stderr.String()
		if out.Len() > 0 {
			errOut = out.String()
		}
	}
	return
}
