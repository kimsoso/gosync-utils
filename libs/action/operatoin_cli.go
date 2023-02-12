// 处理客户端中目录变动相关的函数
package action

import (
	"btsync-utils/libs/utils"
	"log"
	"sort"
)

type byUint []uint

func (a byUint) Len() int           { return len(a) }
func (a byUint) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byUint) Less(i, j int) bool { return a[i] < a[j] }

func PrintStruct(opt Operation) {
	log.Println("ID:", opt.ID)
	log.Println("Path:", opt.Path)
	log.Println("Act:", opt.Act)
	log.Println("IsDir:", opt.IsDir)
	log.Println("Md5:", opt.Md5)
	log.Println("Size:", opt.Size)
}

func (a *Action) SaveOperationCliFromOperation(operation Operation) error {
	if err := a.db.Create(&OperationCli{
		ID:    operation.ID,
		Path:  operation.Path,
		Act:   operation.Act,
		IsDir: operation.IsDir,
		Md5:   operation.Md5,
		Size:  operation.Size,
	}).Error; err != nil {
		return err
	} else {
		a.saveVersion()
		return nil
	}
}
func (a *Action) SaveOperationCli(optcli OperationCli) error {
	if err := a.db.Create(&optcli).Error; err == nil {
		a.saveVersion()
		return nil
	} else {
		return err
	}
}

func (a *Action) LastOperationClis(fromId uint) (out []*OperationCli) {
	out = []*OperationCli{}
	a.db.Find(&out, "id >= ?", fromId)
	return
}

func (a *Action) GetOperationCli(id uint) (out *OperationCli) {
	out = &OperationCli{}
	if a.db.First(out, id).Error == nil {
		return out
	} else {
		return nil
	}
}

func (a *Action) GetOperationClisByIds(ids []uint) (out []*OperationCli) {
	out = []*OperationCli{}
	a.db.Find(&out, "id IN ?", ids)
	return out
}

func (a *Action) LastOperationCliId() uint {
	rowcli := &OperationCli{}
	lastCliId := uint(0)
	if err := a.db.Order("id DESC").First(rowcli).Error; err == nil {
		lastCliId = rowcli.ID
	}
	return lastCliId
}

func (a *Action) LastOperationCliRunningId() uint {
	rowcli := &OperationCli{}
	lastCliId := uint(0)
	if err := a.db.Order("id DESC").First(rowcli).Error; err == nil {
		lastCliId = rowcli.ID
	}

	rowrun := &OperationRunning{}
	lastRunningId := uint(0)

	if err := a.db.Order("id DESC").First(rowrun).Error; err == nil {
		lastRunningId = rowrun.ID
	}
	return utils.If(lastCliId > lastRunningId, lastCliId, lastRunningId)
}

func (a *Action) SaveOperationClisFromRunnings(ids []uint) (saved []uint, err error) {
	a.Lock()
	defer a.Unlock()

	runnings := a.GetRunnings(ids)

	clis := []OperationCli{}
	saved = make([]uint, 0, len(ids))
	for _, running := range runnings {
		saved = append(saved, running.ID)
		clis = append(clis, OperationCli(running))
	}

	sort.Sort(byUint(saved))

	if err = a.db.Create(&clis).Error; err == nil {
		err = a.RemoveRunnings(saved)
	}

	return
}
