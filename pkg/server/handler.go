package server

import (
	"bytes"
	"io"
	"log"

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

	// (example) Reject a connection when PublishingName is empty
	if cmd.PublishingName == "" {
		return errors.New("PublishingName is empty")
	}

	// TODO: Create a new stream context and start transcoding

	return nil
}

func (h *Handler) OnSetDataFrame(timestamp uint32, data *rtmpmsg.NetStreamSetDataFrame) error {
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
	audio.Data = flvBody

	log.Printf("FLV Audio Data: Timestamp = %d, SoundFormat = %+v, SoundRate = %+v, SoundSize = %+v, SoundType = %+v, AACPacketType = %+v, Data length = %+v",
		timestamp,
		audio.SoundFormat,
		audio.SoundRate,
		audio.SoundSize,
		audio.SoundType,
		audio.AACPacketType,
		len(flvBody.Bytes()),
	)

	// TODO: transcribe the audio chunk

	return nil
}

func (h *Handler) OnVideo(timestamp uint32, payload io.Reader) error {
	return nil
}

func (h *Handler) OnClose() {
	log.Printf("OnClose")
}
