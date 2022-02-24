package main

import (
	"context"
	"fmt"
	"github.com/ambazon/network"
	libp2pnet "github.com/libp2p/go-libp2p-core/network"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
)


func main() {

	config, err := network.ParseFlags()
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	node, err := network.NewNode(ctx, config.Seed, 0)
	if err != nil {
		log.Fatal(err)
	}

	node.Host.SetStreamHandler(protocol.ID(config.ProtocolID), func(s libp2pnet.Stream) {
		fmt.Println("Got a new stream!")
	})

	dht, err := network.NewDHT(ctx, node.Host, config.DiscoveryPeers)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Bootstrapping the DHT")
	if err = dht.Bootstrap(ctx); err != nil {
		panic(err)
	}

	// Let's connect to the bootstrap nodes first. They will tell us about the
	// other nodes in the network.
	var wg sync.WaitGroup
	for _, peerAddr := range config.BootstrapPeers {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := node.Host.Connect(ctx, *peerinfo); err != nil {
				log.Println(err)
			} else {
				log.Println("Connection established with bootstrap node:", *peerinfo)
			}
		}()
	}
	wg.Wait()

	go network.Discover(ctx, node.Host, dht, config.Rendezvous)

	run(node.Host, cancel)
}


func run(h host.Host, cancel func()) {
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	<-c

	fmt.Printf("\rExiting...\n")

	cancel()

	if err := h.Close(); err != nil {
		panic(err)
	}
	os.Exit(0)
}