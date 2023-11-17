package example

import (
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"time"
)

// 方法不过是把 this 当成一个参数。
// 分成值接收者和指针接收者两种情形。
// 方法和结构体的定义必须在同一 package 内。
func (v Vertex) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v *Vertex) Scale(f float64) {
	v.X = v.X * f
	v.Y = v.Y * f
}

// Vertex 是 Abser，但 *Vertex 也是 Abser ?
type Abser interface {
	Abs() float64
}

// *Vertex 是 Scaler，但 Vertex 不是 Scaler
type Scaler interface {
	Scale(float64)
}

type Square struct {
	Edge float64
}

func (sq *Square) Scale(f float64) {
	sq.Edge *= f
}

func (sq Square) Area() float64 {
	return sq.Edge * sq.Edge
}

// 静态接口检查，让编译器判断一个结构体是否实现了某个接口
var _ Scaler = (*Square)(nil)

// 任何类型都属于 interface{}
func describe(i interface{}) {
	fmt.Printf("%#v\n", i)
}

// 类型选择样例
// type any = interface{}
func twice(i any) {
	// 可以省略 v :=
	// 如果使用赋值语句，则 v 会是对应 case 中的类型，否则 v 是 i 原本的类型
	switch v := i.(type) {
	case int:
		fmt.Printf("i is int, 2 * i is %v\n", 2*v)
	case string:
		fmt.Printf("i is string, 2 * i is %v\n", v+v)
	default:
		fmt.Printf("IDK the type of i :(\nv is %v\n", v)
	}
}

// 实现 fmt.Stringer 接口的类型可以被 %v 格式化输出
// 注意 IPAddr 和 [4]byte 之间没有等号，否则 IPAddr 相当于别名，而不是新的类型，无法添加方法
type IPAddr [4]byte

func (ipAddr IPAddr) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", ipAddr[0], ipAddr[1], ipAddr[2], ipAddr[3])
}

// 实现 error 接口的类型可以被 %v 格式化输出
type MyError struct {
	When time.Time
	What string
}

func (e *MyError) Error() string {
	return fmt.Sprintf("at %v, %s", e.When, e.What)
}

func somethingMayCauseError() error {
	return &MyError{time.Now(), "Ooops!"}
}

type ErrNegativeSqrt float64

func (err ErrNegativeSqrt) Error() string {
	return fmt.Sprintf("cannot Sqrt negative number: %v", float64(err))
}

func sqrtWithErr(x float64) (float64, error) {
	if x < 0 {
		return 0, ErrNegativeSqrt(x)
	}
	return sqrt(x), nil
}

// io.Reader 的使用和实现
type rot13Reader struct {
	r io.Reader
}

func (rtr *rot13Reader) Read(b []byte) (int, error) {
	n, err := rtr.r.Read(b)
	if err == nil {
		for i := 0; i < n; i++ {
			c := b[i]
			if 'A' <= c && c <= 'Z' {
				b[i] = ((c - 'A' + 13) % 26) + 'A'
			} else if 'a' <= c && c <= 'z' {
				b[i] = ((c - 'a' + 13) % 26) + 'a'
			}
		}
	}
	return n, err
}

func Run6() {
	fmt.Println("--------- Run 6 ---------")
	v := Vertex{1.0, 2.0}
	p := &v
	fmt.Printf("v.Abs(): %v\n", v.Abs())
	// 相当于 (&v).Scale(3.0)
	v.Scale(3.0)
	fmt.Printf("v: %v\n", v)
	// 相当于 (*p).Abs()
	fmt.Printf("p.Abs(): %v\n", p.Abs())
	p.Scale(2.0)
	fmt.Printf("p: %v\n", p)
	// 接口多态性，scaler 底层相当于存储了一个（类型，值）的二元组
	var scaler Scaler = &Vertex{114, 514}
	scaler.Scale(1.5)
	fmt.Printf("scaler: %T, %v\n", scaler, scaler)
	// 底层二元组的值可以是 nil
	var emptyVertex Vertex
	var nilVertex *Vertex
	var abser Abser = emptyVertex
	fmt.Printf("abser: %T, %v\n", abser, abser)
	abser = nilVertex
	fmt.Printf("abser: %T, %v\n", abser, abser)
	// 接口本身可以是 nil，此时不可调用接口中定义的方法，会触发空指针 panic
	var nilAbser Abser
	fmt.Printf("nilAbser: %T, %v\n", nilAbser, nilAbser)
	describe(&Vertex{-1, 2})
	var scaler2 Scaler = &Square{3}
	// 类型断言，若 scaler2 不是 *Square 则会触发 panic
	// 如果括号内的类型不是 Scaler 则无法通过编译
	// 类型断言生成的对象是对被断言者的拷贝，不是引用
	square := scaler2.(*Square)
	fmt.Printf("square.Edge: %v\n", square.Edge)
	// 类型断言形式2，ok 表示断言成功，若断言失败会使 vertex 为 nil
	if vertex, ok := scaler2.(*Vertex); ok {
		fmt.Printf("vertex: %v\n", vertex)
	} else {
		fmt.Println("vertex is nil")
	}
	// 类型选择
	twice(21)
	twice("boing")
	// fmt.Stringer 接口
	hosts := map[string]IPAddr{
		"loopback":  {127, 0, 0, 1},
		"googleDNS": {8, 8, 8, 8},
	}
	for name, ip := range hosts {
		fmt.Printf("%v: %v\n", name, ip)
	}
	// error 接口
	if err := somethingMayCauseError(); err != nil {
		fmt.Println(err)
	}
	f, err := sqrtWithErr(-2)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(f)
	}
	// io.Reader 接口
	r := strings.NewReader("Hello, Reader!")
	b := make([]byte, 8)
	// 每次 Read 从前往后尝试写满 b。
	// n 表示本次写的字节数，若写不满，则后面 8 - n 个字节不会被改变。
	// 当 n == 0 时，err 将为 io.EOF
	for {
		n, err := r.Read(b)
		fmt.Printf("n = %v err = %v b = %v\n", n, err, b)
		fmt.Printf("b[:n] = %q\n", b[:n])
		if err == io.EOF {
			break
		}
	}
	// 测试自行实现的 io.Reader
	s := strings.NewReader("Lbh penpxrq gur pbqr!\n")
	rtr := rot13Reader{s}
	io.Copy(os.Stdout, &rtr)
}
