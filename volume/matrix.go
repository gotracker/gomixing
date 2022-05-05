package volume

// Matrix is an array of Volumes
type Matrix struct {
	StaticMatrix
	Channels int
}

// Apply takes a volume matrix and multiplies it by incoming volumes
func (m Matrix) ApplyToMatrix(mtx Matrix) Matrix {
	if m.Channels == mtx.Channels {
		// simple straight-through
		out := Matrix{
			Channels: m.Channels,
		}
		for i := 0; i < m.Channels; i++ {
			out.StaticMatrix[i] = mtx.StaticMatrix[i].ApplySingle(m.StaticMatrix[i])
		}
		return out
	}

	// more complex applications follow...

	if mtx.Channels == 1 {
		// right (mtx) is mono, so just do direct mono application
		return m.Apply(m.StaticMatrix[0])
	}

	out := Matrix{
		Channels: m.Channels,
	}

	switch m.Channels {
	case 1:
		// left (m) is mono, so do indirect mono application
		mo := mtx.AsMono()
		out.StaticMatrix[0] = m.StaticMatrix[0].ApplySingle(mo[0])
		return out
	case 2:
		// left (m) is stereo, so do indirect stereo application
		st := mtx.AsStereo()
		out.StaticMatrix[0] = m.StaticMatrix[0].ApplySingle(st[0])
		out.StaticMatrix[1] = m.StaticMatrix[1].ApplySingle(st[1])
	case 4:
		// left (m) is quad, so do indirect quad application
		qu := mtx.AsQuad()
		out.StaticMatrix[0] = m.StaticMatrix[0].ApplySingle(qu[0])
		out.StaticMatrix[1] = m.StaticMatrix[1].ApplySingle(qu[1])
		out.StaticMatrix[2] = m.StaticMatrix[2].ApplySingle(qu[2])
		out.StaticMatrix[3] = m.StaticMatrix[3].ApplySingle(qu[3])
	}
	return out
}

func (m Matrix) Apply(vol Volume) Matrix {
	var out Matrix
	out.Channels = m.Channels
	for i := 0; i < m.Channels; i++ {
		out.StaticMatrix[i] = vol.ApplySingle(m.StaticMatrix[i])
	}
	return out
}

func (m *Matrix) Accumulate(in Matrix) {
	if m.Channels == 0 {
		m.Channels = in.Channels
		copy(m.StaticMatrix[:m.Channels], in.StaticMatrix[:m.Channels])
		return
	}

	var dry StaticMatrix
	in.ToChannels(m.Channels, dry[:])
	for i := 0; i < m.Channels; i++ {
		m.StaticMatrix[i] += dry[i]
	}
}

func (m *Matrix) Assign(channels int, data []Volume) {
	m.Channels = channels
	for i := 0; i < channels; i++ {
		m.StaticMatrix[i] = data[i]
	}
}

func (m Matrix) ToChannels(channels int, out []Volume) {
	if m.Channels == channels {
		copy(out, m.StaticMatrix[0:channels])
		return
	}

	switch channels {
	default:
		copy(out, m.StaticMatrix[0:channels])
	case 1:
		mo := m.AsMono()
		copy(out, mo[:])
	case 2:
		st := m.AsStereo()
		copy(out, st[:])
	case 4:
		qu := m.AsQuad()
		copy(out, qu[:])
	}
}

// Sum sums all the elements of the Matrix and returns the resulting Volume
func (m Matrix) Sum() Volume {
	var v Volume
	for i := 0; i < m.Channels; i++ {
		v += m.StaticMatrix[i]
	}
	return v
}

func (m *Matrix) Set(ch int, vol Volume) {
	m.StaticMatrix[ch] = vol
}

func (m Matrix) Get(ch int) Volume {
	return m.StaticMatrix[ch]
}

func (m Matrix) AsMono() [1]Volume {
	var out [1]Volume
	if m.Channels == 1 {
		out[0] = m.StaticMatrix[0]
	} else {
		out[0] = m.Sum() / Volume(m.Channels)
	}
	return out
}

func (m Matrix) AsStereo() [2]Volume {
	switch m.Channels {
	case 1:
		return [2]Volume{m.StaticMatrix[0], m.StaticMatrix[0]}
	case 2:
		return [2]Volume{m.StaticMatrix[0], m.StaticMatrix[1]}
	default:
		return [2]Volume{(m.StaticMatrix[0] + m.StaticMatrix[2]) / 2.0, (m.StaticMatrix[1] + m.StaticMatrix[3]) / 2.0}
	}
}

func (m Matrix) AsQuad() [4]Volume {
	switch m.Channels {
	case 1:
		return [4]Volume{m.StaticMatrix[0], m.StaticMatrix[0], m.StaticMatrix[0], m.StaticMatrix[0]}
	case 2:
		return [4]Volume{m.StaticMatrix[0], m.StaticMatrix[1], m.StaticMatrix[0], m.StaticMatrix[1]}
	default:
		return [4]Volume{m.StaticMatrix[0], m.StaticMatrix[1], m.StaticMatrix[2], m.StaticMatrix[3]}
	}
}
