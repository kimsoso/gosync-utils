package utils

// 取得map的所有键值
func MapKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// 取得map的所有不重复值
func MapUniqueValues[T comparable](in map[interface{}]T) (out []T) {
	mp := make(map[T]uint8, len(in))
	for _, v := range in {
		mp[v] = 0
	}
	return MapKeys(mp)
}
