package sonet

import (
	"btsync-utils/vars"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/fatih/color"
)

// 以应答方式发送文件数据
func SendFileDataRR(conn net.Conn, filePath string, offset, expectSize int64, cbFunc func(sentN int)) (int, error) {
	fd, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}

	defer fd.Close()

	_, err = fd.Seek(offset, 0)
	if err != nil {
		return 0, err
	}

	var packNum uint32 = 0
	buf := make([]byte, vars.BufSize)

	sent := 0
	for {
		n, err := fd.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}
		// check if out of range
		if expectSize > 0 && int64(sent+n) > expectSize {
			n -= (sent + n) - int(expectSize)
		}
		err = SendPack(conn, packNum, buf[:n])
		if err != nil {
			return sent, err
		}

		if remoteNum, _, err := Read(conn); err != nil {
			return sent, err
		} else if remoteNum != packNum {
			return sent, fmt.Errorf("didnt recieve feedback %d", packNum)
		}

		sent += n
		packNum++

		if cbFunc != nil {
			cbFunc(n)
		}
		if expectSize > 0 && sent >= int(expectSize) {
			break
		}
	}
	return sent, nil
}

// 以回复应答的方式接收长数据
func ReceiveLongDataRR(conn net.Conn, callbackFunc func([]byte) (err error, next bool)) error {
	var packNumber uint32 = 0
	for {
		remoteNum, buf, err := Read(conn)
		if err == io.EOF {
			log.Println(color.RedString("long data done!"))
			return nil
		}
		if err != nil {
			return err
		}
		if remoteNum != packNumber {
			return fmt.Errorf(color.RedString("pack is wrong l%d r%d\n", packNumber, remoteNum))
		}

		if err := SendPack(conn, packNumber, []byte{}); err != nil {
			return err
		}

		if err, next := callbackFunc(buf); err != nil || !next {
			return err
		}

		packNumber++
	}
}
