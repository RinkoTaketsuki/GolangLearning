package ext_sort

import (
	"sort"
)

func ExtMergeSort2wayExample(data []uint64) bool {
	type (
		DataPtr = uintptr
		BufPtr  = uintptr
	)

	// 内部存储可以存储的数据个数，这里假设最极端的情况，即内部存储只能放置 2
	// 个元素。如果内部存储放不下 2 个元素，则无法完成该外部排序算法。
	const BUF_LEN BufPtr = 2

	// buf 模拟内部存储，bak 模拟算法需要的额外的外部存储空间，
	// inBak 表示上一次合并的结果存储在 bak 中。
	// read 和 write 分别表示读取和写入外部存储。
	// move 表示把某段存储内容在 data 和 bak 间复制。
	buf := make([]uint64, BUF_LEN)
	bak := make([]uint64, len(data))
	inBak := false
	read := func(bufAddr BufPtr, dataAddr DataPtr, useBak bool) {
		if useBak {
			buf[bufAddr] = bak[dataAddr]
		} else {
			buf[bufAddr] = data[dataAddr]
		}
	}
	write := func(bufAddr BufPtr, dataAddr DataPtr, useBak bool) {
		if useBak {
			bak[dataAddr] = buf[bufAddr]
		} else {
			data[dataAddr] = buf[bufAddr]
		}
	}

	// 排序，按顺序每次从外部存储读取最多 2 个元素，排序后写回外部存储。
	for offset := DataPtr(0); offset <= DataPtr(len(data))-BUF_LEN; offset += BUF_LEN {
		for i := BufPtr(0); i < BUF_LEN; i++ {
			read(i, offset+i, false)
		}
		read(0, offset, false)
		read(1, offset+1, false)
		if buf[1] < buf[0] {
			buf[0], buf[1] = buf[1], buf[0]
		}
		write(0, offset, false)
		write(1, offset+1, false)
	}

	if len(data) <= BUF_LEN {
		return true
	}

	/* 4. 合并，每次选取两个相邻的排好序的分段，合并为一个排好序的分段。*/
	// 迭代子分段长度。
	for segmentLen := uint64(BUF_LEN); segmentLen < DATA_LEN; segmentLen <<= 1 {
		// 迭代子分段对。值得注意的是若 offset2 越界但 offset 1 未越界的情况。
		// 此时 offset1 到 DATA_LEN 这一段的数据一定是有序的，无须处理。
		offset1, offset2 := uint64(0), segmentLen
		for offset2 < DATA_LEN {
			end1, end2 := offset2, min(DATA_LEN, offset2+segmentLen)
			// 合并两个相邻子分段
			nonempty := [BUF_LEN]bool{}
			for i1, i2, io := offset1, offset2, offset1; io < end2; {
				// 读取数据
				if !nonempty[0] && i1 < end1 {
					read(0, i1, inBak)
					nonempty[0] = true
				}
				if !nonempty[1] && i2 < end2 {
					read(1, i2, inBak)
					nonempty[1] = true
				}
				// 比较后写入外部存储，移动数据指针
				if nonempty[0] && nonempty[1] {
					if buf[0] < buf[1] {
						write(0, io, !inBak)
						nonempty[0] = false
						i1++
					} else {
						write(1, io, !inBak)
						nonempty[1] = false
						i2++
					}
				} else if nonempty[0] {
					write(0, io, !inBak)
					nonempty[0] = false
					i1++
				} else { // nonempty[1]
					write(1, io, !inBak)
					nonempty[1] = false
					i2++
				}
				io++
			}
			offset1, offset2 = end2, end2+segmentLen
		}
		// 把输入数组的值同步到输出数组
		if offset1 < DATA_LEN {
			if inBak {
				for ; offset1 < DATA_LEN; offset1 += BUF_LEN {
					copy(buf[:], bak[offset1:])
					copy(data[offset1:], buf[:])
				}
			} else {
				for ; offset1 < DATA_LEN; offset1 += BUF_LEN {
					copy(buf[:], data[offset1:])
					copy(bak[offset1:], buf[:])
				}
			}
		}
		// 一轮合并完毕，交换 data 和 bak
		inBak = !inBak
	}

	return inBak && sort.SliceIsSorted(bak[:], func(i, j int) bool { return bak[i] < bak[j] }) ||
		!inBak && sort.SliceIsSorted(data[:], func(i, j int) bool { return data[i] < data[j] })
}
