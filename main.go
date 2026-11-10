package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/docker/go-connections/sockets"
	"github.com/docker/go-plugins-helpers/volume"
)

const (
	config    = "/etc/linstor/docker-volume.conf"
	plugin    = "linstor"
	socketDir = "/run/docker/plugins"
)

var (
	root = filepath.Join(volume.DefaultDockerRootDirectory, plugin)
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func main() {
	node, err := os.Hostname()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	if err := os.MkdirAll(socketDir, 0755); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	listener, err := sockets.NewUnixSocket(filepath.Join(socketDir, plugin+".sock"), 0)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer listener.Close()

	driver := NewLinstorDriver(config, node, root)
	server := newVolumeServer(driver)
	fmt.Println(server.Serve(listener))
}
