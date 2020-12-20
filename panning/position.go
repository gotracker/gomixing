package panning

import "math"

// Position is stored as polar coordinates
// with Angle of 0 radians being calculated from right
// and >0 rotating counter-clockwise from that point
type Position struct {
	Angle    float32
	Distance float32
}

var (
	// CenterAhead is the position directly ahead of the listener
	CenterAhead = MakeStereoPosition(0.5, 0, 1)
)

// MakeStereoPosition creates a stereo panning position based on a linear interpolation between `leftValue` and `RightValue`
func MakeStereoPosition(value float32, leftValue float32, rightValue float32) Position {
	if leftValue == rightValue {
		panic("leftValue and rightValue should be distinct")
	}
	d := float64(rightValue - leftValue)
	t := (d - float64(value)) / d
	// we're using a 2d rotation matrix to calcuate the left and right channels, so we really want the half angle
	prad := t * math.Pi / 2.0

	return Position{
		Angle:    float32(prad),
		Distance: 1.0,
	}
}
