package diff

type Chunk struct {
	L1From, L1To, L2From, L2To int // [from, to)
	L1Empty, L2Empty           bool
	Conflict                   bool
}

func Do(l1, l2 []string) []Chunk {
	var f func(from1, to1, from2, to2 int) []Chunk

	f = func(from1, to1, from2, to2 int) []Chunk {
		_l1 := l1[from1:to1]
		_l2 := l2[from2:to2]
		if len(_l1) == 0 && len(_l2) == 0 {
			return nil
		}
		if len(_l1) == 0 {
			return []Chunk{
				{
					L1From:   -1,
					L1To:     -1,
					L2From:   from2,
					L2To:     to2,
					L1Empty:  true,
					L2Empty:  false,
					Conflict: false,
				},
			}
		}
		if len(_l2) == 0 {
			return []Chunk{
				{
					L1From:   from1,
					L1To:     to1,
					L2From:   -1,
					L2To:     -1,
					L1Empty:  false,
					L2Empty:  true,
					Conflict: false,
				},
			}
		}
		seq := LCS(_l1, _l2)
		var chunk Chunk
		if seq.Len == 0 {
			chunk = Chunk{
				L1From:   from1,
				L1To:     to1,
				L2From:   from2,
				L2To:     to2,
				L1Empty:  false,
				L2Empty:  false,
				Conflict: true,
			}
			return []Chunk{
				chunk,
			}
		} else {
			chunk = Chunk{
				L1From:   from1 + seq.L1From,
				L1To:     from1 + seq.L1To,
				L2From:   from2 + seq.L2From,
				L2To:     from2 + seq.L2To,
				L1Empty:  false,
				L2Empty:  false,
				Conflict: false,
			}
		}

		result := make([]Chunk, 0)
		result = append(result, f(from1, from1 + seq.L1From, from2, from2 + seq.L2From)...)
		result = append(result, chunk)
		result = append(result, f(from1 + seq.L1To, to1, from2 + seq.L2To, to2)...)
		return result
	}

	return f(0, len(l1), 0, len(l2))
}

type Sequence struct {
	L1, L2                     []string
	L1From, L1To, L2From, L2To int // [from, to)
	Len                        int
}

func LCS(l1, l2 []string) Sequence {
	_l1 := len(l1)
	_l2 := len(l2)
	m := make([][]int, 0)
	for i := 0; i < _l1; i++ {
		m = append(m, make([]int, _l2))
	}
	maxI := 0
	maxJ := 0
	max := 0

	for i := 0; i < _l1; i++ {
		for j := 0; j < _l2; j++ {
			if l1[i] != l2[j] {
				continue
			}
			m[i][j] = 1
			if i-1 >= 0 && j-1 >= 0 {
				m[i][j] += m[i-1][j-1]
			}
			if m[i][j] > max {
				max, maxI, maxJ = m[i][j], i, j
			}
		}
	}
	return Sequence{
		L1:     l1,
		L2:     l2,
		L1From: maxI + 1 - max,
		L1To:   maxI + 1,
		L2From: maxJ + 1 - max,
		L2To:   maxJ + 1,
		Len:    max,
	}
}