package clients

import (
	"context"
	"io"

	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/protocol"
)

// DockerClient is a client interface for interacting with docker
type DockerClient interface {
	PullImage(ctx context.Context, refStr string) error
	CreatePublicNetwork(ctx context.Context, name string) (string, error)
	CreateInternalNetwork(ctx context.Context, name string) (string, error)
	AttachNetwork(ctx context.Context, containerID string, networkID string) error
	GetContainers(ctx context.Context) (DockerContainerList, error)
	StartContainer(ctx context.Context, config DockerContainerConfig) (*DockerContainer, error)
	StopContainer(ctx context.Context, ID string) error
	Prune(ctx context.Context) error
	HasLocalImage(ctx context.Context, ref string) bool
}

// MessageClient receives and publishes messages.
type MessageClient interface {
	Subscribe(subject string, handler interface{})
	Publish(subject string, payload interface{})
}

// AgentClient makes the gRPC requests to evaluate block and txs and receive results.
type AgentClient interface {
	MustDial(config.AgentConfig)
	protocol.AgentClient
	io.Closer
}
