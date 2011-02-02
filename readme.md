This project is intended to make it easier to run stress test benchmarks
against web servers by using http and periodic analysis of the results.
Autobench is great, but can only use static ranges and steps which makes it
less useful for my purposes.

Example usage:

        $ autohttperf -help
        Usage of ./autohttperf: "host1:port1" ...
          -server="localhost": The hostname or IP address of the server
          -cooldown=3: The number of steps to take following an 'error state' (stress only)
          -numerrors=500: The maximum acceptable number of errors to indicate 'stressed' (stress only)
          -stressreqs=false: Perform a request stress test
          -manual=false: Perform a manual benchmark
          -timeout=5: Amount of time before a request is considered unfulfilled
          -port=80: The port on which to bind the server
          -url="/": The URL to be requested
          -numconns=6000: The number of connections to be opened (manual only)
          -stressconn=false: Perform a connection stress test
          -connrate=200: The rate of new connections (connections per second) (manual only)
          -help=false: Display usage information
          -duration=60: The duration of each 'step' of the stress test in seconds (stress only)
          -sleeptime=5: The amount of time (in seconds) to sleep between each round (stress only)
          -requests=5: The number of requests sent per connection (manual only)
        
        autohttperf --server 10.0.0.125 --stressconn worker1.myhost.com:1717 worker2.myhost.com:1717

This is incredibly limited right now, but I am actively using it in order to
benchmark a series of servers from 3 different client machines.  Right now it
doesn't work, but feel free to take a look.
