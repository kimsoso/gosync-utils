package action

import (
	"btsync-utils/libs/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

func (a *Action) NewRunnings(rows []OperationRunning) error {
	return a.db.Create(&rows).Error
}

func (a *Action) NewRunning(row OperationRunning) error {
	return a.db.Create(&row).Error
}

func (a *Action) GetRunning(runningId uint) (out *OperationRunning) {
	out = &OperationRunning{}
	if err := a.db.First(out, runningId).Error; err == nil {
		return out
	} else {
		return nil
	}
}
func (a *Action) GetRunnings() []OperationRunning {
	out := []OperationRunning{}
	a.db.Find(&out)
	return out
}

func (a *Action) GetLastRunnings(fromId uint) []OperationRunning {
	out := make([]OperationRunning, 0)
	a.db.Find(&out, "id >= ?", fromId)
	return out
}

func (a *Action) RemoveRunning(runningId uint) error {
	return a.db.Unscoped().Delete(&OperationRunning{}, "id = ?", runningId).Error
}

// func (a *Action) SaveRunningAndCli(optcli OperationCli) error {
// 	err := a.db.Exec("INSERT INTO operation_clis SELECT * FROM operation_runnings WHERE id = '?';DELETE FROM operation_runnings WHERE id = '?'", optcli.ID, optcli.ID).Error
// 	a.saveVersion()
// 	return err
// }

func (a *Action) SaveRunningAndCli(optcli OperationCli) error {
	a.Lock()
	defer a.Unlock()

	a.SaveOperationCli(optcli)
	a.saveVersion()
	return a.RemoveRunning(optcli.ID)
}

// func (a *Action) SaveRunningAndCli(optcli OperationCli) error {
// 	if err := a.SaveOperationCli(optcli); err != nil {
// 		return err
// 	} else {
// 		a.saveVersion()
// 		return a.RemoveRunning(optcli.ID)
// 	}
// }

// just for company
func (a *Action) saveVersion() {
	lastid := strconv.Itoa(int(a.LastOperationCliId()))
	ioutil.WriteFile(filepath.Join(utils.ExecDir(), VersionFilename), []byte(lastid), os.ModePerm)
}
