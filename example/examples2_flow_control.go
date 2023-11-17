package example

import (
	"fmt"
	"math"
	"runtime"
	"time"
)

// golang 只有 for 循环，没有 while 循环
// 相比 C 语言，for 循环和 if 条件判断没有圆括号，花括号不可省略，花括号不能另起一行写，因为 golang parser 会自动加分号
// 流程控制 for、if、switch 等后面什么都不写，相当于写 true
func forLoop() {
	sum := 0
	// 最常见的 for 循环格式
	for i := 0; i < 10; i++ {
		sum += i
	}
	fmt.Printf("sum 0 to 10: %v\n", sum)

	sum = 1
	// 相当于 while 循环的格式
	for sum <= 1000 {
		sum += sum
	}
	fmt.Printf("minimum 2^n greater than 1000: %v\n", sum)

	sum = 1
	// 相当于 for true
	for {
		if sum > 10 {
			break
		}
		sum += sum
	}
	fmt.Printf("minimum 2^n greater than 10: %v\n", sum)

	// 对于字符串的特殊处理
	for pos, char := range "日本\x80語" { // \x80 是个非法的UTF-8编码
		fmt.Printf("字符 %#U 始于字节位置 %d\n", char, pos)
	}
	// 字符 U+65E5 '日' 始于字节位置 0
	// 字符 U+672C '本' 始于字节位置 3
	// 字符 U+FFFD '�' 始于字节位置 6
	// 字符 U+8A9E '語' 始于字节位置 7
}

// if 语句的条件判断前面还可以执行一条额外语句
// 下面写 else 纯粹为了示范写法，一般这种 return、break 最好不要写多余的 else
func pow(x, n, lim float64) float64 {
	if v := math.Pow(x, n); v < lim {
		return v
	} else {
		fmt.Printf("%g >= %g\n", v, lim)
	}
	return lim
}

func sqrt(x float64) float64 {
	const EPSILON = 0.000001
	z := 1.0
	delta := math.MaxFloat64
	for math.Abs(delta) > EPSILON {
		delta = (z*z - x) / (2 * z)
		z -= delta
	}
	return z
}

func getWindows() string {
	fmt.Println(" *** f is called *** ")
	return "windows"
}

/* switch 语句和 C 语言的 switch 语句相比有如下不同：
 * 1. 类似 if，没有圆括号，判断对象前面可额外执行一条语句
 * 2. 判断对象可以是任意类型
 * 3. case 后所跟条件可以是运行时才能判断的值，如变量或函数返回值
 * 4. 从上到下判断 case，若命中某 case 则下面的 case 不会被计算（如下面命中 linux 时不会调用 getWindows())
 * 5. case 命中后其他 case 中的语句不会被执行，相当于 C 语言的 switch 全部加 break，若要模仿 C 语言 switch 的行为则需要在末尾添加 fallthrough 关键字
 * 6. switch 可省略判断对象，而转为在 case 中写判断条件，相当于 switch true
 */
func switchCase() {
	fmt.Println("Go runs on")
	linux := "linux"
	switch os := runtime.GOOS; os {
	case "darwin":
		fmt.Println("OS X.")
	case linux:
		fmt.Println("Linux.")
	case getWindows():
		fmt.Println("Windows")
	default:
		fmt.Printf("%s.\n", os)
	}

	t := time.Now()
	switch {
	case t.Hour() < 12:
		fmt.Println("Good morning!")
	case t.Hour() < 17:
		fmt.Println("Good afternoon.")
	default:
		fmt.Println("Good evening.")
	}
}

// switch 并不会自动下溯，但 case 可通过逗号分隔来列举相同的处理条件。
func shouldEscape(c byte) bool {
	switch c {
	case ' ', '?', '&', '=', '#', '+', '%':
		return true
	}
	return false
}

func echo(i int) int {
	fmt.Printf("echo is called, i is %v\n", i)
	return i
}

// defer 栈的执行规律：入栈前先对函数参数求值，出栈时执行函数
func deferStack() {
	fmt.Println("defer begin")
	for i := 0; i < 3; i++ {
		defer fmt.Println(echo(i))
	}
	fmt.Println("defer end")
}

func Run2() {
	fmt.Println("--------- Run 2 ---------")
	forLoop()
	fmt.Printf("sqrt(pow(7, 2, 1024)): %v\n", sqrt(pow(7, 2, 1024)))
	switchCase()
	fmt.Printf("shouldEscape('?'): %v\n", shouldEscape('?'))
	deferStack()
}
