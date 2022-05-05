package mixing

import (
	"github.com/gotracker/gomixing/panning"
	"github.com/gotracker/gomixing/volume"
)

// PanMixer is a mixer that's specialized for mixing multichannel audio content
type PanMixer interface {
	GetMixingMatrix(panning.Position) volume.Matrix
	NumChannels() int
}

// GetPanMixer returns the panning mixer that can generate a matrix
// based on input pan value
func GetPanMixer(channels int) PanMixer {
	switch channels {
	case 1:
		return PanMixerMono
	case 2:
		return PanMixerStereo
	case 4:
		return PanMixerQuad
	}

	return nil
}
