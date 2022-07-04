// Command quickstart generates an audio file with the content "Hello, World!".
package main

import (
	"context"
	"io/ioutil"
	"os"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"

	"github.com/alecthomas/kong"
)

func runTTS(text string, languageCode string, voiceName string) ([]byte, error) {
	// Instantiates a client.
	ctx := context.Background()

	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	// Perform the text-to-speech request on the text input with the selected
	// voice parameters and audio file type.
	req := texttospeechpb.SynthesizeSpeechRequest{
		// Set the text input to be synthesized.
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: text},
		},
		// Build the voice request, select the language code ("en-US") and the SSML
		// voice gender ("neutral").
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: languageCode,
			Name:         voiceName,
			// SsmlGender:   texttospeechpb.SsmlVoiceGender_NEUTRAL,
		},
		// Select the type of audio file you want returned.
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_MP3,
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

	LanguageCode string `help:"Language code to use for the synthesis. See full list at: https://cloud.google.com/text-to-speech/docs/voices" default:"en-US"`

	VoiceName string `help:"Voice name to use for the synthesis. Use an empty string to let the GCP API choose. See full list at: https://cloud.google.com/text-to-speech/docs/voices" default:"en-US-Wavenet-F"`
}

func (c *Cli) AfterApply() error {
	// the VoiceName overrides LanguageCode if given to the GCP API.
	// So if a different LanguageCode is used, we reset the VoiceName.
	if c.LanguageCode != "en-US" && c.VoiceName == "en-US-Wavenet-F" {
		c.VoiceName = ""
	}
	return nil
}

var cli Cli

func main() {
	kctx := kong.Parse(&cli)

	input, err := readFile(cli.InputFile)
	kctx.FatalIfErrorf(err)

	audio, err := runTTS(input, cli.LanguageCode, cli.VoiceName)
	kctx.FatalIfErrorf(err)

	err = writeOutput(cli.OutputFile, audio)
	kctx.FatalIfErrorf(err)
}
