package example // 除了 package main，包名最好和其所在目录的名字相同。

// 函数外部所有语句都要以关键字开头，如下面的 package, import, var, const 和 func

import (
	"fmt"
	// 导入的包可以起别名
	sc "strconv"
	// 如果有只需导包而当前文件无须使用的情况，使用空白标识符
	_ "math"
)

// 声明多个同类型变量，初始化为 0
var zeroInt1, zeroInt2 int

// 使用 var (...) 的方式声明变量，string 默认为 ""，bool 默认为 false，complex 默认为 (0+0i)
var (
	emptyString1   string
	false1, false2 bool
	zeroComplex1   complex128
)

// 字面值的类型为 untyped xxx，其可以被精确计算
// 通过 const 可以给 untyped 类型取别名
// 使用 var v = 字面值初始化时，untyped int 会转为 int，untyped float 会转为 float，untyped complex 会转为 complex128
const bigUntypedInt = 0xfffffffffffffffffffffffffff
const smallUntypedInt = bigUntypedInt & 0xffff
const (
	untypedFloat1  = -777777779999.777777779999
	untypedFloat2  = untypedFloat1 + 1.000000000001
	untypedComplex = 777777779999.999999 + 77777777999.77777777999i
)

// 定义变量时可以设定初始值，float32 和 float64 可以用 untyped int 初始化，但整数不能用 untyped float 初始化
var float1, float2 float32 = 0xf, 3.14
var (
	string1                      = "hello"
	complex1, complex2 complex64 = (1 - 2i), (-3 + 4i)
)

// 类型转换可以使用 类型(变量) 的语法，须注意不支持 untyped int 的截断
var toInt = int32(smallUntypedInt)
var toUint16 = uint16(smallUntypedInt)

// 但支持变量的截断
var toUint8 = uint8(toUint16)

// untyped float 转换为 float 变量会损失精度，不能将 untyped float 转换为整数变量
var toFloat32 = float32(untypedFloat1)
var toFloat64 = float64(untypedFloat2)

// untyped complex 转换为 complex 变量会损失精度
var toComplex64 = complex64(untypedComplex)
var toComplex128 = complex128(untypedComplex)

func printObjs() {
	fmt.Printf("zeroInt1: %v\n", zeroInt1)
	fmt.Printf("zeroInt2: %v\n", zeroInt2)
	fmt.Printf("emptyString1: %v\n", emptyString1)
	fmt.Printf("false1: %v\n", false1)
	fmt.Printf("false2: %v\n", false2)
	fmt.Printf("zeroComplex1: %v\n", zeroComplex1)
	fmt.Printf("smallUntypedInt: Default: %T, %v\n", smallUntypedInt, smallUntypedInt)
	fmt.Printf("untypedFloat1: Default: %T, %v\n", untypedFloat1, untypedFloat1)
	fmt.Printf("untypedFloat2: Default: %T, %v\n", untypedFloat2, untypedFloat2)
	fmt.Printf("untypedComplex: Default: %T, %v\n", untypedComplex, untypedComplex)
	fmt.Printf("float1: %v\n", float1)
	fmt.Printf("float2: %v\n", float2)
	fmt.Printf("string1: %v\n", string1)
	fmt.Printf("complex1: %v\n", complex1)
	fmt.Printf("complex2: %v\n", complex2)
	fmt.Printf("toInt: %v\n", toInt)
	fmt.Printf("toUint16: %v\n", toUint16)
	fmt.Printf("toUint8: %v\n", toUint8)
	fmt.Printf("toFloat32: %v\n", toFloat32)
	fmt.Printf("toFloat64: %v\n", toFloat64)
	fmt.Printf("toComplex64: %v\n", toComplex64)
	fmt.Printf("toComplex128: %v\n", toComplex128)

}

// 函数的参数列表有多个相邻的相同类型的值可以合并
func add(x, y int) int {
	return x + y
}

// 函数的返回值可以指定变量名，返回时可以省略返回值
// 相当于在函数开头定义了 var z float64
func sub(x, y float64) (z float64) {
	z = x - y
	return
}

// 函数返回多个值可以是不同类型，函数内部可以用 := 快速定义变量
func getFloatAndString(x int) (float64, string) {
	f := float64(x)
	s := sc.Itoa(x)
	return f, s
}

// 函数返回多个值时也可以省略返回值，此时必须指定每个返回值的变量名
func subAndAddOne(x int) (a, b int) {
	a = x - 1
	b = x + 1
	return
}

// 此处返回 r1 时相当于 r1 = sub(x, y)
func add4(w, x, y, z float64) (r1, r2 float64) {
	r2 = x + z
	return w + y, r2
}

// 闭包范例，体现了函数也是值
func fibonacci() func() int {
	x := 0
	y := 1
	return func() (r int) {
		r = x
		x = y
		y = r + y
		return
	}
}

// 大写字母开头的各种标识符是导出的，即可以被其他 package 引用
func Run1() {
	fmt.Println("--------- Run 1 ---------")
	printObjs()
	// 函数返回多个值时可以直接解包到其他函数的参数，但此时返回列表类型和参数列表类型必须一一对应
	// 变量间的类型转换必须是显式转换
	// _ 用于接收不需要的返回值
	_, s := getFloatAndString(add(subAndAddOne(int(sub(add4(1, 2, 3, 4))))))
	fmt.Printf("getFloatAndString(add(subAndAddOne(int(sub(add4(1, 2, 3, 4)))))): %v\n", s)
	f := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}
	// string 不可更改，可以混合 byte 和 rune。string 支持 + 运算符
	s = "hello" + "世界!"
	fmt.Printf("s: %v\n", s)
}
