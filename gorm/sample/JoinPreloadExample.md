# GORM Join Preload: Nested model example

```go
type User struct {
    gorm.Model
    name        string
    age         string
    birthday    time.Time
    active      bool
    ManagerID   *uint
    Manager     *User
    CompanyID   uint
    Company     Company
}

type Company struct {
    ID      uint
    name    string
}

func fullFetchUsers(users []Users) {
    db.Joins("Manager").Joins("Manager.Company").Find(&users)
}
```

```sql
SELECT 
    "users"."id",
    "users"."created_at",
    "users"."updated_at",
    "users"."deleted_at",
    "users"."name",
    "users"."age",
    "users"."birthday",
    "users"."company_id",
    "users"."manager_id",
    "users"."active",
    "Manager"."id" AS "Manager__id",
    "Manager"."created_at" AS "Manager__created_at",
    "Manager"."updated_at" AS "Manager__updated_at",
    "Manager"."deleted_at" AS "Manager__deleted_at",
    "Manager"."name" AS "Manager__name",
    "Manager"."age" AS "Manager__age",
    "Manager"."birthday" AS "Manager__birthday",
    "Manager"."company_id" AS "Manager__company_id",
    "Manager"."manager_id" AS "Manager__manager_id",
    "Manager"."active" AS "Manager__active",
    "Manager__Company"."id" AS "Manager__Company__id",
    "Manager__Company"."name" AS "Manager__Company__name"
FROM 
    "users"
    LEFT JOIN
    "users" "Manager"
    ON
    "users"."manager_id" = "Manager"."id" AND "Manager"."deleted_at" IS NULL
    LEFT JOIN 
    "companies" "Manager__Company" 
    ON 
    "Manager"."company_id" = "Manager__Company"."id" 
WHERE 
    "users"."deleted_at" IS NULL
```
