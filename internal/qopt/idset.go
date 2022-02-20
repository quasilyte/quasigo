package qopt

type idSet struct {
	ids []uint8
}

func (s *idSet) Reset() {
	s.ids = s.ids[:0]
}

func (s *idSet) Find(id int) int {
	for i, x := range s.ids {
		if int(x) == id {
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
func (s *idSet) Push(id int) {
	s.ids = append(s.ids, uint8(id))
}

func (s *idSet) Add(id int) bool {
	for _, x := range s.ids {
		if int(x) == id {
			return false
		}
	}
	s.ids = append(s.ids, uint8(id))
	return true
}

func (s *idSet) Contains(id int) bool {
	return s.Find(id) != -1
}

func (s *idSet) Remove(id int) {
	ids := s.ids[:0]
	for _, x := range s.ids {
		if int(x) == id {
			continue
		}
		ids = append(ids, x)
	}
	s.ids = ids
}
