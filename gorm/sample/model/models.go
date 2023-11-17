package model

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/RinkoTaketsuki/GolangLearning/gorm/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
默认情况下，模型和关系表之间的映射关系是：

	ID -> 主键
	CreatedAt -> 创建时间
	UpdatedAt -> 更新时间
	结构体名（驼峰）的复数 -> 下划线分隔的表名
	字段名（驼峰） -> 下划线分隔的列名

嵌入字段会被展开，所以为方便起见最好直接嵌入 gorm.Model

下面是该 model 在 MySQL 下对应的表，表名为 users

	+------------+-----------------+------+-----+---------+----------------+
	| Field      | Type            | Null | Key | Default | Extra          |
	+------------+-----------------+------+-----+---------+----------------+
	| id         | bigint unsigned | NO   | PRI | null    | auto_increment |
	| created_at | datetime        | YES  |     | null    |                |
	| updated_at | datetime        | YES  |     | null    |                |
	| deleted_at | datetime        | YES  | MUL | null    |                |
	| name       | varchar(256)    | YES  |     | null    |                |
	| phone      | bigint unsigned | YES  |     | null    |                |
	| birthday   | datetime        | YES  |     | null    |                |
	+------------+-----------------+------+-----+---------+----------------+

	CREATE TABLE `users` (
		`id` bigint unsigned AUTO_INCREMENT,
		`created_at` datetime NULL,
		`updated_at` datetime NULL,
		`deleted_at` datetime NULL,
		`name` varchar(256),
		`phone` bigint unsigned,
		`birthday` datetime NULL,
		PRIMARY KEY (`id`),
		INDEX `idx_users_deleted_at` (`deleted_at`)
	);
*/
type User struct {
	gorm.Model
	Name     string
	Phone    uint
	Birthday time.Time
}

func (u *User) String() string {
	var sb strings.Builder
	sb.WriteString("User:\t")
	sb.WriteString(u.Name)
	sb.WriteByte('\n')
	sb.WriteString(utils.GormModelToString(&u.Model))
	sb.WriteString("Phone:\t")
	sb.WriteString(strconv.Itoa(int(u.Phone)))
	sb.WriteString("\nBirthday:\t")
	sb.WriteString(u.Birthday.String())
	sb.WriteByte('\n')
	return sb.String()
}

/*
Create Hook 相关的执行顺序：

	开始事务
	BeforeSave
	BeforeCreate
	关联前的 save
	插入记录至 db
	关联后的 save
	AfterCreate
	AfterSave
	提交或回滚事务

所有 Hook 返回 error 都会导致事务 Rollback
*/
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.Name == "" {
		return errors.New("try to insert a user with an empty name")
	}
	return nil
}

// 调试输出
func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	utils.PrintBlock("User.AfterCreate")
	utils.PrintStmt(tx.Statement)
	if tx.Statement.RowsAffected != 0 {
		fmt.Printf("Rows affected: %d\n", tx.Statement.RowsAffected)
	}
	if tx.Error != nil {
		fmt.Printf("Error: %s\n", tx.Error.Error())
	}
	fmt.Printf("%s\n", u)
	return nil
}

func (u *User) AfterFind(tx *gorm.DB) (err error) {
	utils.PrintBlock("User.AfterFind")
	utils.PrintStmt(tx.Statement)
	if tx.Statement.RowsAffected != 0 {
		fmt.Printf("Rows returned: %d\n", tx.Statement.RowsAffected)
	}
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			fmt.Println("Error: No data found")
		} else {
			fmt.Printf("Error: %s\n", tx.Error.Error())
		}
	}
	fmt.Printf("%s\n", u)
	return nil
}

/*
Class 用于示范：
 1. 如何使用 GORM 标签控制读写权限，其中 write 权限指同时具有 create 和 update 权限
 2. 如何使用 GORM 标签将普通字段作为嵌入结构体来使用
 3. 如何使用自定义类型字段
 4. 如何设定默认值

须注意的一点是，如果某字段是非指针基础类型（bool、int、string 等），且要写入的值是零值（false、0、"" 等），
且该字段设定了默认值，则该零值会被忽略，转而插入默认值。若要避免这种情况，可将字段类型改为指针或 Scanner 或 Valuer。

有些情况（如 generated column）需要忽略具体的默认值，此时使用 default:(-) 标签

	CREATE TABLE `classes` (
		`id_num` bigint unsigned,
		`id_name` varchar(256),
		`m_read_write42` bigint DEFAULT 42,
		`m_read1` bigint,
		`m_read_create1` GENERATED ALWAYS AS (m_read1 + m_read2),
		`m_read_update` bigint,
		`m_read2` bigint,
		`m_read_create2` bigint,
		`m_create` bigint,
		`created` datetime NULL,
		`updated_timestamp` bigint,
		`created_unix` bigint,
		`updated_milli` bigint,
		`created_nano` bigint,
		`loc` geometry,
		PRIMARY KEY (`id_num`,`id_name`)
	);
*/
type Class struct {
	ThisIsID     MyID         `gorm:"embedded;embeddedPrefix:id_"`
	MReadWrite42 int          `gorm:"<-;default:42"`
	MRead1       int          `gorm:"<-:false"`
	MReadCreate1 int          `gorm:"<-:create"`
	MReadUpdate  int          `gorm:"<-:update"`
	MRead2       int          `gorm:"->;type:bigint AS (m_read_create1 + m_read_create2);default:(-);"`
	MReadCreate2 int          `gorm:"->;<-:create"`
	MCreate      int          `gorm:"->:false;<-:create"`
	MIgnore      int          `gorm:"-"`
	TR           TimeRecorder `gorm:"embedded"`
	Loc          Location
}

// primaryKey 有两个时不会生成自增主键
//
// autoIncrement:false 标签可以显式地使主键不带 AUTO_INCREMENT 属性
type MyID struct {
	Num  uint   `gorm:"primaryKey"`
	Name string `gorm:"primaryKey"`
}

/*
TimeRecorder 用于示范如何使用 GORM 标签和字段类型自定义存储创建时间和更新时间的字段
*/
type TimeRecorder struct {
	Created          time.Time `gorm:"autoCreateTime"`
	UpdatedTimestamp int64     `gorm:"autoUpdateTime"`
	CreatedUnix      int       `gorm:"autoCreateTime"`
	UpdatedMilli     int64     `gorm:"autoUpdateTime:milli"`
	CreatedNano      int64     `gorm:"autoUpdateTime:nano"`
}

type Location struct {
	X, Y int
}

// 定义 SQL 中的数据类型
func (loc Location) GormDataType() string {
	return "geometry"
}

// 定义插入时如何将当前类型的数据转化为 SQL 中的数据
func (loc Location) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	return clause.Expr{
		SQL:  "ST_PointFromText(?)",
		Vars: []any{fmt.Sprintf("POINT(%d %d)", loc.X, loc.Y)},
	}
}

// Citizen 用于示范关联模型
type Citizen struct {
	gorm.Model
	Name     string
	Passport Passport
}

type Passport struct {
	gorm.Model
	Num     uint
	ExpTime time.Time
}
