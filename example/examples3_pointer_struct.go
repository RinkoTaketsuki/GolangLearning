package example

import "fmt"

type Vertex struct {
	X, Y float64
}

var (
	v1 = Vertex{1, 2}  // 创建一个 Vertex 类型的结构体
	v2 = Vertex{X: 1}  // Y:0 被隐式地赋予
	v3 = Vertex{}      // X:0 Y:0
	p  = &Vertex{1, 2} // 创建一个 *Vertex 类型的结构体（指针）
)

func Run3() {
	fmt.Println("--------- Run 3 ---------")
	fmt.Println(v1, p, v2, v3, v1.X, p.Y) // p.Y 等价于 (*p).Y
	// 结构体赋值是值传递
	v4 := v1
	v4.X = 42
	fmt.Printf("v1: %v\n", v1)
}
