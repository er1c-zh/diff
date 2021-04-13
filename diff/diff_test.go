package diff

import (
	"fmt"
	"github.com/fatih/color"
	"strings"
	"testing"
)

func TestDo(t *testing.T) {
	l1 := []string{
		"1",
		"6",
		"bilibili",
		"3",
		"jd",
		"4",
		"5",
	}
	l2 := []string{
		"1",
		"2",
		"bilibili",
		"douyu",
		"3",
		"4",
		"5",
	}
	chunks := Do(l1, l2)

	for _, c := range chunks {
		if c.Conflict {
			color.Red("  <<<<<<<<<<l1\n")
			color.Red("  %s\n", strings.Join(l1[c.L1From:c.L1To], "\n  "))
			color.Red("  ============\n")
			color.Red("  %s\n", strings.Join(l2[c.L2From:c.L2To], "\n  "))
			color.Red("  <<<<<<<<<<l2\n")
		} else {
			if c.L1Empty {
				color.Green("+ %s\n", strings.Join(l2[c.L2From:c.L2To], "\n+ "))
			} else if c.L2Empty {
				color.Green("+ %s\n", strings.Join(l1[c.L1From:c.L1To], "\n+ "))
			} else {
				fmt.Printf("  %s\n", strings.Join(l1[c.L1From:c.L1To], "\n  "))
			}
		}
	}
}


