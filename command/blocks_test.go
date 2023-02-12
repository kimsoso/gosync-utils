// 将多个小的命令组合到一起，形成命令数据块，以避免小文件的传输缓慢的问题!
// 最后数据块不一定固定为输入的blockSize值有可能大，有可能小，由用户使用时配置

package command

import (
	"btsync-utils/libs/action"
	"btsync-utils/libs/config"
	"btsync-utils/libs/dirfile"
	"log"
	"path/filepath"
	"testing"

	c "github.com/smartystreets/goconvey/convey"
)

var (
	conf *config.Conf
	clis []action.OperationCli
	cb   *CommandBlk

	basepath = "/home/soso/go/src/btsync-utils/"
	testdir  = filepath.Join(basepath, "testdata")
	conffile = filepath.Join(basepath, "conf/server.yaml")
)

func init() {
	clis = []action.OperationCli{}
	if files, err := dirfile.GetDirStruct(testdir, testdir); err == nil {
		for i, file := range files {
			clis = append(clis, action.OperationCli{
				ID:    uint(i) + 1,
				Path:  file.Path,
				Act:   uint8(action.ACT_ADD),
				IsDir: file.IsDir,
				Md5:   file.Md5,
				Size:  file.Size,
			})
		}
	} else {
		log.Println("got error:", err)
	}

	conf = config.NewConf(conffile)
}

func Test_tt(t *testing.T) {
	c.Convey("测试的测试", t, func() {
		c.So(100, c.ShouldEqual, 100)
	})
}
