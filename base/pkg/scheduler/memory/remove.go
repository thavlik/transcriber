package memory

func (s *memoryScheduler) Remove(projectID string) error {
	s.l.Lock()
	defer s.l.Unlock()
	delete(s.set, projectID)
	return nil
}
