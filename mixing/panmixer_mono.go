package mixing

import (
	"github.com/gotracker/gomixing/panning"
	"github.com/gotracker/gomixing/volume"
)

// PanMixerMono is a mixer that's specialized for mixing monaural audio content
var PanMixerMono PanMixer = &panMixerMono{}

type panMixerMono volume.Matrix

func (p panMixerMono) GetMixingMatrix(pan panning.Position) volume.Matrix {
	// distance and angle are ignored on mono
	return volume.Matrix{
		StaticMatrix: volume.StaticMatrix{1.0},
		Channels:     1,
	}
}

func (p panMixerMono) NumChannels() int {
	return p.Channels
}
