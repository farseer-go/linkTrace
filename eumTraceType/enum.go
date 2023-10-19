package eumTraceType

type Enum int

const (
	WebApi        Enum = iota // Api
	MqConsumer                // MQ消费
	QueueConsumer             // 本地消费
	FSchedule                 // 调度中心
	Task                      // 本地任务
)
