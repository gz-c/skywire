package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/skycoin/skywire/app"
	"github.com/skycoin/skywire/messages"
	"github.com/skycoin/skywire/nodemanager"
)

func main() {

	//messages.SetDebugLogLevel()
	messages.SetInfoLogLevel()

	var (
		err error
	)

	args := os.Args
	if len(args) < 2 {
		printHelp()
		return
	}

	hopsStr := os.Args[1]

	if hopsStr == "--help" {
		printHelp()
		return
	}

	hops, err := strconv.Atoi(hopsStr)
	if err != nil {
		fmt.Println("\nThe first argument should be a number of hops\n")
		return
	}

	if hops < 1 {
		fmt.Println("\nThe number of hops should be a positive number > 0\n")
		return
	}

	cfg := &nodemanager.NodeManagerConfig{
		Domain:           "mesh.network",
		CtrlAddr:         "127.0.0.1:5999",
		AppTrackerAddr:   "",
		RouteManagerAddr: "",
		LogisticsServer:  "",
	}

	meshnet, err := nodemanager.NewNetwork(cfg)
	if err != nil {
		panic(err)
	}

	defer meshnet.Shutdown()

	clientNode, serverNode := meshnet.CreateSequenceOfNodes(hops+1, 15000)

	serverId := messages.MakeAppId("socksServer0")

	server, err := app.NewSocksServer(serverId, serverNode.AppTalkAddr(), "0.0.0.0:8001")
	if err != nil {
		panic(err)
	}
	defer server.Shutdown()

	client, err := app.NewSocksClient(messages.MakeAppId("client0"), clientNode.AppTalkAddr(), "0.0.0.0:8000")
	if err != nil {
		panic(err)
	}
	defer client.Shutdown()

	err = client.Connect(serverId, serverNode.Id().Hex())
	if err != nil {
		panic(err)
	}

	client.Listen()

}

func printHelp() {
	fmt.Println("\nFORMAT: go run socks.go n , where n is a number of hops")
	fmt.Println("\nUsage example for 10 meshnet hops:")
	fmt.Println("\ngo run socks.go 10")
	fmt.Println("\nNumber of hops should be more than 0\n")
}
