package qopt

type sliceMap[Key comparable, Value any] struct {
	pairs []slotMapPair[Key, Value]
}

type slotMapPair[Key comparable, Value any] struct {
	key Key
	val Value
}

func (m *sliceMap[Key, Value]) Reset() {
	m.pairs = m.pairs[:0]
}

func (m *sliceMap[Key, Value]) FindValue(key Key) (v Value, found bool) {
	for _, p := range m.pairs {
		if p.key == key {
			return p.val, true
		}
	}
	return v, false
}

func (m *sliceMap[Key, Value]) FindIndex(key Key) int {
	for i, p := range m.pairs {
		if p.key == key {
			return i
		}
	}
	return -1
}

func (m *sliceMap[Key, Value]) GetValue(index int) Value {
	i := uint(index)
	if i < uint(len(m.pairs)) {
		return m.pairs[i].val
	}
	var v Value
	return v
}

func (m *sliceMap[Key, Value]) RemoveAt(index int) {
	m.pairs[index] = m.pairs[len(m.pairs)-1]
	m.pairs = m.pairs[:len(m.pairs)-1]
}

func (m *sliceMap[Key, Value]) UpdateValueAt(index int, val Value) {
	m.pairs[index].val = val
}

func (m *sliceMap[Key, Value]) Push(key Key, val Value) {
	m.pairs = append(m.pairs, slotMapPair[Key, Value]{key: key, val: val})
}

func (m *sliceMap[Key, Value]) Add(key Key, val Value) bool {
	for _, p := range m.pairs {
		if p.key == key {
			return false
		}
	}
	m.Push(key, val)
	return true
}

func (m *sliceMap[Key, Value]) Contains(key Key) bool {
	return m.FindIndex(key) != -1
}

func (m *sliceMap[Key, Value]) Remove(key Key) {
	i := m.FindIndex(key)
	if i != -1 {
		m.RemoveAt(i)
	}
}
