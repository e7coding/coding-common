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
)

type IntArr struct {
	arr []int
}

func NewIntArr() *IntArr {
	return NewIntArrSize(0, 0)
}
func NewIntArrSize(size int, cap int) *IntArr {
	return &IntArr{
		arr: make([]int, size, cap),
	}
}

// NewIntArrCopy 返回另一个数组的拷贝
func NewIntArrCopy(arr []int) *IntArr {
	nIntArr := make([]int, len(arr))
	copy(nIntArr, arr)
	return &IntArr{
		arr: nIntArr,
	}
}

// Chunk 将数组拆分成多个size长度的区块
func (aa *IntArr) Chunk(size int) [][]int {
	if size < 1 {
		return nil
	}

	length := len(aa.arr)
	chunks := int(math.Ceil(float64(length) / float64(size)))
	var n [][]int
	for i, end := 0, 0; chunks > 0; chunks-- {
		end = (i + 1) * size
		if end > length {
			end = length
		}
		n = append(n, aa.arr[i*size:end])
		i++
	}

	return n
}

func (aa *IntArr) Pos(val int) (index int) {
	return aa.Search(val)
}

func (aa *IntArr) Index(index int) (val int) {
	return aa.Get(index)
}

func (aa *IntArr) Get(index int) (val int) {

	if index < 0 || index >= len(aa.arr) {
		return
	}
	val = aa.arr[index]

	return
}

func (aa *IntArr) Set(index int, val int) error {

	if index < 0 || index >= len(aa.arr) {
		return jerr.WithCodeF(jcode.NewErrCode(jcode.ParamErr), "index %d out of arr range %d", index, len(aa.arr))
	}
	aa.arr[index] = val

	return nil
}

func (aa *IntArr) SetFrom(arr []int) *IntArr {

	aa.arr = arr

	return aa
}

func (aa *IntArr) Sort(compare func(v1, v2 int) bool) *IntArr {

	sort.Slice(aa.arr, func(i, j int) bool {
		return compare(aa.arr[i], aa.arr[j])
	})

	return aa
}

func (aa *IntArr) InsertBefore(index int, vals ...int) error {

	if index < 0 || index >= len(aa.arr) {
		return jerr.WithCodeF(jcode.NewErrCode(jcode.ParamErr), "index %d out of arr range %d", index, len(aa.arr))
	}
	// 取尾
	tail := append([]int{}, aa.arr[index:]...)
	aa.arr = append(aa.arr[0:index], vals...)
	// 放尾
	aa.arr = append(aa.arr, tail...)

	return nil
}

func (aa *IntArr) InsertAfter(index int, vals ...int) error {

	if index < 0 || index >= len(aa.arr) {
		return jerr.WithCodeF(jcode.NewErrCode(jcode.ParamErr), "index %d out of arr range %d", index, len(aa.arr))
	}
	tail := append([]int{}, aa.arr[index+1:]...)
	aa.arr = append(aa.arr[0:index+1], vals...)
	aa.arr = append(aa.arr, tail...)

	return nil
}

func (aa *IntArr) Remove(index int) {

	aa.removeByIndex(index)

}

func (aa *IntArr) RemoveVal(val int) bool {

	if i := aa.pos(val); i != -1 {
		aa.removeByIndex(i)
		return true
	}

	return false
}

func (aa *IntArr) RemoveVals(vals ...int) {

	for _, val := range vals {
		if i := aa.pos(val); i != -1 {
			aa.removeByIndex(i)
		}
	}

}

func (aa *IntArr) PushLeft(val ...int) *IntArr {

	aa.arr = append(val, aa.arr...)

	return aa
}

func (aa *IntArr) PushRight(val ...int) *IntArr {

	aa.arr = append(aa.arr, val...)

	return aa
}

func (aa *IntArr) PopLeft() (val int) {

	if len(aa.arr) == 0 {
		return 0
	}
	val = aa.arr[0]
	aa.arr = aa.arr[1:]

	return val
}

func (aa *IntArr) PopRight() (val int) {

	index := len(aa.arr) - 1
	if index < 0 {
		return 0
	}
	val = aa.arr[index]
	aa.arr = aa.arr[:index]

	return val
}

func (aa *IntArr) PopLefts(size int) []int {

	if size <= 0 || len(aa.arr) == 0 {
		return nil
	}
	if size >= len(aa.arr) {
		arr := aa.arr
		aa.arr = aa.arr[:0]
		return arr
	}
	val := aa.arr[0:size]
	aa.arr = aa.arr[size:]

	return val
}

func (aa *IntArr) PopRights(size int) []int {

	if size <= 0 || len(aa.arr) == 0 {
		return nil
	}
	index := len(aa.arr) - size
	if index <= 0 {
		arr := aa.arr
		aa.arr = aa.arr[:0]
		return arr
	}
	val := aa.arr[index:]
	aa.arr = aa.arr[:index]

	return val
}

func (aa *IntArr) Append(val ...int) *IntArr {
	aa.PushRight(val...)
	return aa
}

func (aa *IntArr) Len() int {

	length := len(aa.arr)

	return length
}

func (aa *IntArr) Slice() []int {

	arr := make([]int, len(aa.arr))
	copy(arr, aa.arr)

	return arr
}

func (aa *IntArr) Interfaces() []int {
	return aa.Slice()
}

func (aa *IntArr) Clone() (newIntArr *IntArr) {

	arr := make([]int, len(aa.arr))
	copy(arr, aa.arr)

	return &IntArr{arr: arr}
}

func (aa *IntArr) Clear() {

	if len(aa.arr) > 0 {
		aa.arr = make([]int, 0)
	}

}

func (aa *IntArr) Contain(val int) bool {
	return aa.Search(val) != -1
}

func (aa *IntArr) Search(val int) (pos int) {

	pos = aa.pos(val)

	return
}

func (aa *IntArr) Unique() *IntArr {

	if len(aa.arr) == 0 {
		return aa
	}
	var (
		ok     bool
		temp   int
		set    = make(map[int]struct{})
		unique = make([]int, 0, len(aa.arr))
	)
	for i := 0; i < len(aa.arr); i++ {
		temp = aa.arr[i]
		if _, ok = set[temp]; ok {
			continue
		}
		set[temp] = struct{}{}
		unique = append(unique, temp)
	}
	aa.arr = unique

	return aa
}

func (aa *IntArr) LockFunc(f func(arr []int)) *IntArr {

	f(aa.arr)

	return aa
}

func (aa *IntArr) RLockFunc(f func(arr []int)) *IntArr {

	f(aa.arr)

	return aa
}

func (aa *IntArr) Merge(arr ...int) *IntArr {
	return aa.Append(arr...)
}

func (aa *IntArr) Rand() (val int) {

	if len(aa.arr) == 0 {
		return 0
	}
	val = aa.arr[jrand.Intn(len(aa.arr))]

	return
}

func (aa *IntArr) Rands(size int) []int {

	if size <= 0 || len(aa.arr) == 0 {
		return nil
	}
	arr := make([]int, size)
	for i := 0; i < size; i++ {
		arr[i] = aa.arr[jrand.Intn(len(aa.arr))]
	}

	return arr
}

func (aa *IntArr) Reverse() *IntArr {

	for i, j := 0, len(aa.arr)-1; i < j; i, j = i+1, j-1 {
		aa.arr[i], aa.arr[j] = aa.arr[j], aa.arr[i]
	}

	return aa
}

func (aa *IntArr) Join(sep string) string {

	if len(aa.arr) == 0 {
		return ""
	}
	buffer := bytes.NewBuffer(nil)
	for k, v := range aa.arr {
		buffer.WriteString(jconv.String(v))
		if k != len(aa.arr)-1 {
			buffer.WriteString(sep)
		}
	}

	return buffer.String()
}

func (aa *IntArr) Count(val int) (count int) {

	for _, v := range aa.arr {
		if v == val {
			count += 1
		}
	}

	return
}

func (aa *IntArr) String() string {

	buffer := bytes.NewBuffer(nil)
	buffer.WriteByte('[')
	s := ""
	for k, v := range aa.arr {
		s = jconv.String(v)
		if jstr.IsNumeric(s) {
			buffer.WriteString(s)
		} else {
			buffer.WriteString(`"` + jstr.QuoteMeta(s, `"\`) + `"`)
		}
		if k != len(aa.arr)-1 {
			buffer.WriteByte(',')
		}
	}
	buffer.WriteByte(']')

	return buffer.String()
}

func (aa *IntArr) MarshalJSON() (b []byte, err error) {

	b, err = json.Marshal(aa.arr)

	return
}

func (aa *IntArr) UnmarshalJSON(b []byte) error {
	if aa.arr == nil {
		aa.arr = make([]int, 0)
	}

	if err := json.UnmarshalUseNumber(b, &aa.arr); err != nil {
		return err
	}

	return nil
}

func (aa *IntArr) UnmarshalValue(val interface{}) error {

	switch val.(type) {
	case string, []byte:
		return json.UnmarshalUseNumber(jconv.Bytes(val), &aa.arr)
	default:
		aa.arr = jconv.SliceInt(val)
	}
	return nil
}

func (aa *IntArr) Filter(filter func(index int, val int) bool) *IntArr {

	for i := 0; i < len(aa.arr); {
		if filter(i, aa.arr[i]) {
			aa.arr = append(aa.arr[:i], aa.arr[i+1:]...)
		} else {
			i++
		}
	}
	return aa
}

func (aa *IntArr) FilterNil() *IntArr {

	for i := 0; i < len(aa.arr); {
		if empty.IsNil(aa.arr[i]) {
			aa.arr = append(aa.arr[:i], aa.arr[i+1:]...)
		} else {
			i++
		}
	}
	return aa
}

func (aa *IntArr) Walk(f func(val int) int) *IntArr {

	for i, v := range aa.arr {
		aa.arr[i] = f(v)
	}

	return aa
}

func (aa *IntArr) IsEmpty() bool {
	return aa.Len() == 0
}

// DeepCopy implements interface for deep copy of current type.
func (aa *IntArr) DeepCopy() interface{} {
	newSlice := make([]int, len(aa.arr))
	copy(newSlice, aa.arr)
	return &IntArr{arr: newSlice}

}
func (aa *IntArr) pos(val int) int {
	if len(aa.arr) == 0 {
		return -1
	}
	result := -1
	for index, v := range aa.arr {
		if v == val {
			result = index
			break
		}
	}
	return result
}
func (aa *IntArr) removeByIndex(index int) bool {
	if index < 0 || index >= len(aa.arr) {
		return false
	}
	if index == 0 {
		aa.arr = aa.arr[1:]
		return true
	} else if index == len(aa.arr)-1 {
		aa.arr = aa.arr[:index]
		return true
	}
	aa.arr = append(aa.arr[:index], aa.arr[index+1:]...)
	return true
}
