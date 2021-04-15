package diff3

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestDo(t *testing.T) {
	a := []string{
		"1",
		"2",
		"3",
		"6",
		"5",
	}
	b := []string{
		"1",
		"2",
		"3",
		"7",
		"5",
	}
	o := []string{
		"1",
		"2",
		"3",
		"4",
		"5",
	}
	l := Do(a, b, o)
	for _, i := range l {
		j, _ := json.Marshal(i)
		t.Logf("%s", string(j))
	}
	if false {
		return
	}
	for _, i := range l {
		if i.Conflict {
			fmt.Printf(">>>>>>>>>>>>>>>>>a\n")
			fmt.Printf("  %s\n", strings.Join(a[i.L1From:i.L1To], "\n  "))
			fmt.Printf("==================\n")
			fmt.Printf("  %s\n", strings.Join(b[i.L1From:i.L1To], "\n  "))
			fmt.Printf("<<<<<<<<<<<<<<<<<b\n")
		} else {
			if i.UseL1 || i.UseL2 {
				split := "+"
				if (i.UseL1 && i.L2Empty) || (i.UseL2 && i.L1Empty) {
					split = "-"
				}
				if i.UseL1 {
					fmt.Printf(split+" %s\n", strings.Join(a[i.L1From:i.L1To], "\n"+split+" "))
				} else {
					fmt.Printf(split+" %s\n", strings.Join(b[i.L1From:i.L1To], "\n"+split+" "))
				}
			} else {
				fmt.Printf("  %s\n", strings.Join(a[i.L1From:i.L1To], "\n  "))
			}
		}
	}
}

func TestDiff2(t *testing.T) {
	l := Diff2(
		[]string{
			"1",
			"3",
			"6",
			"5",
			"7",
		}, []string{
			"1",
			"2",
			"3",
			"4",
			"5",
		})

	t.Logf("%v", l)

	for _, i := range l {
		t.Logf("%v", i)
	}
}
