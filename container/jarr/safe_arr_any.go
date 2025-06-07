package jarr

import (
	"bytes"
	"github.com/e7coding/coding-common/errs/jcode"
	"github.com/e7coding/coding-common/errs/jerr"
	"github.com/e7coding/coding-common/internal/deepcopy"
	"github.com/e7coding/coding-common/internal/empty"
	"github.com/e7coding/coding-common/internal/json"
	"github.com/e7coding/coding-common/jutil/jconv"
	"github.com/e7coding/coding-common/text/jstr"
	"math"
	"math/rand"
	"sort"
	"sync"
)

type SafeAnyArr struct {
	mu  sync.RWMutex
	arr []interface{}
}

func NewSafeAnyArr() *SafeAnyArr {
	return NewSafeAnyArrSize(0, 0)
}
func NewSafeAnyArrSize(size int, cap int) *SafeAnyArr {
	return &SafeAnyArr{
		arr: make([]interface{}, size, cap),
	}
}

// NewSafeAnyArrCopy 返回另一个数组的拷贝
func NewSafeAnyArrCopy(arr []interface{}) *SafeAnyArr {
	nAnyArr := make([]interface{}, len(arr))
	copy(nAnyArr, arr)
	return &SafeAnyArr{
		arr: nAnyArr,
	}
}

// 1. 基本操作.增、删、改、查

//增
// ----------------------------------------------------------------------

func (saa *SafeAnyArr) Append(val ...interface{}) *SafeAnyArr {
	saa.AppendRight(val...)
	return saa
}
func (saa *SafeAnyArr) AppendLeft(val ...interface{}) *SafeAnyArr {
	saa.mu.Lock()
	saa.arr = append(val, saa.arr...)
	saa.mu.Unlock()
	return saa
}

func (saa *SafeAnyArr) AppendRight(val ...interface{}) *SafeAnyArr {
	saa.mu.Lock()
	saa.arr = append(saa.arr, val...)
	saa.mu.Unlock()
	return saa
}

func (saa *SafeAnyArr) AppendBefore(index int, vals ...interface{}) error {
	saa.mu.Lock()
	defer saa.mu.Unlock()
	if index < 0 || index >= len(saa.arr) {
		return jerr.WithCodeF(jcode.NewErrCode(jcode.ParamErr), "index %d out of arr range %d", index, len(saa.arr))
	}
	tail := append([]interface{}{}, saa.arr[index:]...)
	saa.arr = append(saa.arr[0:index], vals...)
	saa.arr = append(saa.arr, tail...)
	return nil
}

func (saa *SafeAnyArr) AppendAfter(index int, vals ...interface{}) error {
	saa.mu.Lock()
	defer saa.mu.Unlock()
	if index < 0 || index >= len(saa.arr) {
		return jerr.WithCodeF(jcode.NewErrCode(jcode.ParamErr), "index %d out of arr range %d", index, len(saa.arr))
	}
	tail := append([]interface{}{}, saa.arr[index+1:]...)
	saa.arr = append(saa.arr[0:index+1], vals...)
	saa.arr = append(saa.arr, tail...)
	return nil
}

// 删
// ----------------------------------------------------------------------

func (saa *SafeAnyArr) Remove(index int) {
	saa.mu.Lock()
	defer saa.mu.Unlock()
	saa.removeByIndex(index)

}

func (saa *SafeAnyArr) RemoveVal(val interface{}) bool {
	saa.mu.Lock()
	defer saa.mu.Unlock()
	if i := saa.pos(val); i != -1 {
		saa.removeByIndex(i)
		return true
	}
	return false
}

func (saa *SafeAnyArr) RemoveVals(vals ...interface{}) {
	saa.mu.Lock()
	defer saa.mu.Unlock()
	for _, val := range vals {
		if i := saa.pos(val); i != -1 {
			saa.removeByIndex(i)
		}
	}
}

func (saa *SafeAnyArr) Clear() {
	saa.mu.Lock()
	defer saa.mu.Unlock()
	if len(saa.arr) > 0 {
		saa.arr = make([]interface{}, 0)
	}
}

// 改
// ----------------------------------------------------------------------

func (saa *SafeAnyArr) Set(index int, val interface{}) error {
	saa.mu.Lock()
	defer saa.mu.Unlock()
	if index < 0 || index >= len(saa.arr) {
		return jerr.WithCodeF(jcode.NewErrCode(jcode.ParamErr), "index %d out of arr range %d", index, len(saa.arr))
	}
	saa.arr[index] = val
	return nil
}

func (saa *SafeAnyArr) SetFrom(arr []interface{}) *SafeAnyArr {
	saa.mu.Lock()
	saa.arr = arr
	saa.mu.Unlock()
	return saa
}

// 查
// ----------------------------------------------------------------------

func (saa *SafeAnyArr) Pos(val interface{}) (index int) {
	return saa.Find(val)
}

func (saa *SafeAnyArr) Index(index int) (val interface{}) {
	return saa.Get(index)
}

func (saa *SafeAnyArr) Get(index int) (val interface{}) {
	saa.mu.RLock()
	defer saa.mu.RUnlock()
	if index < 0 || index >= len(saa.arr) {
		return
	}
	val = saa.arr[index]
	return
}
func (saa *SafeAnyArr) PopLeft() (val interface{}) {
	saa.mu.Lock()
	defer saa.mu.Unlock()
	if len(saa.arr) == 0 {
		return nil
	}
	val = saa.arr[0]
	saa.arr = saa.arr[1:]
	return val
}

func (saa *SafeAnyArr) PopRight() (val interface{}) {
	saa.mu.Lock()
	defer saa.mu.Unlock()
	index := len(saa.arr) - 1
	if index < 0 {
		return nil
	}
	val = saa.arr[index]
	saa.arr = saa.arr[:index]
	return val
}

func (saa *SafeAnyArr) PopLefts(size int) []interface{} {
	saa.mu.Lock()
	defer saa.mu.Unlock()
	if size <= 0 || len(saa.arr) == 0 {
		return nil
	}
	if size >= len(saa.arr) {
		arr := saa.arr
		saa.arr = saa.arr[:0]
		return arr
	}
	val := saa.arr[0:size]
	saa.arr = saa.arr[size:]
	return val
}

func (saa *SafeAnyArr) PopRights(size int) []interface{} {
	saa.mu.Lock()
	defer saa.mu.Unlock()
	if size <= 0 || len(saa.arr) == 0 {
		return nil
	}
	index := len(saa.arr) - size
	if index <= 0 {
		arr := saa.arr
		saa.arr = saa.arr[:0]
		return arr
	}
	val := saa.arr[index:]
	saa.arr = saa.arr[:index]
	return val
}
func (saa *SafeAnyArr) Find(val interface{}) (pos int) {
	saa.mu.RLock()
	defer saa.mu.RUnlock()
	pos = saa.pos(val)
	return
}

func (saa *SafeAnyArr) Rand() (val interface{}) {
	saa.mu.RLock()
	defer saa.mu.RUnlock()
	if len(saa.arr) == 0 {
		return nil
	}
	val = saa.arr[rand.Intn(len(saa.arr))]
	return
}

func (saa *SafeAnyArr) Rands(size int) []interface{} {
	saa.mu.RLock()
	defer saa.mu.RUnlock()
	if size <= 0 || len(saa.arr) == 0 {
		return nil
	}
	arr := make([]interface{}, size)
	for i := 0; i < size; i++ {
		arr[i] = saa.arr[rand.Intn(len(saa.arr))]
	}
	return arr
}

// Chunk 将数组拆分成多个size长度的区块
func (saa *SafeAnyArr) Chunk(size int) [][]interface{} {
	if size < 1 {
		return nil
	}
	saa.mu.RLock()
	defer saa.mu.RUnlock()
	length := len(saa.arr)
	chunks := int(math.Ceil(float64(length) / float64(size)))
	var n [][]interface{}
	for i, end := 0, 0; chunks > 0; chunks-- {
		end = (i + 1) * size
		if end > length {
			end = length
		}
		n = append(n, saa.arr[i*size:end])
		i++
	}
	return n
}

// Compact 移除假值
func (saa *SafeAnyArr) Compact() *SafeAnyArr {
	saa.mu.Lock()
	defer saa.mu.Unlock()
	for i := 0; i < len(saa.arr); {
		if empty.IsEmpty(saa.arr[i]) {
			saa.arr = append(saa.arr[:i], saa.arr[i+1:]...)
		} else {
			i++
		}
	}
	return saa
}

func (saa *SafeAnyArr) Sort(compare func(v1, v2 interface{}) bool) *SafeAnyArr {
	saa.mu.Lock()
	defer saa.mu.Unlock()
	sort.Slice(saa.arr, func(i, j int) bool {
		return compare(saa.arr[i], saa.arr[j])
	})
	return saa
}

func (saa *SafeAnyArr) Len() int {
	saa.mu.RLock()
	defer saa.mu.RUnlock()
	length := len(saa.arr)
	return length
}

func (saa *SafeAnyArr) Slice() []interface{} {
	saa.mu.RLock()
	defer saa.mu.RUnlock()
	arr := make([]interface{}, len(saa.arr))
	copy(arr, saa.arr)
	return arr
}

func (saa *SafeAnyArr) Interfaces() []interface{} {
	return saa.Slice()
}

func (saa *SafeAnyArr) Clone() (newAnyArr *SafeAnyArr) {
	saa.mu.RLock()
	defer saa.mu.RUnlock()
	arr := make([]interface{}, len(saa.arr))
	copy(arr, saa.arr)
	return &SafeAnyArr{arr: arr}
}

func (saa *SafeAnyArr) Contain(val interface{}) bool {
	return saa.Find(val) != -1
}

func (saa *SafeAnyArr) Uniq() *SafeAnyArr {
	saa.mu.Lock()
	defer saa.mu.Unlock()
	if len(saa.arr) == 0 {
		return saa
	}
	var (
		ok   bool
		tVal interface{}
		set  = make(map[interface{}]struct{})
		uniq = make([]interface{}, 0, len(saa.arr))
	)
	for i := 0; i < len(saa.arr); i++ {
		tVal = saa.arr[i]
		if _, ok = set[tVal]; ok {
			continue
		}
		set[tVal] = struct{}{}
		uniq = append(uniq, tVal)
	}
	saa.arr = uniq
	return saa
}

func (saa *SafeAnyArr) LockFunc(f func(arr []interface{})) *SafeAnyArr {
	saa.mu.Lock()
	defer saa.mu.Unlock()
	f(saa.arr)
	return saa
}

func (saa *SafeAnyArr) RLockFunc(f func(arr []interface{})) *SafeAnyArr {
	saa.mu.RLock()
	defer saa.mu.RUnlock()
	f(saa.arr)
	return saa
}

func (saa *SafeAnyArr) Merge(arr ...interface{}) *SafeAnyArr {
	return saa.Append(arr)
}

func (saa *SafeAnyArr) Reverse() *SafeAnyArr {
	saa.mu.Lock()
	defer saa.mu.Unlock()
	for i, j := 0, len(saa.arr)-1; i < j; i, j = i+1, j-1 {
		saa.arr[i], saa.arr[j] = saa.arr[j], saa.arr[i]
	}
	return saa
}

func (saa *SafeAnyArr) Join(sep string) string {
	saa.mu.RLock()
	defer saa.mu.RUnlock()
	if len(saa.arr) == 0 {
		return ""
	}
	buffer := bytes.NewBuffer(nil)
	for k, v := range saa.arr {
		buffer.WriteString(jconv.String(v))
		if k != len(saa.arr)-1 {
			buffer.WriteString(sep)
		}
	}
	return buffer.String()
}

func (saa *SafeAnyArr) Count(val interface{}) (count int) {
	saa.mu.RLock()
	defer saa.mu.RUnlock()
	for _, v := range saa.arr {
		if v == val {
			count += 1
		}
	}
	return
}

func (saa *SafeAnyArr) String() string {
	saa.mu.RLock()
	defer saa.mu.RUnlock()
	buffer := bytes.NewBuffer(nil)
	buffer.WriteByte('[')
	s := ""
	for k, v := range saa.arr {
		s = jconv.String(v)
		if jstr.IsNumeric(s) {
			buffer.WriteString(s)
		} else {
			buffer.WriteString(`"` + jstr.QuoteMeta(s, `"\`) + `"`)
		}
		if k != len(saa.arr)-1 {
			buffer.WriteByte(',')
		}
	}
	buffer.WriteByte(']')
	return buffer.String()
}

func (saa *SafeAnyArr) MarshalJSON() (b []byte, err error) {
	saa.mu.RLock()
	defer saa.mu.RUnlock()
	b, err = json.Marshal(saa.arr)
	return
}

func (saa *SafeAnyArr) UnmarshalJSON(b []byte) error {
	if saa.arr == nil {
		saa.arr = make([]interface{}, 0)
	}
	saa.mu.Lock()
	defer saa.mu.Unlock()
	if err := json.UnmarshalUseNumber(b, &saa.arr); err != nil {
		return err
	}
	return nil
}

func (saa *SafeAnyArr) UnmarshalValue(val interface{}) error {
	saa.mu.Lock()
	defer saa.mu.Unlock()
	switch val.(type) {
	case string, []byte:
		return json.UnmarshalUseNumber(jconv.Bytes(val), &saa.arr)
	default:
		saa.arr = jconv.SliceAny(val)
	}
	return nil
}

func (saa *SafeAnyArr) Filter(filter func(index int, val interface{}) bool) *SafeAnyArr {
	saa.mu.Lock()
	defer saa.mu.Unlock()
	for i := 0; i < len(saa.arr); {
		if filter(i, saa.arr[i]) {
			saa.arr = append(saa.arr[:i], saa.arr[i+1:]...)
		} else {
			i++
		}
	}
	return saa
}

func (saa *SafeAnyArr) FilterNil() *SafeAnyArr {
	saa.mu.Lock()
	defer saa.mu.Unlock()
	for i := 0; i < len(saa.arr); {
		if empty.IsNil(saa.arr[i]) {
			saa.arr = append(saa.arr[:i], saa.arr[i+1:]...)
		} else {
			i++
		}
	}
	return saa
}

func (saa *SafeAnyArr) Walk(f func(val interface{}) interface{}) *SafeAnyArr {
	saa.mu.Lock()
	defer saa.mu.Unlock()
	for i, v := range saa.arr {
		saa.arr[i] = f(v)
	}
	return saa
}

func (saa *SafeAnyArr) IsEmpty() bool {
	return saa.Len() == 0
}

// DeepCopy implements interface for deep copy of current type.
func (saa *SafeAnyArr) DeepCopy() interface{} {
	saa.mu.RLock()
	defer saa.mu.RUnlock()
	newSlice := make([]interface{}, len(saa.arr))
	for i, v := range saa.arr {
		newSlice[i] = deepcopy.Copy(v)
	}
	return &SafeAnyArr{arr: newSlice}
}
func (saa *SafeAnyArr) pos(val interface{}) int {
	if len(saa.arr) == 0 {
		return -1
	}
	result := -1
	for index, v := range saa.arr {
		if v == val {
			result = index
			break
		}
	}
	return result
}
func (saa *SafeAnyArr) removeByIndex(index int) bool {
	if index < 0 || index >= len(saa.arr) {
		return false
	}
	if index == 0 {
		saa.arr = saa.arr[1:]
		return true
	} else if index == len(saa.arr)-1 {
		saa.arr = saa.arr[:index]
		return true
	}
	saa.arr = append(saa.arr[:index], saa.arr[index+1:]...)
	return true
}
