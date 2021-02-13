package volume

import "math"

// Volume is a mixable volume
type Volume float32

var (
	// VolumeUseInstVol tells the system to use the volume stored on the instrument
	// This is useful for trackers and other musical applications
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
	val := v.withOverflowProtection()
	switch bitsPerSample {
	case 8:
		return int8(val * 128.0)
	case 16:
		return int16(val * 32678.0)
	case 24:
		s := int32(val * 8388608.0)
		return Int24{Hi: int8(s >> 16), Lo: uint16(s & 65535)}
	case 32:
		return int32(val * 2147483648.0)
	}
	return 0
}

// ToIntSample returns a volume as an int32 value ranged to the bits per sample provided
func (v Volume) ToIntSample(bitsPerSample int) int32 {
	val := v.withOverflowProtection()
	switch bitsPerSample {
	case 8:
		return int32(val * 128.0)
	case 16:
		return int32(val * 32678.0)
	case 24:
		return int32(val * 8388608.0)
	case 32:
		return int32(val * 2147483648.0)
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

// Apply takes a volume matrix and multiplies it by incoming volumes
func (m Matrix) Apply(samp ...Volume) Matrix {
	o := make(Matrix, len(m))
	v := Matrix(samp).Sum()
	for i, s := range v.Apply(m...) {
		o[i] = s
	}
	return o
}

// ApplyInSitu takes a volume matrix and multiplies it by incoming volumes
func (m Matrix) ApplyInSitu(samp ...Volume) Matrix {
	v := Matrix(samp).Sum()
	for i, s := range v.Apply(m...) {
		m[i] = s
	}
	return m
}

// Sum sums all the elements of the Matrix and returns the resulting Volume
func (m Matrix) Sum() Volume {
	var v Volume
	for _, s := range m {
		v += s
	}
	return v
}

func (v Volume) withOverflowProtection() float64 {
	val := float64(v)
	if math.Abs(val) <= 1.0 {
		// likely case
		return val
	} else if math.Signbit(val) {
		// overflow, negative
		return -1.0
	}
	// overflow, positive
	return 1.0
}
