package sample

import (
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func GetDB(dbName string) (db *gorm.DB) {
	var err error
	switch {
	case strings.EqualFold(dbName, "mysql"):
		// 可以自定义 MySQL Driver，只需 import 并修改 mysql.Config.DriverName 参数为包名即可。
		// 如 import _ "example.com/my_mysql_driver"，对应的 DriverName 为 "my_mysql_driver"。
		// 也可以使用 sql.DB 初始化，只需将 mysql.Config.Conn 设为该对象的指针，Config 的其他项全设为空即可。
		// DSN 的参数可参考 https://github.com/go-sql-driver/mysql#parameters
		db, err = gorm.Open(mysql.New(mysql.Config{
			// DSN data source name
			DSN: "root:123456@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local",
			// string 类型字段的默认长度
			DefaultStringSize: 256,
			// 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
			DisableDatetimePrecision: true,
			// 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
			DontSupportRenameIndex: true,
			// 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
			DontSupportRenameColumn: true,
			// 根据当前 MySQL 版本自动配置
			SkipInitializeWithVersion: false,
		}), &gorm.Config{})
	case strings.EqualFold(dbName, "postgres") ||
		strings.EqualFold(dbName, "postgresql") ||
		strings.EqualFold(dbName, "postgre") ||
		strings.EqualFold(dbName, "pg"):
		// https://github.com/go-gorm/postgres
		// 与 MySQL 类似，可以自定义 PostgreSQL Driver 或使用 sql.DB 初始化
		db, err = gorm.Open(postgres.New(postgres.Config{
			DSN: "user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai",
			// disables implicit prepared statement usage
			PreferSimpleProtocol: true,
		}), &gorm.Config{})
	case strings.EqualFold(dbName, "sqlite") ||
		strings.EqualFold(dbName, "lite"):
		// github.com/mattn/go-sqlite3
		// 允许使用 cache=shared 参数
		db, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	case strings.EqualFold(dbName, "sqlserver") ||
		strings.EqualFold(dbName, "mssql"):
		// github.com/denisenkom/go-mssqldb
		dsn := "sqlserver://gorm:LoremIpsum86@localhost:9930?database=gorm"
		db, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	default:
		panic("unknown database name")
	}
	if err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
	return
}
