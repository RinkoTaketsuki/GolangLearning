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
	db.First(&firstUser)

	lastUser := model.User{}
	// SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL ORDER BY `users`.`id` DESC LIMIT 1
	db.Last(&lastUser)

	user1 := model.User{}
	// SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL LIMIT 1
	db.Take(&user1)

	user2 := model.User{}
	// 等价于 Take()；若不加 Limit，Find 会进行全表扫描但只返回第一个结果
	db.Limit(1).Find(&user2)
}
