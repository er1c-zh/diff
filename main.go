package main

import "fmt"

func main() {
	seq := LCS("hello", "hello world")

	fmt.Printf("%s %s",
		seq.L1[seq.L1From:seq.L1To],
		seq.L2[seq.L2From:seq.L2To])

}

type Sequence struct {
	L1, L2 string
	L1From, L1To, L2From, L2To int // [from, to)
}

func LCS(l1, l2 string) Sequence {
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
			if i - 1 >= 0 && j - 1 >= 0 {
				m[i][j] += m[i - 1][j - 1]
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
	}
}
