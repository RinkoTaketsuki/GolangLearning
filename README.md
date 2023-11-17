# Golang Learning

好记性不如烂笔头。

## 常用命令

`go run main.go` 可以直接编译并执行 .go 文件

直接执行 `go install` 会编译当前目录下有 `package main` 的文件并将可执行文件放入 `$GOPATH/bin/`

`go install` 后面也可以加 module 名（`$GOPATH/src/` 下的路径名），当前目录需要是 go.mod 所在目录或其子目录

```sh
go install github.com/RinkoTaketsuki/GolangLearning
```

`go install` 后面还可以加 `package main` 所在的 `$GOPATH/src/` 下的相对路径名，当前目录需要是 go.mod 所在目录或其子目录

```sh
go install github.com/RinkoTaketsuki/GolangLearning/grpc_example/greeter_server
go install github.com/RinkoTaketsuki/GolangLearning/grpc_example/greeter_client
```

以上命令生成的可执行文件名取决于所在目录，在 Windows 下如 `GolangLearning.exe`，`greeter_server.exe`，`greeter_client.exe`

`go install` 时如果当前程序 import 了其他 module，且被导入的 module 没有 main，则会自动执行 `go install` 被导入的 module。这次 install 会生成静态库文件，存储在 `$GOPATH/pkg/` 下

`go build` 语法与 `go install` 类似，只不过其仅仅测试能否通过编译，而不生成任何二进制文件

`go test` 语法也与 `go install` 类似，用于执行单元测试

`go get` 后面可以接带域名的 module 名，这样 go 会自动下载源码到 `$GOPATH/src/`，下载静态库到 `$GOPATH/pkg/`。如果之前已经下载过则行为与 `go install` 相同

`go install` 时如果在自己的代码里 import 远程 module 也会做类似的下载行为

`go help xxx` 可以获取 `go xxx` 命令的帮助

`gofmt` 或 `go fmt` 可以格式化代码，这使得所有 Golang 代码风格统一。一个小细节是格式化各种运算符时会用空格体现运算优先级，比如：`x<<8 + y<<16`

`godoc` 根据注释生成文档，且按照注释的原样显示，不做 html 解析等处理

### 注释文档生成规范

建议参照 `$GOROOT/src/fmt`，文档注释都写在 package、func 等上方且注释和声明间没有空行，package 文档使用块注释，一个 package 多个文件时只需在其中一个文件中写 package 文档。

每一个可导出的名称（大写字母开头）都应该写文档注释，成组声明（如 `var(...)`）只需写一个文档注释

`$GOROOT` 是 go 的安装目录，其也有像 `$GOPATH` 的结构，只不过其存放的是标准库

## 命名规范

通常请使用 **驼峰命名**

但 package 名最好是单个单词，需要多个单词时请拆为多级目录，如 `$GOROOT/src/pkg/encoding/base64`，导入时要 `import "encoding/base64"`，其包名默认为 `base64`

Golang 的结构体没有构造函数，但通常在构造 `ring.Ring` 时，会使用 `ring.NewRing()` 作为构造函数，如果 `ring` 里面只有某一个结构体和其方法，则通常将 `ring.New()` 作为构造函数

假定某个结构体有 `member` 成员，则其 Getter 和 Setter 的通常写法为 `obj.Member()` 和 `obj.SetMember(value)`

序列化通常用 `obj.String()` 而不是 `obj.ToString()`，类似的如 `Read()`，`Write()` 等标准库中常用的意义鲜明的命名也最好效仿使用

只包含一个方法的接口通常其接口名为方法名 + `er`，如 `Reader`，`Writter`

常用的异常处理写法示例：

```go
f, err := os.Open(name)
if err != nil {
    return err
}
d, err := f.Stat()
if err != nil {
    f.Close()
    return err
}
codeUsing(f, d)
```

## 作用域与变量重新声明

在满足下列条件时，已被声明的变量 v 可出现在 := 声明中：

- 本次声明与已声明的 v 处于同一作用域中（若 v 已在外层作用域中声明过，则此次声明会创建一个新的变量）

- 在初始化中与其类型相应的值才能赋予 v，且在此次声明中至少另有一个变量是新声明的，如上面例子中的 `d, err := f.Stat()`，函数返回的第二个值必须是 error 类型

```go
var (
    v int
    e error
)
v2, e2 := v, e
v3, e2 := v, e
v4, e2 := v, e
v2, e3 := v, e
v2, v3, v4, e4, e5 := v, v, v, e, e
{
    v2, e2 := v, e
    fmt.Printf("%p %p %p %p %p %p %p\n", &v2, &v3, &v4, &e2, &e3, &e4, &e5)
}
fmt.Printf("%p %p %p %p %p %p %p\n", &v2, &v3, &v4, &e2, &e3, &e4, &e5)
// 0xc0000a6090 0xc0000a6080 0xc0000a6088 0xc0000882c0 0xc000088290 0xc0000882a0 0xc0000882b0
// 0xc0000a6068 0xc0000a6080 0xc0000a6088 0xc000088280 0xc000088290 0xc0000882a0 0xc0000882b0
```

这种特性主要用于大量的 err 判断

单个下划线（空白标识符）不算变量，无须考虑上面的重新声明条件，可以随处多次使用

函数的参数和返回值列表、if 或 for 语句中声明的变量等和其后的大括号属于同一作用域（类似 C 语言）

## 逗号的特殊处理

Golang 没有“逗号运算符”的概念，类似 `a,b,c = x,y,z` 的语句，前后变量数必须一样，且必须是有值的表达式（不能是 `v=1` 或 `v++` 这种表达式）

```go
// 反转 a
for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
    a[i], a[j] = a[j], a[i]
}
```

## 内存分配

`new(类名)` 可认为是 `&类名{}` 的语法糖

`make(类名)` 只可用于 slice、map、和 channel，如果尝试 `new([]int)` 则其返回值是指向 nil slice 的指针（指针本身不是 nil）

对数组的赋值和传参会导致所有数组元素的拷贝，若需要 C 语言那样的数组行为必须使用指针

## 打印

`fmt.Fprint` 系列函数的第一个形参是 `io.Writer` 接口，使用 `os.Stdout` 即可模拟 `fmt.Print` 系列函数

`fmt.Print(obj)` 相当于 `fmt.Printf("%V", obj)`

### 打印格式 %v

```go
type T struct {
    a int
    b float64
    c string
}
t := &T{ 7, -2.35, "abc\tdef" }
var timeZone = map[string]int{
    "UTC":  0*60*60,
    "EST": -5*60*60,
    "CST": -6*60*60,
    "MST": -7*60*60,
    "PST": -8*60*60,
}

fmt.Printf("%v\n", t)           // &{7 -2.35 abc   def}
fmt.Printf("%+v\n", t)          // &{a:7 b:-2.35 c:abc     def}
fmt.Printf("%#v\n", t)          // &main.T{a:7, b:-2.35, c:"abc\tdef"}
fmt.Printf("%v\n", timeZone)    // map[CST:-21600 EST:-18000 MST:-25200 PST:-28800 UTC:0]
fmt.Printf("%+v\n", timeZone)   // map[CST:-21600 EST:-18000 MST:-25200 PST:-28800 UTC:0]
fmt.Printf("%#v\n", timeZone)   // map[string]int{"CST":-21600, "EST":-18000, "MST":-25200, "PST":-28800, "UTC":0}
```

### 打印格式 %q %x

```go
    s := "'hello' `Alice` and \"Bob\""
    b := []byte("'hello' Alice and \"Bob\"")
    fmt.Printf("%q\n", s)   // "'hello' `Alice` and \"Bob\""
    fmt.Printf("%#q\n", s)  // "'hello' `Alice` and \"Bob\""
    fmt.Printf("%x\n", s)   // 2768656c6c6f272060416c6963656020616e642022426f6222
    fmt.Printf("% x\n", s)  // 27 68 65 6c 6c 6f 27 20 60 41 6c 69 63 65 60 20 61 6e 64 20 22 42 6f 62 22
    fmt.Printf("%q\n", b)   // "'hello' Alice and \"Bob\""
    fmt.Printf("%#q\n", b)  // `'hello' Alice and "Bob"`
    fmt.Printf("%x\n", b)   // 2768656c6c6f2720416c69636520616e642022426f6222
    fmt.Printf("% x\n", b)  // 27 68 65 6c 6c 6f 27 20 41 6c 69 63 65 20 61 6e 64 20 22 42 6f 62 22
```

### 实现 fmt.Stringer 时导致无限递归的小细节

```go
type MyString string

func (m *MyString) String() string {
    // return fmt.Sprintf("MyString=%s", m) // 错误：会无限递归
    return fmt.Sprintf("MyString=%s", string(*m)) // 可以：注意转换
}

func main() {
    var s MyString = "hello"
    fmt.Println(&s)
}
```

## iota

```go
type ByteSize float64

const (
    _           = iota // 通过赋予空白标识符来忽略第一个值
    KB ByteSize = 1 << (10 * iota)
    MB
    GB
    TB
    PB
    EB
    ZB
    YB
)
```

## init 执行顺序

```go
// 先执行这个
var (
    home   = os.Getenv("HOME")
    user   = os.Getenv("USER")
    gopath = os.Getenv("GOPATH")
)

// 再执行这个
func init() {
    if user == "" {
        log.Fatal("$USER not set")
    }
    if home == "" {
        home = "/home/" + user
    }
    if gopath == "" {
        gopath = home + "/go"
    }
}

// 最后执行这个
func init() {
    // gopath 可通过命令行中的 --gopath 标记覆盖掉。
    flag.StringVar(&gopath, "gopath", gopath, "override default GOPATH")
}
```

## 值方法和指针方法的行为

```go
type MyInt int

func (m MyInt) f1() {
    fmt.Println(&m)
}

func (m *MyInt) f2() {
    fmt.Println(m)
}

func main() {
    var mi MyInt
    var pmi *MyInt
    fmt.Println(&mi)
    mi.f1() // 编译器会拷贝 mi 到 f1 的 m 参数
    mi.f2() // 编译器自动转换为 (&mi).f2()
    // pmi.f1() // 不符合语法
    pmi.f2() // 编译器同样是拷贝传值，只不过拷贝的是 *MyInt 类型
}

// 0xc00009a000
// 0xc00009a008
// 0xc00009a000
// <nil>
```

## 接口

### 一些实用接口

实现 `io.Writer` 可以作为输出流，可以作为 `fmt.Fprint()` 系列函数的输出

实现 `io.Reader` 可以作为输入流

实现 `fmt.Stringer` 可以自定义转换为字符串时的输出格式

实现 `sort.Interface` 可以对内部元素进行排序，如 `sort.IntSlice`

### any 转 string 的实用方法

```go
var value interface{} // Value 由调用者提供
switch str := value.(type) {
case string:
    return str
case fmt.Stringer:
    return str.String()
}
```

### 接口好实践

如果某接口和实现它的结构体之间有明显的单继承父子类关系，则该结构体的构造函数应该返回父类型

```go
type Hasher interface {
    hash(v interface{}) uint64
}

type Algo1Hasher struct{
    // members
} // 假定实现了 Hasher
type Algo2Hasher struct{
    // members
} // 假定实现了 Hasher

func NewAlgo1Hasher() Hasher {
    // func body
}

func NewAlgo2Hasher() Hasher {
    // func body
}

func f() {
    // 如果用户需要改用其他 Hasher，则只需要在此修改调用的构造函数
    hasher := NewAlgo1Hasher()
    // use hasher
}
```

### 函数也能实现方法

```go
func printArgs(w http.ResponseWriter, req *http.Request) {
    fmt.Println(os.Args)
}

// HandlerFunc 类型实现了 http.Handler 方法，其定义为 type HandlerFunc func(ResponseWriter, *Request)
func serve() {
    // server defination ...
    http.Handle("/args", http.HandlerFunc(printArgs))
}
```

## 内嵌的注意事项

`interface` 之间内嵌需要保证没有同名函数，否则编译失败

`struct` 之间内嵌需要保证同层之间没有同名对象，否则编译失败，不同层之间的同名对象会导致上层覆盖下层

`struct` 内嵌结构体的方法可以直接通过 `obj.Method()` 调用，也可以通过 `obj.EmbeddedStruct.Method()` 调用

## 并发

### 同步 channel 用于等待 goroutine 结束

```go
func main() {
    c := make(chan int)
    go func() {
        // do someting
        c <- 1
    }()
    // do something
    <-c
}
```

### 带缓冲区 channel 用作信号量

```go
var MaxOutstanding int = runtime.NumCPU()

var sem = make(chan int, MaxOutstanding)

type Request struct {
    // defination
}

func Serve(queue chan *Request) {
    for req := range queue {
        sem <- 1
        go func(req *Request) {
            // process req
            <-sem
        }(req)
    }
}

func Serve2(queue chan *Request) {
    for req := range queue {
        req := req // 为该Go协程创建 req 的新实例。省略这一行会导致多个 goroutine 共享同一个 req 值。
        sem <- 1
        go func() {
            // process req
            <-sem
        }()
    }
}
```

### channel 嵌套样例

```go
// runtime.GOMAXPROCS，设置当前最大可用的 CPU 数量，返回的是之前设置的最大可用的 CPU 数量。
// 默认情况下使用 runtime.NumCPU 的值，但是可以被命令行环境变量或者调用此函数并传参正整数修改。传参 0 的话会返回值。
const MaxOutstanding = runtime.GOMAXPROCS(0)

type Request struct {
    args        []int
    f           func([]int) int
    resultChan  chan int
}

var clientRequests chan = make(chan *Request, 128)

func handle(queue chan *Request) {
    for req := range queue {
        req.resultChan <- req.f(req.args)
    }
}

func Serve(clientRequests chan *Request, quit chan bool) {
    // 启动处理程序
    for i := 0; i < MaxOutstanding; i++ {
        go handle(clientRequests)
    }
    <-quit  // 等待通知退出。
}

func sum(a []int) (s int) {
    for _, v := range a {
        s += v
    }
    return
}

func main() {
    request := &Request{[]int{3, 4, 5}, sum, make(chan int)}
    // 发送请求
    clientRequests <- request
    // 等待回应
    fmt.Printf("answer: %d\n", <-request.resultChan)
}
```

### 漏桶缓冲区设计样例

```go
var freeList = make(chan *Buffer, 100)
var serverChan = make(chan *Buffer)

func client() {
    for {
        var b *Buffer
        // 若缓冲区可用就用它，不可用就分配个新的。
        select {
        case b = <-freeList:
            // 获取一个，不做别的。
        default:
            // 非空闲，因此分配一个新的。
            b = new(Buffer)
        }
        load(b)              // 从网络中读取下一条消息。
        serverChan <- b      // 发送至服务器。
    }
}

func server() {
    for {
        b := <-serverChan    // 等待工作。
        process(b)
        // 若缓冲区有空间就重用它。
        select {
        case freeList <- b:
            // 将缓冲区放到空闲列表中，不做别的。
        default:
            // 空闲列表已满，保持就好。
        }
    }
}
```

## 错误

`error` 是 `interface` 不是 `struct`，用户可以自行实现各种各样的 `error`，比如 `os.PathError`

```go
// PathError 记录错误、执行的操作和文件路径
type PathError struct {
    Op string    // "open", "unlink" 等等对文件的操作
    Path string  // 相关文件的路径
    Err error    // 由系统调用返回
}

func (e *PathError) Error() string {
    return e.Op + " " + e.Path + ": " + e.Err.Error()
}
```

错误字符串应尽可能地指明它们的来源以便调试，因为输出错误的位置可能离错误产生的位置非常遥远（C++ 并感）

如果需要像其他语言一样检查指定类型的错误，可以使用类型断言

`panic` 会逐调用堆栈运行 `defer` 栈，然后终止 goroutine 运行

`recover` 会返回 `panic` 传入的内容并终止 `panic` 的回溯。因为 `panic` 后仅会执行 `defer` 栈中的内容，故 `recover` 的调用必须在 `defer` 的函数内部。

### panic 和 recover 搭配用作异常处理的样例

```go
// Error 是解析错误的类型，它满足 error 接口。
type Error string
func (e Error) Error() string {
    return string(e)
}

// error 是 *Regexp 的方法，它通过用一个 Error 
// 触发Panic来报告解析错误。
func (regexp *Regexp) error(err string) {
    panic(Error(err))
}

// Compile 返回该正则表达式解析后的表示。
func Compile(str string) (regexp *Regexp, err error) {
    regexp = new(Regexp)
    // 当发生解析错误时，doParse 会触发 panic
    defer func() {
        if e := recover(); e != nil {
            regexp = nil    // 清理返回值。
            err = e.(Error) // 若不是 RegExp 的 error 方法触发的 panic，将重新触发一个新的 panic（即类型断言失败）。
        }
    }()
    return regexp.doParse(str), nil
}
```

## defer 亦能访问外层函数的局部变量和命名的返回值

```go
func f() int {
    fmt.Println("Exec f")
    i := 1
    fmt.Println("Defer a func")
    defer func() {
        fmt.Println("Exec deferred func, now i is", i)
        i++
        fmt.Println("Ret deffered func, now i is", i)
    }()
    fmt.Println("Ret f, now i is", i)
    return i
}

func g() (i int) {
    fmt.Println("Exec g")
    i = 1
    fmt.Println("Defer a func")
    defer func() {
        fmt.Println("Exec deferred func, now i is", i)
        i++
        fmt.Println("Ret deffered func, now i is", i)
    }()
    fmt.Println("Ret g, now i is", i)
    return i
}

func main() {
    fmt.Println("Call f")
    fmt.Println("Finally f returns:", f())
    fmt.Println("Call g")
    fmt.Println("Finally g returns:", g())
}

// Call f
// Exec f
// Defer a func
// Ret f, now i is 1
// Exec deferred func, now i is 1
// Ret deffered func, now i is 2
// Finally f returns: 1
// Call g
// Exec g
// Defer a func
// Ret g, now i is 1
// Exec deferred func, now i is 1
// Ret deffered func, now i is 2
// Finally g returns: 2
```
