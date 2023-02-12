package file

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

// 创建一个空文件
func TouchFile(filename string) error {
	dir := filepath.Dir(filename)
	if dir != "." {
		os.MkdirAll(dir, os.ModePerm)
	}
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return err
}

// write file with auto create dir
func WriteFile(filename string, data []byte) error {
	dir := filepath.Dir(filename)
	if _, err := os.Stat(dir); err != nil {
		os.MkdirAll(dir, os.ModePerm)
	}
	return os.WriteFile(filename, data, os.ModePerm)
}

// 创建一个0填充的空文件，如果存在则truncate
func ZeroFile(filename string, size int64) error {
	dir := filepath.Dir(filename)
	if dir != "." {
		os.MkdirAll(dir, os.ModePerm)
	}
	if _, err := os.Stat(filename); err != nil {
		if err := TouchFile(filename); err != nil {
			return err
		}
	}
	return os.Truncate(filename, size)
}

// 针对最后一块数据存在不能读满buff的情况，为正常
func ReadRange(filename string, offset, size int64) (data []byte, err error) {
	if fd, err := os.Open(filename); err != nil {
		return nil, err
	} else {
		defer fd.Close()

		if co, err := fd.Seek(offset, 0); err != nil {
			return nil, err
		} else if co != offset {
			return nil, fmt.Errorf("%d is not equeal with offset:%d privided", co, offset)
		} else {
			data = make([]byte, size)
			if rn, err := io.ReadFull(fd, data); err == nil || err == io.ErrUnexpectedEOF {
				return data[:rn], nil
			} else {
				return nil, err
			}
		}
	}
}

func WriteRange(filename string, offset int64, payload []byte) error {
	fd, err := os.OpenFile(filename, os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer fd.Close()

	if ret, err := fd.Seek(offset, 0); err != nil {
		return err
	} else if ret != offset {
		return errors.New("can't move to fixed offset, got:" + strconv.Itoa(int(ret)))
	} else {
		if nn, err := fd.Write(payload); err != nil {
			return err
		} else if nn != len(payload) {
			return errors.New("didn't wrote fixed length bytes:got " + strconv.Itoa(nn))
		} else {
			return nil
		}
	}
}
