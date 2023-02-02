package server

func (s *Server) pumpMessages(msgs <-chan []byte) {
	for {
		select {
		case <-s.ctx.Done():
			return
		case msg, ok := <-msgs:
			if !ok {
				panic("subscription channel closed unexpectedly")
			}
			s.spawn(func() {
				s.broadcastLocal(
					s.ctx,
					msg,
				)
			})
		}
	}
}
