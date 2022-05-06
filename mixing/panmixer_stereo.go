package mixing

import (
	"math"

	"github.com/gotracker/gomixing/panning"
	"github.com/gotracker/gomixing/volume"
)

// PanMixerStereo is a mixer that's specialized for mixing stereo audio content
var PanMixerStereo PanMixer = &panMixerStereo{}

type panMixerStereo struct{}

func (p panMixerStereo) GetMixingMatrix(pan panning.Position) volume.Matrix {
	pangle := float64(pan.Angle)
	s, c := math.Sincos(pangle)
	var d volume.Volume
	if pan.Distance > 0 {
		d = 1 / volume.Volume(pan.Distance*pan.Distance)
	}
	l := d * volume.Volume(s)
	r := d * volume.Volume(c)
	return volume.Matrix{
		StaticMatrix: volume.StaticMatrix{l, r},
		Channels:     2,
	}
}

func (p panMixerStereo) NumChannels() int {
	return 2
}
