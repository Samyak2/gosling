<h1 align=center>
gosling
</h1>

<p align=center>
Natural sounding text-to-speech in the terminal (and more).
</p>


## Pre-requisites

This is *NOT* intended to be a completely-free, pick-up-and-use TTS solution. In fact, it is simply a wrapper around Google's Cloud Text-to-Speech API.

You will need:
- A GCP account with billing enabled.
    - Google gives you [1 million characters free](https://cloud.google.com/text-to-speech/pricing) every month. That's nearly 10 books a month. It's essentially free for personal use.
    - Once you have a GCP account, [enable the TTS API and get a service account](https://cloud.google.com/text-to-speech/docs/before-you-begin).
    - Export service account credentials in your shell. You will need to do this every time you open a new shell. Add it to your shell configuration or make a script to run `gosling` for convenience.
      ```bash
      export GOOGLE_APPLICATION_CREDENTIALS="/path/to/your/service-account.json"
      ```
- Internet connection every time you need some text spoken to you.
- I have only tested this on Linux. Commands for playing audio will be different on other platforms.

## Examples

### Simple text with default options

https://user-images.githubusercontent.com/34161949/178104531-73298a8e-753f-4910-94c6-7cea9a85337a.mp4

### Numbers and punctuation with default options

(the multiple exclamations are something that I have seen other TTSs struggle with):
```
Welcome to gosling!!! It has options such as "Pitch adjustment" in the range -20.0 to 20.0, "Speaking rate/speed" in the range 0.25 to 4.0 and "Volume gain" (in dB) in the range -96.0 to 16.0.
```

https://user-images.githubusercontent.com/34161949/178104603-f8c46b93-4d38-4f71-bdc0-d3d4b3e47b05.mp4

### Other languages

Kannada:

https://user-images.githubusercontent.com/34161949/178105235-19e921c7-355b-4e66-8c3e-e962718002aa.mp4

Check out the [full voice list](https://cloud.google.com/text-to-speech/docs/voices), use `Wavenet` or `Neural2` based voices for better quality.

## Installation

### Pre-built binaries

Go to the [latest release](https://github.com/Samyak2/gosling/releases/latest), scroll down to "Assets" and download the correct file for your platform. Unzip the file and run the `gosling` binary inside:
```bash
./gosling
```

### If you have `go` installed

```bash
go install github.com/Samyak2/gosling@latest
```

## Usage

### Text file

```bash
gosling input.txt output.mp3
```

Play the resulting `output.mp3` file using your audio player.

### Standard input

```bash
echo "hello there" | gosling - output.mp3
```

### Play audio directly

If you have the `play` command, which is usually a part of the `sox` package (`sudo dnf install sox` on Fedora):
```bash
echo "hello there" | gosling - - | play -t mp3 -
```

If you have the `ffplay` command, which is a part of `ffmpeg`:
```bash
echo "hello there" | gosling - - | ffplay -nodisp -autoexit -
```

### Options

`gosling` has a lot of configuration around [language & voice](https://cloud.google.com/text-to-speech/docs/voices), [audio](https://cloud.google.com/text-to-speech/docs/reference/rest/v1/text/synthesize#audioconfig), etc.


See `gosling --help` for all options.
```
Usage: gosling <input-file> <output-file>

Arguments:
  <input-file>     Text file to read from. Use - for standard input.
  <output-file>    Audio file to write to. Use - for standard output.

Flags:
  -h, --help                            Show context-sensitive help.
  -l, --language-code="en-US"           Language code to use for the synthesis. See full list at: https://cloud.google.com/text-to-speech/docs/voices
  -v, --voice-name="en-US-Wavenet-A"    Voice name to use for the synthesis. Use an empty string to let the GCP API choose. See full list at: https://cloud.google.com/text-to-speech/docs/voices
      --pitch=-3                        Pitch adjustment in the range [-20.0, 20.0]. Use a negative number to decrease the pitch. See:
                                        https://cloud.google.com/text-to-speech/docs/reference/rest/v1/text/synthesize#audioconfig
  -r, --speaking-rate=1.0               Speaking rate/speed in the range [0.25, 4.0]. See: https://cloud.google.com/text-to-speech/docs/reference/rest/v1/text/synthesize#audioconfig
      --volume-gain=0.0                 Volume gain (in dB) in the range [-96.0, 16.0]. See: https://cloud.google.com/text-to-speech/docs/reference/rest/v1/text/synthesize#audioconfig
  -s, --[no-]ssml                       Use if text has SSML. Default is plain text. See: https://cloud.google.com/text-to-speech/docs/basics#speech_synthesis_markup_language_ssml_support
      --service-endpoint=STRING         GCP Service Endpoint. You'll need to set this if you want a Neural2 voice. See: https://cloud.google.com/text-to-speech/docs/endpoints.
```

## FAQ

### The voice sounds too robotic

#### WaveNet

By default, on the default language, `gosling` uses a [WaveNet](https://cloud.google.com/text-to-speech/docs/wavenet) based voice model. If you're using a different language, make sure to switch the voice to a WaveNet based one too. Use `--voice-name` for this.

#### Neural2

If WaveNet is not good enough, try using a `Neural2` voice type (search for `Neural2` in the [voice list](https://cloud.google.com/text-to-speech/docs/voices) if you need other languages):
```
gosling input.txt output.mp3 --service-endpoint 'https://us-central1-texttospeech.googleapis.com' -v en-US-Neural2-A
```
TODO: this endpoint is currently timing out for all TTS requests, not sure why.

If Neural2 isn't good enough either, well... you'll have to take this up with Google.

### Why am I getting this error `google: could not find default credentials`?

Either:
- You did not read the [Pre-requisites](#pre-requisites) section.
- You forgot to export the `GOOGLE_APPLICATION_CREDENTIALS` environment variable in your shell.
- Something is wrong with your GCP service account. See [this page](https://cloud.google.com/docs/authentication/production#passing_variable) that is also linked from the error.

### Why don't `--pitch` and `--volume-gain` have short versions?

These options can have negative values and the command-line parser I use behaves weirdly with [negative numbers and short flags](https://github.com/alecthomas/kong/issues/315). I have removed the short versions to avoid making it a pitfall.

### How do I use this with [`foliate`](https://github.com/johnfactotum/foliate)?

I use this script:
```bash
#!/bin/bash
# requires gosling and sox
export GOOGLE_APPLICATION_CREDENTIALS="/path/to/your/service-account.json"
gosling - - | play -t mp3 - &
trap 'kill $!; exit 0' INT
wait
```

Copy and save this to a file and `chmod +x /path/to/foliate-gosling.sh` it.

TODO: this only works with English text. I need to figure out a way to convert `FOLIATE_TTS_LANG_LOWER` to Google's format.

### But why?

When I'm too lazy to read an article, I use Google Assistant's "read me this article" feature on my phone. It's extremely good, especially with text-only articles. I could not find an alternative on desktop (specifically, Linux).

Yes, there are quite a few text-to-speech apps on Linux. Most of them either sound like R2D2 or something from the depths of the void. The only one, that I found, which sounds bearable uses an undocumented Google Translate API (probably a ToS violation?). There are also some pre-trained neural-network based models, but they sound like a person speaking through a very low-bandwidth voice call and they skip over numbers and abbreviations pretending they never existed.

The only text-to-speech that sounded good was Google's. So I thought - "they must have a GCP API for this". And they did. And I hacked this together.

## TODO

- [ ] `speech-dispatcher` support. This will allow using it in Firefox's reader mode, for example.
- [ ] Some pre-processing of raw text - remove extra/unnecessary punctuation, better formatting for numbers, etc.

## License

MIT
