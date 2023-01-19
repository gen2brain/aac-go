package aac_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/gen2brain/aac-go"
	"github.com/youpy/go-wav"
)

func TestEncode(t *testing.T) {
	file, err := os.Open(filepath.Join("testdata", "test.wav"))
	if err != nil {
		t.Fatal(err)
	}

	wr := wav.NewReader(file)
	f, err := wr.Format()
	if err != nil {
		t.Fatal(err)
	}

	buf := bytes.NewBuffer(make([]byte, 0))

	opts := &aac.Options{}
	opts.SampleRate = int(f.SampleRate)
	opts.NumChannels = int(f.NumChannels)

	enc, err := aac.NewEncoder(buf, opts)
	if err != nil {
		t.Fatal(err)
	}

	err = enc.Encode(wr)
	if err != nil {
		t.Error(err)
	}

	err = enc.Close()
	if err != nil {
		t.Error(err)
	}

	err = os.WriteFile(filepath.Join(os.TempDir(), "test.aac"), buf.Bytes(), 0640)
	if err != nil {
		t.Error(err)
	}
}
