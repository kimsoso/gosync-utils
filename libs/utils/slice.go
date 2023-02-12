package utils

// 判断是否含有给出值
func SliceHas[T Ordered](slc []T, value T) (out []int) {
	out = make([]int, 0, len(slc))
	for k, v := range slc {
		if v == value {
			out = append(out, k)
		}
	}
	return
}
