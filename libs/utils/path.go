// 路径处理，例如执行文件位置，相对地址
package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// 执行文件所在目录
func ExecDir() (execdir string) {
	execdir, _ = os.Executable()
	// air 环境下返回工作目录
	if strings.HasSuffix(execdir, "/tmp/main") {
		execdir = execdir[:len(execdir)-9]
	} else if strings.HasPrefix(execdir, "/tmp/") {
		execdir, _ = os.Getwd()
	} else {
		execdir = filepath.Dir(execdir)
	}
	return
}
