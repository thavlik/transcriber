package server

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/transcribestreamingservice"
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
	transcript *transcribestreamingservice.MedicalTranscript,
) error {
	var lastTerm string
	text := transcriber.ConvertTranscript(transcript)
	body, err := json.Marshal(&wsMessage{
		Type: "transcript",
		Payload: map[string]interface{}{
			"text": text,
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(text)
	s.broadcast(body)
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
					body, err := json.Marshal(&wsMessage{
						Type: "ref",
						Payload: map[string]interface{}{
							"matched": matched,
							"terms":   ref.Terms,
							"images":  ref.Images,
						},
					})
					if err != nil {
						panic(err)
					}
					s.broadcast(body) // broadcast the reference material to all websockets
					s.useRef(ref)     // mark the reference material as used
				}
			}
		}
	}
	return nil
}
