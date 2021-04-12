package main

import "fmt"

func main() {
	fmt.Printf("hello")
}

type Item interface {
	EqualTo(Item) bool
	Val() interface{}
}

type Sequence struct {
	L1, L2 []Item
	L1From, L1To, L2From, L2To int
}

func MaxSameSubsequence(l1, l2 []Item) Sequence {
	_l1 := len(l1)
	_l2 := len(l2)
	type seq struct {
		L1From, L1To, L2From, L2To int
	}
	dp := make([][]seq, 0)
	for i := 0; i < len(l1); i++ {
		dp = append(dp, make([]seq, len(l2)))
	}

	dp[0][0] = seq{
		L1From: 0,
		L1To:   0,
		L2From: 0,
		L2To:   0,
	}

	for i := 0; i < _l1; i++ {
		for j := 0; j < _l2; j++ {

		}
	}

	return Sequence{
		L1:     l1,
		L2:     l2,
		L1From: dp[_l1][_l2].L1From,
		L1To:   dp[_l1][_l2].L1To,
		L2From: dp[_l1][_l2].L2From,
		L2To:   dp[_l1][_l2].L2To,
	}
}
