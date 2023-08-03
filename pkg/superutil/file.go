package superutil

import (
	"os"
)

// FileExist 文件是否存在
func FileExist(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	if fi.IsDir() {
		return false
	}
	return true
}
