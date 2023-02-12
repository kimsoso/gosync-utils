package action

import "gorm.io/gorm"

// 目录操作表
type Operation struct {
	gorm.Model
	ID    uint   `gorm:"primaryKey;autoIncrement:true"`
	Path  string `gorm:"index;type:varchar(1024)"`
	Act   uint8  `gorm:"index;default:0"`
	IsDir bool   `gorm:"not null"`
	Md5   string `gorm:"type:char(32)"`
	Size  int64
}

type OperationCli struct {
	gorm.Model
	ID    uint   `gorm:"primaryKey;not null"`
	Path  string `gorm:"index;type:varchar(1024)"`
	Act   uint8  `gorm:"index;default:0"`
	IsDir bool   `gorm:"not null"`
	Md5   string `gorm:"type:char(32)"`
	Size  int64
}

type OperationRunning struct {
	gorm.Model
	ID    uint   `gorm:"primaryKey;not null"`
	Path  string `gorm:"index;type:varchar(1024)"`
	Act   uint8  `gorm:"index;default:0"`
	IsDir bool   `gorm:"not null"`
	Md5   string `gorm:"type:char(32)"`
	Size  int64
}

// 当前目录结构表，ID永远为1
type File struct {
	ID      uint `gorm:"primaryKey;not null"`
	Content string
}
