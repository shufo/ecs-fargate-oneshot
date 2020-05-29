package tasks

import (
	"fmt"
	"os"
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
	container := c.String("container")

	latestDefinition := getLatestTaskDefinition(c, ecsSvc)
	commands := c.Args().Slice()
	service := getService(c, ecsSvc)

	results, err := ecsSvc.RunTask(&ecs.RunTaskInput{
		Count:      aws.Int64(1),
		Cluster:    aws.String(cluster),
		LaunchType: aws.String("FARGATE"),
		NetworkConfiguration: &ecs.NetworkConfiguration{
			AwsvpcConfiguration: service.NetworkConfiguration.AwsvpcConfiguration,
		},
		Overrides: &ecs.TaskOverride{
			ContainerOverrides: []*ecs.ContainerOverride{
				{
					Name:    aws.String(container),
					Command: aws.StringSlice(commands),
				},
			},
		},
		TaskDefinition: latestDefinition,
	})

	if err != nil {
		log.Fatal(err)
	}

	taskArn := *results.Tasks[0].TaskArn

	log.Info("executing tasks...\nPlease wait for tasks to finished")

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
