package main

import "rpc"

type Args struct {
	Host                  string
	Port                  int
	URL                   string
	NumConnections        int
	ConnectionRate        int
	RequestsPerConnection int
	Duration              int
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
	date   int64 // The time the pending call was started
	args   *Args // The arguments passed to the pending call
}

type PerfData struct {
	// These fields MUST be supplied by the implementor, they do not come
	// from the parsed performance data
	BenchmarkId string
	BenchmarkDate int64
	ArgHost string
	ArgPort int
	ArgURL string
	ArgNumConnections int
	ArgConnectionRate int
	ArgRequestsPerConnection int
	ArgDuration int

	// The following fields all come from the parsed data and should not
	// need to be changed.

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
	RepliesPerSecMin, RepliesPerSecAvg, RepliesPerSecMax,
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
