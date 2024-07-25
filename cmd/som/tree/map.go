package tree

type orderedMap[T any] struct {
	items map[string]T
	keys  []string
}

func newOrderedMap[T any]() orderedMap[T] {
	return orderedMap[T]{
		items: make(map[string]T),
		keys:  make([]string, 0),
	}
}

func (m *orderedMap[T]) Set(key string, value T) {
	if _, ok := m.items[key]; !ok {
		m.keys = append(m.keys, key)
	}
	m.items[key] = value
}

func (m *orderedMap[T]) Get(key string) (T, bool) {
	value, ok := m.items[key]
	return value, ok
}

func (m *orderedMap[T]) Keys() []string {
	return m.keys
}

func (m *orderedMap[T]) Len() int {
	return len(m.keys)
}
