package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

const (
	imgNamePrefix  = "docker.io/library"
	dockerUnixSock = "unix:///var/run/docker.sock"
)

// Interface represents docker client
type Interface interface {
	Run(timeout int, args []string)
}

type dockerClient struct {
}

// New returns docker client interface
func New() Interface {
	return &dockerClient{}
}

func (dc dockerClient) Run(timeout int, args []string) {
	clt, err := client.NewClient(dockerUnixSock, "", nil, nil)
	if err != nil {
		log.Println("Fatal error while creating docker client: ", err)
		panic(err)
	}

	ctx := context.Background()
	reader, err := clt.ImagePull(ctx, fmt.Sprintf("%s/%s", imgNamePrefix, args[0]), types.ImagePullOptions{})
	if err != nil {
		log.Printf("Fatal error while pulling docker image: %s, err: %s", fmt.Sprintf("%s/%s", imgNamePrefix, args[0]), err)
		panic(err)
	}

	io.Copy(os.Stdout, reader)

	resp, err := clt.ContainerCreate(ctx, &container.Config{
		Image: args[0],
		Cmd:   args[1:],
		Tty:   false,
	}, nil, nil, "")
	if err != nil {
		log.Println("Fatal error while creating container instance: ", err)
		panic(err)
	}

	if err := clt.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	code, err := clt.ContainerWait(ctx, resp.ID)
	if err != nil {
		if strings.Contains(err.Error(), "context deadline exceeded") {
			log.Println("Exiting as deadline exceeded")
		} else {
			log.Printf("Fatal error while wait: statusCode: %d, err: %v\n", code, err)
		}
	} else {
		log.Println("Executed successfully with exit code: ", code)
	}
}
