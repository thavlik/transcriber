package memory

import "math/rand"

func (s *memoryScheduler) Rand() (string, error) {
	s.l.Lock()
	defer s.l.Unlock()
	n := rand.Intn(len(s.set))
	i := 0
	for k := range s.set {
		if i == n {
			return k, nil
		}
		i++
	}
	return "", nil
}
