package main

import "flag"
import "fmt"
import "log"
import "os"
import "rpc"
import "time"

// Runs a benchmark distributed over a set of clients. Returns a slice of the
// resulting PerfData structures and a boolean flags indicating if all workers
// successfully reported data, i.e. if the benchmark can be trusted.

func RunDistributedBenchmark(workers []*Worker, args *Args) ([]*PerfData, bool) {
	// Fairly simple, just split args up over however many client
	// we're connected to, and perform the benchmark. Don't collate
	// results or anything at the current time.

	// Generate a simple UID based on the current time in nanoseconds.
	nanotime := time.Nanoseconds()
	nanoid := fmt.Sprintf("%#v", nanotime)

	numWorkers := len(workers)
	log.Printf("Distributing benchmark over %d clients", numWorkers)
	log.Printf("Arguments: %#v", args)

	for _, worker := range workers {
		wargs := &Args{
			args.Host,
			args.Port,
			args.URL,
			args.NumConnections / numWorkers,
			args.ConnectionRate / numWorkers,
			args.RequestsPerConnection,
		}

		result := new(Result)

		call := worker.client.Go("HTTPerf.Benchmark", wargs, &result, nil)

		if call.Error != nil {
			log.Printf("[%s] Failed to open connection: %s", worker.id, call.Error)
			worker.args = wargs
			worker.result = nil
			worker.call = nil
		} else {
			log.Printf("[%s] Requested benchmark", worker.id)
			worker.args = wargs
			worker.result = result
			worker.call = call
			worker.date = time.Seconds()
		}
	}

	// Collect the PerfData into a slice
	results := make([]*PerfData, 0, len(workers))
	success := true

	for _, worker := range workers {
		if worker.call == nil {
			// This call was not successful
			success = false
		} else {
			call := <-worker.call.Done
			log.Printf("[%s] Got results", worker.id)
			if call.Error != nil {
				log.Printf("[%s] Error state reported: %s", worker.id, call.Error.String())
				success = false
			} else {
				perfdata, err := ParseResults(worker.result.Stdout, nanoid, worker.date, worker.args)
				if err != nil {
					// Error parsing, report this
					log.Printf("[%s] Error parsing perf data: %s\n", worker.id, err.String())
					success = false
				}
				results = append(results, perfdata)

				if len(worker.result.Stderr) > 0 {
					log.Printf("[%s] Stderr: %s", worker.id, worker.result.Stderr)
				}
			}
		}
	}

	return results, success
}

// Stress test a server for maximum number of connections per second
func StressTestConnections(workers []*Worker) {
	// A list of stress and steps, these should be sequential
	var stressRates = map[int]int{
		0: 100,   // Start benchmarking at rate 100
		100: 100, // Step up 100 rate each round
	}

	// Fetch the starting connection rate from the map
	rate := *startRate
	step := stressRates[rate]
	if step == 0 {
		step = stressRates[0]
	}

	errorState := false
	cooldownSteps := *cooldown

	// Output the TSV header
	WriteTSVHeader(os.Stdout)

	for {
		// Calculate the number of connections to request. Since we're distributing
		// both the rate and the number of connections over several workers, this
		// does not need to take that into account.
		//
		// 10 second duration with 300 connections per second is 3000 connections,
		// regardless of how many clients are used to distribute that load.
		numconns := *testLength * rate

		args := new(Args)
		args.Host = *server
		args.Port = *port
		args.URL = *url
		args.NumConnections = numconns
		args.ConnectionRate = rate
		args.RequestsPerConnection = *requests

		data, ok := RunDistributedBenchmark(workers, args)
		if !ok {
			log.Printf("Stress test for rate %d did not fully succeed", rate)
		}

		WriteTSVParseDataSet(os.Stdout, data)

		// Check if the data set is over the error threshold
		hasErrors := SetHasErrors(data, *numErrors)

		if errorState && !hasErrors {
			log.Printf("Exiting error state, server seems to have recovered")
			errorState = false
			cooldownSteps = *cooldown
		} else if !errorState && hasErrors {
			log.Printf("Entering an error state, will cooldown for %d rounds", cooldownSteps)
			errorState = true
		}

		if errorState {
			cooldownSteps = cooldownSteps - 1
			log.Printf("In an error state with %d rounds to go", cooldownSteps)
		}

		// Stop benchmarking when we've run out of cooldown steps
		if cooldownSteps < 0 {
			break
		}

		// Increment the rate/step accordingly.
		rate = rate + step
		if newStep, ok := stressRates[rate]; ok {
			step = newStep
		}

		log.Printf("Current rate: %d, step: %d", rate, step)

		// Perform any sleep, as directed
		log.Printf("Sleeping for %d seconds", *sleep)
		var sleeptime int64 = int64(*sleep) * 1000000000
		time.Sleep(sleeptime)
		log.Printf("Done sleeping")
	}
}

// Stress test a server for maximum number of requests per second
func StressTestRequests(workers []*Worker) {
}

func RunManualBenchmark(workers []*Worker) {
	args := &Args{
		*server,
		*port,
		*url,
		*numConns,
		*connRate,
		*requests,
	}

	data, ok := RunDistributedBenchmark(workers, args)
	if !ok {
		log.Printf("Manual benchmark did not fully succeed")
	}

	// Write the TSV header
	WriteTSVHeader(os.Stdout)
	// Write out the perf data for each benchmark
	for idx, perfdata := range data {
		if *dumpraw {
			log.Printf("Client %d output: \n%s\n", idx, perfdata.Raw)
		}
		WriteTSVParseData(os.Stdout, perfdata)
	}
}

// General options that every single mode will require
var help *bool = flag.Bool("help", false, "Display usage information")
var server *string = flag.String("server", "localhost", "The hostname or IP address of the server")
var port *int = flag.Int("port", 80, "The port on which to bind the server")
var url *string = flag.String("url", "/", "The URL to be requested")
var timeout *int = flag.Int("timeout", 5, "Amount of time before a request is considered unfulfilled")

// Flags that can be used to turn a mode on or off, these are combined and
// will be executed in the order they are specified here, not the order they
// are specified on the commandline.
var modeStressConn *bool = flag.Bool("stressconn", false, "Perform a connection stress test")
var modeStressReqs *bool = flag.Bool("stressreqs", false, "Perform a request stress test")
var modeManual *bool = flag.Bool("manual", false, "Perform a manual benchmark")

// Manual mode options
var numConns *int = flag.Int("numconns", 6000, "The number of connections to be opened (manual only)")
var connRate *int = flag.Int("connrate", 200, "The rate of new connections (connections per second) (manual only)")
var requests *int = flag.Int("requests", 5, "The number of requests sent per connection (manual only)")

// Stress test options
var numErrors *int = flag.Int("numerrors", 500, "The maximum acceptable number of errors to indicate 'stressed' (stress only)")
var cooldown *int = flag.Int("cooldown", 3, "The number of steps to take following an 'error state' (stress only)")
var testLength *int = flag.Int("duration", 60, "The duration of each 'step' of the stress test in seconds (stress only)")
var sleep *int = flag.Int("sleeptime", 5, "The amount of time (in seconds) to sleep between each round (stress only)")
var startRate *int = flag.Int("startrate", 100, "The connection start rate for the stress test")
var dumpraw *bool = flag.Bool("dumpraw", true, "Dump the raw client output to stderr")

var PrintUsage = func() {
	fmt.Fprintf(os.Stderr, "Usage of %s: \"host1:port1\" ...\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Parse()

	if *help {
		PrintUsage()
		return
	}

	// Build a slice of RPC clients, as specified by the user as arguments
	workers := make([]*Worker, 0, 5)

	for idx, arg := range flag.Args() {
		log.Printf("Opening RPC connection to %s", arg)
		client, err := rpc.DialHTTP("tcp", arg)
		log.Printf("New RPC connection %p", client)

		if err != nil {
			log.Fatalf("Could not connect to client %s: %s", arg, err)
		}

		id := fmt.Sprintf("%s:%d", arg, idx)
		worker := &Worker{arg, id, client, nil, nil, 0, nil}
		workers = append(workers, worker)
	}

	if !*modeStressConn && !*modeStressReqs && !*modeManual {
		log.Fatalf("No mode selected, please supply one of -stressconn, -stressreqs or -manual")
	}

	if *modeManual {
		RunManualBenchmark(workers)
	}

	if *modeStressConn {
		StressTestConnections(workers)
	}

	if *modeStressReqs {
		StressTestRequests(workers)
	}
}
