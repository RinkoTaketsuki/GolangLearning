package utils

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
)

type DB sql.DB

func Slice(args ...any) []any {
	return args
}

// 给 err 加上一些 debug 信息然后直接 log.Fatalf 退出程序
func ErrHandler(e error) {
	if e != nil {
		// 处理特定 SQL 驱动对应的错误的方法举例
		if driverErr, ok := e.(*mysql.MySQLError); ok && driverErr.Number == mysqlerr.ER_ACCESS_DENIED_ERROR {
			e = fmt.Errorf("access denied by MySQL server")
		}
		pc := make([]uintptr, 1)
		runtime.Callers(2, pc)
		f := runtime.FuncForPC(pc[0])
		file, line := f.FileLine(pc[0])
		base := filepath.Base(file)
		runtimeMsg := fmt.Sprintf("[%s: line %d (%s at 0x%x)]", base, line, f.Name(), f.Entry())
		log.Fatalf("%s %s", runtimeMsg, e.Error())
	}
}

// 注：占位符格式说明：
//
//	MySQL：?
//	PG：$N（N 表示第 N 个占位符）
//	SQLite：? 或 $N
//	Oracle：:para_name（用名称标记占位符而不是序数）
func (db *DB) Exec(query string, args ...any) sql.Result {
	sqlDB := (*sql.DB)(db)
	// stmt 内部会记录一个连接信息，其被执行时会优先使用内部记录的连接
	// 如果当前连接不可用（繁忙、中断等），则会在连接池中另找一个可用的连接，并重新 Prepare
	// 当前业务并发量较高时（大量连接繁忙）这种特性会成为性能瓶颈
	stmt, err := sqlDB.Prepare(query)
	ErrHandler(err)
	defer stmt.Close()
	// Exec 或 Query 在正常情况下会 Close stmt
	// 相比不 Prepare 的方式，Prepare 的方式需要和数据库多通信一次
	// 这里不 Prepare 的方式是指直接使用 sql.DB.Exec，且只有一个纯 SQL 参数，没有其他参数
	// 当其他参数不为空时，Exec 内部还是会进行 Prepare
	res, err := stmt.Exec(args...)
	ErrHandler(err)
	return res
}

func (db *DB) ExecMany(queries []string, args [][]any) {
	sqlDB := (*sql.DB)(db)
	// Begin 方法和开启事务的 SQL 不能混用，否则会有一系列的未定义行为
	tx, err := sqlDB.Begin()
	defer tx.Rollback()
	ErrHandler(err)
	for i, q := range queries {
		// tx 也有 Prepare, Query, Exec 等方法，与 sql.DB 的对应方法相同
		// sql.Tx.Prepare 创建的 sql.Stmt 会与事务绑定，而此时使用 sql.DB.Prepare 会导致使用新的连接
		// 其他如 set 变量、创建临时表、更改连接设置等操作获取 Tx 对象时也要使用 tx 而不是 db
		// Tx.Stmt() 方法可以把传入的 stmt 绑定的连接改为当前 tx 的连接，但不建议使用
		// 早期版本的 Golang 必须保证 stmt 比 tx 先关闭，如果要 defer stmt.Close() 需要注意这个问题
		stmt, err := tx.Prepare(q)
		ErrHandler(err)
		_, err = stmt.Exec(args[i]...)
		stmt.Close()
		if err != nil {
			break
		}
	}
	ErrHandler(tx.Commit())
}

// rows 可以多次关闭不会出问题，Close 不会对已关闭的 row 做任何事。
// 因此应该 defer 所有 sql.Rows 的 Close。
func (db *DB) Query(query string, args ...any) *sql.Rows {
	sqlDB := (*sql.DB)(db)
	stmt, err := sqlDB.Prepare(query)
	ErrHandler(err)
	defer stmt.Close()
	// 当不需要返回结果时请勿使用 Query，而应该使用 Exec，因为 sql.Rows 在被 Close 前会一直维持数据库连接
	res, err := stmt.Query(args...)
	ErrHandler(err)
	return res
}

// 读取未知列的表的正确姿势
func (db *DB) QueryResultPrintln(query string, args ...any) {
	var err error

	rows := db.Query(query, args...)
	defer rows.Close()

	// 单独获取列名列表
	// columns, err := rows.Columns()
	// ErrHandler(err)

	// 获取列数据类型列表，其中包含列名列表
	columnTypes, err := rows.ColumnTypes()
	ErrHandler(err)

	maxColNameLen := 0
	maxColTypeNameLen := 0
	for _, ct := range columnTypes {
		colNameLen := len(ct.Name())
		colTypeNameLen := len(ct.DatabaseTypeName())
		if colNameLen > maxColNameLen {
			maxColNameLen = colNameLen
		}
		if colTypeNameLen > maxColTypeNameLen {
			maxColTypeNameLen = colTypeNameLen
		}
	}
	maxColNameLen = (((maxColNameLen + 1) / 8) + 1) * 8
	maxColTypeNameLen = ((maxColTypeNameLen / 8) + 1) * 8

	vals := make([]any, len(columnTypes))
	for i := range vals {
		vals[i] = new(sql.RawBytes) // 使用 %s 输出支持 UTF-8 字符
	}

	for rowCnt := 1; rows.Next(); rowCnt++ {
		ErrHandler(rows.Scan(vals...))
		fmt.Println("----------", rowCnt, "----------")
		for i, val := range vals {
			data := *(val.(*sql.RawBytes))
			if data == nil {
				for n, _ := fmt.Printf("%s:", columnTypes[i].Name()); n < maxColNameLen; n++ {
					fmt.Printf(" ")
				}
				for n, _ := fmt.Printf("%s", columnTypes[i].DatabaseTypeName()); n < maxColTypeNameLen; n++ {
					fmt.Printf(" ")
				}
				fmt.Println("<NULL>")
			} else {
				for n, _ := fmt.Printf("%s:", columnTypes[i].Name()); n < maxColNameLen; n++ {
					fmt.Printf(" ")
				}
				for n, _ := fmt.Printf("%s", columnTypes[i].DatabaseTypeName()); n < maxColTypeNameLen; n++ {
					fmt.Printf(" ")
				}
				fmt.Printf("[ %s ]\n", data)
			}
		}
	}
	ErrHandler(rows.Err())
}

func (db *DB) QueryOneRow(prepare bool, query string, args []any, dest ...any) (has bool) {
	var err error
	var row *sql.Row // sql.Row 对象无须 Close
	sqlDB := (*sql.DB)(db)
	if !prepare {
		row = sqlDB.QueryRow(query, args...)
	} else {
		stmt, err := sqlDB.Prepare(query)
		ErrHandler(err)
		defer stmt.Close()
		row = stmt.QueryRow(args...)
	}
	// 不同于 sql.Rows.Scan，sql.Row.Scan 在其对应的 QueryRow 查询不到结果时也会返回 error
	// 故需要对这种特殊 error 进行特殊处理
	err = row.Scan(dest...)
	if err != nil {
		if err == sql.ErrNoRows {
			has = false
		} else {
			ErrHandler(err)
		}
	} else {
		has = true
	}
	return
}
