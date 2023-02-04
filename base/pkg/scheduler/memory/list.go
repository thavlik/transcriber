package memory

func (s *memoryScheduler) List() ([]string, error) {
	s.l.Lock()
	defer s.l.Unlock()
	n := len(s.set)
	projectIDs := make([]string, n)
	i := 0
	for k := range s.set {
		projectIDs[i] = k
		i++
	}
	return projectIDs, nil
}
