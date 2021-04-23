package diff3

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"strings"
	"testing"
)

func TestDo(t *testing.T) {
	a := []string{
		"1",
		"2",
		"3",
		"4",
		"4",
		"5",
		"6",
		"7",
		"8",
		"9",
	}
	b := []string{
		"1",
		"2",
		"3",
		"4",
		"5",
		"6",
		"7",
		"z",
		"8",
		"9",
		"1",
		"2",
		"3",
		"4",
		"5",
		"6",
		"7",
		"8",
		"9",
	}
	o := []string{
		"1",
		"2",
		"3",
		"4",
		"5",
		"6",
		"7",
		"8",
		"9",
		"1",
		"2",
		"3",
		"4",
		"5",
		"6",
		"7",
		"8",
		"9",
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
			lo := i.OTo - i.OFrom
			la := i.L1To - i.L1From
			lb := i.L2To - i.L2From
			aSame := lo == la
			bSame := lo == lb

			if aSame {
				for idx := 0; idx < lo; idx++ {
					if a[i.L1From + idx] != o[i.OFrom + idx] {
						aSame = false
						break
					}
				}
			}

			if bSame {
				for idx := 0; idx < lo; idx++ {
					if b[i.L2From + idx] != o[i.OFrom + idx] {
						bSame = false
						break
					}
				}
			}

			if !aSame && !bSame {
				fmt.Printf(">>>>>>>>>>>>>>>>>a\n")
				fmt.Printf("  %s\n", strings.Join(a[i.L1From:i.L1To], "\n  "))
				fmt.Printf("==================\n")
				fmt.Printf("  %s\n", strings.Join(b[i.L2From:i.L2To], "\n  "))
				fmt.Printf("<<<<<<<<<<<<<<<<<b\n")
			} else if aSame && bSame {
				panic("all same!")
			} else if aSame || bSame {
				use := a
				from := i.L1From
				to := i.L1To
				line := i.L1To - i.L1From
				anotherLine := i.L2To - i.L2From
				if aSame {
					use = b
					from = i.L2From
					to = i.L2To
					line = i.L2To - i.L2From
					anotherLine = i.L1To - i.L1From
				}

				c := color.New(color.FgCyan)
				prefix := " "
				if anotherLine == 0 {
					prefix = "+"
					c = color.New(color.FgGreen)
				} else if line == 0 {
					if aSame {
						use = a
						from = i.L1From
						to = i.L1To
						line = i.L1To - i.L1From
						anotherLine = i.L2To - i.L2From
					} else {
						use = b
						from = i.L2From
						to = i.L2To
						line = i.L2To - i.L2From
						anotherLine = i.L1To - i.L1From
					}
					prefix = "-"
					c = color.New(color.FgRed)
				}
				for _, s := range use[from:to] {
					_, err := c.Printf("%s %s\n", prefix, s)
					if err != nil {
						t.Error(err)
						return
					}
				}
			} else {
				panic("wtf")
			}
		} else {
			fmt.Printf("  %s\n", strings.Join(a[i.L1From:i.L1To], "\n  "))
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
