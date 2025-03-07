package tts

import (
	"context"
	"fmt"
	"log"
	"sync"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	texttospeechpb "cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	"github.com/npvu1510/crawl-en-vocab/pkg/config"
	"github.com/npvu1510/crawl-en-vocab/pkg/utils"
	"google.golang.org/api/option"
)

type TTSService struct {
	client       *texttospeech.Client
	languageCode string
}

var (
	Service *TTSService
	once    sync.Once
)

func init() {
	once.Do(func() {
		conf := config.MustLoad()

		client, err := texttospeech.NewClient(context.Background(), option.WithCredentialsFile(conf.TTS.CREDENTIALS_FILE_PATH))
		utils.CheckError(err)

		Service = &TTSService{client: client}
	})
}

func GetClient() *texttospeech.Client {
	return Service.client
}

func (s *TTSService) SetLanguageCode(code string) {
	s.languageCode = code
}

func (s *TTSService) TextToSpeech(text string) ([]byte, error) {
	ctx := context.Background()

	//
	req := &texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{
				Text: text,
			},
		},
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: s.languageCode,
			SsmlGender:   texttospeechpb.SsmlVoiceGender_NEUTRAL,
		},
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_MP3,
		},
	}

	//
	resp, err := s.client.SynthesizeSpeech(ctx, req)
	if err != nil {
		log.Fatalf("❌ Error while calling Google TTS API: %v", err)
		return nil, err
	}

	fmt.Printf("✅ TTS for %v successfully!\n", text)
	return resp.AudioContent, nil
}

func (c *TTSService) CloseClient() {
	if Service != nil {
		Service.client.Close()
	}
}
