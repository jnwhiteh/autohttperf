package main

import "reflect"
import "strconv"
import "testing"

var testData = `Maximum connect burst length: 1

Total: connections 10000 requests 10000 replies 10000 test-duration 6.964 s

Connection rate: 1435.9 conn/s (0.7 ms/conn, <=1 concurrent connections)
Connection time [ms]: min 0.2 avg 0.7 max 27.4 median 0.5 stddev 0.7
Connection time [ms]: connect 0.1
Connection length [replies/conn]: 1.000

Request rate: 1435.9 req/s (0.7 ms/req)
Request size [B]: 72.0

Reply rate [replies/s]: min 1444.8 avg 1444.8 max 1444.8 stddev 0.0 (1 samples)
Reply time [ms]: response 0.5 transfer 0.1
Reply size [B]: header 170.0 content 4109.0 footer 2.0 (total 4281.0)
Reply status: 1xx=0 2xx=10000 3xx=0 4xx=0 5xx=0

CPU time [s]: user 1.28 system 5.22 (user 18.4% system 75.0% total 93.5%)
Net I/O: 6101.1 KB/s (50.0*10^6 bps)

Errors: total 0 client-timo 0 socket-timo 0 connrefused 0 connreset 0
Errors: fd-unavail 0 addrunavail 0 ftab-full 0 other 0`

var expectedNums = map[int]float64{
	// Index 0 will contain the entire match
	1: 1,
	2: 10000, 3: 10000, 4: 10000, 5: 6.964,

	6: 1435.9, 7: 0.7, 8: 1,
	9: 0.2, 10: 0.7, 11: 27.4, 12: 0.5, 13: 0.7,
	14: 0.1,
	15: 1.0,

	16: 1435.9, 17: 0.7,
	18: 72.0,

	19: 1444.8, 20: 1444.8, 21: 1444.8, 22: 0.0, 23: 1,
	24: 0.5, 25: 0.1,
	26: 170, 27: 4109, 28: 2.0, 29: 4281,
	30: 0, 31: 10000, 32: 0, 33: 0, 34: 0,

	35: 1.28, 36: 5.22, 37: 18.4, 38: 75.0, 39: 93.5,
	40: 6101.1,

	43: 0, 44: 0, 45: 0, 46: 0, 47: 0,
	48: 0, 49: 0, 50: 0, 51: 0,
}

var expectedStrings = map[int]string{
	0:  testData,
	41: "KB/s",
	42: "50.0*10^6",
}

func TestParseRaw(t *testing.T) {
	results := ParseResultsRaw(testData)
	if len(results) != 52 {
		t.Errorf("Got %d results, was expecting 51", len(results))
	}

	for idx, value := range expectedNums {
		conv, err := strconv.Atof64(results[idx])
		if err != nil {
			t.Errorf("Unable to convert capture %d to a float: %q", idx, results[idx])
		} else {
			if conv != value {
				t.Errorf("Expected %f for result %d, got %f", value, idx, conv)
			}
		}
	}

	for idx, value := range expectedStrings {
		if value != results[idx] {
			t.Errorf("Expected %s for result %d, got %s", value, idx, results[idx])
		}
	}
}

var fields = []string{"Raw", "ConnectionBurstLength", "TotalConnections", "TotalRequests", "TotalReplies", "TestDuration", "ConnectionsPerSecond", "MsPerConnection", "ConcurrentConnections", "ConnectionTimeMin", "ConnectionTimeAvg", "ConnectionTimeMax", "ConnectionTimeMedian", "ConnectionTimeStddev", "ConnectionTimeConnect", "RepliesPerConnection", "RequestsPerSecond", "MsPerRequest", "RequestSize", "RepliesPerSecMin", "RepliesPerSecAvgm", "RepliesPerSecMax", "RepliesPerSecStddev", "RepliesPerSecNumSamples", "ReplyTimeResponse", "ReplyTimeTransfer", "ReplySizeHeader", "ReplySizeContent", "ReplySizeFooter", "ReplySizeTotal", "ReplyStatus_1xx", "ReplyStatus_2xx", "ReplyStatus_3xx", "ReplyStatus_4xx", "ReplyStatus_5xx", "CpuTimeUser", "CpuTimeSystem", "CpuPercUser", "CpuPercSystem", "CpuPercTotal", "NetIOValue", "NetIOUnit", "NetIOBytesPerSecond", "ErrTotal", "ErrClientTimeout", "ErrSocketTimeout", "ErrConnectionRefused", "ErrConnectionReset", "ErrFdUnavail", "ErrAddRunAvail", "ErrFtabFull", "ErrOther"}

func TestParse(t *testing.T) {
	results, err := ParseResults(testData)
	if err != nil {
		t.Errorf("Failed to parse: %s", err.String())
		return
	}

	// Turn the struct into a Type so we can use reflection
	ptr, ok := reflect.NewValue(results).(*reflect.PtrValue)
	if !ok {
		t.Errorf("Could not convert results into a pointer value")
		return
	}

	val, ok := ptr.Elem().(*reflect.StructValue)
	if !ok {
		t.Errorf("Failed when reflecting on struct")
		return
	}

	for idx, expected := range expectedNums {
		// We need to grab the field by name, using the reflect package
		result := val.FieldByName(fields[idx])
		if result == nil {
			t.Errorf("Failed when reflecting field %s", fields[idx])
		}

		fvalue := result.(*reflect.FloatValue)
		fresult := fvalue.Get()

		if expected != fresult {
			t.Errorf("Expected %f for result %s, got %f", expected, fields[idx], fresult)
		}
	}

	for idx, expected := range expectedStrings {
		result := val.FieldByName(fields[idx])
		if result == nil {
			t.Errorf("Failed when reflecting field %s", fields[idx])
		}

		fvalue := result.(*reflect.StringValue)
		sresult := fvalue.Get()

		if expected != sresult {
			t.Errorf("Expected %s for result %s, got %s", expected, fields[idx], sresult)
		}
	}
}
