package main

import "fmt"
import "os"
import "regexp"
import "strconv"

var resultPattern = `Maximum connect burst length: ([0-9]*)

Total: connections ([0-9]*) requests ([0-9]*) replies ([0-9]*) test-duration ([0-9]*\.?[0-9]*) s

Connection rate: ([0-9]*\.?[0-9]*) conn/s \(([0-9]*\.?[0-9]*) ms/conn, <=([0-9]*) concurrent connections\)
Connection time \[ms\]: min ([0-9]*\.?[0-9]*) avg ([0-9]*\.?[0-9]*) max ([0-9]*\.?[0-9]*) median ([0-9]*\.?[0-9]*) stddev ([0-9]*\.?[0-9]*)
Connection time \[ms\]: connect ([0-9]*\.?[0-9]*)
Connection length \[replies/conn\]: ([0-9]*\.?[0-9]*)

Request rate: ([0-9]*\.?[0-9]*) req/s \(([0-9]*\.?[0-9]*) ms/req\)
Request size \[B\]: ([0-9]*\.?[0-9]*)

Reply rate \[replies/s\]: min ([0-9]*\.?[0-9]*) avg ([0-9]*\.?[0-9]*) max ([0-9]*\.?[0-9]*) stddev ([0-9]*\.?[0-9]*) \(([0-9])* samples\)
Reply time \[ms\]: response ([0-9]*\.?[0-9]*) transfer ([0-9]*\.?[0-9]*)
Reply size \[B\]: header ([0-9]*\.?[0-9]*) content ([0-9]*\.?[0-9]*) footer ([0-9]*\.?[0-9]*) \(total ([0-9]*\.?[0-9]*)\)
Reply status: 1xx=([0-9]*) 2xx=([0-9]*) 3xx=([0-9]*) 4xx=([0-9]*) 5xx=([0-9]*)

CPU time \[s\]: user ([0-9]*\.?[0-9]*) system ([0-9]*\.?[0-9]*) \(user ([0-9]*\.?[0-9]*)\% system ([0-9]*\.?[0-9]*)\% total ([0-9]*\.?[0-9]*)\%\)
Net I/O: ([0-9]*\.?[0-9]*) (.*) \((.*) bps\)

Errors: total ([0-9]*) client-timo ([0-9]*) socket-timo ([0-9]*) connrefused ([0-9]*) connreset ([0-9]*)
Errors: fd-unavail ([0-9]*) addrunavail ([0-9]*) ftab-full ([0-9]*) other ([0-9]*)`

var resultRegexp = regexp.MustCompile(resultPattern)

const NUM_RESULTS = 52

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

func ParseResultsRaw(str string) []string {
	return resultRegexp.FindStringSubmatch(str)
}

func ParseResults(str string) (*PerfData, os.Error) {
	results := ParseResultsRaw(str)
	data := new(PerfData)

	var conv float64
	var err os.Error

	data.Raw = results[0]
	if conv, err = strconv.Atof64(results[1]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 1, err.String()))
	}
	data.ConnectionBurstLength = conv
	if conv, err = strconv.Atof64(results[2]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 2, err.String()))
	}
	data.TotalConnections = conv
	if conv, err = strconv.Atof64(results[3]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 3, err.String()))
	}
	data.TotalRequests = conv
	if conv, err = strconv.Atof64(results[4]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 4, err.String()))
	}
	data.TotalReplies = conv
	if conv, err = strconv.Atof64(results[5]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 5, err.String()))
	}
	data.TestDuration = conv
	if conv, err = strconv.Atof64(results[6]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 6, err.String()))
	}
	data.ConnectionsPerSecond = conv
	if conv, err = strconv.Atof64(results[7]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 7, err.String()))
	}
	data.MsPerConnection = conv
	if conv, err = strconv.Atof64(results[8]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 8, err.String()))
	}
	data.ConcurrentConnections = conv
	if conv, err = strconv.Atof64(results[9]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 9, err.String()))
	}
	data.ConnectionTimeMin = conv
	if conv, err = strconv.Atof64(results[10]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 10, err.String()))
	}
	data.ConnectionTimeAvg = conv
	if conv, err = strconv.Atof64(results[11]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 11, err.String()))
	}
	data.ConnectionTimeMax = conv
	if conv, err = strconv.Atof64(results[12]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 12, err.String()))
	}
	data.ConnectionTimeMedian = conv
	if conv, err = strconv.Atof64(results[13]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 13, err.String()))
	}
	data.ConnectionTimeStddev = conv
	if conv, err = strconv.Atof64(results[14]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 14, err.String()))
	}
	data.ConnectionTimeConnect = conv
	if conv, err = strconv.Atof64(results[15]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 15, err.String()))
	}
	data.RepliesPerConnection = conv
	if conv, err = strconv.Atof64(results[16]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 16, err.String()))
	}
	data.RequestsPerSecond = conv
	if conv, err = strconv.Atof64(results[17]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 17, err.String()))
	}
	data.MsPerRequest = conv
	if conv, err = strconv.Atof64(results[18]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 18, err.String()))
	}
	data.RequestSize = conv
	if conv, err = strconv.Atof64(results[19]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 19, err.String()))
	}
	data.RepliesPerSecMin = conv
	if conv, err = strconv.Atof64(results[20]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 20, err.String()))
	}
	data.RepliesPerSecAvgm = conv
	if conv, err = strconv.Atof64(results[21]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 21, err.String()))
	}
	data.RepliesPerSecMax = conv
	if conv, err = strconv.Atof64(results[22]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 22, err.String()))
	}
	data.RepliesPerSecStddev = conv
	if conv, err = strconv.Atof64(results[23]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 23, err.String()))
	}
	data.RepliesPerSecNumSamples = conv
	if conv, err = strconv.Atof64(results[24]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 24, err.String()))
	}
	data.ReplyTimeResponse = conv
	if conv, err = strconv.Atof64(results[25]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 25, err.String()))
	}
	data.ReplyTimeTransfer = conv
	if conv, err = strconv.Atof64(results[26]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 26, err.String()))
	}
	data.ReplySizeHeader = conv
	if conv, err = strconv.Atof64(results[27]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 27, err.String()))
	}
	data.ReplySizeContent = conv
	if conv, err = strconv.Atof64(results[28]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 28, err.String()))
	}
	data.ReplySizeFooter = conv
	if conv, err = strconv.Atof64(results[29]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 29, err.String()))
	}
	data.ReplySizeTotal = conv
	if conv, err = strconv.Atof64(results[30]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 30, err.String()))
	}
	data.ReplyStatus_1xx = conv
	if conv, err = strconv.Atof64(results[31]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 31, err.String()))
	}
	data.ReplyStatus_2xx = conv
	if conv, err = strconv.Atof64(results[32]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 32, err.String()))
	}
	data.ReplyStatus_3xx = conv
	if conv, err = strconv.Atof64(results[33]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 33, err.String()))
	}
	data.ReplyStatus_4xx = conv
	if conv, err = strconv.Atof64(results[34]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 34, err.String()))
	}
	data.ReplyStatus_5xx = conv
	if conv, err = strconv.Atof64(results[35]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 35, err.String()))
	}
	data.CpuTimeUser = conv
	if conv, err = strconv.Atof64(results[36]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 36, err.String()))
	}
	data.CpuTimeSystem = conv
	if conv, err = strconv.Atof64(results[37]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 37, err.String()))
	}
	data.CpuPercUser = conv
	if conv, err = strconv.Atof64(results[38]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 38, err.String()))
	}
	data.CpuPercSystem = conv
	if conv, err = strconv.Atof64(results[39]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 39, err.String()))
	}
	data.CpuPercTotal = conv
	if conv, err = strconv.Atof64(results[40]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 40, err.String()))
	}
	data.NetIOValue = conv
	data.NetIOUnit = results[41]
	data.NetIOBytesPerSecond = results[42]
	if conv, err = strconv.Atof64(results[43]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 43, err.String()))
	}
	data.ErrTotal = conv
	if conv, err = strconv.Atof64(results[44]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 44, err.String()))
	}
	data.ErrClientTimeout = conv
	if conv, err = strconv.Atof64(results[45]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 45, err.String()))
	}
	data.ErrSocketTimeout = conv
	if conv, err = strconv.Atof64(results[46]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 46, err.String()))
	}
	data.ErrConnectionRefused = conv
	if conv, err = strconv.Atof64(results[47]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 47, err.String()))
	}
	data.ErrConnectionReset = conv
	if conv, err = strconv.Atof64(results[48]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 48, err.String()))
	}
	data.ErrFdUnavail = conv
	if conv, err = strconv.Atof64(results[49]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 49, err.String()))
	}
	data.ErrAddRunAvail = conv
	if conv, err = strconv.Atof64(results[50]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 50, err.String()))
	}
	data.ErrFtabFull = conv
	if conv, err = strconv.Atof64(results[51]); err != nil {
		return nil, os.NewError(fmt.Sprintf("Error parsing field %d:%s", 51, err.String()))
	}
	data.ErrOther = conv

	return data, nil
}

/* The following Lua script was used to generate the above code:

local fields = {"Raw", "ConnectionBurstLength", "TotalConnections", "TotalRequests", "TotalReplies", "TestDuration", "ConnectionsPerSecond", "MsPerConnection", "ConcurrentConnections", "ConnectionTimeMin", "ConnectionTimeAvg", "ConnectionTimeMax", "ConnectionTimeMedian", "ConnectionTimeStddev", "ConnectionTimeConnect", "RepliesPerConnection", "RequestsPerSecond", "MsPerRequest", "RequestSize", "RepliesPerSecMin", "RepliesPerSecAvgm", "RepliesPerSecMax", "RepliesPerSecStddev", "RepliesPerSecNumSamples", "ReplyTimeResponse", "ReplyTimeTransfer", "ReplySizeHeader", "ReplySizeContent", "ReplySizeFooter", "ReplySizeTotal", "ReplyStatus_1xx", "ReplyStatus_2xx", "ReplyStatus_3xx", "ReplyStatus_4xx", "ReplyStatus_5xx",
"CpuTimeUser", "CpuTimeSystem", "CpuPercUser", "CpuPercSystem", "CpuPercTotal", "NetIOValue", "NetIOUnit", "NetIOBytesPerSecond", "ErrTotal", "ErrClientTimeout", "ErrSocketTimeout", "ErrConnectionRefused", "ErrConnectionReset", "ErrFdUnavail", "ErrAddRunAvail", "ErrFtabFull", "ErrOther"}

local strings = {Raw = true, NetIOUnit = true, NetIOBytesPerSecond = true}

for idx, field in ipairs(fields) do
	if strings[field] then
		print(string.format("data.%s = results[%d]", field, idx - 1))
	else
		print(string.format("if conv, err = strconv.Atof64(results[%d]); err != nil {", idx - 1))
		print(string.format([[	return nil, os.NewError(fmt.Sprintf("Error parsing field %%d:%%s", %d, err.String()))]], idx - 1))
		print(string.format("}"))
		print(string.format("data.%s = conv", field))
	end
end

*/
