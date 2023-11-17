package sample

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/RinkoTaketsuki/GolangLearning/sql_example/utils"
)

var TableName = `users`

var CreateTableString = fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
	id		INT UNSIGNED NOT NULL UNIQUE AUTO_INCREMENT PRIMARY KEY,
	name	VARCHAR(255) NOT NULL
) ENGINE=InnoDB
`, TableName)

var DropTableString = fmt.Sprintf(`DROP TABLE %s`, TableName)

var ImportDataStrings = [...]string{
	`INSERT INTO users(name) VALUES("Alice")`,
	`INSERT INTO users(name) VALUES("Bob")`,
	`INSERT INTO users(name) VALUES("Catherine")`,
}

func CreateSampleTableAndImportSampleData(db *utils.DB) {
	res := db.Exec(CreateTableString)
	rowsAffected, err := res.RowsAffected()
	utils.ErrHandler(err)
	log.Printf("Create table SUCCESS. %d row(s) is affected.\n", rowsAffected)
	db.ExecMany(ImportDataStrings[:], make([][]any, len(ImportDataStrings)))
	log.Printf("Import data SUCCESS. %d row(s) is inserted.\n", len(ImportDataStrings))
}

func SampleCRUD(db *utils.DB) {
	db.ExecMany([]string{
		`INSERT INTO users VALUES(?,?)`,
		`DELETE FROM users WHERE id = ?`,
	}, [][]any{
		{
			sql.NullInt64{Int64: 0, Valid: false},
			"David",
		},
		{
			2,
		},
	})
	// 需注意 args 中包含 uint64 的情况，如果该值大于 math.MaxInt64，则会产生错误
	rows := db.Query(`SELECT id, name FROM users WHERE id = ?`, 3)
	// Close 也会返回错误，但这种错误我们通常无法处理，故一般忽略或 panic
	defer rows.Close()
	var (
		id   int
		name string
	)
	// Next 返回 false 时会自动 Close rows。
	// 已经迭代完最后一行时再次调用 Next 会返回 false，因此下面的循环在不出错时会迭代所有行。
	// Next 发生错误也会导致 Next 返回 false，且此时会调用 rows.Close()。
	// 即便如此也最好显式关闭 rows，像上面的 defer 语句那样处理，多次关闭无害。
	// 由于 Next 的这些特性，在迭代处理 rows 是请谨慎使用 break 或 continue 等流程控制语句。
	for rows.Next() {
		// 必须先 Next 再 Scan，即便是第一行
		// Scan 执行时，假设数据库某一列的类型是 VARCHAR，但其中存储的值都能转换成整数，我们也可以传入 *int
		// 即 Scan 会对扫描到的数据尝试转换成传入的指针类型对应的数据
		// 需要处理可能为 NULL 的列时，scan 需要输入 sql.NullXXX 类型，如果未提供可以仿照 sql.NullString 实现自己的可为 NULL 的类型
		// 更好的处理方法是在 SQL 中使用 COALESCE 函数
		utils.ErrHandler(rows.Scan(&id, &name))
		fmt.Println(id, name)
	}
	// 如果上面的迭代过程中 Next 因为错误而返回了 false，则 Err 方法会返回 Next 中生成的错误
	utils.ErrHandler(rows.Err())
}

func DropSampleTable(db *utils.DB) {
	res := db.Exec(DropTableString)
	rowsAffected, err := res.RowsAffected()
	utils.ErrHandler(err)
	log.Printf("Drop Table Success. %d row(s) is affected.\n", rowsAffected)
}
