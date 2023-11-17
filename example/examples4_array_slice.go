package example

import (
	"fmt"
	"strings"
)

func arrayAndSlice() {
	// 数组的定义 1
	var a [2]string
	// 数组元素的赋值
	a[1] = "World"
	// 数组元素的引用
	fmt.Printf("a[0]: %v\n", a[0])
	fmt.Printf("a[1]: %v\n", a[1])

	// 数组的定义 2
	primes := [6]int{2, 3, 5, 7, 11, 13}

	// 数组的定义 3
	arr := [...]string{"aaa", "bbb", "ccc"}
	fmt.Printf("arr: %T, %v\n", arr, arr)

	// 切片的定义 1（左闭右开）
	var primes1to4 []int = primes[1:4]
	fmt.Printf("primes1to4: %T, %v\n", primes1to4, primes1to4)

	// 切片是底层数组的引用
	primes1to4[0] = 17
	fmt.Printf("primes: %T, %v\n", primes, primes)

	// 切片的定义 2（相当于声明了一个隐含数组的 [:] 切片）
	boolValues := []bool{true, false, true, true, false, true}
	fmt.Printf("boolValues: %v\n", boolValues)

	// 切片的定义 3（花括号可嵌套）
	structValues := []struct {
		i int
		b bool
	}{
		{2, true},
		{3, false},
		{5, true},
		{7, true},
		{11, false},
		{13, true},
	}
	fmt.Printf("structValues: %v\n", structValues)

	// 切片的定义 4（使用 make）
	ms1 := make([]int, 5)    // 相当于引用一个长度为 5 的数组，且 ms1 的 len 和 cap 均为 5
	ms2 := make([]int, 3, 5) // 相当于引用一个长度为 5 的数组，且 ms1 的 len 为 3，cap 为 5，引用底层数组的前 3 个元素
	fmt.Printf("make([]int, 5): %v\n", ms1)
	fmt.Printf("make([]int, 3, 5): %v\n", ms2)
}

// 即便 s 是 nil 也能正常工作
func printLenAndCap(s []int) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}

// 切片的长度就是它所包含的元素个数。
// 切片的容量是从它的第一个元素开始数，到其底层数组元素末尾的个数。(注意不是底层数组的长度)
func sliceLenAndCap() {
	s := []int{2, 3, 5, 7, 11, 13}
	printLenAndCap(s)
	// 截取切片使其长度为 0
	s = s[:0]
	printLenAndCap(s)
	// 拓展其长度，注意底层数组始终只有一个，且不能扩展到超出底层数组的范围，否则会有 panic: runtime error
	s = s[:4]
	printLenAndCap(s)
	// 舍弃前两个值
	s = s[2:]
	printLenAndCap(s)
}

// 可以声明长度为 0 的数组，以及长度为 0 的切片。
// 注意区分 nil 和长度为 0 的切片。
func sliceZeroValue() {
	var a [0]int
	fmt.Printf("var a [0]int, a: %T, %v\n", a, a)
	s := a[:]
	fmt.Printf("a[:]: len=%d cap=%d %v\n", len(s), cap(s), s)
	var s2 []int
	if nil == s2 {
		fmt.Println("var s2 []int, s2 is nil")
	}
}

// 切片的元素也可以是切片
func slice2d() {
	board := [][]string{
		{"_", "_", "_"},
		{"_", "_", "_"},
		{"_", "_", "_"},
	}

	// 两个玩家轮流打上 X 和 O
	board[0][0] = "X"
	board[2][2] = "O"
	board[1][2] = "X"
	board[1][0] = "O"
	board[0][2] = "X"

	for i := 0; i < len(board); i++ {
		fmt.Printf("%s\n", strings.Join(board[i], " "))
	}
}

func sliceAppend() {
	var s []int
	printLenAndCap(s)

	// 添加一个空切片
	s = append(s, 0)
	printLenAndCap(s)

	// 这个切片会按需增长
	s = append(s, 1)
	printLenAndCap(s)

	// 可以一次性添加多个元素
	s = append(s, 2, 3, 4)
	printLenAndCap(s)
}

func arrayAndSliceForRange() {
	var pow1 = [8]int{1, 2, 4, 8, 16, 32, 64, 128}
	// 通常形式的 for range
	for i, v := range pow1 {
		fmt.Printf("2**%d = %d\n", i, v)
	}
	pow2 := make([]int, 10)
	// 可以省略第二个值，也可以写成 i, _
	for i := range pow2 {
		pow2[i] = 1 << uint(i) // == 2**i
	}
	// 省略第一个值需要显式地使用 _
	for _, value := range pow2 {
		fmt.Printf("%d\n", value)
	}
}

// 二维切片的分配、引用和 for range
func iota2d(length, width int) (s [][]int) {
	s = make([][]int, length)
	prev := 0
	// 注意固定行长的二维切片的分配技巧，不适用于 s 中的某一行会 append 的情况
	_s := make([]int, length*width)
	for i := range s {
		s[i], _s = _s[:width], _s[width:]
		s[i][0] = prev
		prev++
		for j := range s[i][:len(s[i])-1] {
			s[i][j+1] = s[i][j] + 1
		}
	}
	return
}

// break 会终止 switch 的执行，若 switch 在 for 循环内，且需要在 switch 内终止 for 循环，需要在 for 循环前面加标签。
// 加标签亦可用于跳出多层循环。
// continue 也有类似的跳转用法，但只能用于 for。
func findValInSlice(slice []int, val int) {
	found := false
loop:
	for _, elem := range slice {
		switch elem {
		case val:
			found = true
			break loop
		}
	}
	if found {
		fmt.Printf("%v is found\n", val)
	}
}

// 其实该函数可以用直接返回 s + ";" 的方式实现。
// 主要是为了体现 []byte 和 string 间的相互转换和 append 的用法。
// []rune 也可实现类似转换，但会把 string 中的每一个字符都转换成 32 位
// 如果用 []byte 转换含非 ASCII 字符的 string 会导致字符被拆开
func addSemiColon(s string) string {
	arr := []byte(s)
	arr = append(arr, ';')
	return string(arr)
}

// 函数参数的最后一项可以加 ...，但必须指定类型，函数内部会将其解释为 slice
// 如果 nums 没有实参则会解释为 nil
func addInts(addf func(int, int) int, nums ...int) int {
	fmt.Println(nums)
	if nums == nil {
		return 0
	}
	ans := nums[0]
	for i := 1; i < len(nums); i++ {
		ans = addf(ans, nums[i])
	}
	return ans
}

func Min(a ...int) int {
	min := int(^uint(0) >> 1) // 最大的 int
	for _, i := range a {
		if i < min {
			min = i
		}
	}
	return min
}

func printSlice(name string, s []int) {
	fmt.Printf("%v: %v, is nil: %v, length: %v, capacity: %v\n", name, s, s == nil, len(s), cap(s))
}

func zeroValues() {
	var len0ArrayZV [0]int
	var len3ArrayZV [3]int
	fmt.Printf("len0ArrayZV: %v, length: %v\n", len0ArrayZV, len(len0ArrayZV))
	fmt.Printf("len3ArrayZV: %v, length: %v\n", len3ArrayZV, len(len3ArrayZV))
	var sliceZV []int
	var len0slice = make([]int, 0)
	var len0ArrayZVCP = len0ArrayZV[:]
	var len3ArrayZVCP = len3ArrayZV[0:0]
	var sliceZVCP = sliceZV[:0]
	printSlice("sliceZV", sliceZV)
	printSlice("len0slice", len0slice)
	printSlice("len0ArrayZVCP", len0ArrayZVCP)
	printSlice("len3ArrayZVCP", len3ArrayZVCP)
	printSlice("sliceZVCP", sliceZVCP)
}

func Run4() {
	fmt.Println("--------- Run 4 ---------")
	arrayAndSlice()
	sliceLenAndCap()
	sliceZeroValue()
	slice2d()
	sliceAppend()
	arrayAndSliceForRange()
	s2d := iota2d(5, 3)
	fmt.Printf("iota2d(5): %v\n", s2d)
	findValInSlice(s2d[1], 2)
	fmt.Printf("addSemiColon(\"sample\"): %v\n", addSemiColon("sample"))
	fmt.Printf("addInts(func(a int, b int) int { return a + b }, 1, 2, 3): %v\n", addInts(func(a int, b int) int { return a + b }, 1, 2, 3))
	fmt.Printf("addInts(func(i1, i2 int) int { return i1 * i2 }): %v\n", addInts(func(i1, i2 int) int { return i1 * i2 }))
	nums := []int{3, 5, 9, 1}
	// slice 解包
	fmt.Printf("Min(nums...): %v\n", Min(nums...))
	zeroValues()
}
