package subcli

import (
	"testing"
)

func compare(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	if (a == nil) != (b == nil) {
		return false
	}

	b = b[:len(a)] // this line is the key
	for i, v := range a {
		if v != b[i] { // here is no bounds checking for b[i]
			return false
		}
	}

	return true
}

func testingData() *SubCli {
	s := New(Program{"testing", "0.0", "This is a test"})

	s.AddCmd(SubCommand{"test4", "test4 desc", "test4 help", nil, nil})
	s.AddCmd(SubCommand{"test2", "test2 desc", "test2 help", nil, nil})
	s.AddCmd(SubCommand{"test3", "test3 desc", "test3 help", nil, nil})
	s.AddCmd(SubCommand{"test1", "test1 desc", "test1 help", nil, nil})
	return s
}

func TestSubLess(t *testing.T) {
	subs := subCommands{
		SubCommand{"test1", "test1 desc", "test1 help", nil, nil},
		SubCommand{"test2", "test2 desc", "test2 help", nil, nil},
		SubCommand{"test3", "test3 desc", "test3 help", nil, nil},
		SubCommand{"test4", "test4 desc", "test4 help", nil, nil},
	}

	err := false
	chLess := func(i, j int, ans bool) {
		if subs.Less(i, j) != ans {
			t.Logf("Subless: [%s < %s] expected: %v got: %v",
				subs[i].Command,
				subs[j].Command,
				subs.Less(i, j),
				ans,
			)
			err = true
		}
	}
	chLess(0, 3, true)
	chLess(3, 1, false)
	chLess(2, 3, true)
	chLess(0, 1, true)
	chLess(1, 0, false)
	if err {
		t.Error("subless errors")
	}

}
