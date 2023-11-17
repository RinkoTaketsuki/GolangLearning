package main

import (
	"database/sql"
	"database/sql/driver"

	"github.com/RinkoTaketsuki/GolangLearning/gorm/sample"
	"gorm.io/gorm"
)

type Consumer interface {
	Consume()
}

type Producer interface {
	Produce()
}

type A struct{}

func (A) Consume()  {}
func (*A) Produce() {}

var _ Consumer = A{}
var _ Consumer = &A{}
var _ Producer = &A{}

func (*A) Scan(value interface{}) error {
	return nil
}

func (*A) Value() (driver.Value, error) {
	return nil, nil
}

var _ driver.Valuer = &A{}
var _ sql.Scanner = &A{}

func main() {
	sample.CreateExamples(&gorm.Session{SkipHooks: false, DryRun: false})
	sample.FindExamples(&gorm.Session{SkipHooks: false, DryRun: false})
}
