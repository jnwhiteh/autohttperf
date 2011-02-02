package main

import "flag"
import "fmt"
import "exec"
import "http"
import "io/ioutil"
import "log"
import "net"
import "os"
import "rpc"

type Args struct {
	Host                  string
	Port                  int
	URL                   string
	NumConnections        int
	ConnectionRate        int
	RequestsPerConnection int
}

type Result struct {
	Stdout     string
	Stderr     string
	ExitStatus int
}

type HTTPerf int

const (
	ERR_EXECNOTFOUND = "Could not find the 'httperf' executable: %s"
	ERR_RUNFAILED    = "Failed to run command: %s"
	ERR_WAIT         = "Failed when waiting on pid %d"
	ERR_NOTEXITED    = "Command did not properly exit: %s"
	ERR_READOUT      = "Could not read stdout: %s"
	ERR_READERR      = "Could not read stderr: %s"
)

func (h *HTTPerf) Benchmark(args *Args, result *Result) os.Error {
	// Try to find the 'httpperf' command, which must exist in the PATH
	// of the current user/environment.

	perfexec, err := exec.LookPath("httperf")
	if err != nil {
		return os.NewError(fmt.Sprintf(ERR_EXECNOTFOUND, err.String()))
	}

	// Build the httperf commandline and build a result to return
	argv := []string{
		perfexec,
		"--server", args.Host,
		"--port", fmt.Sprintf("%d", args.Port),
		"--uri", args.URL,
		"--num-conns", fmt.Sprintf("%d", args.NumConnections),
		"--rate", fmt.Sprintf("%d", args.ConnectionRate),
		"--num-calls", fmt.Sprintf("%d", args.RequestsPerConnection),
		"--hog",
	}

	log.Printf("++ [%p] Running benchmark of %s on port %d", args, args.Host, args.Port)
	log.Printf("   [%p] Input arguments: %#v", args, args)
	log.Printf("   [%p] Commandline arguments: %#v", args, argv)

	cmd, err := exec.Run(argv[0], argv, nil, "", exec.DevNull, exec.Pipe, exec.Pipe)
	if err != nil {
		return os.NewError(fmt.Sprintf(ERR_RUNFAILED, err.String()))
	}

	defer cmd.Close()

	log.Printf("   [%p] Process successfully started with PID: %d", args, cmd.Pid)

	output, err := ioutil.ReadAll(cmd.Stdout)
	if err != nil {
		return os.NewError(fmt.Sprintf(ERR_READOUT, err.String()))
	}
	errout, err := ioutil.ReadAll(cmd.Stderr)
	if err != nil {
		return os.NewError(fmt.Sprintf(ERR_READERR, err.String()))
	}

	log.Printf("   [%p] Finished reading stdout and stderr", args)

	w, err := cmd.Wait(0)

	log.Printf("-- [%p] Command joined and finished", args)

	if err != nil {
		return os.NewError(fmt.Sprintf(ERR_WAIT, cmd.Pid))
	} else if !w.Exited() {
		return os.NewError(fmt.Sprintf(ERR_NOTEXITED, w.String()))
	}

	result.Stdout = string(output)
	result.Stderr = string(errout)
	result.ExitStatus = int(w.WaitStatus)

	return nil
}

var host *string = flag.String("host", "*", "The host on which to bind the server")
var port *int = flag.Int("port", 1717, "The port on which to bind the server")

func main() {
	flag.Parse()

	httperf := new(HTTPerf)
	rpc.Register(httperf)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", fmt.Sprintf("%s:%d", *host, *port))
	if e != nil {
		log.Exit("listen error:", e)
	}

	log.Printf("Now listening for requests on %s:%d", *host, *port)
	http.Serve(l, nil)
}
