package main

import "fmt"
import "io"
import "log"
import "strings"
import "reflect"

// The 'Raw' field is omitted here, since all of the data is already included
var fieldNames = []string{"BenchmarkId", "BenchmarkDate", "ArgHost", "ArgPort", "ArgURL", "ArgNumConnections", "ArgConnectionRate", "ArgRequestsPerConnection", "ConnectionBurstLength", "TotalConnections", "TotalRequests", "TotalReplies", "TestDuration", "ConnectionsPerSecond", "MsPerConnection", "ConcurrentConnections", "ConnectionTimeMin", "ConnectionTimeAvg", "ConnectionTimeMax", "ConnectionTimeMedian", "ConnectionTimeStddev", "ConnectionTimeConnect", "RepliesPerConnection", "RequestsPerSecond", "MsPerRequest", "RequestSize", "RepliesPerSecMin", "RepliesPerSecAvgm", "RepliesPerSecMax", "RepliesPerSecStddev", "RepliesPerSecNumSamples", "ReplyTimeResponse", "ReplyTimeTransfer", "ReplySizeHeader", "ReplySizeContent", "ReplySizeFooter", "ReplySizeTotal", "ReplyStatus_1xx", "ReplyStatus_2xx", "ReplyStatus_3xx", "ReplyStatus_4xx", "ReplyStatus_5xx", "CpuTimeUser", "CpuTimeSystem", "CpuPercUser", "CpuPercSystem", "CpuPercTotal", "NetIOValue", "NetIOUnit", "NetIOBytesPerSecond", "ErrTotal", "ErrClientTimeout", "ErrSocketTimeout", "ErrConnectionRefused", "ErrConnectionReset", "ErrFdUnavail", "ErrAddRunAvail", "ErrFtabFull", "ErrOther"}

// Write a CSV header to the given writer including each of the field names
// above, and an optional list of additional column names specified. In the
// resulting file, the optional columns are listed first.
func WriteTSVHeader(w io.Writer) {
	numColumns := len(fieldNames)
	columns := make([]string, 0, numColumns)

	columns = append(columns, fieldNames...)

	io.WriteString(w, strings.Join(columns, ","))
	io.WriteString(w, "\n")
}

func WriteTSVParseDataSet(w io.Writer, data []*PerfData) {
	for _, result := range data {
		WriteTSVParseData(w, result)
	}
}

func WriteTSVParseData(w io.Writer, data *PerfData) {
	numColumns := len(fieldNames)
	columns := make([]string, 0, numColumns)

	// Turn the struct into a Type so we can use reflection
	ptr, ok := reflect.NewValue(data).(*reflect.PtrValue)
	if !ok {
		log.Exitf("Could not convert results into a pointer value")
		return
	}

	val, ok := ptr.Elem().(*reflect.StructValue)
	if !ok {
		log.Exitf("Failed when reflecting on struct")
		return
	}

	// Move through every field, fetching the value by name and adding
	// it to the columns slice

	for _, field := range fieldNames {
		column := val.FieldByName(field)
		if column == nil {
			log.Exitf("Failed when reflecting field %s", field)
		}

		switch t := column.(type) {
			case *reflect.StringValue:
				svalue := column.(*reflect.StringValue)
				scolumn := svalue.Get()
				columns = append(columns, scolumn)
			case *reflect.FloatValue:
				fvalue := column.(*reflect.FloatValue)
				fcolumn := fvalue.Get()
				columns = append(columns, fmt.Sprintf("%#v", fcolumn))
			case *reflect.IntValue:
				ivalue := column.(*reflect.IntValue)
				icolumn := ivalue.Get()
				columns = append(columns, fmt.Sprintf("%#v", icolumn))
			default:
				log.Exitf("Got a field we cannot handle: %s", field)
		}
	}

	io.WriteString(w, strings.Join(columns, ","))
	io.WriteString(w, "\n")
}

func SetHasErrors(perfdata []*PerfData) bool {
	for _, data := range perfdata {
		if data.ErrTotal > 0 {
			return true
		}
	}

	return false
}
