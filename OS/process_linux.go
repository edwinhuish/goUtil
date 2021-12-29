// +build linux

package thuOS
func UserHomeDir() string {
	return os.Getenv("HOME")
}
