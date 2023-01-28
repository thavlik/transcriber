package server

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/transcribestreamingservice"
	"github.com/thavlik/transcriber/pkg/refmat"

	"go.uber.org/zap"
)

func (s *server) checkTerm(term string) *refmat.ReferenceMaterial {
	term = strings.ToLower(term)
	ref, ok := s.refs[term]
	if ok {
		return ref
	}
	return nil
}

func (s *server) handleTranscript(
	transcript *transcribestreamingservice.MedicalTranscript,
) error {
	var lastTerm string
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
				ref := s.checkTerm(term)
				if ref == nil && lastTerm != "" || (ref != nil && s.isRefUsed(ref)) {
					// see if a compound word will match a new reference material
					matched = fmt.Sprintf("%s %s", lastTerm, term)
					ref = s.checkTerm(matched)
					//if ref == nil {
					//	s.log.Debug("no reference material found",
					//		zap.String("search", matched))
					//}
				}
				lastTerm = term
				if ref != nil && !s.isRefUsed(ref) {
					s.log.Debug("found reference material",
						zap.String("matched", matched),
						zap.Strings("terms", ref.Terms))
					body, err := json.Marshal(ref)
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
