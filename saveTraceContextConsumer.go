package linkTrace

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"strconv"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/snc"
	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/linkTrace/eumTraceType"
)

// FopsServer fops地址
var FopsServer string

// SaveTraceContextConsumer 上传链路记录到FOPS中心
func SaveTraceContextConsumer(subscribeName string, lstMessage collections.ListAny, remainingCount int) {
	if traceContext := trace.CurTraceContext.Get(); traceContext != nil {
		traceContext.Ignore()
	}
	lstTraceContext := collections.NewList[TraceContext]()
	lstMessage.Foreach(func(item *any) {
		// 上下文
		dto := (*item).(TraceContext)
		if len(dto.List) == 0 && dto.TraceType != eumTraceType.WebApi {
			return
		}
		dto.printLog()

		// 链路超过200条，则丢弃
		if len(dto.List) > 200 {
			dto.List = dto.List[0:200]
			switch dto.TraceType {
			case eumTraceType.WebApi:
				flog.Warningf("【%s链路追踪】链路明细超过200条，TraceId:%s，耗时：%s，%s", dto.TraceType.ToString(), flog.Green(dto.TraceId), flog.Red(dto.UseTs.String()), dto.WebContext.WebPath)
			case eumTraceType.MqConsumer, eumTraceType.QueueConsumer, eumTraceType.EventConsumer:
				flog.Warningf("【%s链路追踪】链路明细超过200条，TraceId:%s，耗时：%s，%s", dto.TraceType.ToString(), flog.Green(dto.TraceId), flog.Red(dto.UseTs.String()), dto.ConsumerContext.ConsumerQueueName)
			case eumTraceType.Task, eumTraceType.FSchedule:
				flog.Warningf("【%s链路追踪】链路明细超过200条，TraceId:%s，耗时：%s，%s %s", dto.TraceType.ToString(), flog.Green(dto.TraceId), flog.Red(dto.UseTs.String()), dto.TaskContext.TaskName, dto.TaskContext.TaskGroupName)
			default:
				flog.Warningf("【%s链路追踪】链路明细超过200条，TraceId:%s，耗时：%s", dto.TraceType.ToString(), flog.Green(dto.TraceId), flog.Red(dto.UseTs.String()))
			}
		}
		lstTraceContext.Add(dto)
	})
	if err := uploadTrace(lstTraceContext); err != nil {
		exception.ThrowRefuseException(err.Error())
	}
}

type UploadTraceRequest struct {
	List any
}

// UploadTrace 上传链路记录
func uploadTrace(lstTraceContext any) error {
	bodyByte, _ := snc.Marshal(UploadTraceRequest{List: lstTraceContext})
	url := configure.GetFopsServer() + "linkTrace/upload"
	newRequest, _ := http.NewRequest("POST", url, bytes.NewReader(bodyByte))
	newRequest.Header.Set("Content-Type", "application/json")
	// 链路追踪
	if traceContext := container.Resolve[trace.IManager]().GetCurTrace(); traceContext != nil {
		newRequest.Header.Set("Trace-Id", traceContext.GetTraceId())
		newRequest.Header.Set("Trace-Level", strconv.Itoa(traceContext.GetTraceLevel()))
		newRequest.Header.Set("Trace-App-Name", core.AppName)
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // 不验证 HTTPS 证书
			},
		},
	}
	rsp, err := client.Do(newRequest)
	if err != nil {
		return fmt.Errorf("上传链路记录到FOPS失败：%s", err.Error())
	}

	apiRsp := core.NewApiResponseByReader[any](rsp.Body)
	if apiRsp.StatusCode != 200 {
		return fmt.Errorf("上传链路记录到FOPS失败（%v）：%s", apiRsp.StatusCode, apiRsp.StatusMessage)
	}
	return err
}
