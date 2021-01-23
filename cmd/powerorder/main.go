package main

import (
	"log"
	"powerorder/app/router"
	"runtime/debug"
	"time"

	"powerorder/version"

	bs "github.com/tal-tech/hera/bootstrap"
	"github.com/tal-tech/hera/ginhttp"
	"github.com/tal-tech/xtools/confutil"
	"github.com/tal-tech/xtools/flagutil"

	"github.com/spf13/cast"
)

func main() {
	//show version
	ver := flagutil.GetVersion()
	if *ver {
		version.Version()
		return
	}

	//init conf
	confutil.InitConfig()

	defer recovery()
	s := ginhttp.NewServer()
	engine := s.GetGinEngine()

	//Add middleware
	//You can customize the middleware according to your actual needs
	engine.Use()

	router.RegisterRouter(engine)

	//Front hook for service startup
	s.AddBeforeServerStartFunc(
		bs.InitLoggerWithConf(),
		bs.InitTraceLogger("Powerorder", "1.0"),
		s.InitConfig(),
	)

	//Exec hook Funcs before the service to closing
	s.AddAfterServerStopFunc(bs.CloseLogger())

	er := s.Serve()
	if er != nil {
		log.Printf("Server stop err:%v", er)
	} else {
		log.Printf("Server exit")
	}
}

func recovery() {
	//panic cause program exit quickly，Some logs may not have time to be written to disk
	time.Sleep(time.Second)

	if rec := recover(); rec != nil {
		log.Printf("Panic Panic occur")
		if err, ok := rec.(error); ok {
			log.Printf("PanicRecover Unhandled error: %v\n stack:%v", err.Error(), cast.ToString(debug.Stack()))
		} else {
			log.Printf("PanicRecover Panic: %v\n stack:%v", rec, cast.ToString(debug.Stack()))
		}
	}
}
