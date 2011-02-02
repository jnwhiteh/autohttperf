package main

import "rpc"

type Args struct {
	Host                  string
	Port                  int
	URL                   string
	NumConnections        int
	ConnectionRate        int
	RequestsPerConnection int
	Hog                   bool
}

type Result struct {
	Stdout     string
	Stderr     string
	ExitStatus int
}

type Worker struct {
	addr   string // The address of the RPC worker client
	id     string // A string UID for this worker
	client *rpc.Client
	result *Result   // The pending RPC result
	call   *rpc.Call // The pending RPC call result
}

type PerfData struct {
	Raw string
	ConnectionBurstLength,
	TotalConnections, TotalRequests, TotalReplies, TestDuration,
	ConnectionsPerSecond, MsPerConnection, ConcurrentConnections,
	ConnectionTimeMin, ConnectionTimeAvg, ConnectionTimeMax,
	ConnectionTimeMedian, ConnectionTimeStddev,
	ConnectionTimeConnect,
	RepliesPerConnection,
	RequestsPerSecond, MsPerRequest,
	RequestSize,
	RepliesPerSecMin, RepliesPerSecAvgm, RepliesPerSecMax,
	RepliesPerSecStddev, RepliesPerSecNumSamples,
	ReplyTimeResponse, ReplyTimeTransfer,
	ReplySizeHeader, ReplySizeContent, ReplySizeFooter, ReplySizeTotal,
	ReplyStatus_1xx, ReplyStatus_2xx, ReplyStatus_3xx, ReplyStatus_4xx, ReplyStatus_5xx,
	CpuTimeUser, CpuTimeSystem, CpuPercUser, CpuPercSystem, CpuPercTotal,
	NetIOValue float64
	NetIOUnit, NetIOBytesPerSecond string
	ErrTotal, ErrClientTimeout, ErrSocketTimeout, ErrConnectionRefused,
	ErrConnectionReset, ErrFdUnavail, ErrAddRunAvail, ErrFtabFull, ErrOther float64
}
