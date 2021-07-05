module github.com/portainer/docker-compose-wrapper

go 1.15

require (
	github.com/compose-spec/compose-go v0.0.0-20210702130122-154903ab827c
	github.com/docker/cli v20.10.7+incompatible
	github.com/docker/compose-cli v1.0.18-0.20210703141553-1da8be257bb0
)

replace github.com/jaguilar/vt100 => github.com/tonistiigi/vt100 v0.0.0-20190402012908-ad4c4a574305
