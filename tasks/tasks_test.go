package tasks

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
)

func TestGetTaskOverride(t *testing.T) {
	cases := []struct {
		container      string
		commands       []string
		cpu            int64
		memory         int64
		taskDefinition *ecs.TaskDefinition
		expect         *ecs.TaskOverride
	}{
		{
			container: "app",
			commands:  []string{"echo", "test"},
			cpu:       256,
			memory:    512,
			taskDefinition: &ecs.TaskDefinition{
				ContainerDefinitions: []*ecs.ContainerDefinition{
					{
						Cpu:    aws.Int64(256),
						Memory: aws.Int64(512),
						Name:   aws.String("app"),
					},
					{
						Cpu:    aws.Int64(256),
						Memory: aws.Int64(512),
						Name:   aws.String("nginx"),
					},
				},
			},
			expect: &ecs.TaskOverride{
				ContainerOverrides: []*ecs.ContainerOverride{
					{
						Cpu:     aws.Int64(256),
						Memory:  aws.Int64(512),
						Name:    aws.String("app"),
						Command: aws.StringSlice([]string{"echo", "test"}),
					},
					{
						Cpu:    aws.Int64(0),
						Memory: aws.Int64(0),
						Name:   aws.String("nginx"),
					},
				},
				Cpu:    aws.String("256"),
				Memory: aws.String("512"),
			},
		},
		{
			container: "app",
			commands:  []string{"echo", "test"},
			cpu:       1024,
			memory:    2048,
			taskDefinition: &ecs.TaskDefinition{
				ContainerDefinitions: []*ecs.ContainerDefinition{
					{
						Cpu:    aws.Int64(256),
						Memory: aws.Int64(512),
						Name:   aws.String("app"),
					},
					{
						Cpu:    aws.Int64(256),
						Memory: aws.Int64(512),
						Name:   aws.String("nginx"),
					},
					{
						Cpu:    aws.Int64(256),
						Memory: aws.Int64(512),
						Name:   aws.String("php-fpm"),
					},
				},
			},
			expect: &ecs.TaskOverride{
				ContainerOverrides: []*ecs.ContainerOverride{
					{
						Cpu:     aws.Int64(1024),
						Memory:  aws.Int64(2048),
						Name:    aws.String("app"),
						Command: aws.StringSlice([]string{"echo", "test"}),
					},
					{
						Cpu:    aws.Int64(0),
						Memory: aws.Int64(0),
						Name:   aws.String("nginx"),
					},
					{
						Cpu:    aws.Int64(0),
						Memory: aws.Int64(0),
						Name:   aws.String("php-fpm"),
					},
				},
				Cpu:    aws.String("1024"),
				Memory: aws.String("2048"),
			},
		},
	}

	for _, c := range cases {
		result := GetTaskOverride(&GetTaskOverrideInput{
			container:      c.container,
			commands:       c.commands,
			cpu:            c.cpu,
			memory:         c.memory,
			taskDefinition: c.taskDefinition,
		})

		if !reflect.DeepEqual(result, c.expect) {
			t.Fail()
		}
	}
}
