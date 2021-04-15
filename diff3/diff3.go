package diff3

const (
	nullIdx = -1
)

type Section struct {
	L1From, L1To, L2From, L2To int  // [from, to)
	L1Empty, L2Empty           bool // 任意一个是否为空
	Conflict                   bool // 是否是real-conflict
	UseL1, UseL2               bool // 是否是使用某一种代替另一个
}

type Chunk struct {
	L1Idx, L2Idx     int
	L1Empty, L2Empty bool
	Conflict         bool
}

func Do(a, b, o []string) []Section {
	ao := getDiff2Matcher(a, o)
	bo := getDiff2Matcher(b, o)

	result := make([]Section, 0)
	ia, ib, io := 0, 0, 0
	waitConflict := true
	done := false

	for {
		_ia, _ib, _io := ia, ib, io

		if waitConflict {
			for {
				ca, ok := ao.oMap[io]
				if !ok {
					done = true
					break
				}
				cb, ok := bo.oMap[io]
				if !ok {
					done = true
					break
				}
				if ca.Conflict || cb.Conflict {
					break
				} else {
					ia++
					ib++
					io++
				}
			}
		} else {
			for {
				ca, ok := ao.oMap[io]
				if !ok {
					done = true
					break
				}
				cb, ok := bo.oMap[io]
				if !ok {
					done = true
					break
				}
				if !ca.Conflict && !cb.Conflict {
					ia = ca.L1Idx
					ib = cb.L1Idx
					break
				} else {
					io++
				}
			}
		}

		if done {
			ia = len(a)
			ib = len(b)
			io = len(o)
		}

		if waitConflict {
			result = append(result, Section{
				L1From:   _ia,
				L1To:     ia,
				L2From:   _ib,
				L2To:     ib,
				L1Empty:  ia == _ia,
				L2Empty:  ib == _ib,
				Conflict: false,
			})
		} else {
			__ia, __ib := _ia, _ib
			realConflict := false
			useA := true
			for _io < io {
				_realConflict := realConflict
				_useA := useA

				ca, ok := ao.oMap[_io]
				if !ok {
					panic("fail 1")
				}
				cb, ok := bo.oMap[_io]
				if !ok {
					panic("fail 2")
				}
				__ia = ca.L1Idx
				__ib = cb.L2Idx
				if ca.Conflict && cb.Conflict {
					if !realConflict {
						_realConflict = true
						goto merge
					}
					goto _continue
				} else if ca.Conflict || cb.Conflict {
					if realConflict {
						_realConflict = false
						_useA = ca.Conflict
						goto merge
					}
					// free-conflict
					if (ca.Conflict && useA) || (cb.Conflict && !useA) {
						goto _continue
					} else {
						_useA = !_useA
						goto merge
					}
				} else {
					// unexpected status
					panic("fail 3")
				}
			merge:
				// todo
				if _ia == __ia && _ib == __ib {
					// 没有数据，直接切换模式
					goto next
				}
				if realConflict {
					result = append(result, Section{
						L1From:   _ia,
						L1To:     __ia,
						L2From:   _ib,
						L2To:     __ib,
						L1Empty:  _ia == __ia,
						L2Empty:  _ib == __ib,
						Conflict: true,
						UseL1:    false,
						UseL2:    false,
					})
				} else {
					if useA {
						result = append(result, Section{
							L1From:   _ia,
							L1To:     __ia,
							L2From:   _ib,
							L2To:     __ib,
							L1Empty:  _ia == __ia,
							L2Empty:  _ib == __ib,
							Conflict: false,
							UseL1:    true,
							UseL2:    false,
						})
					} else {
						result = append(result, Section{
							L1From:   _ia,
							L1To:     __ia,
							L2From:   _ib,
							L2To:     __ib,
							L1Empty:  _ia == __ia,
							L2Empty:  _ib == __ib,
							Conflict: false,
							UseL1:    false,
							UseL2:    true,
						})
					}
				}
				_io++
			next:
				_ia, _ib = __ia, __ib
				realConflict = _realConflict
				useA = _useA
				continue
			_continue:
				if _io+1 == io {
					__ia, __ib = ia, ib
					goto merge
				}
				_io++
			}
		}
		waitConflict = !waitConflict
		_ia, _ib, _io = ia, ib, io
		if done {
			break
		} else {
			continue
		}
	}

	return result
}

type Matcher interface {
	Match(i, j int) bool
}

type matcher struct {
	oMap map[int]Chunk
	aMap map[int]Chunk
}

func getDiff2Matcher(a, o []string) matcher {
	ao := Diff2(a, o)
	m := matcher{
		oMap: map[int]Chunk{},
		aMap: map[int]Chunk{},
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

func Diff2(l1, l2 []string) []Chunk {
	var f func(l1From, l1To, l2From, l2To int) []Chunk
	f = func(l1From, l1To, l2From, l2To int) []Chunk {
		_l1 := l1[l1From:l1To]
		_l2 := l2[l2From:l2To]
		if len(_l1) == 0 && len(_l2) == 0 {
			return nil
		}
		if len(_l1) == 0 {
			result := make([]Chunk, 0, l2To-l2From)
			for i := l2From; i < l2To; i++ {
				result = append(result, Chunk{
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
			result := make([]Chunk, 0, l1To-l1From)
			for i := l1From; i < l1To; i++ {
				result = append(result, Chunk{
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
			result := make([]Chunk, 0, l1To-l1From+l2To-l2From)
			for i := l1From; i < l1To; i++ {
				result = append(result, Chunk{
					L1Idx:    i,
					L2Idx:    nullIdx,
					L1Empty:  false,
					L2Empty:  true,
					Conflict: true,
				})
			}
			for i := l2From; i < l2To; i++ {
				result = append(result, Chunk{
					L1Idx:    nullIdx,
					L2Idx:    i,
					L1Empty:  true,
					L2Empty:  false,
					Conflict: true,
				})
			}
			return result
		}
		result := make([]Chunk, 0, l1To-l1From+l2To-l2From)
		result = append(result, f(l1From, l1From+seq.L1From, l2From, l2From+seq.L2From)...)
		for i := 0; i < seq.Len; i++ {
			result = append(result, Chunk{
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
