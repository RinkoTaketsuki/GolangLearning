package utils

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"gorm.io/gorm"
)

func PrintBlock(str string) {
	width := 0
	for _, r := range str {
		if utf8.RuneLen(r) >= 2 {
			width += 2
		} else {
			width++
		}
	}
	var sb strings.Builder
	bar := "+" + strings.Repeat("-", width+2) + "+\n"
	in := "| " + str + " |\n"
	sb.WriteString(bar)
	sb.WriteString(in)
	sb.WriteString(bar)
	fmt.Print(sb.String())
}

func PrintStmt(stmt *gorm.Statement) {
	if stmt == nil {
		return
	}
	fmt.Printf("SQL: %s\n", stmt.SQL.String())
	for i, v := range stmt.Vars {
		fmt.Printf("Var %d:\t %#v\n", i, v)
	}
}

func GormModelToString(m *gorm.Model) string {
	var sb strings.Builder
	sb.WriteString("GORM model:\t")
	sb.WriteString(strconv.Itoa(int(m.ID)))
	sb.WriteString("\n  Created at:\t")
	sb.WriteString(m.CreatedAt.String())
	sb.WriteString("\n  Updated at:\t")
	sb.WriteString(m.UpdatedAt.String())
	sb.WriteString("\n  Deleted at:\t")
	if m.DeletedAt.Valid {
		sb.WriteString(m.DeletedAt.Time.String())
		sb.WriteByte('\n')
	} else {
		sb.WriteString("NULL\n")
	}
	return sb.String()
}
