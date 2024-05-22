build:
	# build binary for linux
	GOOS=linux GOARCH=amd64 go build \
		-tags="containers_image_openpgp exclude_graphdriver_devicemapper exclude_graphdriver_btrfs" \
		-o=bin/harboradapter \
		github.com/vmware/carbon-black-adapter-for-harbor/cmd/harboradapter

docker-build publish check-release-var: release_version ?=

check-release-var:
ifndef release_version
	$(error release_version is required to publish a release)
endif

docker-build:
	docker build . -t cbartifactory/harbor_adapter:$(release_version)

publish: check-release-var docker-build
	docker push cbartifactory/harbor_adapter:$(release_version)
