## aac-go
[![Build Status](https://github.com/gen2brain/aac-go/actions/workflows/build.yml/badge.svg)](https://github.com/gen2brain/aac-go/actions)
[![GoDoc](https://godoc.org/github.com/gen2brain/aac-go?status.svg)](https://godoc.org/github.com/gen2brain/aac-go) 
[![Go Report Card](https://goreportcard.com/badge/github.com/gen2brain/aac-go?branch=master)](https://goreportcard.com/report/github.com/gen2brain/aac-go) 

`aac-go` provides AAC codec encoder based on [VisualOn AAC encoder](https://github.com/mstorsjo/vo-aacenc) library.

### Installation

    go get -u github.com/gen2brain/aac-go

### Examples

See [micgrab](https://github.com/gen2brain/aac-go/blob/master/examples/micgrab/micgrab.go) example.

### Usage

```go
package main

import (
	"bytes"
	"os"

	"github.com/gen2brain/aac-go"
	"github.com/youpy/go-wav"
)

func main() {
	file, err := os.Open("test.wav")
	if err != nil {
		panic(err)
	}

	wreader := wav.NewReader(file)
	f, err := wreader.Format()
	if err != nil {
		panic(err)
	}

	buf := bytes.NewBuffer(make([]byte, 0))

	opts := &aac.Options{}
	opts.SampleRate = int(f.SampleRate)
	opts.NumChannels = int(f.NumChannels)

	enc, err := aac.NewEncoder(buf, opts)
	if err != nil {
		panic(err)
	}

	err = enc.Encode(wreader)
	if err != nil {
		panic(err)
	}

	err = enc.Close()
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("test.aac", buf.Bytes(), 0644)
	if err != nil {
		panic(err)
	}
}
```

## More

For H.264 encoder see [x264-go](https://github.com/gen2brain/x264-go).
