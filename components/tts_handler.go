package linux_screen_reader

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/faiface/beep/wav"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"

	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
)

type ttsHandle struct {
	mu   sync.Mutex
	init bool
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

	f, err := t.request(text)
	if err != nil {
		return err
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	streamer, format, err := wav.Decode(f)
	if err != nil {
		fmt.Println("occurrence error data: " + text)
		fmt.Println("catch error : ", err)
		return err
	}

	defer func() {
		if err = streamer.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	if !t.init {
		err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
		if err != nil {
			log.Println(text)
			log.Println(err)
			return err
		}
		t.init = true
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		time.Sleep(1 * time.Second)
		wg.Done()
	})))
	wg.Wait()

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
