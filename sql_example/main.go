package main

import (
	"database/sql"

	"github.com/RinkoTaketsuki/GolangLearning/sql_example/utils"

	// 匿名 import 需要使用的驱动，一是为了执行其中的 init 函数，二是将驱动注册到 database/sql 包
	// 该驱动不支持多语句、多结果和存储过程
	_ "github.com/go-sql-driver/mysql"
)

/*
连接池意味着在单个数据库上执行两个连续的语句可能会打开两个连接并分别执行它们。
例如，LOCK TABLES 后跟 INSERT 可能会阻塞，这是因为 INSERT 位于一个不持有表锁的连接上。
尝试使用 USE 等语句修改连接状态时，本次 Exec 或 Query 修改的是当前连接的状态，后面执行的语句不一定会使用这个连接。
且 USE 对应的连接被回收到连接池时还会污染其他语句执行时的环境。
同理，BEGIN、COMMIT、ROLLBACK 等语句也不应被当作 Query 或 Exec 直接使用。

连接是在需要并且池中没有可用的空闲连接时才创建的。默认情况下，连接数量没有限制。
如果尝试一次执行很多操作，则可以创建任意数量的连接。这可能导致数据库返回错误，例如「连接过多」。

在 Go 1.1 或更高版本中，您可以使用 db.SetMaxIdleConns(N) 来限制池中空闲连接的数量。不过，这并不限制池的大小。
在 Go 1.2.1 或更高版本中，您可以使用 db.SetMaxOpenConns(N) 来限制与数据库的打开连接数。
不幸的是，https://groups.google.com/d/msg/golang-dev/jOTqHxI09ns/x79ajll-ab4J 说明 db.SetMaxOpenConns(N) 在 1.2 中无法安全使用。

连接回收速度非常快。使用 db.SetMaxIdleConns(N) 设置大量空闲连接可以减少这种搅动，并有助于保持连接可重复使用。
但长时间保持空闲连接会导致问题 (例如 Microsoft Azure 上的 MySQL 中的 https://github.com/go-sql-driver/mysql/issues/257)。
如果由于连接空闲时间太长而导致连接超时，请尝试 db.SetMaxIdleConns(0)。

您还可以通过设置 db.SetConnMaxLifetime(duration) 指定连接可重用的最长时间，因为重用长时间的连接可能会导致网络问题。
这会延迟关闭未使用的连接，即可能会延迟关闭过期的连接。

但需要注意 Tx 是与连接绑定的，即事务中所有语句在同一连接内执行。
由此造成的问题是，如果事务过程中某个 Rows 没有被关闭却尝试 Query 另一个结果会产生错误。
*/
func main() {
	// Open 并不会创建到数据库服务器的连接，连接行为会推迟到第一次需要连接时
	// 请勿频繁进行 Open 和 Close，这容易导致一些网络故障，应该在自己的业务逻辑中传递 *sql.DB
	// sql 内置的连接错误处理机制通常可以使我们无须考虑连接中断的问题，一般会重试 10 次连接
	db, err := sql.Open("mysql", "sysdba:transwarp@tcp(localhost:15307)/test")
	utils.ErrHandler(err)
	defer db.Close()
	// 如果需要立即建立连接可使用 Ping()，该函数会检查输入的 dataSourceName 是否可用
	utils.ErrHandler(db.Ping())
	utilsDB := (*utils.DB)(db)
	// sample.CreateSampleTableAndImportSampleData(utilsDB)
	// sample.SampleCRUD(utilsDB)
	// sample.DropSampleTable(utilsDB)
	// utilsDB.QueryResultPrintln("select * from mysql.user where host = ?", "%")
	utilsDB.QueryResultPrintln("select ConCat(?, 'AbC')", "DeF")
}
