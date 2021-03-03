module github.com/portainer/docker-compose-wrapper

go 1.15

require (
	github.com/containerd/containerd v1.4.3 // indirect
	github.com/portainer/libcompose v0.5.3
	github.com/stretchr/testify v1.7.0 // indirect
)

replace github.com/docker/docker => github.com/docker/engine v1.4.2-0.20200204220554-5f6d6f3f2203
