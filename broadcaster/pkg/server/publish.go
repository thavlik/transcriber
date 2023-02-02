package server

import "go.uber.org/zap"

const channelName = "bc"

// publish sends a message to all connected clients across all servers
func (s *Server) publish(msg string) {
	if s.redisClient == nil {
		// no redis, broadcast locally
		go s.broadcastLocal(s.ctx, []byte(msg))
		return
	}
	// we will receive the message from redis, so don't broadcast locally
	if _, err := s.redisClient.Publish(
		s.ctx,
		channelName,
		msg,
	).Result(); err != nil {
		s.log.Error("failed to publish on redis", zap.Error(err))
		// TODO: should we panic?
		return
	}
}
