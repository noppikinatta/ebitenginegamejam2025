// Copyright 2022 noppikinatta
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package asset

import (
	"bytes"
	_ "embed"
	"errors"
	"io"
	"log"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

const (
	bgmLength int64 = 21054898
)

//go:embed sound/bgm.ogg
var bgm []byte

//go:embed sound/explosion.wav
var seExplosion []byte

const sampleRate int = 48000

var context *audio.Context

func init() {
	context = audio.NewContext(sampleRate)
	soundCache = map[Sound]*audio.Player{}
}

type Sound int

const (
	BGM Sound = iota
	SEExplosion
)

func LoadSounds() error {
	ss := []struct {
		Resource []byte
		Sound    Sound
		FileType fileType
		Volume   float64
	}{
		{seExplosion, SEExplosion, fileTypeWav, 0.125},
		{bgm, BGM, fileTypeOgg, 0.25},
	}

	for _, s := range ss {
		err := load(s.Resource, s.Sound, s.FileType, s.Volume)
		if err != nil {
			return err
		}
	}

	return nil
}

type fileType int

const (
	fileTypeWav fileType = iota
	fileTypeMp3
	fileTypeOgg
)

func load(resource []byte, sound Sound, ftype fileType, vol float64) error {
	var s io.ReadSeeker
	var err error

	switch ftype {
	case fileTypeWav:
		s, err = wav.DecodeWithSampleRate(sampleRate, bytes.NewReader(resource))
		if err != nil {
			return err
		}
	case fileTypeMp3:
		s, err = mp3.DecodeWithSampleRate(sampleRate, bytes.NewReader(resource))
		if err != nil {
			return err
		}
	case fileTypeOgg:
		s, err = vorbis.DecodeWithSampleRate(sampleRate, bytes.NewReader(resource))
		if err != nil {
			return err
		}
	default:
		return errors.New("not supported filetype")
	}

	// BGM loops
	if sound == BGM {
		s = audio.NewInfiniteLoop(s, bgmLength) //int64(len(resource)))
	}

	p, err := context.NewPlayer(s)
	if err != nil {
		return err
	}
	p.SetVolume(vol)
	soundCache[sound] = p

	return nil
}

var soundCache map[Sound]*audio.Player

func PlaySound(s Sound) {
	p := soundCache[s]
	err := p.Rewind()
	if err != nil {
		log.Println(err)
	}
	p.Play()
}

func StopSound(s Sound) {
	p := soundCache[s]
	p.Pause()
}
