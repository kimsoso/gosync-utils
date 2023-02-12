/*
 * @Author: soso
 * @Date: 2022-01-26 10:56:33
 * @LastEditTime: 2022-03-03 17:51:06
 * @LastEditors: Please set LastEditors
 * @Description: sonet传输包
 * @FilePath: /go-mesh-sync/go-utils/utils/sonet/consts.go
 */
package sonet

import (
	"btsync-utils/vars"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"os"
	"time"

	"github.com/vmihailenco/msgpack"
)

const (
	timeoutDruation = time.Second * 10
	retryTime       = 3
)

func WriteCommandWithValue(conn net.Conn, command uint8, value interface{}) error {
	return Write(conn, BinaryCommandWithValue(command, value))
}

func WriteCommandWithPayload(conn net.Conn, command uint8, payload interface{}) error {
	return Write(conn, BinaryCommandWithPayload(command, payload))
}

func ReadCommandWithPayload(conn net.Conn) (cmd uint8, payload []byte, err error) {
	if _, data, err := Read(conn); err != nil {
		return 0, nil, err
	} else {
		reader := bytes.NewReader(data)
		binary.Read(reader, binary.BigEndian, &cmd)
		payload = data[1:]
		return cmd, payload, nil
	}
}

func ReadCommandWithValue(conn net.Conn, value interface{}) (cmd uint8, err error) {
	if _, data, err := Read(conn); err != nil {
		return 0, err
	} else {
		reader := bytes.NewReader(data)
		binary.Read(reader, binary.BigEndian, &cmd)
		if value != nil {
			binary.Read(reader, binary.BigEndian, value)
		}
		return cmd, nil
	}
}

// 二进制命令,uint8
func BinaryCommandWithPayload(command uint8, payload interface{}) []byte {
	buf := bytes.NewBuffer([]byte{})

	binary.Write(buf, binary.BigEndian, command)
	if payload != nil {
		data, _ := msgpack.Marshal(payload)
		buf.Write(data)
	}

	return buf.Bytes()
}

func BinaryCommandWithValue(command uint8, value interface{}) []byte {
	buf := bytes.NewBuffer([]byte{})

	binary.Write(buf, binary.BigEndian, command)
	if value != nil {
		binary.Write(buf, binary.BigEndian, value)
	}

	return buf.Bytes()
}

/**
 * @description: 读取命令数据
 * @param {net.Conn} conn
 * @return {*}
 */
func Read(conn net.Conn) (uint32, []byte, error) {
	msg, err := ReceivePack(conn)
	if err != nil {
		return 0, nil, err
	}
	return msg.GetMsgId(), msg.GetData(), nil
}

/**
 * @description: 写入自拼写命令
 * @param {net.Conn} conn
 * @param {[]byte} command
 * @return {*}
 */
func Write(conn net.Conn, command []byte) error {
	return SendPack(conn, 0, command)
}

// 发送文件
func SendFileData(conn net.Conn, filePath string, offset, expectSize int64, cbFunc func(sentN int)) (int, error) {
	fd, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer fd.Close()

	_, err = fd.Seek(offset, 0)
	if err != nil {
		return 0, err
	}

	buf := make([]byte, vars.BufSize)
	var packNum uint32 = 0
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
		sent += n
		packNum += 1
		if cbFunc != nil {
			cbFunc(n)
		}
		if expectSize > 0 && sent >= int(expectSize) {
			break
		}
	}
	return sent, nil
}

/**
 * @description: 接收长数据
 * @param {net.Conn} conn
 * @param {func} callbackFunc
 * @return {error}
 */
func ReceiveLongData(conn net.Conn, callbackFunc func([]byte) (err error, next bool)) error {
	var packNumber uint32 = 0
	for {
		_n, buf, err := Read(conn)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if _n != packNumber {
			return errors.New("pack is wrong")
		}

		if err, next := callbackFunc(buf); err != nil || !next {
			return err
		}

		packNumber += 1
	}
}
