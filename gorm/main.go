package main

import (
	"github.com/RinkoTaketsuki/GolangLearning/gorm/sample"
	"gorm.io/gorm"
)

func main() {
	sample.CreateExamples(&gorm.Session{SkipHooks: false, DryRun: false})
	sample.FindExamples(&gorm.Session{SkipHooks: false, DryRun: false})
}
