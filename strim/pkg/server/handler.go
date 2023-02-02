package server

import (
	"bytes"
	"context"
	"io"
	"log"
	"sync"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/transcriber/pkg/source"
	"github.com/thavlik/transcriber/transcriber/pkg/source/aac"
	flvtag "github.com/yutopp/go-flv/tag"
	"github.com/yutopp/go-rtmp"
	rtmpmsg "github.com/yutopp/go-rtmp/message"
	"go.uber.org/zap"
)

var _ rtmp.Handler = (*Handler)(nil)

type CheckStreamKey func(string) error

// Handler An RTMP connection handler
type Handler struct {
	rtmp.DefaultHandler
	wg             *sync.WaitGroup
	source         *aac.AACSource
	ctx            context.Context
	cancel         context.CancelFunc
	newSource      chan<- source.Source
	checkStreamKey CheckStreamKey
	log            *zap.Logger
}

func NewHandler(
	ctx context.Context,
	newSource chan<- source.Source,
	checkStreamKey CheckStreamKey,
	wg *sync.WaitGroup,
	log *zap.Logger,
) *Handler {
	wg.Add(1)
	ctx, cancel := context.WithCancel(ctx)
	return &Handler{
		ctx:            ctx,
		cancel:         cancel,
		newSource:      newSource,
		checkStreamKey: checkStreamKey,
		wg:             wg,
		log:            log,
	}
}

func (h *Handler) OnServe(conn *rtmp.Conn) {
}

func (h *Handler) OnConnect(timestamp uint32, cmd *rtmpmsg.NetConnectionConnect) error {
	log.Printf("OnConnect: %#v", cmd)
	return nil
}

func (h *Handler) OnCreateStream(timestamp uint32, cmd *rtmpmsg.NetConnectionCreateStream) error {
	log.Printf("OnCreateStream: %#v", cmd)
	return nil
}

func (h *Handler) OnPublish(
	_ *rtmp.StreamContext,
	timestamp uint32,
	cmd *rtmpmsg.NetStreamPublish,
) error {
	log.Printf("OnPublish [%s]", cmd.PublishingType)
	// cmd.PublishingName is the stream secret key in OBS
	// use this value to determine which clients should
	// receive the transcription over websocket
	if cmd.PublishingName == "" {
		return errors.New("PublishingName is empty")
	}
	if err := h.checkStreamKey(cmd.PublishingName); err != nil {
		return errors.Wrap(err, "invalid stream key")
	}
	return nil
}

func (h *Handler) OnSetDataFrame(timestamp uint32, data *rtmpmsg.NetStreamSetDataFrame) error {
	log.Printf("OnSetDataFrame")
	return nil
}

func (h *Handler) OnAudio(timestamp uint32, payload io.Reader) error {
	var audio flvtag.AudioData
	if err := flvtag.DecodeAudioData(payload, &audio); err != nil {
		return err
	}
	if audio.SoundType != flvtag.SoundTypeStereo {
		h.log.Warn("audio stream is not stereo", zap.Int("type", int(audio.SoundType)))
		return errors.Errorf("only stereo sound is supported")
	}
	if audio.SoundFormat != flvtag.SoundFormatAAC {
		h.log.Warn("unsupported audio format", zap.Int("format", int(audio.SoundFormat)))
		return errors.Errorf("only AAC sound is supported")
	}
	if audio.SoundSize != flvtag.SoundSize16Bit {
		h.log.Warn("unsupported audio size", zap.Int("size", int(audio.SoundSize)))
		return errors.Errorf("only 16-bit audio is supported")
	}
	flvBody := new(bytes.Buffer)
	if _, err := io.Copy(flvBody, audio.Data); err != nil {
		return err
	}
	data := flvBody.Bytes()
	if h.source == nil {
		var sampleRate int64
		switch audio.SoundRate {
		case flvtag.SoundRate5_5kHz:
			return errors.Errorf("sample rate of 5kHz is less than minimum supported rate of 8kHz")
		case flvtag.SoundRate11kHz:
			sampleRate = 11025
		case flvtag.SoundRate22kHz:
			sampleRate = 22050
		case flvtag.SoundRate44kHz:
			sampleRate = 44100
		default:
			return errors.Errorf("invalid sound rate: %d", audio.SoundRate)
		}
		var err error
		h.source, err = aac.NewAACSource(
			h.ctx,
			sampleRate,
			audio.SoundType == flvtag.SoundTypeStereo,
			h.log,
		)
		if err != nil {
			h.log.Warn("failed to initialize audio decoder", zap.Error(err))
			return errors.Wrap(err, "failed to initialize audio decoder")
		}
		// notify the server that we have a new source
		select {
		case <-h.ctx.Done():
			return h.ctx.Err()
		case h.newSource <- h.source:
		}
	}
	switch audio.AACPacketType {
	case flvtag.AACPacketTypeSequenceHeader:
		// received codec information
		//h.log.Debug("got aac sequence header", zap.ByteString("data", data))
		//cfg := &aac.AudioSpecificConfig{}
		//if err := cfg.UnmarshalBinary(data); err != nil {
		//	h.log.Warn("failed to parse audio specific config", zap.Error(err))
		//	return errors.Wrap(err, "failed to parse audio specific config")
		//}
		//h.log.Debug("decoded aac sequence header", zap.String("asc", cfg.String()))
		if err := h.source.InitSeqHeader(data); err != nil {
			h.log.Warn("failed to initialize audio decoder", zap.Error(err))
			return errors.Wrap(err, "failed to initialize audio decoder")
		}
		return nil
	case flvtag.AACPacketTypeRaw:
		// received audio data
		if n, err := h.source.Write(data); err != nil {
			h.log.Warn("failed to decode audio", zap.Error(err))
			return errors.Wrap(err, "failed to decode audio")
		} else if n != len(data) {
			h.log.Warn("failed to write full frame to decoder",
				zap.Int("n", n),
				zap.Int("len", len(data)))
			return errors.Errorf("failed to decode audio: %d != %d", n, len(data))
		}
		return nil
	default:
		return errors.Errorf("invalid AAC packet type: %d", audio.AACPacketType)
	}
}

func (h *Handler) OnVideo(timestamp uint32, payload io.Reader) error {
	return nil
}

func (h *Handler) OnClose() {
	defer h.wg.Done()
	h.log.Debug("OnClose")
	h.cancel()
	// TODO: properly cancel transcription
	if h.source != nil {
		h.source.Close()
		h.source = nil
	}
}
