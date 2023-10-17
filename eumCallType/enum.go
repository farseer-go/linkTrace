package eumCallType

type Enum int

const (
	HttpClient    Enum = iota // HttpClient
	GrpcClient                // GrpcClient
	Database                  // Database
	Redis                     // Redis
	Mq                        // Mq
	Elasticsearch             // Elasticsearch
	Hand                      // Hand
	KeyLocation               // KeyLocation
)

func (receiver Enum) ToString() string {
	switch receiver {
	case HttpClient:
		return "HttpClient"
	case GrpcClient:
		return "GrpcClient"
	case Database:
		return "Database"
	case Redis:
		return "Redis"
	case Mq:
		return "Mq"
	case Elasticsearch:
		return "Elasticsearch"
		return "Mq"
	case Hand:
		return "Hand"
		return "Mq"
	case KeyLocation:
		return "KeyLocation"
	}
	return ""
}
