package volume

// StaticMatrix is an array of Volumes
type StaticMatrix [4]Volume

func (m StaticMatrix) Channels() int {
	return 4
}

// Apply takes a volume matrix and multiplies it by incoming volumes
func (m StaticMatrix) ApplyToMatrix(mtx StaticMatrix) StaticMatrix {
	return StaticMatrix{
		mtx[0].ApplySingle(m[0]),
		mtx[1].ApplySingle(m[1]),
		mtx[2].ApplySingle(m[2]),
		mtx[3].ApplySingle(m[3]),
	}
}

func (m StaticMatrix) Apply(vol Volume) StaticMatrix {
	return StaticMatrix{
		vol.ApplySingle(m[0]),
		vol.ApplySingle(m[1]),
		vol.ApplySingle(m[2]),
		vol.ApplySingle(m[3]),
	}
}

func (m *StaticMatrix) Accumulate(in StaticMatrix) {
	m[0] += in[0]
	m[1] += in[1]
	m[2] += in[2]
	m[3] += in[3]
}

func (m *StaticMatrix) Assign(channels int, data []Volume) {
	switch channels {
	case 1:
		v := data[0]
		m[0] = v
		m[1] = v
		m[2] = v
		m[3] = v
	case 2:
		l := data[0]
		r := data[1]
		m[0] = l
		m[1] = r
		m[2] = l
		m[3] = r
	case 3:
		l := data[0]
		r := data[1]
		cr := data[2]
		m[0] = l
		m[1] = r
		m[2] = cr
		m[3] = cr
	case 4:
		m[0] = data[0]
		m[1] = data[1]
		m[2] = data[2]
		m[3] = data[3]
	}
}

func (m StaticMatrix) ToChannels(channels int) []Volume {
	switch channels {
	default:
		return nil
	case 1:
		return []Volume{m.Sum() / 4.0}
	case 2:
		l := (m[0] + m[2]) / 2.0
		r := (m[1] + m[3]) / 2.0
		return []Volume{l, r}
	case 4:
		return m[:]
	}
}

// Sum sums all the elements of the StaticMatrix and returns the resulting Volume
func (m StaticMatrix) Sum() Volume {
	return m[0] + m[1] + m[2] + m[3]
}

func (m *StaticMatrix) Set(ch int, vol Volume) {
	m[ch] = vol
}

func (m StaticMatrix) Get(ch int) Volume {
	return m[ch]
}

func (m StaticMatrix) AsMono() [1]Volume {
	return [1]Volume{(m[0] + m[1] + m[2] + m[3]) / 4.0}
}

func (m StaticMatrix) AsStereo() [2]Volume {
	return [2]Volume{(m[0] + m[2]) / 2.0, (m[1] + m[3]) / 2.0}
}

func (m StaticMatrix) AsQuad() [4]Volume {
	return m
}
