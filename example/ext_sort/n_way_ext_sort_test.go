package ext_sort

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"
)

func TestExtMergeSort2way(t *testing.T) {
	testcases := [][]uint64{
		{},
		{42},
		{0, 0},
		{1, 2},
		{2, 1},
	}
	for i := 0; i < 100; i++ {
		dataLenBig, err := rand.Int(rand.Reader, big.NewInt(1000))
		if err != nil {
			t.Error(err)
		}
		dataLen := dataLenBig.Uint64()
		data := make([]uint64, dataLen)
		for i := range data {
			val, err := rand.Int(rand.Reader, big.NewInt(100000))
			if err != nil {
				panic(err)
			}
			data[i] = val.Uint64()
		}
		testcases = append(testcases, data)
	}
	for i, testcase := range testcases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ExtMergeSort2way(testcase)
			if !sort.SliceIsSorted(testcase, func(i, j int) bool {
				return testcase[i] < testcase[j]
			}) {
				t.FailNow()
			}
		})
	}
}

func WriteFixedLengthRandomNumbersToFile(f *os.File, num int, upperBound *big.Int) error {
	var sb strings.Builder
	for range num {
		i, err := rand.Int(rand.Reader, upperBound)
		if err != nil {
			return err
		}
		sb.WriteString(fmt.Sprintf("%08d", i.Uint64()))
	}
	_, err := f.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	fmt.Println(sb.String())
	_, err = f.WriteString(sb.String())
	return err
}

func TestExtMergeSortNWay(t *testing.T) {
	file, err := os.Create("numbers.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	type Case struct {
		num        int
		upperBound *big.Int
		n          int
	}
	testcases := []Case{
		{0, big.NewInt(10), 2},
		{1, big.NewInt(10), 2},
		{2, big.NewInt(10), 2},
		{3, big.NewInt(10), 2},
		{4, big.NewInt(10), 2},
	}
	for range 100 {
		numBig, err := rand.Int(rand.Reader, big.NewInt(16))
		if err != nil {
			t.Fatal(err)
		}
		num := int(numBig.Int64())
		upperBound, err := rand.Int(rand.Reader, big.NewInt(16))
		if err != nil {
			t.Fatal(err)
		}
		upperBound.Add(upperBound, big.NewInt(1))
		n, err := rand.Int(rand.Reader, big.NewInt(8))
		if err != nil {
			t.Fatal(err)
		}
		if n.Cmp(big.NewInt(2)) < 0 {
			continue
		}
		testcases = append(testcases, Case{num, upperBound, int(n.Int64())})
	}
	for _, testcase := range testcases {
		if !t.Run(
			fmt.Sprintf("n: %d, number of elements: %d, element range: [0, %d)",
				testcase.n, testcase.num, testcase.upperBound.Int64()), func(t *testing.T) {
				WriteFixedLengthRandomNumbersToFile(file, testcase.num, testcase.upperBound)
				if err := ExtMergeSortNWay(file, 8, func(b1, b2 []byte) bool {
					return strings.Compare(string(b1), string(b2)) < 0
				}, testcase.n); err != nil {
					t.Fatal(err)
				}
				file.Seek(0, io.SeekStart)
				buf := make([]byte, 8*testcase.num)
				file.Read(buf)
				nums := make([]int, testcase.num)
				for i := range nums {
					nums[i], err = strconv.Atoi(string(buf[8*i : 8*(i+1)]))
					if err != nil {
						t.Fatal(err)
					}
				}
				if !sort.IntsAreSorted(nums) {
					t.FailNow()
				}
			}) {
			t.FailNow()
		}
	}
}

func TestXXX(t *testing.T) {
	file, err := os.Create("numbers.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	file.WriteString("5648581")
	err = ExtMergeSortNWay(file, 1, func(b1, b2 []byte) bool { return b1[0] < b2[0] }, 6)
	if err != nil {
		t.Fatal(err)
	}
}
