package aac

import (
	"fmt"

	"github.com/pkg/errors"
)

// The AAC object type in RAW AAC frame.
// Refer to @doc ISO_IEC_14496-3-AAC-2001.pdf, @page 23, @section 1.5.1.1 Audio object type definition
type ObjectType uint8

const (
	ObjectTypeForbidden ObjectType = iota

	ObjectTypeMain
	ObjectTypeLC
	ObjectTypeSSR

	ObjectTypeHE   ObjectType = 5  // HE=LC+SBR
	ObjectTypeHEv2 ObjectType = 29 // HEv2=LC+SBR+PS
)

func (o ObjectType) String() string {
	switch o {
	case ObjectTypeForbidden:
		return "Forbidden"
	case ObjectTypeMain:
		return "Main"
	case ObjectTypeLC:
		return "LC"
	case ObjectTypeSSR:
		return "SSR"
	case ObjectTypeHE:
		return "HE"
	case ObjectTypeHEv2:
		return "HEv2"
	default:
		return "Unknown"
	}
}

// The aac sample rate index.
// Refer to @doc ISO_IEC_13818-7-AAC-2004.pdf, @page 46, @section Table 35 – Sampling frequency
type SampleRateIndex uint8

const (
	SampleRateIndex96kHz SampleRateIndex = iota
	SampleRateIndex88kHz
	SampleRateIndex64kHz
	SampleRateIndex48kHz
	SampleRateIndex44kHz
	SampleRateIndex32kHz
	SampleRateIndex24kHz
	SampleRateIndex22kHz
	SampleRateIndex16kHz
	SampleRateIndex12kHz
	SampleRateIndex11kHz
	SampleRateIndex8kHz
	SampleRateIndex7kHz
	SampleRateIndexReserved0
	SampleRateIndexReserved1
	SampleRateIndexReserved2
	SampleRateIndexReserved3
	SampleRateIndexForbidden
)

func (v SampleRateIndex) String() string {
	switch v {
	case SampleRateIndex96kHz:
		return "96kHz"
	case SampleRateIndex88kHz:
		return "88kHz"
	case SampleRateIndex64kHz:
		return "64kHz"
	case SampleRateIndex48kHz:
		return "48kHz"
	case SampleRateIndex44kHz:
		return "44kHz"
	case SampleRateIndex32kHz:
		return "32kHz"
	case SampleRateIndex24kHz:
		return "24kHz"
	case SampleRateIndex22kHz:
		return "22kHz"
	case SampleRateIndex16kHz:
		return "16kHz"
	case SampleRateIndex12kHz:
		return "12kHz"
	case SampleRateIndex11kHz:
		return "11kHz"
	case SampleRateIndex8kHz:
		return "8kHz"
	case SampleRateIndex7kHz:
		return "7kHz"
	case SampleRateIndexReserved0, SampleRateIndexReserved1, SampleRateIndexReserved2, SampleRateIndexReserved3:
		return "Reserved"
	default:
		return "Forbidden"
	}
}

func (v SampleRateIndex) ToHz() int {
	aacSR := []int{
		96000, 88200, 64000, 48000,
		44100, 32000, 24000, 22050,
		16000, 12000, 11025, 8000,
		7350, 0, 0, 0,
		/* To avoid overflow by forbidden */
		0,
	}
	return aacSR[v]
}

// The aac channel.
// Refer to @doc ISO_IEC_13818-7-AAC-2004.pdf, @page 72, @section Table 42 – Implicit speaker mapping
type Channels uint8

const (
	ChannelForbidden Channels = iota
	// center front speaker
	// FFMPEG: mono           FC
	ChannelMono
	// left, right front speakers
	// FFMPEG: stereo         FL+FR
	ChannelStereo
	// center front speaker, left, right front speakers
	// FFMPEG: 2.1            FL+FR+LFE
	// FFMPEG: 3.0            FL+FR+FC
	// FFMPEG: 3.0(back)      FL+FR+BC
	Channel3
	// center front speaker, left, right center front speakers, rear surround
	// FFMPEG: 4.0            FL+FR+FC+BC
	// FFMPEG: quad           FL+FR+BL+BR
	// FFMPEG: quad(side)     FL+FR+SL+SR
	// FFMPEG: 3.1            FL+FR+FC+LFE
	Channel4
	// center front speaker, left, right front speakers, left surround, right surround rear speakers
	// FFMPEG: 5.0            FL+FR+FC+BL+BR
	// FFMPEG: 5.0(side)      FL+FR+FC+SL+SR
	// FFMPEG: 4.1            FL+FR+FC+LFE+BC
	Channel5
	// center front speaker, left, right front speakers, left surround, right surround rear speakers,
	// front low frequency effects speaker
	// FFMPEG: 5.1            FL+FR+FC+LFE+BL+BR
	// FFMPEG: 5.1(side)      FL+FR+FC+LFE+SL+SR
	// FFMPEG: 6.0            FL+FR+FC+BC+SL+SR
	// FFMPEG: 6.0(front)     FL+FR+FLC+FRC+SL+SR
	// FFMPEG: hexagonal      FL+FR+FC+BL+BR+BC
	Channel5_1 // speakers: 6
	// center front speaker, left, right center front speakers, left, right outside front speakers,
	// left surround, right surround rear speakers, front low frequency effects speaker
	// FFMPEG: 7.1            FL+FR+FC+LFE+BL+BR+SL+SR
	// FFMPEG: 7.1(wide)      FL+FR+FC+LFE+BL+BR+FLC+FRC
	// FFMPEG: 7.1(wide-side) FL+FR+FC+LFE+FLC+FRC+SL+SR
	Channel7_1 // speakers: 7
	// FFMPEG: 6.1            FL+FR+FC+LFE+BC+SL+SR
	// FFMPEG: 6.1(back)      FL+FR+FC+LFE+BL+BR+BC
	// FFMPEG: 6.1(front)     FL+FR+LFE+FLC+FRC+SL+SR
	// FFMPEG: 7.0            FL+FR+FC+BL+BR+SL+SR
	// FFMPEG: 7.0(front)     FL+FR+FC+FLC+FRC+SL+SR
)

func (v Channels) String() string {
	switch v {
	case ChannelMono:
		return "Mono(FC)"
	case ChannelStereo:
		return "Stereo(FL+FR)"
	case Channel3:
		return "FL+FR+FC"
	case Channel4:
		return "FL+FR+FC+BC"
	case Channel5:
		return "FL+FR+FC+SL+SR"
	case Channel5_1:
		return "FL+FR+FC+LFE+SL+SR"
	case Channel7_1:
		return "FL+FR+FC+LFE+BL+BR+SL+SR"
	default:
		return "Forbidden"
	}
}

type AudioSpecificConfig struct {
	Object     ObjectType      // AAC object type.
	SampleRate SampleRateIndex // AAC sample rate, not the FLV sampling rate.
	Channels   Channels        // AAC channel configuration.
}

func (v *AudioSpecificConfig) validate() (err error) {
	switch v.Object {
	case ObjectTypeMain, ObjectTypeLC, ObjectTypeSSR, ObjectTypeHE, ObjectTypeHEv2:
	default:
		return errors.Errorf("invalid object %#x", uint8(v.Object))
	}

	if v.SampleRate < SampleRateIndex88kHz || v.SampleRate > SampleRateIndex7kHz {
		return errors.Errorf("invalid sample-rate %#x", uint8(v.SampleRate))
	}

	if v.Channels < ChannelMono || v.Channels > Channel7_1 {
		return errors.Errorf("invalid channels %#x", uint8(v.Channels))
	}
	return
}

func (v *AudioSpecificConfig) UnmarshalBinary(data []byte) (err error) {
	// AudioSpecificConfig
	// Refer to @doc ISO_IEC_14496-3-AAC-2001.pdf, @page 33, @section 1.6.2.1 AudioSpecificConfig
	//
	// only need to decode the first 2bytes:
	// audioObjectType, 5bits.
	// samplingFrequencyIndex, aac_sample_rate, 4bits.
	// channelConfiguration, aac_channels, 4bits
	//
	// @see SrsAacTransmuxer::write_audio
	if len(data) < 2 {
		return errors.Errorf("requires 2 but only %v bytes", len(data))
	}

	t0, t1 := uint8(data[0]), uint8(data[1])

	v.Object = ObjectType((t0 >> 3) & 0x1f)
	v.SampleRate = SampleRateIndex(((t0 << 1) & 0x0e) | ((t1 >> 7) & 0x01))
	v.Channels = Channels((t1 >> 3) & 0x0f)

	return v.validate()
}

func (v *AudioSpecificConfig) MarshalBinary() (data []byte, err error) {
	if err = v.validate(); err != nil {
		return
	}

	// AudioSpecificConfig
	// Refer to @doc ISO_IEC_14496-3-AAC-2001.pdf, @page 33, @section 1.6.2.1 AudioSpecificConfig
	//
	// only need to decode the first 2bytes:
	// audioObjectType, 5bits.
	// samplingFrequencyIndex, aac_sample_rate, 4bits.
	// channelConfiguration, aac_channels, 4bits
	return []byte{
		byte(byte(v.Object)&0x1f)<<3 | byte(byte(v.SampleRate)&0x0e)>>1,
		byte(byte(v.SampleRate)&0x01)<<7 | byte(byte(v.Channels)&0x0f)<<3,
	}, nil
}

func (v *AudioSpecificConfig) String() string {
	return fmt.Sprintf(
		"AudioSpecificConfig(object=%s, sample-rate=%s, channels=%s)",
		v.Object.String(),
		v.SampleRate.String(),
		v.Channels.String(),
	)
}
