package logs

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/shufo/ecs-fargate-oneshot/utils"
	"github.com/urfave/cli/v2"
)

type ShowLogsInput struct {
	Ctx       *cli.Context
	Sess      *session.Session
	LogConfig *ecs.LogConfiguration
	TaskID    string
}

type LogConfiguration struct {
	LogGroupName    *string
	LogStreamPrefix *string
	Container       string
	TaskID          string
}

func ShowLogs(input *ShowLogsInput) {
	logsSvc := cloudwatchlogs.New(input.Sess)

	logConfig := LogConfiguration{
		LogGroupName:    input.LogConfig.Options["awslogs-group"],
		LogStreamPrefix: input.LogConfig.Options["awslogs-stream-prefix"],
		Container:       input.Ctx.String("container"),
		TaskID:          input.TaskID,
	}

	logStreamName := getLogStreamName(logConfig)

	resp, err := getLogEvents(&GetLogEventsInput{
		svc:       logsSvc,
		ctx:       input.Ctx,
		logConfig: logConfig,
		taskID:    input.TaskID,
	})

	if err != nil {
		fmt.Println("Got error getting log events:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("Event messages for stream " + logStreamName + " in log group LOG-GROUP-NAME: " + *input.LogConfig.Options["awslogs-group"])

	gotToken := ""
	nextToken := ""

	for {
		gotToken = nextToken
		nextToken = *resp.NextForwardToken

		if gotToken == nextToken {
			break
		}

		for _, event := range resp.Events {
			fmt.Println("  ", *event.Message)
		}

		resp, _ = getLogEvents(&GetLogEventsInput{
			svc:       logsSvc,
			ctx:       input.Ctx,
			logConfig: logConfig,
			taskID:    input.TaskID,
			nextToken: nextToken,
		})
	}
}

func RunShowLogsWithTaskId(c *cli.Context) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	ecsSvc := ecs.New(sess)

	taskID := c.String("task-id")

	// get task id from passed argments if flag not exists
	if taskID == "" {
		taskID = c.Args().First()
	}

	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	// read from pipe if exists
	if info.Mode()&os.ModeNamedPipe != 0 {
		reader := bufio.NewReader(os.Stdin)
		var output []rune

		for {
			input, _, err := reader.ReadRune()
			if err != nil && err == io.EOF {
				break
			}
			output = append(output, input)
		}

		var args string

		for j := 0; j < len(output); j++ {
			args = args + string(output[j])
		}

		taskID = strings.TrimSpace(args)
	}

	taskDefinition := utils.GetTaskDefinitionFromTaskID(c.String("cluster"), taskID)
	logConfiguration := utils.GetLogConfigurationFromTaskDefinition(ecsSvc, &taskDefinition, c.String("container"))

	ShowLogs(&ShowLogsInput{
		Ctx:       c,
		Sess:      sess,
		LogConfig: logConfiguration,
		TaskID:    taskID,
	})
	return nil
}

func getLogStreamName(config LogConfiguration) string {
	if *config.LogStreamPrefix == "" {
		return strings.Join([]string{config.Container, config.TaskID}, "/")
	}

	return strings.Join([]string{*config.LogStreamPrefix, config.Container, config.TaskID}, "/")
}

type GetLogEventsInput struct {
	svc       *cloudwatchlogs.CloudWatchLogs
	ctx       *cli.Context
	logConfig LogConfiguration
	taskID    string
	nextToken string
}

func getLogEvents(i *GetLogEventsInput) (*cloudwatchlogs.GetLogEventsOutput, error) {
	logStreamName := getLogStreamName(i.logConfig)

	param := cloudwatchlogs.GetLogEventsInput{
		Limit:         aws.Int64(100),
		LogGroupName:  i.logConfig.LogGroupName,
		LogStreamName: aws.String(logStreamName),
		StartFromHead: aws.Bool(true),
	}

	if i.nextToken != "" {
		param.NextToken = aws.String(i.nextToken)
	}

	return i.svc.GetLogEvents(&param)
}
