/*
 * Copyright 2016 Charith Ellawala
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package gcgrpcpool

import (
	"fmt"
	"sync"

	log "github.com/Sirupsen/logrus"

	"github.com/charithe/gcgrpcpool/gcgrpc"
	"github.com/golang/groupcache"
	"github.com/golang/groupcache/consistenthash"
	pb "github.com/golang/groupcache/groupcachepb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const defaultReplicas = 50

type GRPCPool struct {
	self        string
	opts        GRPCPoolOptions
	mu          sync.Mutex
	peers       *consistenthash.Map
	grpcGetters map[string]*grpcGetter
}

type GRPCPoolOptions struct {
	Replicas        int
	HashFn          consistenthash.Hash
	PeerDialOptions []grpc.DialOption
}

func NewGRPCPool(self string, server *grpc.Server) *GRPCPool {
	return NewGRPCPoolOptions(self, server, nil)
}

var grpcPoolCreated bool

func NewGRPCPoolOptions(self string, server *grpc.Server, opts *GRPCPoolOptions) *GRPCPool {
	if grpcPoolCreated {
		panic("NewGRPCPool must be called only once")
	}

	grpcPoolCreated = true

	pool := &GRPCPool{
		self:        self,
		grpcGetters: make(map[string]*grpcGetter),
	}

	if opts != nil {
		pool.opts = *opts
	}

	if pool.opts.Replicas == 0 {
		pool.opts.Replicas = defaultReplicas
	}

	if pool.opts.PeerDialOptions == nil {
		pool.opts.PeerDialOptions = []grpc.DialOption{grpc.WithInsecure()}
	}

	pool.peers = consistenthash.New(pool.opts.Replicas, pool.opts.HashFn)
	groupcache.RegisterPeerPicker(func() groupcache.PeerPicker { return pool })
	gcgrpc.RegisterPeerServer(server, pool)
	return pool
}

func (gp *GRPCPool) Set(peers ...string) {
	gp.mu.Lock()
	defer gp.mu.Unlock()
	gp.peers = consistenthash.New(gp.opts.Replicas, gp.opts.HashFn)
	tempGetters := make(map[string]*grpcGetter, len(peers))
	for _, peer := range peers {
		if getter, exists := gp.grpcGetters[peer]; exists == true {
			tempGetters[peer] = getter
			delete(gp.grpcGetters, peer)
		} else {
			getter, err := newGRPCGetter(peer, gp.opts.PeerDialOptions...)
			if err != nil {
				log.Warnf("Failed to open connection to [%s] : %v", peer, err)
			} else {
				tempGetters[peer] = getter
				gp.peers.Add(peer)
			}
		}
	}

	for p, g := range gp.grpcGetters {
		g.close()
		delete(gp.grpcGetters, p)
	}

	gp.grpcGetters = tempGetters
}

func (gp *GRPCPool) PickPeer(key string) (groupcache.ProtoGetter, bool) {
	gp.mu.Lock()
	defer gp.mu.Unlock()

	if gp.peers.IsEmpty() {
		return nil, false
	}

	if peer := gp.peers.Get(key); peer != gp.self {
		return gp.grpcGetters[peer], true
	}
	return nil, false
}

func (gp *GRPCPool) Retrieve(ctx context.Context, req *gcgrpc.RetrieveRequest) (*gcgrpc.RetrieveResponse, error) {
	group := groupcache.GetGroup(req.Group)
	if group == nil {
		return nil, fmt.Errorf("Unable to find group [%s]", req.Group)
	}
	group.Stats.ServerRequests.Add(1)
	var value []byte
	err := group.Get(ctx, req.Key, groupcache.AllocatingByteSliceSink(&value))
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve [%s]: %v", req, err)
	}

	return &gcgrpc.RetrieveResponse{Value: value}, nil
}

func (gp *GRPCPool) AddPeers(ctx context.Context, peers *gcgrpc.Peers) (*gcgrpc.Ack, error) {
	gp.mu.Lock()
	defer gp.mu.Unlock()
	for _, peer := range peers.PeerAddr {
		if _, exists := gp.grpcGetters[peer]; exists != true {
			getter, err := newGRPCGetter(peer, gp.opts.PeerDialOptions...)
			if err != nil {
				log.Warnf("Failed to open connection to [%s]: %v", peer, err)
			} else {
				log.Infof("Adding peer [%s]", peer)
				gp.grpcGetters[peer] = getter
				gp.peers.Add(peer)
			}
		}
	}
	return &gcgrpc.Ack{}, nil

}
func (gp *GRPCPool) RemovePeers(ctx context.Context, peers *gcgrpc.Peers) (*gcgrpc.Ack, error) {
	return &gcgrpc.Ack{}, nil
}
func (gp *GRPCPool) SetPeers(ctx context.Context, peers *gcgrpc.Peers) (*gcgrpc.Ack, error) {
	return &gcgrpc.Ack{}, nil
}

type grpcGetter struct {
	address string
	conn    *grpc.ClientConn
}

func newGRPCGetter(address string, dialOpts ...grpc.DialOption) (*grpcGetter, error) {
	conn, err := grpc.Dial(address, dialOpts...)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to [%s]: %v", address, err)
	}
	return &grpcGetter{address: address, conn: conn}, nil
}

func (g *grpcGetter) Get(ctx groupcache.Context, in *pb.GetRequest, out *pb.GetResponse) error {
	client := gcgrpc.NewPeerClient(g.conn)
	resp, err := client.Retrieve(context.Background(), &gcgrpc.RetrieveRequest{Group: *in.Group, Key: *in.Key})
	if err != nil {
		return fmt.Errorf("Failed to GET [%s]: %v", in, err)
	}

	out.Value = resp.Value
	return nil
}

func (g *grpcGetter) close() {
	if g.conn != nil {
		g.conn.Close()
	}
}
