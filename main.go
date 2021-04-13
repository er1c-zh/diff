package main

import (
	"flag"
	"fmt"
	"github.com/er1c-zh/diff/diff"
	"github.com/fatih/color"
	"io/ioutil"
	"strings"
)

var (
	fileArgs string
)

func _flag() {
	flag.StringVar(&fileArgs, "f", "", "file1,file2")
}

func main() {
	_flag()
	flag.Parse()
	if fileArgs == "" {
		flag.Usage()
		return
	}
	files := strings.Split(fileArgs, ",")
	if len(files) != 2 {
		flag.Usage()
		return
	}

	var l1, l2 []string

	parseFile := func(p string, t *[]string) error {
		color.Cyan("file: %s\n", p)
		b, err := ioutil.ReadFile(p)
		if err != nil {
			return err
		}
		*t = strings.Split(string(b), "\n")
		return nil
	}

	for _, s := range []struct{
		p string
		t *[]string
	}{
		{files[0], &l1},
		{files[1], &l2},
	}{
		err := parseFile(s.p, s.t)
		if err != nil {
			fmt.Printf("Err: %s\n", err.Error())
			flag.Usage()
			return
		}
	}

	chunks := diff.Do(l1, l2)

	for _, c := range chunks {
		if c.Conflict {
			color.Red(">>>>>>>>>>>>file 1\n")
			color.Red("  %s\n", strings.Join(l1[c.L1From:c.L1To], "\n  "))
			color.Red("==============\n")
			color.Red("  %s\n", strings.Join(l2[c.L2From:c.L2To], "\n  "))
			color.Red("<<<<<<<<<<<<file 2\n")
		} else {
			if c.L1Empty {
				color.Green("+ %s\n", strings.Join(l2[c.L2From:c.L2To], "\n+ "))
			} else if c.L2Empty {
				color.Green("+ %s\n", strings.Join(l1[c.L1From:c.L1To], "\n+ "))
			} else {
				color.Cyan("= %s\n", strings.Join(l1[c.L1From:c.L1To], "\n= "))
			}
		}
	}
}


