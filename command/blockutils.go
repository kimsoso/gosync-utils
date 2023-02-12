package command

import (
	"btsync-utils/libs/utils"
	"encoding/binary"
	"math"
)

const (
	startString = "btsync start"
	endString   = "btsync end"
)

func (c *CommandBlk) blockcount(size int64) uint16 {
	return uint16(math.Ceil(float64(size) / float64(c.blocksize)))
}

func startBytes() []byte {
	return append(binary.BigEndian.AppendUint16(make([]byte, 0, 2), 12), []byte(startString)...)
}

func endBytes() []byte {
	return append(binary.BigEndian.AppendUint16(make([]byte, 0, 2), 10), []byte(endString)...)
}

// 判断是否是块开始
func IsStart(in []byte) bool {
	return binary.BigEndian.Uint16(in) == 12 && utils.B2S(in[2:14]) == startString
}

// 判断是否是块结束
func IsEnd(in []byte) bool {
	return binary.BigEndian.Uint16(in) == 10 && utils.B2S(in[2:12]) == endString
}
