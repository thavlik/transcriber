package memory

func (s *memoryScheduler) Add(projectIDs ...string) error {
	s.l.Lock()
	for _, projectID := range projectIDs {
		s.set[projectID] = struct{}{}
		notify := s.notify
		for _, ch := range notify {
			go func(ch chan<- struct{}) { ch <- struct{}{} }(ch)
		}
	}
	s.l.Unlock()
	return nil
}
