package client

import (
	"os"
	"testing"
)

type tCase struct {
	filename   string
	shouldFail bool
}

func TestDefSyntax(t *testing.T) {
	path := "/root/go/src/github.com/plarun/scheduler/internal/client/etc/test/def/"
	testCases := []tCase{
		// {path + "def1.txt", false},
		// {path + "def2.txt", false},
		{path + "def3.txt", false},
		{path + "def4.txt", false},
		{path + "def5.txt", false},
		{path + "def6.txt", false},
		{path + "def7.txt", false},
	}

	for tc, testCase := range testCases {
		file, err := os.Open(testCase.filename)
		if err != nil {
			panic("file not found: " + testCase.filename)
		}

		syntax := newDefinition(file)
		err = syntax.Parse()

		if err != nil && !testCase.shouldFail {
			t.Errorf("case %d: Error: %v", tc+1, err)
		}

		if err == nil && testCase.shouldFail {
			t.Errorf("case %d: should fail", tc+1)
		}
	}
}
