package tasks

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/briandowns/spinner"
	"github.com/shufo/ecs-fargate-oneshot/logs"
	"github.com/shufo/ecs-fargate-oneshot/utils"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func RunTask(c *cli.Context) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	ecsSvc := ecs.New(sess)

	cluster := c.String("cluster")
	latestDefinition := getLatestTaskDefinition(c, ecsSvc)
	service := getService(c, ecsSvc)

	taskDefinition := describeTaskDefinition(c, ecsSvc, *latestDefinition)

	overrides := GetTaskOverride(&GetTaskOverrideInput{
		container:      c.String("container"),
		commands:       c.Args().Slice(),
		cpu:            c.Int64("cpu"),
		memory:         c.Int64("memory"),
		taskDefinition: taskDefinition,
	})

	log.Debugln("Task Overrides: ", overrides)

	results, err := ecsSvc.RunTask(&ecs.RunTaskInput{
		Count:      aws.Int64(1),
		Cluster:    aws.String(cluster),
		LaunchType: aws.String("FARGATE"),
		NetworkConfiguration: &ecs.NetworkConfiguration{
			AwsvpcConfiguration: service.NetworkConfiguration.AwsvpcConfiguration,
		},
		Overrides:      overrides,
		TaskDefinition: latestDefinition,
	})

	if err != nil {
		log.Fatal(err)
	}

	taskArn := *results.Tasks[0].TaskArn

	log.Info("executing tasks...\nPlease wait for tasks to be finished")

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)

	if c.Bool("verbose") || c.Bool("progress") {
		s.Start()
	}

	err = ecsSvc.WaitUntilTasksStopped(&ecs.DescribeTasksInput{
		Cluster: aws.String(cluster),
		Tasks:   aws.StringSlice([]string{taskArn}),
	})

	if err != nil {
		log.Fatal(err)
	}

	if c.Bool("verbose") || c.Bool("progress") {
		s.Stop()
	}

	taskID := utils.ExtractTaskId(taskArn)

	log.Infoln("INFO: task finished")
	log.Infoln("INFO: taskId: ", taskID)

	if c.Bool("show-cloudwatch-logs") {
		logConfiguration := utils.GetLogConfigurationFromTaskDefinition(ecsSvc, latestDefinition, c.String("container"))
		logs.ShowLogs(&logs.ShowLogsInput{
			Ctx:       c,
			Sess:      sess,
			LogConfig: logConfiguration,
			TaskID:    taskID,
		})
	}

	fmt.Println(taskID)
}

func getLatestTaskDefinition(c *cli.Context, svc *ecs.ECS) *string {
	taskDefinition := c.String("task-definition")
	sort := "DESC"

	taskDefinitions, err := svc.ListTaskDefinitions(&ecs.ListTaskDefinitionsInput{
		FamilyPrefix: aws.String(taskDefinition),
		Sort:         aws.String(sort),
	})

	if err != nil {
		log.Fatal(err)
	}

	if len(taskDefinitions.TaskDefinitionArns) == 0 {
		fmt.Println("There is no enough task definitions for task definition: ", taskDefinition)
		os.Exit(1)
	}

	return taskDefinitions.TaskDefinitionArns[0]
}

func getService(c *cli.Context, svc *ecs.ECS) *ecs.Service {
	res, err := svc.DescribeServices(&ecs.DescribeServicesInput{
		Cluster:  aws.String(c.String("cluster")),
		Services: aws.StringSlice([]string{c.String("service")}),
	})

	if err != nil {
		log.Fatal(err)
	}

	return res.Services[0]
}

func describeTaskDefinition(c *cli.Context, svc *ecs.ECS, taskDefinition string) *ecs.TaskDefinition {
	res, err := svc.DescribeTaskDefinition(&ecs.DescribeTaskDefinitionInput{
		TaskDefinition: &taskDefinition,
	})

	if err != nil {
		log.Fatal(err)
	}

	return res.TaskDefinition
}

type GetTaskOverrideInput struct {
	container      string
	commands       []string
	cpu            int64
	memory         int64
	taskDefinition *ecs.TaskDefinition
}

func GetTaskOverride(input *GetTaskOverrideInput) *ecs.TaskOverride {
	container := input.container
	commands := input.commands
	cpu := input.cpu
	memory := input.memory
	taskDefinition := input.taskDefinition

	containerOverrides := []*ecs.ContainerOverride{
		{
			Name:    aws.String(container),
			Command: aws.StringSlice(commands),
		},
	}

	if cpu != 0 && memory != 0 {
		containerOverrides[0].Cpu = aws.Int64(cpu)
		containerOverrides[0].Memory = aws.Int64(memory)

		for i := 0; i < len(taskDefinition.ContainerDefinitions); i++ {
			containerDefinition := taskDefinition.ContainerDefinitions[i]
			if *containerDefinition.Name != container {
				taskDefinition.ContainerDefinitions[i].Cpu = aws.Int64(0)
				taskDefinition.ContainerDefinitions[i].Memory = aws.Int64(0)
			}

		}

		for i := 0; i < len(taskDefinition.ContainerDefinitions); i++ {
			containerDefinition := taskDefinition.ContainerDefinitions[i]
			if *containerDefinition.Name != container {
				containerOverrides = append(containerOverrides, &ecs.ContainerOverride{
					Name:   containerDefinition.Name,
					Cpu:    containerDefinition.Cpu,
					Memory: containerDefinition.Memory,
				})
			}
		}
	}

	overrides := &ecs.TaskOverride{
		ContainerOverrides: containerOverrides,
	}

	if cpu != 0 && memory != 0 {
		overrides.Cpu = aws.String(strconv.FormatInt(cpu, 10))
		overrides.Memory = aws.String(strconv.FormatInt(memory, 10))
	}

	return overrides
}
