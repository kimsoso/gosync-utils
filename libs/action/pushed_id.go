// 处理服务器端向客户端同步所需的相关操作
package action

// 同步记录表，对每个子网同步到的ID
type Sync struct {
	Subnet   string `gorm:"primaryKey;type:varchar(24)"`
	PushedID uint   `gorm:""`
}

// 记录每个子网已经同步的最后ID
func (act *Action) SetPushedID(subnet string, doneID uint) error {
	if act.db.First(&Sync{}, "subnet = ?", subnet).Error == nil {
		return act.db.Model(Sync{}).Where("subnet = ?", subnet).Update("pushed_id", doneID).Error
	} else {
		return act.db.Create(&Sync{
			Subnet:   subnet,
			PushedID: doneID,
		}).Error
	}
}
