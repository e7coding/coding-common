package jarr

import (
	"bytes"
	"github.com/e7coding/coding-common/errs/jcode"
	"github.com/e7coding/coding-common/errs/jerr"
	"github.com/e7coding/coding-common/internal/empty"
	"github.com/e7coding/coding-common/internal/json"
	"github.com/e7coding/coding-common/jutil/jconv"
	"github.com/e7coding/coding-common/jutil/jrand"
	"github.com/e7coding/coding-common/text/jstr"
	"math"
	"sort"
	"sync"
)

type SafeIntArr struct {
	mu  sync.RWMutex
	arr []int
}

func NewSafeIntArr() *SafeIntArr {
	return NewSafeIntArrSize(0, 0)
}
func NewSafeIntArrSize(size int, cap int) *SafeIntArr {
	return &SafeIntArr{
		arr: make([]int, size, cap),
	}
}

// NewSafeIntArrCopy 返回另一个数组的拷贝
func NewSafeIntArrCopy(arr []int) *SafeIntArr {
	nIntArr := make([]int, len(arr))
	copy(nIntArr, arr)
	return &SafeIntArr{
		arr: nIntArr,
	}
}

// 1. 基本操作.增、删、改、查

//增
// ----------------------------------------------------------------------

func (sia *SafeIntArr) Append(val ...int) *SafeIntArr {
	sia.AppendRight(val...)
	return sia
}
func (sia *SafeIntArr) AppendLeft(val ...int) *SafeIntArr {
	sia.mu.Lock()
	defer sia.mu.Unlock()
	sia.arr = append(val, sia.arr...)

	return sia
}

func (sia *SafeIntArr) AppendRight(val ...int) *SafeIntArr {
	sia.mu.Lock()
	defer sia.mu.Unlock()
	sia.arr = append(sia.arr, val...)

	return sia
}

func (sia *SafeIntArr) AppendBefore(index int, vals ...int) error {
	sia.mu.Lock()
	defer sia.mu.Unlock()
	if index < 0 || index >= len(sia.arr) {
		return jerr.WithCodeF(jcode.NewErrCode(jcode.ParamErr), "index %d out of arr range %d", index, len(sia.arr))
	}
	tail := append([]int{}, sia.arr[index:]...)
	sia.arr = append(sia.arr[0:index], vals...)
	sia.arr = append(sia.arr, tail...)

	return nil
}

func (sia *SafeIntArr) AppendAfter(index int, vals ...int) error {
	sia.mu.Lock()
	defer sia.mu.Unlock()
	if index < 0 || index >= len(sia.arr) {
		return jerr.WithCodeF(jcode.NewErrCode(jcode.ParamErr), "index %d out of arr range %d", index, len(sia.arr))
	}
	tail := append([]int{}, sia.arr[index+1:]...)
	sia.arr = append(sia.arr[0:index+1], vals...)
	sia.arr = append(sia.arr, tail...)

	return nil
}

// 删
// ----------------------------------------------------------------------

func (sia *SafeIntArr) Remove(index int) {
	sia.mu.Lock()
	defer sia.mu.Unlock()
	sia.removeByIndex(index)

}

func (sia *SafeIntArr) RemoveVal(val int) bool {
	sia.mu.Lock()
	defer sia.mu.Unlock()
	if i := sia.pos(val); i != -1 {
		sia.removeByIndex(i)
		return true
	}

	return false
}

func (sia *SafeIntArr) RemoveVals(vals ...int) {
	sia.mu.Lock()
	defer sia.mu.Unlock()
	for _, val := range vals {
		if i := sia.pos(val); i != -1 {
			sia.removeByIndex(i)
		}
	}

}

func (sia *SafeIntArr) Clear() {
	sia.mu.Lock()
	defer sia.mu.Unlock()
	if len(sia.arr) > 0 {
		sia.arr = make([]int, 0)
	}

}

// 改
// ----------------------------------------------------------------------

func (sia *SafeIntArr) Set(index int, val int) error {
	sia.mu.Lock()
	defer sia.mu.Unlock()
	if index < 0 || index >= len(sia.arr) {
		return jerr.WithCodeF(jcode.NewErrCode(jcode.ParamErr), "index %d out of arr range %d", index, len(sia.arr))
	}
	sia.arr[index] = val

	return nil
}

func (sia *SafeIntArr) SetFrom(arr []int) *SafeIntArr {
	sia.mu.Lock()
	defer sia.mu.Unlock()
	sia.arr = arr

	return sia
}

// 查
// ----------------------------------------------------------------------

func (sia *SafeIntArr) Pos(val int) (index int) {
	return sia.Find(val)
}

func (sia *SafeIntArr) Index(index int) (val int) {
	return sia.Get(index)
}

func (sia *SafeIntArr) Get(index int) (val int) {
	sia.mu.RLock()
	defer sia.mu.RUnlock()
	if index < 0 || index >= len(sia.arr) {
		return
	}
	val = sia.arr[index]

	return
}
func (sia *SafeIntArr) PopLeft() (val int) {
	sia.mu.Lock()
	defer sia.mu.Unlock()
	if len(sia.arr) == 0 {
		return 0
	}
	val = sia.arr[0]
	sia.arr = sia.arr[1:]

	return val
}

func (sia *SafeIntArr) PopRight() (val int) {
	sia.mu.Lock()
	defer sia.mu.Unlock()
	index := len(sia.arr) - 1
	if index < 0 {
		return 0
	}
	val = sia.arr[index]
	sia.arr = sia.arr[:index]

	return val
}

func (sia *SafeIntArr) PopLefts(size int) []int {
	sia.mu.Lock()
	defer sia.mu.Unlock()
	if size <= 0 || len(sia.arr) == 0 {
		return nil
	}
	if size >= len(sia.arr) {
		arr := sia.arr
		sia.arr = sia.arr[:0]
		return arr
	}
	val := sia.arr[0:size]
	sia.arr = sia.arr[size:]

	return val
}

func (sia *SafeIntArr) PopRights(size int) []int {
	sia.mu.Lock()
	defer sia.mu.Unlock()
	if size <= 0 || len(sia.arr) == 0 {
		return nil
	}
	index := len(sia.arr) - size
	if index <= 0 {
		arr := sia.arr
		sia.arr = sia.arr[:0]
		return arr
	}
	val := sia.arr[index:]
	sia.arr = sia.arr[:index]

	return val
}
func (sia *SafeIntArr) Find(val int) (pos int) {
	sia.mu.RLock()
	defer sia.mu.RUnlock()
	pos = sia.pos(val)

	return
}

func (sia *SafeIntArr) Rand() (val int) {
	sia.mu.RLock()
	defer sia.mu.RUnlock()
	if len(sia.arr) == 0 {
		return 0
	}
	val = sia.arr[jrand.Intn(len(sia.arr))]

	return
}

func (sia *SafeIntArr) Rands(size int) []int {
	sia.mu.RLock()
	defer sia.mu.RUnlock()

	if size <= 0 || len(sia.arr) == 0 {
		return nil
	}
	arr := make([]int, size)
	for i := 0; i < size; i++ {
		arr[i] = sia.arr[jrand.Intn(len(sia.arr))]
	}

	return arr
}

// Chunk 将数组拆分成多个size长度的区块
func (sia *SafeIntArr) Chunk(size int) [][]int {
	if size < 1 {
		return nil
	}
	sia.mu.RLock()
	defer sia.mu.RUnlock()
	length := len(sia.arr)
	chunks := int(math.Ceil(float64(length) / float64(size)))
	var n [][]int
	for i, end := 0, 0; chunks > 0; chunks-- {
		end = (i + 1) * size
		if end > length {
			end = length
		}
		n = append(n, sia.arr[i*size:end])
		i++
	}

	return n
}

// Compact 移除假值
func (sia *SafeIntArr) Compact() *SafeIntArr {
	sia.mu.Lock()
	defer sia.mu.Unlock()
	for i := 0; i < len(sia.arr); {
		if empty.IsEmpty(sia.arr[i]) {
			sia.arr = append(sia.arr[:i], sia.arr[i+1:]...)
		} else {
			i++
		}
	}

	return sia
}

func (sia *SafeIntArr) Sort(compare func(v1, v2 int) bool) *SafeIntArr {
	sia.mu.Lock()
	defer sia.mu.Unlock()
	sort.Slice(sia.arr, func(i, j int) bool {
		return compare(sia.arr[i], sia.arr[j])
	})

	return sia
}

func (sia *SafeIntArr) Len() int {
	sia.mu.RLock()
	defer sia.mu.RUnlock()
	length := len(sia.arr)

	return length
}

func (sia *SafeIntArr) Slice() []int {
	sia.mu.RLock()
	defer sia.mu.RUnlock()
	arr := make([]int, len(sia.arr))
	copy(arr, sia.arr)

	return arr
}

func (sia *SafeIntArr) Interfaces() []int {
	return sia.Slice()
}

func (sia *SafeIntArr) Clone() (newIntArr *SafeIntArr) {
	sia.mu.RLock()
	defer sia.mu.RUnlock()
	arr := make([]int, len(sia.arr))
	copy(arr, sia.arr)

	return &SafeIntArr{arr: arr}
}

func (sia *SafeIntArr) Contain(val int) bool {
	return sia.Find(val) != -1
}

func (sia *SafeIntArr) Uniq() *SafeIntArr {
	sia.mu.Lock()
	defer sia.mu.Unlock()
	if len(sia.arr) == 0 {
		return sia
	}
	var (
		ok   bool
		tVal int
		set  = make(map[int]struct{})
		uniq = make([]int, 0, len(sia.arr))
	)
	for i := 0; i < len(sia.arr); i++ {
		tVal = sia.arr[i]
		if _, ok = set[tVal]; ok {
			continue
		}
		set[tVal] = struct{}{}
		uniq = append(uniq, tVal)
	}
	sia.arr = uniq

	return sia
}

func (sia *SafeIntArr) LockFunc(f func(arr []int)) *SafeIntArr {
	sia.mu.Lock()
	defer sia.mu.Unlock()
	f(sia.arr)

	return sia
}

func (sia *SafeIntArr) RLockFunc(f func(arr []int)) *SafeIntArr {
	sia.mu.RLock()
	defer sia.mu.RUnlock()
	f(sia.arr)

	return sia
}

func (sia *SafeIntArr) Merge(arr ...int) *SafeIntArr {
	return sia.Append(arr...)
}

func (sia *SafeIntArr) Reverse() *SafeIntArr {
	sia.mu.Lock()
	defer sia.mu.Unlock()
	for i, j := 0, len(sia.arr)-1; i < j; i, j = i+1, j-1 {
		sia.arr[i], sia.arr[j] = sia.arr[j], sia.arr[i]
	}

	return sia
}

func (sia *SafeIntArr) Join(sep string) string {
	sia.mu.RLock()
	defer sia.mu.RUnlock()
	if len(sia.arr) == 0 {
		return ""
	}
	buffer := bytes.NewBuffer(nil)
	for k, v := range sia.arr {
		buffer.WriteString(jconv.String(v))
		if k != len(sia.arr)-1 {
			buffer.WriteString(sep)
		}
	}

	return buffer.String()
}

func (sia *SafeIntArr) Count(val int) (count int) {
	sia.mu.RLock()
	defer sia.mu.RUnlock()
	for _, v := range sia.arr {
		if v == val {
			count += 1
		}
	}

	return
}

func (sia *SafeIntArr) String() string {
	sia.mu.RLock()
	defer sia.mu.RUnlock()

	buffer := bytes.NewBuffer(nil)
	buffer.WriteByte('[')
	s := ""
	for k, v := range sia.arr {
		s = jconv.String(v)
		if jstr.IsNumeric(s) {
			buffer.WriteString(s)
		} else {
			buffer.WriteString(`"` + jstr.QuoteMeta(s, `"\`) + `"`)
		}
		if k != len(sia.arr)-1 {
			buffer.WriteByte(',')
		}
	}
	buffer.WriteByte(']')

	return buffer.String()
}

func (sia *SafeIntArr) MarshalJSON() (b []byte, err error) {
	sia.mu.RLock()
	defer sia.mu.RUnlock()
	b, err = json.Marshal(sia.arr)

	return
}

func (sia *SafeIntArr) UnmarshalJSON(b []byte) error {
	if sia.arr == nil {
		sia.arr = make([]int, 0)
	}
	sia.mu.Lock()
	defer sia.mu.Unlock()
	if err := json.UnmarshalUseNumber(b, &sia.arr); err != nil {
		return err
	}

	return nil
}

func (sia *SafeIntArr) UnmarshalValue(val interface{}) error {
	sia.mu.Lock()
	defer sia.mu.Unlock()
	switch val.(type) {
	case string, []byte:
		return json.UnmarshalUseNumber(jconv.Bytes(val), &sia.arr)
	default:
		sia.arr = jconv.SliceInt(val)
	}

	return nil
}

func (sia *SafeIntArr) Filter(filter func(index int, val int) bool) *SafeIntArr {
	sia.mu.Lock()
	defer sia.mu.Unlock()
	for i := 0; i < len(sia.arr); {
		if filter(i, sia.arr[i]) {
			sia.arr = append(sia.arr[:i], sia.arr[i+1:]...)
		} else {
			i++
		}
	}

	return sia
}

func (sia *SafeIntArr) FilterNil() *SafeIntArr {
	sia.mu.Lock()
	defer sia.mu.Unlock()
	for i := 0; i < len(sia.arr); {
		if empty.IsNil(sia.arr[i]) {
			sia.arr = append(sia.arr[:i], sia.arr[i+1:]...)
		} else {
			i++
		}
	}

	return sia
}

func (sia *SafeIntArr) Walk(f func(val int) int) *SafeIntArr {
	sia.mu.Lock()
	defer sia.mu.Unlock()
	for i, v := range sia.arr {
		sia.arr[i] = f(v)
	}

	return sia
}

func (sia *SafeIntArr) IsEmpty() bool {
	return sia.Len() == 0
}

// DeepCopy implements interface for deep copy of current type.
func (sia *SafeIntArr) DeepCopy() interface{} {
	sia.mu.RLock()
	defer sia.mu.RUnlock()
	newSlice := make([]int, len(sia.arr))
	copy(newSlice, sia.arr)

	return &IntArr{arr: newSlice}
}
func (sia *SafeIntArr) pos(val int) int {
	if len(sia.arr) == 0 {
		return -1
	}
	result := -1
	for index, v := range sia.arr {
		if v == val {
			result = index
			break
		}
	}
	return result
}
func (sia *SafeIntArr) removeByIndex(index int) bool {
	if index < 0 || index >= len(sia.arr) {
		return false
	}
	if index == 0 {
		sia.arr = sia.arr[1:]
		return true
	} else if index == len(sia.arr)-1 {
		sia.arr = sia.arr[:index]
		return true
	}
	sia.arr = append(sia.arr[:index], sia.arr[index+1:]...)
	return true
}
