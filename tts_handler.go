package google_tts_screenreader

import (
	"bytes"
	"context"
	"log"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"

	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
)

type ttsHandle struct {
}
type ttsHandler interface {
	play(text string) error
}

func (t *ttsHandle) play(text string) error {

	f, err := t.request(text)

	streamer, format, err := wav.Decode(f)
	if err != nil {
		log.Println(text)

		// todo エラーファイルを収集すする
		//f, err := ioutil.TempFile("/tmp/", "tts-error-data")
		//if err == nil {
		//	// todo エラーファイルを収集する
		//	// bytes.ioutil.WriteFile(f.Name(), f, 644)
		//}

		log.Println(err)
	}
	defer func() {
		_ = streamer.Close()
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
