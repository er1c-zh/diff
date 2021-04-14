package diff3

import (
	"encoding/json"
	"testing"
)

func TestDo(t *testing.T) {
	l := Do([]string{
		"1",
		"2",
		"3",
		"4",
		"5",
	}, []string{
		"1",
		"2",
		"3",
		"5",
		"5",
	}, []string{
		"1",
		"2",
		"3",
		"4",
		"5",
	})
	for _, i := range l {
		j, _ := json.Marshal(i)
		t.Logf("%s", string(j))
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
