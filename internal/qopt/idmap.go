package qopt

type idMap struct {
	kv []uint16
}

func (s *idMap) Reset() {
	s.kv = s.kv[:0]
}

func (s *idMap) FindValue(id uint8) (uint8, bool) {
	for _, x := range s.kv {
		if uint8(x>>8) == id {
			return uint8(x), true
		}
	}
	return 0, false
}

func (s *idMap) FindIndex(id uint8) int {
	for i, x := range s.kv {
		if uint8(x>>8) == id {
			return i
		}
	}
	return -1
}

func (s *idMap) GetValue(index int) uint8 {
	i := uint(index)
	if i < uint(len(s.kv)) {
		return uint8(s.kv[i])
	}
	return 0
}

func (s *idMap) RemoveAt(index int) {
	s.kv[index] = s.kv[len(s.kv)-1]
	s.kv = s.kv[:len(s.kv)-1]
}

func (s *idMap) UpdateValueAt(index, val int) {
	s.kv[index] = (s.kv[index] & 0xff00) | uint16(val)
}

func (s *idMap) Push(id, val uint8) {
	s.kv = append(s.kv, ((uint16(id)<<8)&0xff00)|uint16(val))
}

func (s *idMap) Add(id, val uint8) bool {
	for _, x := range s.kv {
		if uint8(x>>8) == id {
			return false
		}
	}
	s.Push(id, val)
	return true
}

func (s *idMap) Contains(id uint8) bool {
	return s.FindIndex(id) != -1
}

func (s *idMap) Remove(id uint8) {
	kv := s.kv[:0]
	for _, x := range s.kv {
		if uint8(x>>8) == id {
			continue
		}
		kv = append(kv, x)
	}
	s.kv = kv
}
