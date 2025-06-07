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
)

type AnyArr struct {
	arr []interface{}
}

func NewAnyArr() *AnyArr {
	return NewAnyArrSize(0, 0)
}
func NewAnyArrSize(size int, cap int) *AnyArr {
	return &AnyArr{
		arr: make([]interface{}, size, cap),
	}
}

// NewAnyArrCopy 返回另一个数组的拷贝
func NewAnyArrCopy(arr []interface{}) *AnyArr {
	nAnyArr := make([]interface{}, len(arr))
	copy(nAnyArr, arr)
	return &AnyArr{
		arr: nAnyArr,
	}
}

// Chunk 将数组拆分成多个size长度的区块
func (aa *AnyArr) Chunk(size int) [][]interface{} {
	if size < 1 {
		return nil
	}

	length := len(aa.arr)
	chunks := int(math.Ceil(float64(length) / float64(size)))
	var n [][]interface{}
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

// Compact 移除假值
func (aa *AnyArr) Compact() *AnyArr {

	for i := 0; i < len(aa.arr); {
		if empty.IsEmpty(aa.arr[i]) {
			aa.arr = append(aa.arr[:i], aa.arr[i+1:]...)
		} else {
			i++
		}
	}

	return aa
}

func (aa *AnyArr) Pos(val interface{}) (index int) {
	return aa.Search(val)
}

func (aa *AnyArr) Index(index int) (val interface{}) {
	return aa.Get(index)
}

func (aa *AnyArr) Get(index int) (val interface{}) {

	if index < 0 || index >= len(aa.arr) {
		return
	}
	val = aa.arr[index]

	return
}

func (aa *AnyArr) Set(index int, val interface{}) error {

	if index < 0 || index >= len(aa.arr) {
		return jerr.WithCodeF(jcode.NewErrCode(jcode.ParamErr), "index %d out of arr range %d", index, len(aa.arr))
	}
	aa.arr[index] = val

	return nil
}

func (aa *AnyArr) SetFrom(arr []interface{}) *AnyArr {

	aa.arr = arr

	return aa
}

func (aa *AnyArr) Sort(compare func(v1, v2 interface{}) bool) *AnyArr {

	sort.Slice(aa.arr, func(i, j int) bool {
		return compare(aa.arr[i], aa.arr[j])
	})

	return aa
}

func (aa *AnyArr) InsertBefore(index int, vals ...interface{}) error {

	if index < 0 || index >= len(aa.arr) {
		return jerr.WithCodeF(jcode.NewErrCode(jcode.ParamErr), "index %d out of arr range %d", index, len(aa.arr))
	}
	// 取尾
	tail := append([]interface{}{}, aa.arr[index:]...)
	aa.arr = append(aa.arr[0:index], vals...)
	// 放尾
	aa.arr = append(aa.arr, tail...)

	return nil
}

func (aa *AnyArr) InsertAfter(index int, vals ...interface{}) error {

	if index < 0 || index >= len(aa.arr) {
		return jerr.WithCodeF(jcode.NewErrCode(jcode.ParamErr), "index %d out of arr range %d", index, len(aa.arr))
	}
	tail := append([]interface{}{}, aa.arr[index+1:]...)
	aa.arr = append(aa.arr[0:index+1], vals...)
	aa.arr = append(aa.arr, tail...)

	return nil
}

func (aa *AnyArr) Remove(index int) {

	aa.removeByIndex(index)

}

func (aa *AnyArr) RemoveVal(val interface{}) bool {

	if i := aa.pos(val); i != -1 {
		aa.removeByIndex(i)
		return true
	}

	return false
}

func (aa *AnyArr) RemoveVals(vals ...interface{}) {

	for _, val := range vals {
		if i := aa.pos(val); i != -1 {
			aa.removeByIndex(i)
		}
	}

}

func (aa *AnyArr) PushLeft(val ...interface{}) *AnyArr {

	aa.arr = append(val, aa.arr...)

	return aa
}

func (aa *AnyArr) PushRight(val ...interface{}) *AnyArr {

	aa.arr = append(aa.arr, val...)

	return aa
}

func (aa *AnyArr) PopLeft() (val interface{}) {

	if len(aa.arr) == 0 {
		return nil
	}
	val = aa.arr[0]
	aa.arr = aa.arr[1:]

	return val
}

func (aa *AnyArr) PopRight() (val interface{}) {

	index := len(aa.arr) - 1
	if index < 0 {
		return nil
	}
	val = aa.arr[index]
	aa.arr = aa.arr[:index]

	return val
}

func (aa *AnyArr) PopLefts(size int) []interface{} {

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

func (aa *AnyArr) PopRights(size int) []interface{} {

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

func (aa *AnyArr) Append(val ...interface{}) *AnyArr {
	aa.PushRight(val...)
	return aa
}

func (aa *AnyArr) Len() int {

	length := len(aa.arr)

	return length
}

func (aa *AnyArr) Slice() []interface{} {

	arr := make([]interface{}, len(aa.arr))
	copy(arr, aa.arr)

	return arr
}

func (aa *AnyArr) Interfaces() []interface{} {
	return aa.Slice()
}

func (aa *AnyArr) Clone() (newAnyArr *AnyArr) {

	arr := make([]interface{}, len(aa.arr))
	copy(arr, aa.arr)

	return &AnyArr{arr: arr}
}

func (aa *AnyArr) Clear() {

	if len(aa.arr) > 0 {
		aa.arr = make([]interface{}, 0)
	}

}

func (aa *AnyArr) Contain(val interface{}) bool {
	return aa.Search(val) != -1
}

func (aa *AnyArr) Search(val interface{}) (pos int) {

	pos = aa.pos(val)

	return
}

func (aa *AnyArr) Unique() *AnyArr {

	if len(aa.arr) == 0 {
		return aa
	}
	var (
		ok     bool
		temp   interface{}
		set    = make(map[interface{}]struct{})
		unique = make([]interface{}, 0, len(aa.arr))
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

func (aa *AnyArr) LockFunc(f func(arr []interface{})) *AnyArr {

	f(aa.arr)

	return aa
}

func (aa *AnyArr) RLockFunc(f func(arr []interface{})) *AnyArr {

	f(aa.arr)

	return aa
}

func (aa *AnyArr) Merge(arr ...interface{}) *AnyArr {
	return aa.Append(arr)
}

func (aa *AnyArr) Rand() (val interface{}) {

	if len(aa.arr) == 0 {
		return nil
	}
	val = aa.arr[rand.Intn(len(aa.arr))]

	return
}

func (aa *AnyArr) Rands(size int) []interface{} {

	if size <= 0 || len(aa.arr) == 0 {
		return nil
	}
	arr := make([]interface{}, size)
	for i := 0; i < size; i++ {
		arr[i] = aa.arr[rand.Intn(len(aa.arr))]
	}

	return arr
}

func (aa *AnyArr) Reverse() *AnyArr {

	for i, j := 0, len(aa.arr)-1; i < j; i, j = i+1, j-1 {
		aa.arr[i], aa.arr[j] = aa.arr[j], aa.arr[i]
	}

	return aa
}

func (aa *AnyArr) Join(sep string) string {

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

func (aa *AnyArr) Count(val interface{}) (count int) {

	for _, v := range aa.arr {
		if v == val {
			count += 1
		}
	}

	return
}

func (aa *AnyArr) String() string {

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

func (aa *AnyArr) MarshalJSON() (b []byte, err error) {

	b, err = json.Marshal(aa.arr)

	return
}

func (aa *AnyArr) UnmarshalJSON(b []byte) error {
	if aa.arr == nil {
		aa.arr = make([]interface{}, 0)
	}

	if err := json.UnmarshalUseNumber(b, &aa.arr); err != nil {
		return err
	}

	return nil
}

func (aa *AnyArr) UnmarshalValue(val interface{}) error {

	switch val.(type) {
	case string, []byte:
		return json.UnmarshalUseNumber(jconv.Bytes(val), &aa.arr)
	default:
		aa.arr = jconv.SliceAny(val)
	}

	return nil
}

func (aa *AnyArr) Filter(filter func(index int, val interface{}) bool) *AnyArr {

	for i := 0; i < len(aa.arr); {
		if filter(i, aa.arr[i]) {
			aa.arr = append(aa.arr[:i], aa.arr[i+1:]...)
		} else {
			i++
		}
	}
	return aa
}

func (aa *AnyArr) FilterNil() *AnyArr {

	for i := 0; i < len(aa.arr); {
		if empty.IsNil(aa.arr[i]) {
			aa.arr = append(aa.arr[:i], aa.arr[i+1:]...)
		} else {
			i++
		}
	}
	return aa
}

func (aa *AnyArr) Walk(f func(val interface{}) interface{}) *AnyArr {

	for i, v := range aa.arr {
		aa.arr[i] = f(v)
	}

	return aa
}

func (aa *AnyArr) IsEmpty() bool {
	return aa.Len() == 0
}

// DeepCopy implements interface for deep copy of current type.
func (aa *AnyArr) DeepCopy() interface{} {

	newSlice := make([]interface{}, len(aa.arr))
	for i, v := range aa.arr {
		newSlice[i] = deepcopy.Copy(v)
	}

	return &AnyArr{arr: newSlice}
}
func (aa *AnyArr) pos(val interface{}) int {
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
func (aa *AnyArr) removeByIndex(index int) bool {
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
