package main

import (
	"encoding/json"
	"fmt"
	"github.com/er1c-zh/diff/diff3"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
	"strings"
)

func usage() {
	fmt.Printf("usage: diff3 file_a file_old file_b\n")
}

func main() {
	args := os.Args
	if len(args) != 4 {
		fmt.Printf("Fail: invalid params.\n")
		usage()
		return
	}
	var a, b, o []string
	files := args[1:]
	for path, ptr := range map[string]*[]string{
		files[0]: &a,
		files[1]: &o,
		files[2]: &b,
	} {
		tmp, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Printf("Fail: can't read file %s: %s\n", files[0], err.Error())
			usage()
			return
		}
		*ptr = strings.Split(string(tmp), "\n")
	}

	l := diff3.Do(a, b, o)

	for _, i := range l {
		j, _ := json.Marshal(i)
		fmt.Printf("%s\n", string(j))
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
				f := color.HiRed
				f(">>>>>>>>>>>>>>>>>a\n")
				f("  %s\n", strings.Join(a[i.L1From:i.L1To], "\n  "))
				f("==================\n")
				f("  %s\n", strings.Join(b[i.L2From:i.L2To], "\n  "))
				f("<<<<<<<<<<<<<<<<<b\n")
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

				f := color.Cyan
				prefix := " "
				if anotherLine == 0 {
					prefix = "+"
					f = color.Green
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
					f = color.Red
				}
				for _, s := range use[from:to] {
					f("%s %s\n", prefix, s)
				}
			} else {
				panic("wtf")
			}
		} else {
			fmt.Printf("  %s\n", strings.Join(a[i.L1From:i.L1To], "\n  "))
		}
	}
}
