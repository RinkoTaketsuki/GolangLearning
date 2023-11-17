package example

import (
	"fmt"
	"strings"
)

type Coordinate struct {
	Lat, Long float64
}

// 类似 slice，单纯声明会使其为 nil
var nilMap map[string]Coordinate

// map 的初始化形式
var exampleMap = map[string]Coordinate{
	"Bell Labs": {40.68433, -74.39967},
	"Google":    {37.42202, -122.08408},
}

func mapUsage() {
	fmt.Printf("exampleMap: %v\n", exampleMap)

	if nil == nilMap {
		fmt.Println("This map is nil, you should call \"make\" to allocate it.")
	}
	nilMap = make(map[string]Coordinate)
	// 赋值
	nilMap["Bell Labs"] = Coordinate{
		40.68433, -74.39967,
	}
	// 引用
	fmt.Println(nilMap["Bell Labs"])

	m := make(map[string]int)
	// 赋值
	m["Answer"] = 42
	fmt.Println("The value:", m["Answer"])
	// 修改
	m["Answer"] = 48
	fmt.Println("The value:", m["Answer"])
	// 删除，此时直接引用会返回默认零值，key 参数不存在时也不会产生错误
	delete(m, "Answer")
	fmt.Println("The value:", m["Answer"])
	// 用第二个返回值检查 key 是否存在
	v, ok := m["Answer"]
	fmt.Println("The value:", v, "Present?", ok)
}

func wordCount(s string) (m map[string]int) {
	words := strings.Fields(s)
	m = make(map[string]int)
	for _, word := range words {
		m[word]++
	}
	return
}

func Run5() {
	fmt.Println("--------- Run 5 ---------")
	mapUsage()
	fmt.Printf("wordCount(\"I ate a donut. Then I ate another donut.\"): %v\n", wordCount("I ate a donut. Then I ate another donut."))
}
