package main

import (
	"context"
	"os"
	"path/filepath"
	"sync"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/scribe/pkg/source"
	"github.com/thavlik/transcriber/scribe/pkg/source/mp3"
	"github.com/thavlik/transcriber/scribe/pkg/source/wav"
	"github.com/thavlik/transcriber/scribe/pkg/transcribe"
)

var testTranscribeArgs struct {
	specialty string
}

var testTranscribeCmd = &cobra.Command{
	Use:   "transcribe",
	Short: "test Amazon Transcribe with a WAV file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := args[0]
		f, err := os.Open(input)
		if err != nil {
			return err
		}
		defer f.Close()
		var src source.Source
		switch filepath.Ext(input) {
		case ".wav":
			src, err = wav.NewWavSource(cmd.Context(), f)
			if err != nil {
				return errors.Wrap(err, "wav.NewWavSource")
			}
		case ".mp3":
			src, err = mp3.NewMP3Source(cmd.Context(), f)
			if err != nil {
				return errors.Wrap(err, "mp3.NewMP3Source")
			}
		default:
			return errors.New("unsupported file type")
		}
		wg := new(sync.WaitGroup)
		wg.Add(1)
		defer wg.Wait()
		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()
		transcripts := make(chan *transcribe.Transcript, 16)
		go func() {
			defer wg.Done()
			transcribe.PrintTranscripts(ctx, transcripts)
		}()
		return transcribe.Transcribe(
			ctx,
			src,
			testTranscribeArgs.specialty,
			transcripts,
			base.DefaultLog,
		)
	},
}

func init() {
	testCmd.AddCommand(testTranscribeCmd)
	testTranscribeCmd.Flags().StringVarP(
		&testTranscribeArgs.specialty,
		"specialty",
		"s",
		"",
		"if set, the medical model with the given specialty is used for transcription",
	)
}
