package server

import "go.uber.org/zap"

const channelName = "bc"

// publish sends a message to all connected clients across all servers
func (s *Server) publish(msg []byte) {
	// we will receive the message from redis, so don't broadcast locally
	if err := s.pub.Publish(
		s.ctx,
		channelName,
		msg,
	); err != nil {
		s.log.Error("failed to publish", zap.Error(err))
		// TODO: should we panic?
		return
	}
}
