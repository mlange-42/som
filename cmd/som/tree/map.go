package tree

type orderedMap[T any] struct {
	items map[string]T
	keys  []string
}

func (m *orderedMap[T]) Add(key string, value T) {
	m.items[key] = value
	m.keys = append(m.keys, key)
}

func (m *orderedMap[T]) Get(key string) (T, bool) {
	value, ok := m.items[key]
	return value, ok
}

func (m *orderedMap[T]) Keys() []string {
	return m.keys
}
