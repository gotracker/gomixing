package volume

import "math"

// Volume is a mixable volume
type Volume float32

var (
	// VolumeUseInstVol tells the system to use the volume stored on the instrument
	VolumeUseInstVol = Volume(math.Inf(-1))
)

// Matrix is an array of Volumes
type Matrix []Volume

// Int24 is an approximation of a 24-bit integer
type Int24 struct {
	Hi int8
	Lo uint16
}

// ToSample returns a volume as a typed value supporting the bits per sample provided
func (v Volume) ToSample(bitsPerSample int) interface{} {
	switch bitsPerSample {
	case 8:
		return int8(v * 128.0)
	case 16:
		return int16(v * 32678.0)
	case 24:
		s := int32(v * 8388608.0)
		return Int24{Hi: int8(s >> 16), Lo: uint16(s & 65535)}
	case 32:
		return int32(v * 2147483648.0)
	}
	return 0
}

// ToIntSample returns a volume as an int32 value ranged to the bits per sample provided
func (v Volume) ToIntSample(bitsPerSample int) int32 {
	switch bitsPerSample {
	case 8:
		return int32(v * 128.0)
	case 16:
		return int32(v * 32678.0)
	case 24:
		return int32(v * 8388608.0)
	case 32:
		return int32(v * 2147483648.0)
	}
	return 0
}

// Apply multiplies the volume to 0..n samples, then returns an array of the results
func (v Volume) Apply(samp ...Volume) Matrix {
	o := make(Matrix, len(samp))
	for i, s := range samp {
		o[i] = s * v
	}
	return o
}

// Apply takes a volume matrix and multiplies it my incoming volumes
func (m Matrix) Apply(samp ...Volume) Matrix {
	o := make(Matrix, len(m))
	for _, s := range samp {
		for i, v := range s.Apply(m...) {
			o[i] += v
		}
	}
	return o
}
