package server

import (
	"context"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/define/pkg/diseasecache"
	"github.com/thavlik/transcriber/define/pkg/storage"
	"go.uber.org/zap"
)

func Entry(
	ctx context.Context,
	serverOpts *base.ServerOptions,
	storage storage.Storage,
	diseaseCache diseasecache.DiseaseCache,
	openAISecretKey string,
	log *zap.Logger,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	s := NewServer(
		ctx,
		storage,
		diseaseCache,
		gpt3.NewClient(
			openAISecretKey,
			gpt3.WithDefaultEngine(gpt3.TextDavinci003Engine),
		),
		log,
	)
	defer s.ShutDown()

	s.spawn(func() {
		base.RunMetrics(
			ctx,
			serverOpts.MetricsPort,
			log,
		)
	})
	base.SignalReady(log)

	return s.ListenAndServe(
		serverOpts.Port,
	)
}
