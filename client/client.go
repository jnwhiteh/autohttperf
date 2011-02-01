package main

import "flag"
import "fmt"
import "log"
import "os"
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

// Runs a benchmark distributed over a set of clients
func RunDistributedBenchmark(workers []*Worker, args *Args) {
	// Fairly simple, just split args up over however many client
	// we're connected to, and perform the benchmark. Don't collate
	// results or anything at the current time.

	numWorkers := len(workers)
	log.Printf("Distributing benchmark over %d clients", numWorkers)
	log.Printf("Arguments: %#v", args)

	for _, worker := range workers {
		wargs := &Args{
			args.Host,
			args.Port,
			args.URL,
			args.NumConnections / numWorkers,
			args.ConnectionRate,
			args.RequestsPerConnection,
			args.Hog,
		}

		result := new(Result)

		call := worker.client.Go("HTTPerf.Benchmark", wargs, &result, nil)

		if call.Error != nil {
			log.Printf("[%s] Failed to open connection: %s", worker.id, call.Error)
			worker.result = nil
			worker.call = nil
		} else {
			log.Printf("[%s] Requested benchmark", worker.id)
			worker.result = result
			worker.call = call
		}
	}

	for _, worker := range workers {
		if worker.call != nil {
			call := <-worker.call.Done
			log.Printf("[%s] Got results", worker.id)
			if call.Error != nil {
				log.Printf("[%s] Error state reported: %s", worker.id, call.Error.String())
			}
		}
	}
}

// Stress test a server for maximum number of connections per second
func StressTestConnections(workers []*Worker) {
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
		true,
	}
	RunDistributedBenchmark(workers, args)
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
var numErrors *int = flag.Int("numerrors", 500, "The maximum acceptable number of errors to indicate 'stressed'")
var cooldown *int = flag.Int("cooldown", 3, "The number of steps to take following an 'error state'")

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
			log.Exitf("Could not connect to client %s: %s", arg, err)
		}

		id := fmt.Sprintf("%s:%d", arg, idx)
		worker := &Worker{arg, id, client, nil, nil}
		workers = append(workers, worker)
	}

	if !*modeStressConn && !*modeStressReqs && !*modeManual {
		log.Exitf("No mode selected, please supply one of -stressconn, -stressreqs or -manual")
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
