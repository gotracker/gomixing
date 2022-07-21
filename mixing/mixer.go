package mixing

import (
	"github.com/gotracker/gomixing/sampling"
	"github.com/gotracker/gomixing/volume"
)

// Mixer is a manager for mixing multiple single- and multi-channel samples into a single multi-channel output stream
type Mixer struct {
	Channels      int
	BitsPerSample int
}

// NewMixBuffer returns a mixer buffer with a number of channels
// of preallocated sample data
func (m Mixer) NewMixBuffer(samples int) MixBuffer {
	return make(MixBuffer, samples)
}

// GetDefaultMixerVolume returns the default mixer volume value based on the number of mixed channels
// not to be confused with the number of output channels
func GetDefaultMixerVolume(numMixedChannels int) volume.Volume {
	return 1.0 / volume.Volume(numMixedChannels)
}

// Flatten will to a final saturation mix of all the row's channel data into a single output buffer
func (m Mixer) Flatten(panmixer PanMixer, samplesLen int, row []ChannelData, mixerVolume volume.Volume, sampleFormat sampling.Format) []byte {
	data := m.NewMixBuffer(samplesLen)
	formatter := sampling.GetFormatter(sampleFormat)
	for _, rdata := range row {
		for _, cdata := range rdata {
			if cdata.Flush != nil {
				cdata.Flush()
			}
			if len(cdata.Data) > 0 {
				volMtx := panmixer.GetMixingMatrix(cdata.Pan).Apply(cdata.Volume)
				data.Add(cdata.Pos, &cdata.Data, volMtx)
			}
		}
	}
	return data.ToRenderData(samplesLen, m.Channels, mixerVolume, formatter)
}

// FlattenToInts runs a flatten on the channel data into separate channel data of int32 variety
// these int32s still respect the BitsPerSample size
func (m Mixer) FlattenToInts(panmixer PanMixer, samplesLen int, row []ChannelData, mixerVolume volume.Volume) [][]int32 {
	data := m.NewMixBuffer(samplesLen)
	for _, rdata := range row {
		for _, cdata := range rdata {
			if cdata.Flush != nil {
				cdata.Flush()
			}
			if len(cdata.Data) > 0 {
				volMtx := panmixer.GetMixingMatrix(cdata.Pan).Apply(cdata.Volume)
				data.Add(cdata.Pos, &cdata.Data, volMtx)
			}
		}
	}
	return data.ToIntStream(panmixer.NumChannels(), samplesLen, m.BitsPerSample, mixerVolume)
}

// FlattenTo will to a final saturation mix of all the row's channel data into a single output buffer
func (m Mixer) FlattenTo(resultBuffers [][]byte, panmixer PanMixer, samplesLen int, row []ChannelData, mixerVolume volume.Volume, sampleFormat sampling.Format) {
	data := m.NewMixBuffer(samplesLen)
	formatter := sampling.GetFormatter(sampleFormat)
	for _, rdata := range row {
		for _, cdata := range rdata {
			if cdata.Flush != nil {
				cdata.Flush()
			}
			if len(cdata.Data) > 0 {
				volMtx := panmixer.GetMixingMatrix(cdata.Pan).Apply(cdata.Volume)
				data.Add(cdata.Pos, &cdata.Data, volMtx)
			}
		}
	}
	data.ToRenderDataWithBufs(resultBuffers, samplesLen, mixerVolume, formatter)
}
