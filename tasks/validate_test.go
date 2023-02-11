package tasks

import (
	"testing"
)

func TestValidateCpuAndMemory(t *testing.T) {
	cases := []struct {
		cpu    uint64
		memory uint64
		expect bool
	}{
		{
			cpu:    0,
			memory: 0,
			expect: true,
		},
		{
			cpu:    128,
			memory: 128,
			expect: false,
		},
		{
			cpu:    256,
			memory: 512,
			expect: true,
		},
		{
			cpu:    1024,
			memory: 2048,
			expect: true,
		},
	}

	for _, c := range cases {
		result := ValidateCombinationOfCpuAndMemory(c.cpu, c.memory)

		if result != c.expect {
			t.Fail()
		}
	}
}
