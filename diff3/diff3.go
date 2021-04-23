package diff3

const (
	nullIdx = -1
)

type Chunk struct {
	L1From, L1To, L2From, L2To int // [from, to)
	OFrom, OTo                 int
	L1Empty, L2Empty           bool // 任意一个是否为空
	Conflict                   bool // 是否是real-conflict
}

type MatcherItem struct {
	L1Idx, L2Idx     int
	L1Empty, L2Empty bool
	Conflict         bool
}

func Do(a, b, o []string) []Chunk {
	ao := getDiff2Matcher(a, o)
	bo := getDiff2Matcher(b, o)

	result := make([]Chunk, 0)
	ia, ib, io := 0, 0, 0
	stable := true

	i := 0

	for {
		if len(o) <= io+i {
			stable = true
			goto done
		}
		if ao.Match(ia+i, io+i) && bo.Match(ib+i, io+i) {
			i++
			continue
		}
		// find conflict
		if i == 0 {
			// 没有stable chunk
			_i := 0
			for {
				if io+_i >= len(o) {
					stable = false
					goto done
				}
				_ia, aOk := ao.Find(io + _i)
				_ib, bOk := bo.Find(io + _i)
				if aOk && bOk {
					// find stable
					result = append(result, Chunk{
						L1From:   ia,
						L1To:     _ia,
						L2From:   ib,
						L2To:     _ib,
						OFrom:    io,
						OTo:      io + _i,
						L1Empty:  _ia == _ib,
						L2Empty:  _ia == _ib,
						Conflict: true,
					})
					ia = _ia
					ib = _ib
					io += _i
					i = 0
					break
				}
				_i += 1
			}
		} else {
			// 输出stable chunk
			result = append(result, Chunk{
				L1From:   ia,
				L1To:     ia + i,
				L2From:   ib,
				L2To:     ib + i,
				OFrom:    io,
				OTo:      io + i,
				L1Empty:  i == 0,
				L2Empty:  i == 0,
				Conflict: false,
			})
			ia += i
			ib += i
			io += i
			i = 0 // reset
			continue
		}
	}

done:
	result = append(result, Chunk{
		L1From:   ia,
		L1To:     len(a),
		L2From:   ib,
		L2To:     len(b),
		OFrom:    io,
		OTo:      len(o),
		L1Empty:  ia == len(a),
		L2Empty:  ib == len(b),
		Conflict: !stable,
	})

	return result
}

type Matcher interface {
	Match(i, j int) bool
}

type matcher struct {
	oMap map[int]MatcherItem
	aMap map[int]MatcherItem
}

func (m matcher) Match(i1, io int) bool {
	r := false
	if item, ok := m.oMap[io]; ok &&
		!item.Conflict &&
		item.L1Idx == i1 {
		return true
	}
	return r
}

func (m matcher) Find(io int) (int, bool) {
	if item, ok := m.oMap[io]; ok &&
		!item.Conflict {
		return item.L1Idx, true
	}
	return 0, false
}

func getDiff2Matcher(a, o []string) matcher {
	ao := Diff2(a, o)
	m := matcher{
		oMap: map[int]MatcherItem{},
		aMap: map[int]MatcherItem{},
	}
	for _, c := range ao {
		if c.L1Idx != nullIdx {
			m.aMap[c.L1Idx] = c
		}
		if c.L2Idx != nullIdx {
			m.oMap[c.L2Idx] = c
		}
	}
	return m
}

func Diff2(l1, l2 []string) []MatcherItem {
	var f func(l1From, l1To, l2From, l2To int) []MatcherItem
	f = func(l1From, l1To, l2From, l2To int) []MatcherItem {
		_l1 := l1[l1From:l1To]
		_l2 := l2[l2From:l2To]
		if len(_l1) == 0 && len(_l2) == 0 {
			return nil
		}
		if len(_l1) == 0 {
			result := make([]MatcherItem, 0, l2To-l2From)
			for i := l2From; i < l2To; i++ {
				result = append(result, MatcherItem{
					L1Idx:    nullIdx,
					L2Idx:    i,
					L1Empty:  true,
					L2Empty:  false,
					Conflict: true,
				})
			}
			return result
		}
		if len(_l2) == 0 {
			result := make([]MatcherItem, 0, l1To-l1From)
			for i := l1From; i < l1To; i++ {
				result = append(result, MatcherItem{
					L1Idx:    i,
					L2Idx:    nullIdx,
					L1Empty:  false,
					L2Empty:  true,
					Conflict: true,
				})
			}
			return result
		}

		seq := LCS(_l1, _l2)
		if seq.Len == 0 {
			result := make([]MatcherItem, 0, l1To-l1From+l2To-l2From)
			for i := l1From; i < l1To; i++ {
				result = append(result, MatcherItem{
					L1Idx:    i,
					L2Idx:    nullIdx,
					L1Empty:  false,
					L2Empty:  true,
					Conflict: true,
				})
			}
			for i := l2From; i < l2To; i++ {
				result = append(result, MatcherItem{
					L1Idx:    nullIdx,
					L2Idx:    i,
					L1Empty:  true,
					L2Empty:  false,
					Conflict: true,
				})
			}
			return result
		}
		result := make([]MatcherItem, 0, l1To-l1From+l2To-l2From)
		result = append(result, f(l1From, l1From+seq.L1From, l2From, l2From+seq.L2From)...)
		for i := 0; i < seq.Len; i++ {
			result = append(result, MatcherItem{
				L1Idx:    l1From + seq.L1From + i,
				L2Idx:    l2From + seq.L2From + i,
				L1Empty:  false,
				L2Empty:  false,
				Conflict: false,
			})
		}
		result = append(result, f(l1From+seq.L1To, l1To, l2From+seq.L2To, l2To)...)
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
