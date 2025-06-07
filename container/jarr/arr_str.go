package jarr

import (
	"bytes"
	"github.com/e7coding/coding-common/errs/jerr"
	"math"
	"sort"
	"strings"

	"github.com/e7coding/coding-common/internal/json"
	"github.com/e7coding/coding-common/jutil/jconv"
	"github.com/e7coding/coding-common/jutil/jrand"
	"github.com/e7coding/coding-common/text/jstr"
)

// StrArr 字符串数组
type StrArr struct {
	array []string
}

// NewStrArr 创建并返回一个空字符串数组
func NewStrArr() *StrArr {
	return NewStrArrSize(0, 0)
}

// NewStrArrSize 根据指定的长度和容量创建并返回字符串数组
func NewStrArrSize(size int, cap int) *StrArr {
	return &StrArr{
		array: make([]string, size, cap),
	}
}

// NewStrArrFrom 从给定切片创建并返回字符串数组
func NewStrArrFrom(array []string) *StrArr {
	return &StrArr{
		array: array,
	}
}

// NewStrArrCopy 复制给定切片并创建新的字符串数组
func NewStrArrCopy(array []string) *StrArr {
	newArray := make([]string, len(array))
	copy(newArray, array)
	return &StrArr{
		array: newArray,
	}
}

// At 根据索引获取元素，越界时返回空字符串
func (sa *StrArr) At(index int) (value string) {
	value = sa.Get(index)
	return
}

// Get 根据索引获取元素和存在标志
func (sa *StrArr) Get(index int) (value string) {
	if index < 0 || index >= len(sa.array) {
		return ""
	}
	return sa.array[index]
}

// SetAt 在指定索引设置值，越界时返回错误
func (sa *StrArr) SetAt(index int, value string) error {
	if index < 0 || index >= len(sa.array) {
		return jerr.WithMsgF("index %d out of array range %d", index, len(sa.array))
	}
	sa.array[index] = value
	return nil
}

// Update 用新切片替换底层数组
func (sa *StrArr) Update(array []string) *StrArr {
	sa.array = array
	return sa
}

// Replace 用给定切片前段替换当前数组
func (sa *StrArr) Replace(array []string) *StrArr {
	maxs := len(array)
	if maxs > len(sa.array) {
		maxs = len(sa.array)
	}
	for i := 0; i < maxs; i++ {
		sa.array[i] = array[i]
	}
	return sa
}

// Sum 将所有元素按整数转换后求和
func (sa *StrArr) Sum() (sum int) {
	for _, v := range sa.array {
		sum += jconv.Int(v)
	}
	return
}

// Sort 对数组进行升序（默认）或降序排序
func (sa *StrArr) Sort(reverse ...bool) *StrArr {
	if len(reverse) > 0 && reverse[0] {
		sort.Slice(sa.array, func(i, j int) bool {
			return strings.Compare(sa.array[i], sa.array[j]) >= 0
		})
	} else {
		sort.Strings(sa.array)
	}
	return sa
}

// SortFunc 使用自定义比较函数排序
func (sa *StrArr) SortFunc(less func(v1, v2 string) bool) *StrArr {

	sort.Slice(sa.array, func(i, j int) bool {
		return less(sa.array[i], sa.array[j])
	})
	return sa
}

func (sa *StrArr) ByFunc(f func(array []string)) *StrArr {
	f(sa.array)
	return sa
}

// AppendBefore 在指定索引前插入元素
func (sa *StrArr) AppendBefore(index int, values ...string) error {

	if index < 0 || index >= len(sa.array) {
		return jerr.WithMsgF("index %d out of array range %d", index, len(sa.array))
	}
	rear := append([]string{}, sa.array[index:]...)
	sa.array = append(sa.array[0:index], values...)
	sa.array = append(sa.array, rear...)
	return nil
}

// AppendAfter 在指定索引后插入元素
func (sa *StrArr) AppendAfter(index int, values ...string) error {

	if index < 0 || index >= len(sa.array) {
		return jerr.WithMsgF("index %d out of array range %d", index, len(sa.array))
	}
	rear := append([]string{}, sa.array[index+1:]...)
	sa.array = append(sa.array[0:index+1], values...)
	sa.array = append(sa.array, rear...)
	return nil
}

// RemoveAt 删除指定索引的元素并返回该值，不存在时 found=false
func (sa *StrArr) RemoveAt(index int) (value string, found bool) {

	return sa.doRemove(index)
}

func (sa *StrArr) doRemove(index int) (val string, found bool) {
	if index < 0 || index >= len(sa.array) {
		return "", false
	}
	// Determine array boundaries when deleting to improve deletion efficiency.
	if index == 0 {
		val = sa.array[0]
		sa.array = sa.array[1:]
		return val, true
	} else if index == len(sa.array)-1 {
		val = sa.array[index]
		sa.array = sa.array[:index]
		return val, true
	}
	// If it is a non-boundary delete,
	// it will involve the creation of an array,
	// then the deletion is less efficient.
	val = sa.array[index]
	sa.array = append(sa.array[:index], sa.array[index+1:]...)
	return val, true
}

// RemoveByVal 根据值删除第一个匹配元素，返回是否删除成功.
func (sa *StrArr) RemoveByVal(value string) bool {
	if i := sa.IndexOf(value); i != -1 {
		_, found := sa.RemoveAt(i)
		return found
	}
	return false
}

// RemoveByVals 批量根据值删除元素
func (sa *StrArr) RemoveByVals(values ...string) {

	for _, value := range values {
		if i := sa.index(value); i != -1 {
			sa.doRemove(i)
		}
	}
}

// Prepend 向数组左侧添加元素
func (sa *StrArr) Prepend(value ...string) *StrArr {
	sa.array = append(value, sa.array...)
	return sa
}

// Append 向数组右侧添加元素
func (sa *StrArr) Append(vals ...string) *StrArr {
	sa.array = append(sa.array, vals...)
	return sa
}

// PopFront 弹出并返回左侧第一个元素，空时 found=false
func (sa *StrArr) PopFront() (string, bool) {
	if len(sa.array) == 0 {
		return "", false
	}
	v := sa.array[0]
	sa.array = sa.array[1:]
	return v, true
}

// PopRight pops and returns an item from the end of array.
// Note that if the array is empty, the `found` is false.
func (sa *StrArr) PopRight() (value string, found bool) {

	index := len(sa.array) - 1
	if index < 0 {
		return "", false
	}
	value = sa.array[index]
	sa.array = sa.array[:index]
	return value, true
}

// PopBack 弹出并返回右侧最后一个元素，空时 found=false
func (sa *StrArr) PopBack() (string, bool) {
	idx := len(sa.array) - 1
	if idx < 0 {
		return "", false
	}
	v := sa.array[idx]
	sa.array = sa.array[:idx]
	return v, true
}

// PopRandom 随机弹出并返回一个元素
func (sa *StrArr) PopRandom() (string, bool) {
	if len(sa.array) == 0 {
		return "", false
	}
	return sa.doRemove(jrand.Intn(len(sa.array)))
}

// PopRandoms 随机弹出并返回多个元素
func (sa *StrArr) PopRandoms(size int) []string {
	if size <= 0 || len(sa.array) == 0 {
		return nil
	}
	if size >= len(sa.array) {
		size = len(sa.array)
	}
	out := make([]string, size)
	for i := 0; i < size; i++ {
		out[i], _ = sa.doRemove(jrand.Intn(len(sa.array)))
	}
	return out
}

// PopFronts 弹出并返回左侧多个元素
func (sa *StrArr) PopFronts(size int) []string {
	if size <= 0 || len(sa.array) == 0 {
		return nil
	}
	if size >= len(sa.array) {
		out := sa.array
		sa.array = sa.array[:0]
		return out
	}
	out := sa.array[:size]
	sa.array = sa.array[size:]
	return out
}

// PopBacks 弹出并返回右侧多个元素
func (sa *StrArr) PopBacks(size int) []string {
	if size <= 0 || len(sa.array) == 0 {
		return nil
	}
	start := len(sa.array) - size
	if start <= 0 {
		out := sa.array
		sa.array = sa.array[:0]
		return out
	}
	out := sa.array[start:]
	sa.array = sa.array[:start]
	return out
}

// SliceRange 按区间提取子切片，超界自动剪裁
func (sa *StrArr) SliceRange(start int, end ...int) []string {
	e := len(sa.array)
	if len(end) > 0 && end[0] < e {
		e = end[0]
	}
	if start > e {
		return nil
	}
	if start < 0 {
		start = 0
	}
	return sa.array[start:e]
}

// SubSlice 按偏移量和长度提取子切片，支持负数偏移和长度
func (sa *StrArr) SubSlice(offset int, length ...int) []string {
	size := len(sa.array)
	if len(length) > 0 {
		size = length[0]
	}
	if offset > len(sa.array) {
		return nil
	}
	if offset < 0 {
		offset += len(sa.array)
		if offset < 0 {
			return nil
		}
	}
	if size < 0 {
		offset += size
		size = -size
		if offset < 0 {
			return nil
		}
	}
	end := offset + size
	if end > len(sa.array) {
		end = len(sa.array)
	}
	return sa.array[offset:end]
}

// Len 返回数组长度
func (sa *StrArr) Len() int {
	return len(sa.array)
}

// Slice 返回底层切片（不拷贝）
func (sa *StrArr) Slice() []string {
	return sa.array
}

// ToInterfaces 将字符串数组转换为 []interface{}
func (sa *StrArr) ToInterfaces() []interface{} {
	out := make([]interface{}, len(sa.array))
	for i, v := range sa.array {
		out[i] = v
	}
	return out
}

// Clone 深拷贝当前数组并返回新实例
func (sa *StrArr) Clone() *StrArr {
	newArr := make([]string, len(sa.array))
	copy(newArr, sa.array)
	return NewStrArrFrom(newArr)
}

// Clear 清空所有元素
func (sa *StrArr) Clear() *StrArr {
	sa.array = make([]string, 0)
	return sa
}

// Contains 检查是否包含指定值
func (sa *StrArr) Contains(val string) bool {
	return sa.IndexOf(val) != -1
}

// ContainsIgnoreCase 忽略大小写检查是否包含指定值
func (sa *StrArr) ContainsIgnoreCase(val string) bool {
	for _, v := range sa.array {
		if strings.EqualFold(v, val) {
			return true
		}
	}
	return false
}

// IndexOf 返回指定值的索引，未找到返回 -1
func (sa *StrArr) IndexOf(val string) int {
	return sa.index(val)
}

// index 不加锁情况下查找值的索引
func (sa *StrArr) index(val string) int {
	for i, v := range sa.array {
		if v == val {
			return i
		}
	}
	return -1
}

// Uniq 去重并保留原有顺序
func (sa *StrArr) Uniq() *StrArr {
	seen := make(map[string]struct{}, len(sa.array))
	out := make([]string, 0, len(sa.array))
	for _, v := range sa.array {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			out = append(out, v)
		}
	}
	sa.array = out
	return sa
}

// Merge 将任意类型切片合并到当前数组
func (sa *StrArr) Merge(src interface{}) *StrArr {
	return sa.Append(jconv.Strings(src)...)
}

// Fill 从指定索引开始，用相同值填充若干元素
func (sa *StrArr) Fill(start, num int, val string) error {
	if start < 0 || start > len(sa.array) {
		return jerr.WithMsgF("index %d out of array range %d", start, len(sa.array))
	}
	for i := start; i < start+num; i++ {
		if i >= len(sa.array) {
			sa.array = append(sa.array, val)
		} else {
			sa.array[i] = val
		}
	}
	return nil
}

// Chunk 按固定大小拆分成多个切片
func (sa *StrArr) Chunk(size int) [][]string {
	if size < 1 {
		return nil
	}
	length := len(sa.array)
	count := int(math.Ceil(float64(length) / float64(size)))
	out := make([][]string, 0, count)
	for i := 0; i < count; i++ {
		start := i * size
		end := start + size
		if end > length {
			end = length
		}
		out = append(out, sa.array[start:end])
	}
	return out
}

// Pad 在左右两端填充元素至指定长度
func (sa *StrArr) Pad(size int, val string) *StrArr {
	if size == 0 || (size > 0 && size <= len(sa.array)) || (size < 0 && -size <= len(sa.array)) {
		return sa
	}
	n := size
	if n < 0 {
		n = -n
	}
	n -= len(sa.array)
	tmp := make([]string, n)
	for i := range tmp {
		tmp[i] = val
	}
	if size > 0 {
		sa.array = append(sa.array, tmp...)
	} else {
		sa.array = append(tmp, sa.array...)
	}
	return sa
}

// Random 随机返回一个元素（不删除）
func (sa *StrArr) Random() (string, bool) {
	if len(sa.array) == 0 {
		return "", false
	}
	return sa.array[jrand.Intn(len(sa.array))], true
}

// Randoms 随机返回多个元素（不删除）
func (sa *StrArr) Randoms(size int) []string {
	if size <= 0 || len(sa.array) == 0 {
		return nil
	}
	out := make([]string, size)
	for i := 0; i < size; i++ {
		out[i] = sa.array[jrand.Intn(len(sa.array))]
	}
	return out
}

// Shuffle 随机打乱数组
func (sa *StrArr) Shuffle() *StrArr {
	for i, j := range jrand.Perm(len(sa.array)) {
		sa.array[i], sa.array[j] = sa.array[j], sa.array[i]
	}
	return sa
}

// Reverse 反转数组顺序
func (sa *StrArr) Reverse() *StrArr {
	for i, j := 0, len(sa.array)-1; i < j; i, j = i+1, j-1 {
		sa.array[i], sa.array[j] = sa.array[j], sa.array[i]
	}
	return sa
}

// Join 用给定分隔符拼接成字符串
func (sa *StrArr) Join(glue string) string {
	if len(sa.array) == 0 {
		return ""
	}
	buf := bytes.NewBuffer(nil)
	for i, v := range sa.array {
		buf.WriteString(v)
		if i < len(sa.array)-1 {
			buf.WriteString(glue)
		}
	}
	return buf.String()
}

// Count 统计每个值出现的次数
func (sa *StrArr) Count() map[string]int {
	m := make(map[string]int, len(sa.array))
	for _, v := range sa.array {
		m[v]++
	}
	return m
}

// ForEach 升序遍历，回调返回 false 时提前停止
func (sa *StrArr) ForEach(f func(idx int, val string) bool) {
	for i, v := range sa.array {
		if !f(i, v) {
			break
		}
	}
}

// ForEachAsc 同 ForEach
func (sa *StrArr) ForEachAsc(f func(idx int, val string) bool) {
	sa.ForEach(f)
}

// ForEachDesc 降序遍历，回调返回 false 时提前停止
func (sa *StrArr) ForEachDesc(f func(idx int, val string) bool) {
	for i := len(sa.array) - 1; i >= 0; i-- {
		if !f(i, sa.array[i]) {
			break
		}
	}
}

// String 返回类似 json.Marshal 的字符串表示
func (sa *StrArr) String() string {
	if sa == nil {
		return ""
	}
	buf := bytes.NewBuffer(nil)
	buf.WriteByte('[')
	for i, v := range sa.array {
		buf.WriteString(`"` + jstr.QuoteMeta(v, `"\`) + `"`)
		if i < len(sa.array)-1 {
			buf.WriteByte(',')
		}
	}
	buf.WriteByte(']')
	return buf.String()
}

// MarshalJSON 实现 json.Marshaler 接口
func (sa *StrArr) MarshalJSON() ([]byte, error) {
	return json.Marshal(sa.array)
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (sa *StrArr) UnmarshalJSON(b []byte) error {
	if sa.array == nil {
		sa.array = make([]string, 0)
	}
	return json.UnmarshalUseNumber(b, &sa.array)
}

// UnmarshalValue 将多种类型值解析为字符串数组
func (sa *StrArr) UnmarshalValue(val interface{}) error {
	switch val.(type) {
	case string, []byte:
		return json.UnmarshalUseNumber(jconv.Bytes(val), &sa.array)
	default:
		sa.array = jconv.SliceStr(val)
		return nil
	}
}

// Filter 依据回调函数过滤元素，返回剩余数组
func (sa *StrArr) Filter(f func(idx int, val string) bool) *StrArr {
	for i := 0; i < len(sa.array); {
		if f(i, sa.array[i]) {
			sa.array = append(sa.array[:i], sa.array[i+1:]...)
		} else {
			i++
		}
	}
	return sa
}

// FilterEmpty 移除所有空字符串元素
func (sa *StrArr) FilterEmpty() *StrArr {
	return sa.Filter(func(_ int, v string) bool { return v == "" })
}

// Walk 将用户提供的函数“f”应用于数组的每个元素。
func (sa *StrArr) Walk(f func(value string) string) *StrArr {

	for i, v := range sa.array {
		sa.array[i] = f(v)
	}
	return sa
}

// IsEmpty 判断数组是否为空
func (sa *StrArr) IsEmpty() bool {
	return len(sa.array) == 0
}

// DeepCopy 深度复制当前实例并返回新对象
func (sa *StrArr) DeepCopy() interface{} {
	if sa == nil {
		return nil
	}
	newArr := make([]string, len(sa.array))
	copy(newArr, sa.array)
	return NewStrArrFrom(newArr)
}
