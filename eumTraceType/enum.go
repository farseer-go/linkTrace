package eumTraceType

type Enum int

const (
	WebApi    Enum = iota // Api
	FSchedule             // 调度中心
	Consumer              // MQ消费
	Job                   // 本地任务
)
