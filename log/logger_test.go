package log_test

import (
	"context"
	bplusc "gitlab.intelligentb.com/devops/bplus/context"
	log "gitlab.intelligentb.com/devops/bplus/log"
	"testing"
)

/*
func TestError(t *testing.T) {
	ctx := context.TODO()
	ctx = bplusc.Add(ctx,bplusc.TraceID,"TRACEID123")
	out := capturer.CaptureOutput(func() {
		log.Error(ctx, "message")
	})
	fmt.Fprintf(os.Stderr,"out is %s\n",out)
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(out),&m); if err != nil {
		t.Errorf("Cannot unmarshal the json from log. Err=%s\n",err.Error())
		t.Fail()
		return
	}
	assert.Equal(t,m["TraceID"],"TRACEID123")
}
*/
func TestErrorWithFields(t *testing.T) {
	ctx := context.TODO()
	ctx = bplusc.Add(ctx,bplusc.TraceID,"TRACEID123")
	log.ErrorWithFields(ctx,map[string]string{
		"foo":"bar",
		}, "message",
	)
}

func TestWarn(t *testing.T) {
	ctx := context.TODO()
	ctx = bplusc.Add(ctx,bplusc.TraceID,"TRACEID123")
	log.Warn(ctx,"message")
}

func TestWarnWithFields(t *testing.T) {
	ctx := context.TODO()
	ctx = bplusc.Add(ctx,bplusc.TraceID,"TRACEID123")
	log.WarnWithFields(ctx,map[string]string{
		"foo":"bar",
	}, "message",
	)
}

func TestInfo(t *testing.T) {
	ctx := context.TODO()
	ctx = bplusc.Add(ctx,bplusc.TraceID,"TRACEID123")
	log.Info(ctx,"message")
}

func TestInfoWithFields(t *testing.T) {
	ctx := context.TODO()
	ctx = bplusc.Add(ctx,bplusc.TraceID,"TRACEID123")
	log.InfoWithFields(ctx,map[string]string{
		"foo":"bar",
	}, "message",
	)
}

