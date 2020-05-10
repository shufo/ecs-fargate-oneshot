package utils

import (
	"testing"
)

func TestExtractTaskId(t *testing.T) {
	cases := []struct {
		param  string
		expect string
	}{
		{
			param:  "foo/bar",
			expect: "bar",
		},
		{
			param:  "aws:arn:~~/dea2ccc0-da13-462d-88ce-f65aa7764d98",
			expect: "dea2ccc0-da13-462d-88ce-f65aa7764d98",
		},
	}

	for _, c := range cases {
		taskID := ExtractTaskId(c.param)

		if taskID != c.expect {
			t.Fail()
		}
	}
}
