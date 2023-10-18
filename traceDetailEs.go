package linkTrace

import (
	"fmt"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/linkTrace/eumCallType"
)

type TraceDetailEs struct {
	TraceDetail
	IndexName   string // 索引名称
	AliasesName string // 别名
}

func (receiver *TraceDetailEs) GetTraceDetail() *TraceDetail {
	return &receiver.TraceDetail
}

func (receiver *TraceDetailEs) ToString() string {
	return fmt.Sprintf("[%s]耗时：%s，%s IndexName=%s，AliasesName=%s", flog.Yellow(receiver.CallType.ToString()), flog.Red(receiver.UseTs.String()), receiver.CallMethod, receiver.IndexName, receiver.AliasesName)
}

// TraceElasticsearch Elasticsearch埋点
func TraceElasticsearch(method string, IndexName string, AliasesName string) *TraceDetailEs {
	detail := &TraceDetailEs{
		TraceDetail: newTraceDetail(eumCallType.Elasticsearch),
		IndexName:   IndexName,
		AliasesName: AliasesName,
	}
	add(detail)
	return detail
}
