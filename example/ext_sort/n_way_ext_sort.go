package ext_sort

import (
	"container/heap"
	"errors"
	"io"
	"os"
	"sort"
)

func ExtMergeSort2way(data []uint64) {
	type (
		DataPtr = uintptr
		BufPtr  = uintptr
	)

	// 内部存储可以存储的数据个数，这里假设最极端的情况，即内部存储只能放置 2
	// 个元素。如果内部存储放不下 2 个元素，则无法完成该外部排序算法。
	const BUF_LEN BufPtr = 2

	dataLen := DataPtr(len(data))

	// buf 模拟内部存储，bak 模拟算法需要的额外的外部存储空间。
	// inBak 表示上一次合并的结果存储在 bak 中。
	// read 和 write 分别表示读取和写入外部存储。
	buf := make([]uint64, BUF_LEN)
	bak := make([]uint64, dataLen)
	inBak := false
	read := func(buf []uint64, dataAddr DataPtr, useBak bool) {
		if useBak {
			copy(buf, bak[dataAddr:])
		} else {
			copy(buf, data[dataAddr:])
		}
	}
	write := func(buf []uint64, dataAddr DataPtr, useBak bool) {
		if useBak {
			copy(bak[dataAddr:], buf)
		} else {
			copy(data[dataAddr:], buf)
		}
	}

	// 排序，按顺序每次从外部存储读取最多 2 个元素，排序后写回外部存储。
	for offset := DataPtr(0); offset+BUF_LEN <= dataLen; offset += BUF_LEN {
		read(buf, offset, false)
		if buf[1] < buf[0] {
			buf[0], buf[1] = buf[1], buf[0]
		}
		write(buf, offset, false)
	}

	// 如果内部存储足以容纳所有外部存储的数据，则上述排序已将所有数据排序完毕。
	if dataLen <= BUF_LEN {
		return
	}

	// 合并，每次选取两个相邻的排好序的分段，合并为一个排好序的分段。
	// 迭代 segmentLen。该变量表示当前循环轮次是将每两个相邻的长度为 segmentLen 的有
	// 序区间合并为一个 2 * segmentLen 的有序区间。
	// 第奇数次迭代中，输入数据为 data，输出数据为 bak。
	// 第偶数次迭代中，输入数据为 bak，输出数据为 data。
	// 举例：第 1 次迭代中，data 为 [1, 4, 2, 3, 2]。
	// 则当前循环轮次结束后，data 不变，bak 为 [1, 2, 3, 4, 2]。
	for segmentLen := BUF_LEN; segmentLen < dataLen; segmentLen <<= 1 {
		// 迭代有序区间对。offset1 为左边有序区间的首元素下标。offset2 为右边有序区间
		// 的首元素下标。值得注意的是 offset2 越界但 offset 1 未越界的情况。此时下标
		// [offset1, dataLen) 这一段的数据一定是有序的，将这段数据直接拷贝到输出数据
		// 的对应区域即可。
		offset1, offset2 := DataPtr(0), segmentLen
		for offset2 < dataLen {
			// [offset1, end1) 为左边的有序区间，[offset2, end2) 为右边的有序区间。
			// 由下面的定义可知这两个区间长度并不相同。
			end1, end2 := offset2, min(dataLen, offset2+segmentLen)
			// nonempty[i] 表示目前 buf[i] 中存储着未输出到输出数据区域的数据。
			nonempty := [BUF_LEN]bool{}
			// [offset1, i1) 为左边有序区间已被读取的数据。
			// [offset2, i2) 为右边有序区间已被读取的数据。
			// [offset1, io) 为已写入输出区域的数据。
			// 循环结束后，i1 == end1，i2 == end2，io == end2。
			for i1, i2, io := offset1, offset2, offset1; io < end2; {
				// 若 buf[0] 为空，则将 i1 中的数据读入
				if !nonempty[0] && i1 < end1 {
					read(buf[0:1], i1, inBak)
					nonempty[0] = true
				}
				// 若 buf[1] 为空，则将 i2 中的数据读入
				if !nonempty[1] && i2 < end2 {
					read(buf[1:2], i2, inBak)
					nonempty[1] = true
				}
				// 若 buf[0] 和 buf[1] 均非空，则比较后将较小者写入输出区域。
				// 若 buf[0] 和 buf[1] 仅有一个非空，直接将这个数据写入输出区域。
				if nonempty[0] && nonempty[1] {
					if buf[0] < buf[1] {
						write(buf[0:1], io, !inBak)
						nonempty[0] = false
						i1++
					} else {
						write(buf[1:2], io, !inBak)
						nonempty[1] = false
						i2++
					}
				} else if nonempty[0] {
					write(buf[0:1], io, !inBak)
					nonempty[0] = false
					i1++
				} else { // nonempty[1]
					write(buf[1:2], io, !inBak)
					nonempty[1] = false
					i2++
				}
				io++
			}
			// 处理下一组有序区间对。
			offset1, offset2 = end2, end2+segmentLen
		}
		// 处理上述的 offset2 越界但 offset 1 未越界的情况。
		if offset1 < dataLen {
			if inBak {
				for ; offset1 < dataLen; offset1 += BUF_LEN {
					copy(buf[:], bak[offset1:])
					copy(data[offset1:], buf[:])
				}
			} else {
				for ; offset1 < dataLen; offset1 += BUF_LEN {
					copy(buf[:], data[offset1:])
					copy(bak[offset1:], buf[:])
				}
			}
		}
		// 一轮合并完毕，交换 data 和 bak
		inBak = !inBak
	}
	if inBak {
		copy(data, bak)
	}
}

type Elem2Idx struct {
	e []byte
	i int
}

type MyHeap struct {
	s        []Elem2Idx
	lt       func([]byte, []byte) bool
	buf      []byte
	elemSize int
}

func (h *MyHeap) Len() int {
	return len(h.s)
}

func (h *MyHeap) Less(i, j int) bool {
	return h.lt(h.s[i].e, h.s[j].e)
}

func (h *MyHeap) Swap(i, j int) {
	h.s[i], h.s[j] = h.s[j], h.s[i]
}

func (h *MyHeap) Push(x any) {
	l := len(h.s)
	if l == cap(h.s) {
		panic("try to push an element into a full heap")
	}
	h.s = h.s[:l+1]
	h.s[l] = x.(Elem2Idx)
}

func (h *MyHeap) Pop() any {
	l := len(h.s) - 1
	if l < 0 {
		panic("try to pop an element from an empty heap")
	}
	e2i := h.s[l]
	h.s = h.s[:l]
	return e2i
}

type MySlice struct {
	s  [][]byte
	lt func(i, j int) bool
}

func (s *MySlice) Len() int {
	return len(s.s)
}

func (s *MySlice) Less(i, j int) bool {
	return s.lt(i, j)
}

func (s *MySlice) Swap(i, j int) {
	if i == j {
		return
	}
	si, sj := s.s[i], s.s[j]
	l := len(si)
	for k := range l {
		si[k] ^= sj[k]
		sj[k] ^= si[k]
		si[k] ^= sj[k]
	}
}

// Elem should be with a fixed size in ext-memory.
// Parser read an Elem from data. If data reach the end (like EOF), the parser
// should do nothing to the output and return (true, nil).
// The err return param means exceptions.
func ExtMergeSortNWay(
	data io.ReadWriteSeeker,
	elemSize int,
	lt func([]byte, []byte) bool,
	n int,
) error {
	if elemSize < 1 || n < 2 {
		return errors.New("wrong parameters")
	}

	// 用于读写数据的内部存储。
	buf := make([]byte, n*elemSize)

	// elems 将 buf 以 elemSize 进行分割，形成二级 slice。
	elems := make([][]byte, n)
	for i := range elems {
		elems[i] = buf[i*elemSize : (i+1)*elemSize]
	}

	// 算法需要的额外的外部存储空间。
	bak, err := os.CreateTemp("", "ext_merge_sort*")
	if err != nil {
		return err
	}

	// 表示上一次合并的结果存储在 bak 中。
	inBak := false

	// 元素总数。
	totalElemNum := 0

	// 第一次排序，按顺序每次从外部存储读取最多 n 个元素，排序后写回外部存储。同时统计
	// 元素总数。若 data 中的总字节数不能被 elemSize 整除，则会忽略多余的字节。
	data.Seek(0, io.SeekStart)
	for {
		byteNum, err := data.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		elemNum := byteNum / elemSize
		totalElemNum += elemNum
		sort.Sort(&MySlice{elems[:elemNum], func(i, j int) bool { return lt(elems[i], elems[j]) }})
		data.Seek(-int64(byteNum), io.SeekCurrent)
		_, err = data.Write(buf[:byteNum])
		if err != nil {
			return err
		}
	}

	// 如果内部存储足以容纳所有外部存储的数据，则上述排序已将所有数据排序完毕。
	if totalElemNum <= n {
		return nil
	}

	// 合并，每次选取 n 个相邻的排好序的分段，合并为一个排好序的分段。
	// 迭代 segmentLen。该变量表示当前循环轮次是将每 n 个相邻的长度为 segmentLen 的有
	// 序区间合并为一个 n * segmentLen 的有序区间。
	// 第奇数次迭代中，输入数据为 data，输出数据为 bak。
	// 第偶数次迭代中，输入数据为 bak，输出数据为 data。

	idxs := make([]int, n) // 输入 segments 的输入下标。
	ends := make([]int, n) // 输入 segments 的结束下标。

	myHeap := &MyHeap{
		s:        make([]Elem2Idx, 0, n),
		lt:       lt,
		buf:      buf,
		elemSize: elemSize,
	}

	// 从 stream 的 offset 位置读取字节到 elem
	read := func(elem []byte, offset int, fromBak bool) error {
		var stream io.ReadWriteSeeker
		if fromBak {
			stream = bak
		} else {
			stream = data
		}
		stream.Seek(int64(offset), io.SeekStart)
		_, err := stream.Read(elem)
		return err
	}

	// 向 stream 的 offset 位置写入 elem
	write := func(elem []byte, offset int, toBak bool) error {
		var stream io.ReadWriteSeeker
		if toBak {
			stream = bak
		} else {
			stream = data
		}
		stream.Seek(int64(offset), io.SeekStart)
		_, err := stream.Write(elem)
		return err
	}

	for segmentLen := n; segmentLen < totalElemNum; segmentLen *= n {
		// 第一组 n * segmentLen 的 idxs
		idxs[0] = 0
		for i := 1; i < len(idxs); i++ {
			idxs[i] = idxs[i-1] + segmentLen
		}
		offset := 0 // 输出 segment 的输出位置。
		// 迭代多个有序区间组。
		for idxs[0] < totalElemNum {
			// 依据上次迭代产生的新的 idxs 初始化 ends 和 offsets
			for i := range ends {
				ends[i] = idxs[i] + segmentLen
			}
			// 初始化 myHeap
			for i, beg := range idxs {
				if beg >= totalElemNum {
					break
				}
				e2i := Elem2Idx{
					e: buf[i*elemSize : (i+1)*elemSize],
					i: i,
				}
				if err := read(e2i.e, idxs[i]*elemSize, inBak); err != nil {
					return err
				}
				idxs[i]++
				heap.Push(myHeap, e2i)
			}
			// 通过读写 myHeap 实现合并
			for myHeap.Len() > 0 {
				e2i := heap.Pop(myHeap).(Elem2Idx)
				i := e2i.i
				if err := write(e2i.e, offset, !inBak); err != nil {
					return err
				}
				offset += elemSize
				if idxs[i] < ends[i] && idxs[i] < totalElemNum {
					if err := read(e2i.e, idxs[i]*elemSize, inBak); err != nil {
						return err
					}
					idxs[i]++
					heap.Push(myHeap, e2i)
				}
			}
			// 处理下一组 n * segmentLen
			for i := range idxs {
				idxs[i] = ends[i] + (n-1)*segmentLen
			}
		}
		// 一轮合并完毕，交换 data 和 bak
		inBak = !inBak
	}

	if inBak {
		bak.Seek(0, io.SeekStart)
		data.Seek(0, io.SeekStart)
		for {
			n, err := bak.Read(buf)
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			_, err = data.Write(buf[:n])
			if err != nil {
				return err
			}
		}
	}

	if err := bak.Close(); err != nil {
		return err
	}
	return os.Remove(bak.Name())
}
