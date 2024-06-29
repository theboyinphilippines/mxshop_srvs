package main

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

//func main() {
//	cfg := jaegercfg.Configuration{
//		//采样器设置
//		Sampler: &jaegercfg.SamplerConfig{
//			Type:  jaeger.SamplerTypeConst,
//			Param: 1,
//		},
//		//jaeger agent设置
//		Reporter: &jaegercfg.ReporterConfig{
//			LogSpans:           true,
//			LocalAgentHostPort: "192.168.0.101:6831",
//		},
//		ServiceName: "mxshop",
//	}
//	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
//	if err != nil {
//		panic(any(err))
//	}
//	defer closer.Close()
//
//	//1发送一个简单的span
//	//span := tracer.StartSpan("go-grpc-web")
//	//time.Sleep(time.Second)
//	//defer span.Finish()
//
//
//	//2发送一个多级嵌套的span
//	parentSpan := tracer.StartSpan("main")
//	span := tracer.StartSpan("funcA",opentracing.ChildOf(parentSpan.Context()))
//	time.Sleep(time.Second)
//	span.Finish()
//
//	time.Sleep(500* time.Millisecond)
//
//	span2 := tracer.StartSpan("funcB",opentracing.ChildOf(parentSpan.Context()))
//	time.Sleep(2*time.Second)
//	span2.Finish()
//
//	parentSpan.Finish()
//}

func main() {
	if len(os.Args) != 2 {
		panic("ERROR: Expecting one argument")
	}

	tracer, cloer := Init("hello-world")
	defer cloer.Close()
	opentracing.SetGlobalTracer(tracer)

	helloTo := os.Args[1]

	//创建parent span
	parentSpan := tracer.StartSpan("say-hello")
	parentSpan.SetTag("hello-to", helloTo)
	defer parentSpan.Finish()

	//创建一个新的ctx，将parent span的信息与context关联
	ctx := opentracing.ContextWithSpan(context.Background(), parentSpan)

	//将新的新的ctx传入，创建一个子span，父span是ctx中的span
	helloStr := formatString(ctx, helloTo)

	printHello(ctx, helloStr)

	//helloStr := fmt.Sprintf("Hello, %s!", helloTo)

}

func formatString(ctx context.Context, helloTo string) string {
	//创建子span，关联到parent span，从parent span的相关联的ctx中提取到子span
	span, _ := opentracing.StartSpanFromContext(ctx, "formatString")
	defer span.Finish()

	v := url.Values{}
	v.Set("helloTo", helloTo)
	url := "http://localhost:8081/format?" + v.Encode()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err.Error())
	}
	xhttp := &http.Client{}
	resp, err := xhttp.Do(req)
	if err != nil {
		panic(err.Error())
	}

	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)

	helloStr := string(buf)

	span.LogFields(
		log.String("event", "string-format"),
		log.String("value", helloStr),
	)

	return helloStr
}

func printHello(ctx context.Context, helloStr string) {
	span, _ := opentracing.StartSpanFromContext(ctx, "printHello")
	defer span.Finish()

	v := url.Values{}
	v.Set("helloStr", helloStr)
	url := "http://localhost:8082/publish?" + v.Encode()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err.Error())
	}

	xhttp := &http.Client{}
	_, err = xhttp.Do(req)
	if err != nil {
		panic(err.Error())
	}

}

// Init returns an instance of Jaeger Tracer that samples 100% of traces and logs all spans to stdout.
func Init(service string) (opentracing.Tracer, io.Closer) {
	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}
	tracer, closer, err := cfg.New(service, config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	return tracer, closer
}
