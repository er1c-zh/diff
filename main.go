package main

import (
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
		if i.Conflict {
			color.Red(">>>>>>>>>>>>>>>>>a\n")
			color.Red("  %s\n", strings.Join(a[i.L1From:i.L1To], "\n  "))
			color.Red("==================\n")
			color.Red("  %s\n", strings.Join(b[i.L2From:i.L2To], "\n  "))
			color.Red("<<<<<<<<<<<<<<<<<b\n")
		} else {
			if (i.UseL1 || i.UseL2) && (i.L1Empty || i.L2Empty) {
				f := color.Green
				split := "+"
				if (i.UseL1 && i.L1Empty) || (i.UseL2 && i.L2Empty) {
					split = "-"
					f = color.Red
				} else if i.UseL1 {
					f(split+" %s\n", strings.Join(a[i.L1From:i.L1To], "\n"+split+" "))
				} else {
					f(split+" %s\n", strings.Join(b[i.L2From:i.L2To], "\n"+split+" "))
				}
			} else {
				if i.UseL1 {
					color.Cyan("  %s\n", strings.Join(a[i.L1From:i.L1To], "\n  "))
				} else if i.UseL2 {
					color.Cyan("  %s\n", strings.Join(b[i.L2From:i.L2To], "\n  "))
				} else {
					fmt.Printf("  %s\n", strings.Join(a[i.L1From:i.L1To], "\n  "))
				}
			}
		}
	}
}
