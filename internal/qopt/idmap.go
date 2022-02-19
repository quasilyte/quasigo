package qopt

type idMap struct {
	kv []uint16
}

func (s *idMap) Reset() {
	s.kv = s.kv[:0]
}

func (s *idMap) Find(id int) int {
	for i, x := range s.kv {
		if int(x>>8) == id {
			return i
		}
	}
	return -1
}

func (s *idMap) GetValue(index int) uint8 {
	i := uint(index)
	if i < uint(len(s.kv)) {
		pair := s.kv[i]
		return uint8(pair & 0xff)
	}
	return 0
}

func (s *idMap) RemoveAt(index int) {
	s.kv[index] = s.kv[len(s.kv)-1]
	s.kv = s.kv[:len(s.kv)-1]
}

func (s *idMap) Push(id, val int) {
	s.kv = append(s.kv, uint16(((id<<8)&0xff00)|(val&0xff)))
}

func (s *idMap) Add(id, val int) bool {
	for _, x := range s.kv {
		if int(x>>8) == id {
			return false
		}
	}
	s.Push(id, val)
	return true
}

func (s *idMap) Contains(id int) bool {
	return s.Find(id) != -1
}

func (s *idMap) Remove(id int) {
	s.kv = s.kv[:0]
	for _, x := range s.kv {
		if int(x>>8) == id {
			continue
		}
		s.kv = append(s.kv, x)
	}
}
