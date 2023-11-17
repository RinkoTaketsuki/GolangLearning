package example

import (
	"fmt"
	"sync"
	"time"
)

func sum(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	c <- sum // 将和送入 c
}

func produce(wg *sync.WaitGroup, ch chan int) {
	if wg != nil {
		defer wg.Done()
	}
	defer close(ch)
	for i := 0; i < 20; i++ {
		fmt.Printf("Put %v into ch\n", i)
		ch <- i
	}
}

func consume(wg *sync.WaitGroup, ch chan int) {
	if wg != nil {
		defer wg.Done()
	}
	for i := 0; i < 20; i++ {
		fmt.Println(<-ch)
	}
}

func fibonacciCh(n int, c chan int) {
	defer close(c)
	x, y := 0, 1
	for i := 0; i < n; i++ {
		c <- x
		x, y = y, x+y
	}
}

// select 随机选择一个未阻塞的 case 执行，若全阻塞则阻塞
func fibonacciSelectCh(c chan int, quit chan bool) {
	x, y := 0, 1
	for {
		select {
		case c <- x:
			x, y = y, x+y
		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}

// SafeCounter 的并发使用是安全的。
type SafeCounter struct {
	v   map[string]int
	mux sync.Mutex
}

// Inc 增加给定 key 的计数器的值。
func (c *SafeCounter) Inc(key string) {
	c.mux.Lock()
	// Lock 之后同一时刻只有一个 goroutine 能访问 c.v
	c.v[key]++
	c.mux.Unlock()
}

// Value 返回给定 key 的计数器的当前值。
func (c *SafeCounter) Value(key string) int {
	c.mux.Lock()
	// Lock 之后同一时刻只有一个 goroutine 能访问 c.v
	defer c.mux.Unlock()
	return c.v[key]
}

func Run7() {
	// t, _ := time.ParseDuration("3s")
	s := []int{7, 2, 8, -9, 4, 0}
	// 无缓冲 channel，发送方和接收方都准备好前二者会一直阻塞
	c := make(chan int)
	go sum(s[:len(s)/2], c)
	go sum(s[len(s)/2:], c)
	x, y := <-c, <-c // 从 c 中接收
	fmt.Println(x, y, x+y)
	// 有缓冲 channel，仅当信道的缓冲区填满后，向其发送数据时才会阻塞。当缓冲区为空时，接受方会阻塞。
	// 使用 sync.WaitGroup 是防止主协程提前结束
	c2 := make(chan int, 5)
	var wg sync.WaitGroup
	wg.Add(2)
	go produce(&wg, c2)
	go consume(&wg, c2)
	wg.Wait()
	// channel 接收的第二个参数。
	// 若未关闭，则正常读取和阻塞，且 ok 总为 true；
	// 若已关闭，则缓冲区有值时读取，且 ok 为 true，无值时读取零值，且 ok 为 false。
	// 不能向已关闭信道发送数据。
	close(c)
	c = make(chan int, 2)
	c <- 42
	c <- 10
	v, ok := <-c
	fmt.Printf("recv: %v, ok: %v\n", v, ok)
	close(c)
	v, ok = <-c
	fmt.Printf("recv: %v, ok: %v\n", v, ok)
	v, ok = <-c
	fmt.Printf("recv: %v, ok: %v\n", v, ok)
	// range channel 会在 channel 关闭时结束，除了这种情况通常不需要手动关闭 channel
	c = make(chan int, 10)
	go fibonacciCh(2*cap(c), c)
	for i := range c {
		fmt.Println(i)
	}
	// select 语句用法
	c = make(chan int)
	quit := make(chan bool)
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(<-c)
		}
		quit <- true
	}()
	fibonacciSelectCh(c, quit)
	// 若 select 的非 default case 均在阻塞，则执行 default
	ticker := time.NewTicker(100 * time.Millisecond)
	boom := time.After(500 * time.Millisecond)
loop:
	for {
		select {
		case <-ticker.C:
			fmt.Println("tick.")
		case <-boom:
			fmt.Println("BOOM!")
			break loop
		default:
			fmt.Println("    .")
			time.Sleep(50 * time.Millisecond)
		}
	}
	ticker.Stop()
	// Mutex 的使用
	sc := SafeCounter{v: make(map[string]int)}
	for i := 0; i < 1000; i++ {
		go sc.Inc("somekey")
	}
	time.Sleep(time.Second * 3)
	fmt.Println(sc.Value("somekey"))
	// 并发爬虫示范
	if !us.testAndInsert("https://golang.org/") {
		go Crawl("https://golang.org/", 4, fetcher)
	}
	time.Sleep(time.Second * 10)
}

type Fetcher interface {
	// Fetch 返回 URL 的 body 内容，并且将在这个页面上找到的 URL 放到一个 slice 中。
	Fetch(url string) (body string, urls []string, err error)
}

type urlSet struct {
	mu sync.Mutex
	s  map[string]bool
}

func (us *urlSet) testAndInsert(url string) bool {
	defer us.mu.Unlock()
	us.mu.Lock()
	if !us.s[url] {
		us.s[url] = true
		return false
	}
	return true
}

var us *urlSet = &urlSet{s: make(map[string]bool)}

// Crawl 使用 fetcher 从某个 URL 开始递归的爬取页面，直到达到最大深度。
func Crawl(url string, depth int, fetcher Fetcher) {
	if depth <= 0 {
		return
	}
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("found: %s %q\n", url, body)
	for _, u := range urls {
		if !us.testAndInsert(u) {
			go Crawl(u, depth-1, fetcher)
		}
	}
}

// fakeFetcher 是返回若干结果的 Fetcher。
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher 是填充后的 fakeFetcher。
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}
