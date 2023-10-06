package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"time"

	speech "cloud.google.com/go/speech/apiv1"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
)

func extractAudio(videoFile string, audioFile string) error {
	cmd := exec.Command("ffmpeg", "-i", videoFile, "-q:a", "0", "-map", "a", audioFile, "-y")
	return cmd.Run()
}

func transcribeAudio(audioFile string) ([]*speechpb.SpeechRecognitionResult, error) {
	ctx := context.Background()
	client, err := speech.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(audioFile)
	if err != nil {
		return nil, err
	}

	resp, err := client.Recognize(ctx, &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_LINEAR16,
			SampleRateHertz: 16000,
			LanguageCode:    "en-US",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: data},
		},
	})
	if err != nil {
		return nil, err
	}

	return resp.Results, nil
}

func writeSRT(results []*speechpb.SpeechRecognitionResult, srtFile string) error {
	file, err := os.Create(srtFile)
	if err != nil {
		return err
	}
	defer file.Close()

	index := 1
	for _, result := range results {
		for _, alt := range result.Alternatives {
			if len(alt.Words) == 0 {
				continue
			}

			startTime := alt.Words[0].StartTime
			endTime := alt.Words[len(alt.Words)-1].EndTime

			start := time.Duration(startTime.GetSeconds())*time.Second + time.Duration(startTime.GetNanos())*time.Nanosecond
			end := time.Duration(endTime.GetSeconds())*time.Second + time.Duration(endTime.GetNanos())*time.Nanosecond

			entry := fmt.Sprintf("%d\n%s --> %s\n%s\n\n", index, start, end, alt.Transcript)
			_, err := file.WriteString(entry)
			if err != nil {
				return err
			}
			index++
		}
	}
	return nil
}

func main() {
	var videoFile string
	flag.StringVar(&videoFile, "video", "", "Path to the video file")
	flag.Parse()

	if videoFile == "" {
		log.Fatalf("Please provide a video file path using the -video flag.")
	}

	usr, err := user.Current()
	if err != nil {
		log.Fatalf("Failed to get current user: %v", err)
	}

	audioFile := filepath.Join(usr.HomeDir, "tmp", "audio.wav")
	srtFileName := filepath.Base(videoFile)
	srtFile := filepath.Join(usr.HomeDir, "tmp", srtFileName+".srt")

	err = extractAudio(videoFile, audioFile)
	if err != nil {
		log.Fatalf("Failed to extract audio: %v", err)
	}

	results, err := transcribeAudio(audioFile)
	if err != nil {
		log.Fatalf("Failed to transcribe audio: %v", err)
	}

	err = writeSRT(results, srtFile)
	if err != nil {
		log.Fatalf("Failed to write SRT file: %v", err)
	}

	fmt.Printf("Subtitles generated successfully and saved to %s\n", srtFile)
}
