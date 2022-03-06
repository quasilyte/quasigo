package qopt

type idSet struct {
	ids []uint8
}

func (s *idSet) Reset() {
	s.ids = s.ids[:0]
}

func (s *idSet) Find(id uint8) int {
	for i, x := range s.ids {
		if x == id {
			return i
		}
	}
	return -1
}

func (s *idSet) RemoveAt(index int) {
	s.ids[index] = s.ids[len(s.ids)-1]
	s.ids = s.ids[:len(s.ids)-1]
}

// Push adds an ID without any checks.
//
// This method should only be used if Find/Contains reported
// that there are no such ID inside this set before.
func (s *idSet) Push(id uint8) {
	s.ids = append(s.ids, id)
}

func (s *idSet) Add(id uint8) bool {
	for _, x := range s.ids {
		if x == id {
			return false
		}
	}
	s.ids = append(s.ids, id)
	return true
}

func (s *idSet) Contains(id uint8) bool {
	return s.Find(id) != -1
}

func (s *idSet) Remove(id uint8) {
	ids := s.ids[:0]
	for _, x := range s.ids {
		if x == id {
			continue
		}
		ids = append(ids, x)
	}
	s.ids = ids
}
