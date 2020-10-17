package linux_screen_reader

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
)

type ttsHandle struct {
	mu sync.Mutex
}
type ttsHandler interface {
	play(text string) error
}

func NewTTSHandle() *ttsHandle {
	return &ttsHandle{
		mu: sync.Mutex{},
	}
}

func (t *ttsHandle) play(text string) error {

	t.mu.Lock()
	defer t.mu.Unlock()

	f, err := t.request(text)
	if err != nil {
		return err
	}

	streamer, format, err := wav.Decode(f)
	if err != nil {
		e := err
		fmt.Println("occurrence error data: " + text)
		fmt.Println("catch error : ", err)

		tmpf, err := ioutil.TempFile("/tmp/", "tts-error-data")
		if err == nil {
			_ = ioutil.WriteFile(tmpf.Name(), f.Bytes(), 644)
			_ = tmpf.Close()
		}
		return e
	}

	defer func() {
		if err = streamer.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	done := make(chan bool)

	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		log.Println(text)
		log.Println(err)
		return err
	}

	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))
	<-done

	return nil
}

func (t *ttsHandle) request(text string) (*bytes.Buffer, error) {

	ctx := context.Background()

	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	req := texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: text},
		},
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: "ja-JP",
			SsmlGender:   texttospeechpb.SsmlVoiceGender_NEUTRAL,
		},
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_LINEAR16,
			SpeakingRate:  1.12,
		},
	}

	resp, err := client.SynthesizeSpeech(ctx, &req)
	if err != nil {
		log.Println(text)
		log.Println(err)
		return nil, err
	}

	return bytes.NewBuffer(resp.AudioContent), nil
}
