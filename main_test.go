package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

var out io.Reader = os.Stdout

func TestCLI(t *testing.T) {
	for _, test := range []struct {
		Args   []string
		Output string
	}{
		{
			Args:   []string{"./ecs-fargate-oneshot", "--help"},
			Output: "--cluster",
		},
		{
			Args:   []string{"./ecs-fargate-oneshot", "--cluster", "app", "--service", "app", "run"},
			Output: "--task-definition",
		},
	} {
		t.Run("CLI test", func(t *testing.T) {
			os.Args = test.Args
			out = bytes.NewBuffer(nil)
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// run
			run(os.Args)

			// read stdout
			outC := make(chan string)
			go func() {
				var buf bytes.Buffer
				io.Copy(&buf, r)
				outC <- buf.String()
			}()

			// back to normal state
			w.Close()
			os.Stdout = old
			actual := <-outC

			fmt.Println("actual: ", actual)

			if strings.Contains(actual, test.Output) == false {
				fmt.Println(actual, test.Output)
				t.Errorf("expected %s, but got %s", test.Output, actual)
			}
		})
	}
}
