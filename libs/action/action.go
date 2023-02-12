package action

import (
	"sync"

	"gorm.io/gorm"
)

type Action struct {
	sync.Mutex
	db *gorm.DB
}

func New(db *gorm.DB) *Action {
	return &Action{
		db: db,
	}
}
