package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type RpcServer struct {
}

func (t *RpcServer) Stop(_ string, _ *struct{}) error {
	StopCurrentWorkflow()
	return nil
}

func StartRpcServer() {
	rpc.Register(new(RpcServer))
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}
