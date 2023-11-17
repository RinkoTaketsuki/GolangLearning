package example

import (
	"fmt"
	"sync"
	"sync/atomic"
)

/*
	乐观锁的原理（用 C 语言解释）：

1. 假定要被保护的变量类型为 DataType，变量名为 data，有：

	Datatype data;

2. 定义新类型 DataTypeWithVersion，有：

	typedef struct DataTypeWithVersion {
		DataType data,
		unsigned long long version
	} DataTypeWithVersion;

3. 初始化：

	DataTypeWithVersion dv = { data, 0 };

4. 读取时直接读取结构体的 data 部分，尝试写入时需要以下的流程，需要 CAS 技术：

	int try_write(DataTypeWithVersion *dv, DataType *newData) {
		DataTypeWithVersion oldDv = *dv;
		DataTypeWithVersion newDv = { *newData, dv.version + 1 };
		if compare_and_swap(&dv, &oldDv, &newDv)
			// success
		else
			// failure
	}

5. 具体业务需要使用死循环（自旋）不断尝试上述流程，直到成功写入
*/
type OptimisticAtomicValue struct {
	v uint64
}

func NewOptimisticAtomicValue(init uint32) *OptimisticAtomicValue {
	return &OptimisticAtomicValue{uint64(init)}
}

func (oav *OptimisticAtomicValue) getVersion() uint32 {
	return uint32(oav.v >> 32)
}

func (oav *OptimisticAtomicValue) setVersion(ver uint32) {
	oav.v &= (uint64(ver)<<32 + 0xffffffff)
}

func (oav *OptimisticAtomicValue) Load() uint32 {
	return uint32(oav.v & 0xffffffff)
}

func (oav *OptimisticAtomicValue) TryStore(newValue uint32) (success bool) {
	oldOav := &OptimisticAtomicValue{v: oav.v}
	newOav := NewOptimisticAtomicValue(newValue)
	newOav.setVersion(oldOav.getVersion() + 1)
	return atomic.CompareAndSwapUint64(&oav.v, oldOav.v, newOav.v)
}

func (oav *OptimisticAtomicValue) TryModify(mod func(uint32) uint32) (success bool) {
	oldOav := &OptimisticAtomicValue{v: oav.v}
	newOav := NewOptimisticAtomicValue(mod(oldOav.Load()))
	newOav.setVersion(oldOav.getVersion() + 1)
	return atomic.CompareAndSwapUint64(&oav.v, oldOav.v, newOav.v)
}

// 场景：存取款，账户初始余额为 0。有 16 个用户并行地存取款，其中 10 个用户存 20
// 次 100 元，另外 6 个取 30 次 50 元，预期最终余额为 11000 元
func OavExample() {
	oav := NewOptimisticAtomicValue(uint32(0))
	var wg sync.WaitGroup
	wg.Add(16)
	deposit := func() {
		defer wg.Done()
		for i := 0; i < 20; i++ {
			for !oav.TryModify(func(i uint32) uint32 {
				return i + 100
			}) {
			}
		}
	}
	withdraw := func() {
		defer wg.Done()
		for i := 0; i < 30; i++ {
			for !oav.TryModify(func(i uint32) uint32 {
				return i - 50
			}) {
			}
		}
	}
	for i := 0; i < 10; i++ {
		go deposit()
	}
	for i := 0; i < 6; i++ {
		go withdraw()
	}
	wg.Wait()
	final := oav.Load()
	if final == 11000 {
		fmt.Println("OK")
	} else {
		fmt.Printf("Failed! Final amount is %d\n", final)
	}
}
