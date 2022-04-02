package qcompile

import "fmt"

type scope struct {
	vars   []string
	depths []int
}

func (s *scope) NumLiveVars() int {
	return len(s.vars)
}

func (s *scope) NumLevels() int {
	return len(s.depths)
}

func (s *scope) Enter() {
	s.depths = append(s.depths, 0)
}

func (s *scope) Leave() int {
	depth := s.depths[len(s.depths)-1]
	s.depths = s.depths[:len(s.depths)-1]
	s.vars = s.vars[:len(s.vars)-depth]
	return depth
}

func (s *scope) PushVar(name string) {
	s.depths[len(s.depths)-1]++
	s.vars = append(s.vars, name)
}

func (s *scope) LookupInCurrent(name string) int {
	num := s.depths[len(s.depths)-1]
	for i, x := range s.vars[len(s.vars)-num:] {
		if x == name {
			fmt.Printf("vars=%q num=%d name=%q result=%d\n", s.vars, num, name, len(s.vars)-num+i)
			return len(s.vars) - num + i
		}
	}
	return -1
}

func (s *scope) Lookup(name string) int {
	for i := len(s.vars) - 1; i >= 0; i-- {
		if s.vars[i] == name {
			return i
		}
	}
	return -1
}
