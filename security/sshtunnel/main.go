package main

import (
	"fmt"
	"log"
	"os"
	"time"

	tl "github.com/lib-golang-test/security/sshtunnel/tunnel"
	"golang.org/x/crypto/ssh"
)

func main() {
	// Setup the tunnel, but do not yet start it yet.
	remote := "ec2-user@jumpbox.us-east-1.mydomain.com"
	authMethod := ssh.Password("password")
	// authMethod := ssh.PrivateKeyFile("path/to/private/key.pem")

	destination := "dqrsdfdssdfx.us-east-1.redshift.amazonaws.com:5439"
	localPort := "8080"
	tunnel, err := tl.NewSSHTunnel(remote, authMethod, destination, localPort)
	if err != nil {
		fmt.Printf("NewSSHTunnel: %s \n", err.Error())
		return
	}
	tunnel.Log = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
	go tunnel.Start()
	time.Sleep(100 * time.Millisecond)

	// NewSSHTunnel will bind to a random port so that you can have
	// multiple SSH tunnels available. The port is available through:
	//   tunnel.Local.Port

	// You can use any normal Go code to connect to the destination
	// server through localhost. You may need to use 127.0.0.1 for
	// some libraries.
	//
	// Here is an example of connecting to a PostgreSQL server:
	// conn := fmt.Sprintf("host=127.0.0.1 port=%d username=foo", tunnel.Local.Port)
	// db, err := sql.Open("postgres", conn)
	// ...
}
