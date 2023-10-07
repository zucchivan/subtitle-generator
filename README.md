# Subtitle Generator

## Overview

This is a simple command-line tool that automates the process of generating subtitle files in SRT format. It uses Google's Speech-to-Text API to transcribe the audio from a given video file and then creates an SRT file with the transcribed text.

## Pre-requisites

- Go 1.21 or higher
- FFmpeg installed and available in your PATH
- Google Cloud credentials set up for Speech-to-Text API

## Google Cloud Credentials Setup

1. Go to the [Google Cloud Console](https://console.cloud.google.com/).
2. Create a new project or select an existing one.
3. Navigate to "APIs & Services" > "Credentials".
4. Create a new service account and download the JSON key file.
5. Set the `GOOGLE_APPLICATION_CREDENTIALS` environment variable to the path of the downloaded JSON key file:

    ```bash
    export GOOGLE_APPLICATION_CREDENTIALS="/path/to/your-service-account-file.json"
    ```

    (Add this line to your `.bashrc` or `.zshrc` to make it permanent)

## Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/zucchivan/subtitle-generator.git
    ```

2. Navigate to the project `cmd` directory:

    ```bash
    cd subtitle-generator/cmd
    ```

3. Install the required Go packages:

    ```bash
    go mod download
    ```

4. Build the application:

    ```bash
    go build -o subtitle-generator
    ```

## Usage

To generate an SRT file for a video, run the following command:

```bash
./subtitle-generator -video=/path/to/your-video.mp4
```

This will generate an SRT file under ~/tmp/ with the same name as the input video file.
