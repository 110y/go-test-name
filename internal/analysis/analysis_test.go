package analysis_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/110y/go-test-name/internal/analysis"
)

func TestGetTestInfo(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		path     string
		pos      int
		expected *analysis.TestInfo
	}{
		"should return TestInfo includes expected TestFuncName and SubTestNames": {
			path: "package1/file1_test.go",
			pos:  178,
			expected: &analysis.TestInfo{
				TestFuncName: "TestFoo",
				SubTestNames: []string{"sub test name"},
			},
		},
		"should return TestInfo includes expected TestFuncName and SubTestNames with escaping regexp meta characters": {
			path: "package1/file1_test.go",
			pos:  229,
			expected: &analysis.TestInfo{
				TestFuncName: "TestFoo",
				SubTestNames: []string{`sub test name includes \(regexp meta characters\)`},
			},
		},
		"should return TestInfo includes expected TestFuncName and SubTestNames for ": {
			path: "package1/file1_test.go",
			pos:  676,
			expected: &analysis.TestInfo{
				TestFuncName: "TestBar",
				SubTestNames: []string{`sub test name`},
			},
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			path := fmt.Sprintf("testdata/%s", test.path)
			actual, err := analysis.GetTestInfo(context.Background(), path, test.pos)
			if err != nil {
				t.Fatalf("error: %s\n", err.Error())
			}

			if diff := cmp.Diff(test.expected, actual); diff != "" {
				t.Errorf("\n(-expected, +actual)\n%s", diff)
			}
		})
	}
}
