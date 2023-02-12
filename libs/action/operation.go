// 处理目录变动相关的函数
package action

import "log"

type ACT_STAT uint8

const (
	ACT_ADD ACT_STAT = 0 + iota
	ACT_UPDATE
	ACT_DELETE
)

// 新建一个操作
func (act *Action) NewOperation(path string, isDir bool, action uint8, md5 string, size int64) *Operation {
	row := &Operation{
		Path:  path,
		Act:   action,
		IsDir: isDir,
		Md5:   md5,
		Size:  size,
	}

	if err := act.db.Create(row).Error; err == nil {
		return row
	} else {
		log.Println("save db operation error:", err)
		return nil
	}
}

func (act *Action) SaveOperations(in []*Operation) (out []*Operation) {
	if err := act.db.Create(in).Error; err != nil {
		return nil
	} else {
		return in
	}
}

func (act *Action) LastId() uint {
	lastrow := &Operation{}
	if err := act.db.Last(lastrow).Error; err != nil {
		return 0
	} else {
		return lastrow.ID
	}
}

// 获取操作列表
func (act *Action) LastOperations(fromId uint, order string) []Operation {
	rows := []Operation{}
	exec := act.db
	exec = exec.Order("id " + order)
	exec.Find(&rows, "id >= ?", fromId)

	return rows
}

// 获取一个操作
func (act *Action) GetOperation(id uint) *Operation {
	row := &Operation{}
	if err := act.db.First(row, id).Error; err == nil {
		return row
	} else {
		return nil
	}
}

// 获取列表以提供id列表的方式
func (a *Action) GetOperationsById(ids []uint) (out []*Operation) {
	out = []*Operation{}
	a.db.Find(&out, ids)
	return out
}
