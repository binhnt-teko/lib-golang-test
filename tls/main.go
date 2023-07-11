package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/blcvn/lib-golang-test/tls/spring"

	pb "github.com/bbengfort/sping/echo"
	"github.com/urfave/cli"
)

// Default values for various options.
const (
	DefaultPort  = uint(3264)
	DefaultPings = uint(8)
	DefaultDelay = int64(100)
)

func signalHandler() {
	// Make signal channel and register notifiers for Interupt and Terminate
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	signal.Notify(sigchan, syscall.SIGTERM)

	// Block until we receive a signal on the channel
	<-sigchan

	// Log the shutdown
	log.Println("shutting down!")
	os.Exit(0)
}

func main() {

	// Create the command line application
	app := cli.NewApp()
	app.Name = "sping"
	app.Usage = "implements simple, secure ping with grpc"

	// Describe the commands in the app
	app.Commands = []cli.Command{
		{
			Name:   "serve",
			Usage:  "run the sping server",
			Action: startServer,
			Flags: []cli.Flag{
				cli.UintFlag{
					Name:  "p, port",
					Usage: "specify the port to listen on",
					Value: DefaultPort,
				},
				cli.StringFlag{
					Name:  "n, name",
					Usage: "specify the name of the client",
				},
			},
		},
		{
			Name:   "echo",
			Usage:  "run the sping client",
			Action: startClient,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "n, name",
					Usage: "specify the name of the client",
				},
				cli.UintFlag{
					Name:  "p, port",
					Usage: "specify the port to ping to",
					Value: DefaultPort,
				},
				cli.UintFlag{
					Name:  "l, limit",
					Usage: "specify the max number of pings to send",
					Value: DefaultPings,
				},
				cli.Int64Flag{
					Name:  "d, delay",
					Usage: "the delay between pings in milliseconds",
					Value: DefaultDelay,
				},
			},
		},
	}

	// Run the application
	app.Run(os.Args)

}

// NewClient returns a ping client with the specified options.
func NewClient(dailer spring.Dailer, address, name string, delay int64, limit uint) *spring.PingClient {
	conn, err := dailer(address)
	if err != nil {
		log.Fatalln(err)
	}
	return &spring.PingClient{
		Name:             name,
		Delay:            time.Duration(delay) * time.Millisecond,
		Limit:            limit,
		Connection:       conn,
		SecurePingClient: pb.NewSecurePingClient(conn),
	}
}

// NewServer returns a ping server with specified options.
func NewServer() *spring.PingServer {
	fmt.Println("Start NewServer ")

	return new(spring.PingServer)
}

// Run the ping server
func startServer(c *cli.Context) error {
	fmt.Println("Start server with port: ", c.Uint("port"))
	go signalHandler()

	server := NewServer()
	err := server.Serve(c.Uint("port"))

	if err != nil {
		fmt.Println("Cannot start server ")
		return cli.NewExitError(err.Error(), 1)
	}
	fmt.Println("startServer finished ")

	return nil
}

// Run the ping client
func startClient(c *cli.Context) error {
	go signalHandler()
	var err error

	// Get the addr to ping to with the associated port
	if c.NArg() != 1 {
		return cli.NewExitError("specify an address to ping to", 1)
	}
	addr := fmt.Sprintf("%s:%d", c.Args()[0], c.Uint("port"))
	fmt.Println("Start client with addr: ", addr)

	// Get the hostname if no name is specified
	name := c.String("name")
	if name == "" {
		if name, err = os.Hostname(); err != nil {
			return cli.NewExitError("no hostname for the pinger", 1)
		}
	}

	fmt.Println("Start client with addr: ", addr, name, c.Int64("delay"), c.Uint("limit"))

	// Create the client to start pinging to.
	client := NewClient(spring.MutualTLS, addr, name, c.Int64("delay"), c.Uint("limit"))
	defer client.Connection.Close()

	if err = client.Run(); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}
