package main

import "rpc"

type Args struct {
	Host string
	Port int
	URL string
	NumConnections int
	ConnectionRate int
	RequestsPerConnection int
	Hog bool
}

type Result struct {
	Stdout string
	Stderr string
	ExitStatus int
}

type Worker struct {
	addr string				// The address of the RPC worker client
	id string				// A string UID for this worker
	client *rpc.Client
	result *Result			// The pending RPC result
	call *rpc.Call			// The pending RPC call result
}
