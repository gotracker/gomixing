package mixing

import (
	"math"

	"github.com/gotracker/gomixing/panning"
	"github.com/gotracker/gomixing/volume"
)

// PanMixer is a mixer that's specialized for mixing multichannel audio content
type PanMixer interface {
	GetMixingMatrix(panning.Position) volume.Matrix
	Channels() int
}

var (
	// PanMixerMono is a mixer that's specialized for mixing monaural audio content
	PanMixerMono PanMixer = &panMixerMono{}

	// PanMixerStereo is a mixer that's specialized for mixing stereo audio content
	PanMixerStereo PanMixer = &panMixerStereo{}

	// PanMixerQuad is a mixer that's specialized for mixing quadraphonic audio content
	PanMixerQuad PanMixer = &panMixerQuad{}
)

type panMixerMono struct {
	PanMixer
}

func (p panMixerMono) GetMixingMatrix(pan panning.Position) volume.Matrix {
	// distance and angle are ignored on mono
	return volume.Matrix{1.0}
}

func (p panMixerMono) Channels() int {
	return 1
}

type panMixerStereo struct {
	PanMixer
}

func (p panMixerStereo) GetMixingMatrix(pan panning.Position) volume.Matrix {
	pangle := float64(pan.Angle)
	s, c := math.Sincos(pangle)
	var d volume.Volume
	if pan.Distance > 0 {
		d = 1 / volume.Volume(pan.Distance*pan.Distance)
	}
	l := d * volume.Volume(s)
	r := d * volume.Volume(c)
	return volume.Matrix{l, r}
}

func (p panMixerStereo) Channels() int {
	return 2
}

type panMixerQuad struct {
	PanMixer
}

func (p panMixerQuad) GetMixingMatrix(pan panning.Position) volume.Matrix {
	pangle := float64(pan.Angle)
	sf, cf := math.Sincos(pangle)
	sr, cr := math.Sin(pangle+math.Pi/2.0), math.Cos(pangle-math.Pi/2.0)
	var d volume.Volume
	if pan.Distance > 0 {
		d = 1 / volume.Volume(pan.Distance*pan.Distance)
	}
	lf := d * volume.Volume(sf)
	rf := d * volume.Volume(cf)
	lr := d * volume.Volume(cr)
	rr := d * volume.Volume(sr)
	return volume.Matrix{lf, rf, lr, rr}
}

func (p panMixerQuad) Channels() int {
	return 4
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
