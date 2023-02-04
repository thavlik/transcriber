package memory

func (s *memoryScheduler) Notify() <-chan struct{} {
	ch := make(chan struct{}, 1)
	s.l.Lock()
	defer s.l.Unlock()
	s.notify = append(s.notify, ch)
	return ch
}
