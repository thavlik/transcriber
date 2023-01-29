package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/transcribestreamingservice"
	"github.com/thavlik/transcriber/pkg/comprehend"
	"github.com/thavlik/transcriber/pkg/refmat"
	"github.com/thavlik/transcriber/pkg/transcriber"

	"go.uber.org/zap"
)

func (s *server) checkTerm(term string) []*refmat.ReferenceMaterial {
	term = strings.ToLower(term)
	refs, ok := s.refs[term]
	if ok {
		return refs
	}
	return nil
}

func (s *server) handleTranscript(
	ctx context.Context,
	transcript *transcribestreamingservice.MedicalTranscript,
) error {
	text := transcriber.ConvertTranscript(transcript)
	go s.broadcastMessage(
		ctx,
		"transcript",
		map[string]interface{}{
			"text": text,
		})
	go func() {
		entities, err := comprehend.Comprehend(
			ctx,
			text,
			//[]string{"OTHER"},
			nil,
			s.log,
		)
		if err != nil {
			s.log.Error("comprehend error", zap.Error(err))
			return
		}
		if len(entities) == 0 {
			return
		}
		s.broadcastMessage(
			ctx,
			"keyterms",
			map[string]interface{}{
				"entities": entities,
			})
		first := entities[0]
		top := first
		for _, entity := range entities[1:] {
			if entity.Score > top.Score {
				top = entity
			}
		}
		s.log.Debug("comprehended entities",
			zap.Int("count", len(entities)),
			zap.String("first.Text", first.Text),
			zap.Float64("first.Score", first.Score),
			zap.String("top.Text", top.Text),
			zap.Float64("top.Score", top.Score),
		)
	}()

	var lastTerm string
	fmt.Println(text)
	for _, result := range transcript.Results {
		for _, alt := range result.Alternatives {
			for _, item := range alt.Items {
				if aws.StringValue(item.Type) == "punctuation" {
					continue
				}
				term := aws.StringValue(item.Content)
				if term == "" {
					continue
				}
				matched := term
				refs := s.checkTerm(term)
				if lastTerm != "" {
					if refs == nil || s.areRefsUsed(refs) {
						// see if a compound word will match a reference material
						matched = fmt.Sprintf("%s %s", lastTerm, term)
						refs = s.checkTerm(matched)
						//if ref == nil {
						//	s.log.Debug("no reference material found",
						//		zap.String("search", matched))
						//}
					}
				}
				lastTerm = term
				for _, ref := range refs {
					if s.isRefUsed(ref) {
						continue
					}
					s.log.Debug("detected reference material",
						zap.String("matched", matched),
						zap.Strings("terms", ref.Terms))
					// broadcast the reference material to all websockets
					s.broadcastMessage(
						ctx,
						"ref",
						map[string]interface{}{
							"matched": ref.Terms[0],
							"terms":   ref.Terms,
							"images":  ref.Images,
						})
					s.useRef(ref) // mark the reference material as used
				}
			}
		}
	}
	return nil
}
