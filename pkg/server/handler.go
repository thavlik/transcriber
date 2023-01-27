package server

import (
	"bytes"
	"io"
	"log"

	"github.com/izern/go-fdkaac/fdkaac"
	"github.com/pkg/errors"
	flvtag "github.com/yutopp/go-flv/tag"
	"github.com/yutopp/go-rtmp"
	rtmpmsg "github.com/yutopp/go-rtmp/message"
	"go.uber.org/zap"
)

var _ rtmp.Handler = (*Handler)(nil)

// Handler An RTMP connection handler
type Handler struct {
	rtmp.DefaultHandler
	dec *fdkaac.AacDecoder
	log *zap.Logger
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

func (h *Handler) OnPublish(_ *rtmp.StreamContext, timestamp uint32, cmd *rtmpmsg.NetStreamPublish) error {
	log.Printf("OnPublish: %#v", cmd)
	// cmd.PublishingName is the stream secret key in OBS
	// use this value to determine which clients should
	// receive the transcription over websocket
	if cmd.PublishingName == "" {
		return errors.New("PublishingName is empty")
	}
	if h.dec != nil {
		return errors.New("decoder already exists, did the client publish twice?")
	}
	h.dec = fdkaac.NewAacDecoder()
	// TODO: Create a new stream context and start transcoding
	return nil
}

func (h *Handler) OnSetDataFrame(timestamp uint32, data *rtmpmsg.NetStreamSetDataFrame) error {
	log.Printf("OnSetDataFrame: %#v", data)
	return nil
}

func (h *Handler) OnAudio(timestamp uint32, payload io.Reader) error {
	var audio flvtag.AudioData
	if err := flvtag.DecodeAudioData(payload, &audio); err != nil {
		return err
	}

	flvBody := new(bytes.Buffer)
	if _, err := io.Copy(flvBody, audio.Data); err != nil {
		return err
	}

	log.Printf("FLV Audio Data: Timestamp = %d, SoundFormat = %+v, SoundRate = %+v, SoundSize = %+v, SoundType = %+v, AACPacketType = %+v, Data length = %+v",
		timestamp,
		audio.SoundFormat,
		audio.SoundRate,
		audio.SoundSize,
		audio.SoundType,
		audio.AACPacketType,
		len(flvBody.Bytes()),
	)

	pcm, err := h.dec.Decode(flvBody.Bytes())
	if err != nil {
		return errors.Wrap(err, "failed to decode aac audio")
	} else if pcm == nil {
		return nil
	}

	// TODO: write the pcm data to the stream source

	return nil
}

func (h *Handler) OnVideo(timestamp uint32, payload io.Reader) error {
	return nil
}

func (h *Handler) OnClose() {
	h.log.Debug("OnClose")
	_ = h.dec.Close()
	h.dec = nil
}
