/*
 * @Author: soso
 * @Date: 2022-01-27 17:52:43
 * @LastEditTime: 2022-02-16 15:30:36
 * @LastEditors: Please set LastEditors
 * @Description: 发包，收包
 * @FilePath: /sync-client/go-utils/utils/sonet/net.go
 */
package sonet

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

// 发送数据包
func SendPack(conn net.Conn, msgId uint32, data []byte) (err error) {
	dp := NewDataPack()
	msg := &Message{
		Id:      msgId,
		DataLen: uint32(len(data)),
		Data:    data,
	}
	sendData, err := dp.Pack(msg)
	if err != nil {
		return
	}
	retry := 0
RETRY:
	if err = conn.SetWriteDeadline(time.Now().Add(timeoutDruation)); err != nil {
		return err
	}
	_, err = conn.Write(sendData)
	if err != nil && os.IsTimeout(err) {
		retry++
		if retry <= RetryTimes {
			goto RETRY
		}
	}

	return
}

// 接收数据包
func ReceivePack(conn net.Conn) (data *Message, err error) {
	dp := NewDataPack()
	headData := make([]byte, dp.GetHeadLen())

	if err = conn.SetReadDeadline(time.Now().Add(timeoutDruation)); err != nil {
		return nil, err
	}
	_, err = io.ReadFull(conn, headData)
	if err != nil {
		return
	}
	msgHead, err := dp.Unpack(headData)
	if err != nil {
		return nil, err
	}
	if msgHead.GetDataLen() > 0 {
		msg := msgHead.(*Message)
		msg.Data = make([]byte, msg.GetDataLen())

		retry := 0
	RETRY:
		if err = conn.SetReadDeadline(time.Now().Add(timeoutDruation)); err != nil {
			return nil, err
		}
		if _, err = io.ReadFull(conn, msg.Data); err != nil && os.IsTimeout(err) {
			retry++
			if retry <= RetryTimes {
				goto RETRY
			}
		}

		if err != nil {
			return nil, err
		} else {
			return msg, nil
		}
	}

	return msgHead.(*Message), nil
}

// 带超时的连接
func TimeoutConn(ip string, port int) (conn net.Conn, err error) {
	dialer := net.Dialer{Timeout: timeoutDruation}
	return dialer.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
}
