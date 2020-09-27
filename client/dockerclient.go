package client

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

type Interface interface {
	Run(ctx context.Context)
}

type dockerClient struct {
}

func New() Interface {
	return &dockerClient
}

func (dc dockerClient) Run(ctx context.Context, args []string) {
	if ctx == nil {
		ctx = context.Background()
	}

	clt, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Println("Fatal error while creating docker client: ", err)
		panic(err)
	}

	reader, err := clt.ImagePull(ctx, args[0])
	if err != nil {
		log.Println("Fatal error while pulling docker image: ", err)
		panic(err)
	}

	io.Copy(os.Stdout, reader)

	resp, err := clt.ContainerCreate(ctx, &container.Config{
		Image: args[0],
		Cmd:   args[1:],
		Tty:   false,
	}, nil, nil, nil, "")
	if err != nil {
		log.Println("Fatal error while creating container instance: ", err)
		panic(err)
	}

	if err := clt.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	chStatus, chErr := clt.ContainerWait(ctx, resp.ID)
	select {
	case err <- chErr:
		if err != nil {
			panic(err)
		}
	case <-chStatus:
	}

	out, err := clt.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
}
