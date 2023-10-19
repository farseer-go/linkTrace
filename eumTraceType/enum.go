package eumTraceType

type Enum int

const (
	WebApi        Enum = iota // Api
	FSchedule                 // 调度中心
	MqConsumer                // MQ消费
	QueueConsumer             // 本地消费
	Job                       // 本地任务
)
