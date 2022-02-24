package network

import (
	"context"
	"crypto/rand"
	"fmt"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"io"

	mrand "math/rand"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/multiformats/go-multiaddr"
)

type Node struct {
	Ctx context.Context
	Host host.Host
	ActivePeers map[string]string
	KadDHT *dht.IpfsDHT
	Discovery *discovery.RoutingDiscovery
	PubSub *pubsub.PubSub

}

func NewNode(ctx context.Context, seed int64, port int) (Node, error) {

	// If the seed is zero, use real cryptographic randomness. Otherwise, use a
	// deterministic randomness source to make generated keys stay the same
	// across multiple runs
	var r io.Reader
	if seed == 0 {
		r = rand.Reader
	} else {
		r = mrand.New(mrand.NewSource(seed))
	}

	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		return Node{}, err
	}

	addr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port))

	host, err := libp2p.New(
		libp2p.ListenAddrs(addr),
		libp2p.Identity(priv),
	)

	if err != nil {
		return Node{}, err
	}

	return Node{
		Host:        host,
		ActivePeers: make(map[string]string),
	}, nil
}

func (host *Node) Call(peer Node) {

}

func (host *Node) MultiCalls(peers []Node) {

}
