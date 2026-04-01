package worker

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	db "pxsemic.com/simplebank/db/sqlc"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	Start() error
	ProcessorTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	store  db.Store
	server *asynq.Server
}

func NewRedisTaskProcessor(db db.Store, redisOpt asynq.RedisClientOpt) TaskProcessor {

	logger := NewLogger()
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Queues: map[string]int{
				QueueCritical: 5,
				QueueDefault:  15,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Error().Err(err).Str("type", task.Type()).
					Bytes("payload", task.Payload()).Msg("process task failed")
			}),
			Logger: logger,
		},
	)
	redis.SetLogger(logger) // 替换go-redis包 logger为我们自定义的logger

	return &RedisTaskProcessor{
		store:  db,
		server: server,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessorTaskSendVerifyEmail)
	return processor.server.Start(mux)
}
