module github.com/vmware/carbon-black-adapter-for-harbor

go 1.16

require (
	github.com/docker/distribution v2.7.1+incompatible
	github.com/gin-gonic/gin v1.7.0
	github.com/google/uuid v1.3.0
	github.com/sirupsen/logrus v1.8.1
	github.com/vmware/carbon-black-cloud-container-cli v1.0.1
)

replace (
	github.com/containerd/containerd => github.com/containerd/containerd v1.5.9
	github.com/go-restruct/restruct => github.com/go-restruct/restruct v1.2.0-alpha
	github.com/opencontainers/runc => github.com/opencontainers/runc v1.0.3
)
