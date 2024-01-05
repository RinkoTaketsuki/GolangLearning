package sample

import (
	"time"

	"github.com/RinkoTaketsuki/GolangLearning/gorm/sample/model"
	"github.com/RinkoTaketsuki/GolangLearning/gorm/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CreateExamples(config *gorm.Session) {
	// gorm.DB.Session 可用于给 DB 添加设定
	db := GetDB("MySQL").Session(config)

	// 建表
	db.AutoMigrate(new(model.User), new(model.Class))

	utils.PrintBlock("*** Insert one object ***")

	user := model.User{Name: "Alice", Phone: 110, Birthday: time.Now().AddDate(-8, 0, 0)}

	// INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`name`,`phone`,`birthday`) VALUES (?,?,?,?,?,?)
	// 当 user.ID 为空时，Create 后 ID 会被置为自增主键的值，默认初始值为 1
	// 当 user.CreatedAt 为空时，Create 后 CreatedAt 会被置为 time.Now()
	// 当 user.UpdatedAt 为空时，Create 后 UpdatedAt 会被置为 time.Now()
	// 其他属性为空时会被设定对应的零值
	db.Create(&user)

	// INSERT INTO `users` (`created_at`,`updated_at`,`name`,`phone`) VALUES (?,?,?,?)
	// CreatedAt 和 UpdatedAt 会被隐式 Select
	db.Select("Name", "Phone", "CreatedAt").Create(&user)

	// Omit 不存在隐含选择，Omit 剩下的列便是所有要 INSERT 的列
	// 如果不 Omit ID 会导致 INSERT 中包含 ID 而发生错误
	// INSERT INTO `users` (`updated_at`,`deleted_at`,`birthday`) VALUES (?,?,?)
	db.Omit("ID", "Name", "Phone", "CreatedAt").Create(&user)

	utils.PrintBlock("*** Insert some objects ***")
	var users = []model.User{
		{Name: "Catherine", Birthday: time.Date(1999, 12, 31, 23, 59, 59, 0, time.Local)},
		{Name: "David", Birthday: time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local)},
		{Name: "Elizabeth", Birthday: time.Date(1919, 8, 10, 11, 45, 14, 893, time.Local)}}

	// 多行 INSERT，Create Hooks 会被多次调用（每个 model 调用一次），但生成的 SQL 只有一句
	db.Create(&users)

	// 清除 User 对象中的 gorm.Model 信息
	for i := range users {
		users[i].Model = gorm.Model{}
	}

	// BATCH INSERT，Create Hooks 同样会被多次调用，但生成的 SQL 一次最多插入 batchSize 行数据
	// 生成 len(users) / batchSize 句 SQL，除最后一句外，其他 SQL 均插入 batchSize 行数据
	db.CreateInBatches(&users, 2)

	utils.PrintBlock("*** Insert one or more maps ***")

	// 须注意的是，插入 map 引起的表中的自增主键变化不会体现在任何一个 User 对象中，
	// 所以当 map 与模型混用时要特别注意自增主键的问题
	// 通过 map 插入时只会插入给定的列，没有任何隐式 Select 或隐式 Omit 的列
	userModel := db.Model(new(model.User))
	userModel.Create(map[string]any{
		"Name":     "Bob",
		"Phone":    911,
		"Birthday": time.Date(1900, 12, 25, 11, 11, 11, 0, time.Local),
	})

	// 插入多个 map
	userModel.Create([]map[string]any{
		{
			"Name":     "Floyd",
			"Birthday": time.Now().Add(-100000 * time.Hour),
		},
		{
			"Name":     "Gauss",
			"Birthday": time.Now().Add(-150000 * time.Hour),
		},
	})

	utils.PrintBlock("*** Insert a object with self-defined type members ***")

	// 如果模型中某个非嵌入成员包含自定义的类型，该类型需要实现
	// schema.GormDataTypeInterface 和 gorm.Valuer 这两个接口
	// 下面分别使用模型和 map 两种方式插入自定义类型的数据
	obj1 := model.Class{
		ThisIsID: model.MyID{
			Name: "Obj1",
		},
		Loc: model.Location{X: 10, Y: 20},
	}
	db.Create(&obj1)
	classModel := db.Model(model.Class{})
	// 此处主键如果没有默认值必须全部指定
	classModel.Create(map[string]any{
		"Num":  2,
		"Name": "Obj2",
		"Loc":  clause.Expr{SQL: "ST_PointFromText(?)", Vars: []interface{}{"POINT(100 200)"}},
	})
	utils.PrintStmt(classModel.Statement)

	// utils.PrintBlock("*** Insert associated objects ***")

	// citizen := model.Citizen{
	// 	Name: "John",
	// 	Passport: model.Passport{
	// 		Num:     114514,
	// 		ExpTime: time.Now().AddDate(10, 0, 0),
	// 	},
	// }
	// db.Create(&citizen)
	// db.Omit("Passport").Create(citizen)
	// // omit all associations
	// db.Omit(clause.Associations).Create(&user)
	// time.Sleep(time.Second * 3)

	utils.PrintBlock("*** Upsert ***")

	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&users)
	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.Assignments(map[string]any{"phone": 114514}),
	}).Create(&users)
	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.Assignments(map[string]any{"phone": gorm.Expr("GREATEST(phone, VALUES(phone))")}),
	}).Create(&users)
	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "phone"}),
	}).Create(&users)

	// UPDATE 不包括 ID 和 CreatedAt
	db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&users)
}

func FindExamples(config *gorm.Session) {
	// gorm.DB.Session 可用于给 DB 添加设定
	db := GetDB("MySQL").Session(config)

	// 建表
	db.AutoMigrate(new(model.User), new(model.Class))

	utils.PrintBlock("*** Select one object ***")

	firstUser := model.User{}
	// SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1
	// 模型有主键约束时，使用主键排序，否则使用第一个列排序
	db.First(&firstUser)

	lastUser := model.User{}
	// SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL ORDER BY `users`.`id` DESC LIMIT 1
	// 模型有主键约束时，使用主键排序，否则使用第一个列排序
	db.Last(&lastUser)

	user1 := model.User{}
	// SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL LIMIT 1
	db.Take(&user1)

	user2 := model.User{}
	// 等价于 Take()；若不加 Limit，Find 会进行全表扫描但只返回第一个结果
	db.Limit(1).Find(&user2)

	// 使用 map 接收 First、Last、Take 的结果，基于 Model 方法传入的模型排序
	firstUserMap := make(map[string]any)
	db.Model(new(model.User)).First(&firstUserMap)
	lastUserMap := make(map[string]any)
	db.Model(new(model.User)).Last(&lastUserMap)
	userMap1 := make(map[string]any)
	db.Model(new(model.User)).Take(&userMap1)

	// 指定 Table 而不指定 Model 时，Take 可以正常工作，但 First 和 Last 会无法排序
	userMap2 := make(map[string]any)
	db.Table("users").Take(&userMap2)

	// 指定主键查找
	id10User1 := model.User{}
	db.First(&id10User1, 10)
	id10User2 := model.User{}
	db.First(&id10User2, "10")
	id10User3 := model.User{Model: gorm.Model{ID: 10}}
	db.First(&id10User3)
	id10User4 := model.User{}
	db.Model(&model.User{Model: gorm.Model{ID: 10}}).First(&id10User4)

	utils.PrintBlock("*** Select all objects ***")

	allUsers1 := make([]model.User, 0)
	db.Find(&allUsers1)

	allUsers2 := make([]model.User, 0)
	// 与 CreateExamples 示例中的 Select 含义类似
	db.Select("name", "phone").Find(&allUsers2)

	allUsers3 := make([]model.User, 0)
	db.Select("COALESCE(phone,?)", 999).Find(&allUsers3)

	allUsers4 := make([]model.User, 0)
	db.Limit(-1).Find(&allUsers4)

	orderedAllUsers1 := make([]model.User, 0)
	db.Order("phone desc, name").Find(&orderedAllUsers1)

	orderedAllUsers2 := make([]model.User, 0)
	db.Order("phone desc").Order("name").Find(&orderedAllUsers2)

	orderedAllUsers3 := make([]model.User, 0)
	db.Clauses(clause.OrderBy{
		Expression: clause.Expr{
			SQL:                "FIELD(id,?)",
			Vars:               []interface{}{[]int{1, 2, 3}},
			WithoutParentheses: true,
		},
	}).Find(&orderedAllUsers3)

	utils.PrintBlock("*** Select filtered objects ***")

	id123Users := make([]model.User, 0)
	// 使用 IN 表达式查询
	db.Where([]int{1, 2, 3}).Find(&id123Users)

	filteredUsers1 := make([]model.User, 0)
	// 如果 db 指定了参数非零的 Model 或 Find 内模型包含非零参数，这些参数会以 AND 运算符与 WHERE 条件连接
	// 多个 Where 方法可以串联，同样使用 AND 运算符连接
	// 以下 Where 方法中输入的参数同样可作为 Find、First、Last、Take 的 conds 参数
	// Not 方法与 Where 类似，只不过会在整个条件上加 NOT 运算符
	// Or 方法与 Where 类似，只不过与其他 Where 方法串联时使用 OR 运算符
	// Having 方法与 Where 类似
	db.Where("name LIKE ? OR phone = ?", "%Alice%", 110).Find(&filteredUsers1)

	filteredUsers2 := make([]model.User, 0)
	// Where 接收结构体指针参数
	db.Where(&model.User{Phone: 110}).Find(&filteredUsers2)

	filteredUsers3 := make([]model.User, 0)
	// Where 接收 map 参数，相比结构体，map 可以查询某一列为零的情况
	db.Where(map[string]any{"phone": 110}).Find(&filteredUsers3)

	filteredUsers4 := make([]model.User, 0)
	// Where 接收结构体指针参数时，可以指定查询哪几列，可用于查询零值
	db.Where(&model.User{Name: "Alice"}, "name", "phone").Find(&filteredUsers4)

	threeUsers := make([]model.User, 0)
	db.Offset(2).Limit(3).Find(&threeUsers)
}
