package linkTrace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/fs/trace"
	"net/http"
)

// FopsServer fops地址
var FopsServer string

// SaveTraceContextConsumer 上传到FOPS中心
func SaveTraceContextConsumer(subscribeName string, lstMessage collections.ListAny, remainingCount int) {
	trace.CurTraceContext.Get().Ignore()
	lstTraceContext := collections.NewList[TraceContext]()
	lstMessage.Foreach(func(item *any) {
		// 上下文
		dto := (*item).(TraceContext)
		lstTraceContext.Add(dto)
	})
	if err := upload(lstTraceContext); err != nil {
		exception.ThrowRefuseException(err.Error())
	}
	return
}

type UploadRequest struct {
	List collections.List[TraceContext]
}

func upload(lstTraceContext collections.List[TraceContext]) error {
	bodyByte, _ := json.Marshal(UploadRequest{List: lstTraceContext})
	url := fmt.Sprintf("%slinkTrace/upload", FopsServer)
	newRequest, _ := http.NewRequest("POST", url, bytes.NewReader(bodyByte))
	newRequest.Header.Set("Content-Type", "application/json")
	// 链路追踪
	if traceContext := container.Resolve[trace.IManager]().GetCurTrace(); traceContext != nil {
		newRequest.Header.Set("Trace-Id", parse.ToString(traceContext.GetTraceId()))
		newRequest.Header.Set("Trace-App-Name", fs.AppName)
	}
	client := &http.Client{}
	rsp, err := client.Do(newRequest)
	if err != nil {
		return err
	}

	apiRsp := core.NewApiResponseByReader[any](rsp.Body)
	if apiRsp.StatusCode != 200 {
		return fmt.Errorf("上传链路记录到%s失败（%v）：%s", url, rsp.StatusCode, apiRsp.StatusMessage)
	}

	return err
}
