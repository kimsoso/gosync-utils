// 将多个小的命令组合到一起，形成命令数据块，以避免小文件的传输缓慢的问题!
// 最后数据块不一定固定为输入的blockSize值有可能大，有可能小，由用户使用时配置

package command

import (
	"btsync-utils/libs/action"
	"btsync-utils/libs/config"
	"btsync-utils/libs/pool"
	"btsync-utils/libs/sonet"
	"btsync-utils/libs/utils"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/fatih/color"
	"github.com/vmihailenco/msgpack"
)

type CommandBlk struct {
	conf *config.Conf

	blocksize int

	currentOptCli action.OperationCli
}

func NewCommandBlk(conf *config.Conf, blockSize int) *CommandBlk {
	return &CommandBlk{
		conf:      conf,
		blocksize: blockSize,
	}
}

// 组合命令数据
// main: [uint16(optcli length)][optcli content][uint8,0=不存在,1=存在文件数据块,2=文件已删除，3=大文件][uint16(block count)]
// blocks: [uint16(blockRealsize)][block content]

func (c *CommandBlk) SendCommands(conn net.Conn, optclis []action.OperationCli) {
	if len(optclis) == 0 {
		return
	}

	dataBlock := make([]byte, 0, c.blocksize*2)

	if c.startSend(conn) != nil {
		return
	}
	defer c.stopSend(conn)

	packNum := 0
	clids := make([]uint, 0, len(optclis))
	tmpClis := make([]action.OperationCli, 0, len(optclis))
	for i := 0; i < len(optclis); i++ {
		clids = append(clids, optclis[i].ID)
		tmpClis = append(tmpClis, optclis[i])

		dataBlock = append(dataBlock, c.PackOptcli(optclis[i])...)

		if len(dataBlock) > c.blocksize {
			c.logClis(clids, tmpClis)
			if err := c.sendBlockOfCommands(conn, dataBlock, packNum); err != nil {
				log.Println("send block error:", err)
				return
			}

			packNum++

			clids = clids[:0]
			tmpClis = tmpClis[:0]
			dataBlock = dataBlock[:0]
		}
	}

	c.logClis(clids, tmpClis)

	if len(dataBlock) > 0 {
		c.sendBlockOfCommands(conn, dataBlock, packNum)
	}
}

func (c *CommandBlk) logClis(clis []uint, optclis []action.OperationCli) {
	log.Println("this block include these ids:", color.RedString("%v", clis))
	for _, cli := range optclis {
		log.Println(cli.ID, utils.If(cli.IsDir, "D", "F"), cli.Path)
	}
}

func (c *CommandBlk) startSend(conn net.Conn) error {
	if err := sonet.WriteCommandWithValue(conn, CBLOCKS_START, nil); err != nil {
		log.Println("send command blocks start error:", err)
		return err
	} else {
		if _, cmd, err := sonet.Read(conn); err != nil {
			log.Println("send command blocks start feedback error:", err)
			return err
		} else if cmd[0] != OK {
			log.Println("send command blocks start recieve wrong reply:", cmd[0])
			return errors.New("rev wrong")
		}
	}
	return nil
}

func (c *CommandBlk) stopSend(conn net.Conn) {
	if err := sonet.Write(conn, endBytes()); err != nil {
		log.Println("send command blocks fininsh error:", err)
	}
}

func (c *CommandBlk) sendBlockOfCommands(conn net.Conn, cmds []byte, packNum int) error {
	if err := sonet.SendPack(conn, uint32(packNum), cmds); err != nil {
		return err

	} else {
	ReReadSignal:
		if remoteNum, rdata, err := sonet.Read(conn); err != nil {
			return err

		} else if remoteNum != uint32(packNum) {
			return fmt.Errorf("send command blocks mismatch packnum: %d, %d", packNum, remoteNum)

		} else if rdata[0] == CBLOCKS_TRANSFILE {
			offset := binary.BigEndian.Uint64(rdata[1:])
			realpath := c.conf.RealFilepath(utils.B2S(rdata[9:]))

			if _, err := sonet.SendFileDataRR(conn, realpath, int64(offset), 0, nil); err != nil {
				return err
			} else {
				goto ReReadSignal
			}

		} else if rdata[0] != OK {
			return fmt.Errorf("send command blocks recieve wrong reply:%d", rdata[0])
		}
	}
	return nil
}

func (c *CommandBlk) PackOptcli(optcli action.OperationCli) (out []byte) {
	c.currentOptCli = optcli

	// writer
	wt := bytes.NewBuffer(make([]byte, 0, c.blocksize*2))
	// optcli content
	cdata, _ := msgpack.Marshal(c.currentOptCli)
	// cli pack length
	binary.Write(wt, binary.BigEndian, uint16(len(cdata)))
	// pack
	wt.Write(cdata)

	// 块标志,[uint8,0=不存在,1=存在文件数据块,2=文件已删除，3=大文件][uint16(block count)]
	blockbit := uint8(0)
	// file data blocks count
	blockscount := utils.If((c.currentOptCli.Act == uint8(action.ACT_ADD) || c.currentOptCli.Act == uint8(action.ACT_UPDATE)) && !optcli.IsDir && optcli.Size > 0, c.blockcount(optcli.Size), 0)

	if blockscount > 0 {
		if _, err := os.Stat(c.conf.RealFilepath(optcli.Path)); err != nil {
			blockbit = 2
		} else if blockscount == 1 {
			blockbit = 1
		} else {
			blockbit = 3
		}
	}
	binary.Write(wt, binary.BigEndian, blockbit)

	if blockbit == 1 {
		filedata, _ := os.ReadFile(c.conf.RealFilepath(c.currentOptCli.Path))
		binary.Write(wt, binary.BigEndian, uint16(len(filedata)))
		wt.Write(filedata)
	}

	return wt.Bytes()
}

// 解包网络数据
// main: [uint16(optcli length)][optcli content][uint8,0=不存在,1=存在文件数据块,2=文件已删除，3=大文件][uint16(block count)]
// blocks: [uint16(blockRealsize)][block content]
func UnpackBlock(
	blockdata []byte,
	currentIndex int,
	processCli func(cli *action.OperationCli, blockbit uint8, filedata []byte) error) (nextIndex int, err error) {
	cliLen := binary.BigEndian.Uint16(blockdata[currentIndex:])
	nextIndex = currentIndex + 2 + int(cliLen)

	cli := pool.OperationCliPool.Get().(*action.OperationCli)
	defer pool.OperationCliPool.Put(cli)

	if err = msgpack.Unmarshal(blockdata[currentIndex+2:nextIndex], cli); err != nil {
		return
	}

	blockbit := blockdata[nextIndex]
	nextIndex++

	var filedata []byte
	if blockbit == 1 {
		blocklen := binary.BigEndian.Uint16(blockdata[nextIndex:])
		nextIndex += 2
		filedata = blockdata[nextIndex : nextIndex+int(blocklen)]
		nextIndex += int(blocklen)
	} else {
		filedata = make([]byte, 0)
	}

	if processCli != nil {
		err = processCli(cli, blockbit, filedata)
	}

	return nextIndex, err
}
