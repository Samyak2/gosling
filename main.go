// Command quickstart generates an audio file with the content "Hello, World!".
package main

import (
	"context"
	"io/ioutil"
	"os"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"google.golang.org/api/option"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"

	"github.com/alecthomas/kong"
)

func RunTTS(text string, options Cli) ([]byte, error) {
	// Instantiates a client.
	ctx := context.Background()

	client, err := texttospeech.NewClient(ctx, option.WithEndpoint(options.ServiceEndpoint))
	if err != nil {
		return nil, err
	}
	defer client.Close()

	var input *texttospeechpb.SynthesisInput

	if options.Ssml {
		input = &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Ssml{Ssml: text},
		}
	} else {
		input = &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: text},
		}
	}

	// Perform the text-to-speech request on the text input with the selected
	// voice parameters and audio file type.
	req := texttospeechpb.SynthesizeSpeechRequest{
		Input: input,
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: options.LanguageCode,
			Name:         options.VoiceName,
			// SsmlGender:   texttospeechpb.SsmlVoiceGender_NEUTRAL,
		},
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_MP3,
			Pitch:         options.Pitch,
			SpeakingRate:  options.SpeakingRate,
			VolumeGainDb:  options.VolumeGain,
		},
	}

	resp, err := client.SynthesizeSpeech(ctx, &req)
	if err != nil {
		return nil, err
	}

	// The resp's AudioContent is binary.
	return resp.AudioContent, nil
}

func readFile(filename string) (string, error) {
	if filename == "-" {
		data, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return "", err
		}
		return string(data), err
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), err
}

func writeOutput(filename string, data []byte) error {
	if filename == "-" {
		_, err := os.Stdout.Write(data)
		return err
	}
	return ioutil.WriteFile(filename, data, 0644)
}

type Cli struct {
	InputFile string `help:"Text file to read from. Use - for standard input." arg:""`

	OutputFile string `help:"Audio file to write to. Use - for standard output." arg:""`

	LanguageCode string `short:"l" help:"Language code to use for the synthesis. See full list at: https://cloud.google.com/text-to-speech/docs/voices" default:"en-US"`

	VoiceName string `short:"v" help:"Voice name to use for the synthesis. Use an empty string to let the GCP API choose. See full list at: https://cloud.google.com/text-to-speech/docs/voices" default:"en-US-Wavenet-A"`

	Pitch float64 `help:"Pitch adjustment in the range [-20.0, 20.0]. Use a negative number to decrease the pitch. See: https://cloud.google.com/text-to-speech/docs/reference/rest/v1/text/synthesize#audioconfig" default:"-3"`

	SpeakingRate float64 `short:"r" help:"Speaking rate/speed in the range [0.25, 4.0]. See: https://cloud.google.com/text-to-speech/docs/reference/rest/v1/text/synthesize#audioconfig" default:"1.0"`

	VolumeGain float64 `help:"Volume gain (in dB) in the range [-96.0, 16.0]. See: https://cloud.google.com/text-to-speech/docs/reference/rest/v1/text/synthesize#audioconfig" default:"0.0"`

	Ssml bool `short:"s" help:"Use if text has SSML. Default is plain text. See: https://cloud.google.com/text-to-speech/docs/basics#speech_synthesis_markup_language_ssml_support" negatable:"" default:"false"`

	ServiceEndpoint string `help:"GCP Service Endpoint. You'll need to set this if you want a Neural2 voice. See: https://cloud.google.com/text-to-speech/docs/endpoints." optional:""`
}

func (c *Cli) AfterApply() error {
	// the VoiceName overrides LanguageCode if given to the GCP API.
	// So if a different LanguageCode is used, we reset the VoiceName.
	if c.LanguageCode != "en-US" && c.VoiceName == "en-US-Wavenet-A" {
		c.VoiceName = ""
	}
	return nil
}

var cli Cli

func main() {
	kctx := kong.Parse(&cli)

	input, err := readFile(cli.InputFile)
	kctx.FatalIfErrorf(err)

	audio, err := RunTTS(input, cli)
	kctx.FatalIfErrorf(err)

	err = writeOutput(cli.OutputFile, audio)
	kctx.FatalIfErrorf(err)
}
