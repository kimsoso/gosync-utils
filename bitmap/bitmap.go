package bitmap

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"sync"
)

// 位图
type BitMap struct {
	sync.RWMutex
	bits []byte
	max  uint32
}

// 初始化一个BitMap
// 一个byte有8位,可代表8个数字,取余后加1为存放最大数所需的容量
func NewBitMap(max uint32) *BitMap {
	bits := make([]byte, (max>>3)+1)
	return &BitMap{bits: bits, max: max}
}

// 初始化一个空的BitMap，以备反序列化填充
func NewEmptyBitMap() *BitMap {
	return &BitMap{}
}

// 添加一个数字到位图
// 计算添加数字在数组中的索引index,一个索引可以存放8个数字
// 计算存放到索引下的第几个位置,一共0-7个位置
// 原索引下的内容与1左移到指定位置后做或运算
func (b *BitMap) Add(num uint) *BitMap {
	b.Lock()
	defer b.Unlock()

	if num >= uint(b.Max()) {
		panic("input exceed range of bitmap max:" + strconv.Itoa(int(b.Max()-1)))
	}
	index := num >> 3
	pos := num & 0x07
	b.bits[index] |= 1 << pos
	return b
}

// 判断一个数字是否在位图
// 找到数字所在的位置,然后做与运算
func (b *BitMap) IsExist(num uint) bool {
	b.RLock()
	defer b.RUnlock()

	index := num >> 3
	pos := num & 0x07
	return b.bits[index]&(1<<pos) != 0
}

// 删除一个数字在位图
// 找到数字所在的位置取反,然后与索引下的数字做与运算
func (b *BitMap) Remove(num uint) *BitMap {
	b.Lock()
	defer b.Unlock()

	if num >= uint(b.Max()) {
		panic("input exceed range of bitmap max:" + strconv.Itoa(int(b.Max()-1)))
	}
	index := num >> 3
	pos := num & 0x07
	b.bits[index] = b.bits[index] & ^(1 << pos)
	return b
}

// 是不是全位空
func (b *BitMap) IsEmpty() bool {
	b.RLock()
	defer b.RUnlock()

	for i := 0; i < len(b.bits); i++ {
		if b.bits[i] != 0 {
			return false
		}
	}
	return true
}

// 是不是全位设置
func (b *BitMap) IsFull() bool {
	b.RLock()
	defer b.RUnlock()

	len := b.max >> 3
	for i := 0; i < int(len); i++ {
		if b.bits[i] != 255 {
			return false
		}
	}
	var a byte = 255
	return b.bits[len] == ^(a << (b.max & 0x07))
}

// 全位设置 1
func (b *BitMap) AddFull() *BitMap {
	b.Lock()
	defer b.Unlock()

	len := b.max >> 3
	for i := 0; i < int(len); i++ {
		b.bits[i] = 255
	}
	var a byte = 255
	b.bits[len] = ^(a << (b.max & 0x07))
	return b
}

// 全位设置 0
func (b *BitMap) RemoveAll() *BitMap {
	b.Lock()
	defer b.Unlock()

	len := b.max>>3 + 1
	for i := 0; i < int(len); i++ {
		b.bits[i] = 0
	}
	return b
}

// bytes length
func (b *BitMap) Len() int {
	return len(b.bits)
}

// 位图的最大数字
func (b *BitMap) Max() uint32 {
	return b.max
}

func (b *BitMap) String() string {
	b.RLock()
	defer b.RUnlock()

	return fmt.Sprint(b.bits)
}

// 序列化
func (b *BitMap) Serialize() []byte {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, b.max)
	binary.Write(buf, binary.LittleEndian, b.bits)
	return buf.Bytes()
}

// 反序列化
func (b *BitMap) UnSerialize(in []byte) error {
	b.Lock()
	defer b.Unlock()

	buf := bytes.NewBuffer(in)

	var max uint32
	err := binary.Read(buf, binary.LittleEndian, &max)
	if err != nil {
		return err
	}

	b.max = max
	bits := make([]byte, (max>>3)+1)

	binary.Read(buf, binary.LittleEndian, &bits)
	if err != nil {
		return err
	}

	b.bits = bits
	return nil
}
