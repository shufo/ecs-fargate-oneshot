# ecs-fargate-oneshot

Executes oneshot task on ECS (fargate)

## Installation

Using go get

```bash
$ go get -u github.com/shufo/ecs-fargate-oneshot
```

or download manually

TBD

### resource permissions

You must have these IAM permissions to execute tasks

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "ExecuteOneShotTask",
      "Effect": "Allow",
      "Action": [
        "ecs:DescribeServices",
        "ecs:DescribeTasks",
        "ecs:DescribeTaskDefinition",
        "ecs:RunTask",
        "ecs:ListTaskDefinitions",
        "logs:GetLogEvents"
      ],
      "Resource": "*"
    }
  ]
}
```

## Usage

### tasks

Executes task (with showing logs)

```bash
$ ecs-fargate-oneshot -v \
    --cluster cluster-name \
    --service service-name \
    run --task-definition app \
    --container app \
    --show-cloudwatch-logs \
    echo "foo bar"
```

NOTE: If you would like to show logs, you must define log configuration for cloudwatch logs on ecs task definition.

```json
...
"logConfiguration": {
    "logDriver": "awslogs",
    "options": {
        "awslogs-group" : "/fargate/service/app",
        "awslogs-region": "us-east-1",
        "awslogs-stream-prefix": "app"
    }
}
```

or if you would not require logs

```bash
$ ecs-fargate-oneshot \
    --cluster cluster-name \
    --service service-name \
    run --task-definition app \
    --container app \
    echo "foo bar"
```

### logs

You can showing logs after task execution

```bash
# run task without logs
$ ecs-fargate-oneshot \
    --cluster cluster-name \
    --service service-name \
    run --task-definition app \
    --container app \
    echo "foo bar"
# output =>
287771bd-92f4-407c-870f-7a480b94cbc7 # this is taskId

# show up logs by passing task id from execution output
$ ecs-fargate-oneshot \
    --cluster cluster-name \
    --service service-name \
    logs --container app \
    287771bd-92f4-407c-870f-7a480b94cbc7
# output =>
Event messages for stream app/app/287771bd-92f4-407c-870f-7a480b94cbc7 in log group LOG-GROUP-NAME: /fargat
e/service/service-name
   foo bar

# or you can pass the task id from stdin
$ ecs-fargate-oneshot --cluster cluster-name --service service-name run --task-definition app --container app echo "foo bar" | ecs-fargate-oneshot --cluster cluster-name --service service-name logs --container app
```

## Options

### Global

|             option |       description | default | required |
| -----------------: | ----------------: | ------: | -------- |
|  `--cluster`, `-c` |      Cluster name |    `""` | yes      |
|  `--service`, `-s` |      Service name |    `""` | yes      |
|  `--verbose`, `-v` | Show verbose logs | `false` | no       |
| `--progress`, `-p` |     Show progress | `false` | no       |
|    `--help,`, `-h` |         Show help | `false` | no       |

### run sub command

|                          option |                  description | default | required |
| ------------------------------: | ---------------------------: | ------: | -------- |
|       `--task-definition`, `-t` |         Task definition name |    `""` | yes      |
|             `--container`, `-n` |               Container name |    `""` | yes      |
| `--show-cloudwatch-logs,`, `-l` | Show logs on cloudwatch logs | `false` | no       |

### logs sub command

|              option |    description | default | required |
| ------------------: | -------------: | ------: | -------- |
| `--container`, `-n` | Container name |    `""` | yes      |
|   `--task-id`, `-t` |  A ECS task id |    `""` | no       |

## Contributing

1.  Fork it
2.  Create your feature branch (`git checkout -b my-new-feature`)
3.  Commit your changes (`git commit -am 'Add some feature'`)
4.  Push to the branch (`git push origin my-new-feature`)
5.  Create new Pull Request

## LICENSE

MIT

## Troubleshoot

- If you encounter the error about the aws API execution

check if the aws environment variable is properly set

```
$ env | grep aws -i

export AWS_DEFAULT_REGION=
export AWS_ACCESS_KEY_ID=
export AWS_SECRET_ACCESS_KEY=
```

## TODO

- [ ] Add option to execute tasks without waiting task status