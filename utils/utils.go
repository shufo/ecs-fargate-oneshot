package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	log "github.com/sirupsen/logrus"
)

func GetTaskDefinitionFromTaskID(cluster string, taskID string) string {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := ecs.New(sess)

	resp, err := svc.DescribeTasks(&ecs.DescribeTasksInput{Cluster: aws.String(cluster), Tasks: aws.StringSlice([]string{taskID})})

	if err != nil {
		log.Fatal(err)
	}

	if len(resp.Tasks) == 0 {
		fmt.Println("There is no finished tasks for taskId: ", taskID)
		os.Exit(1)
	}

	return *resp.Tasks[0].TaskDefinitionArn
}

func GetLogConfigurationFromTaskDefinition(svc *ecs.ECS, taskDefinitionArn *string, container string) *ecs.LogConfiguration {
	res, err := svc.DescribeTaskDefinition(&ecs.DescribeTaskDefinitionInput{
		TaskDefinition: taskDefinitionArn,
	})

	if err != nil {
		log.Fatal(err)
	}

	for _, d := range res.TaskDefinition.ContainerDefinitions {
		if *d.Name != container {
			continue
		}

		if *d.LogConfiguration.LogDriver != "awslogs" {
			continue
		}

		return d.LogConfiguration
	}

	log.Fatal("There is no matched container name. Please check your container definitions.")

	return nil
}

func ExtractTaskId(taskArn string) string {
	if !strings.Contains(taskArn, "/") {
		log.Fatal("seems passed parameter is not task ARN")
	}

	return strings.Split(taskArn, "/")[2]
}
