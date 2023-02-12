package dao

import (
	"btsync-utils/libs/action"
	"btsync-utils/libs/utils"
	"log"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	// "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	dbPath = filepath.Join(utils.ExecDir(), "opts.db")
)

func NewDB() *gorm.DB {
	log.Println("sqlite db file is:", dbPath)

	_, errExist := os.Stat(dbPath)

	if db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}); err != nil {
		log.Panic("初始化数据库错误,", err)
		return nil
	} else {
		if sqlDB, err := db.DB(); err != nil {
			log.Fatalln("sqlite sqldb error:", err)
		} else {
			sqlDB.SetMaxIdleConns(1)    //最大空闲连接数
			sqlDB.SetMaxOpenConns(1)    //最大连接数
			sqlDB.SetConnMaxLifetime(0) //设置连接空闲超时，不超时
			sqlDB.SetConnMaxIdleTime(0) //设置连接空闲超时，不超时
		}

		if os.IsNotExist(errExist) {
			initAllTables(db)
		}
		return db
	}
}

// 初始化所有表结构
func initAllTables(gdb *gorm.DB) {
	log.Println("init all tables in db")

	NewTable(gdb, "operations", action.Operation{}, false)
	NewTable(gdb, "operationClis", action.OperationCli{}, false)
	NewTable(gdb, "operationRunnings", action.OperationRunning{}, false)
	NewTable(gdb, "files", action.File{}, false)
}

// 新建数据库表
func NewTable(db *gorm.DB, tableName string, tableStruct interface{}, dropExist bool) {
	if dropExist && db.Migrator().HasTable(tableName) {
		db.Migrator().DropTable(tableStruct)
	}
	db.Migrator().CreateTable(tableStruct)
}
