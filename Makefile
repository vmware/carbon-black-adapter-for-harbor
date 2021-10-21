build:
	# build binary for linux
	GOOS=linux GOARCH=amd64 go build \
		-tags="containers_image_openpgp exclude_graphdriver_devicemapper exclude_graphdriver_btrfs" \
		-o=bin/harboradapter \
		github.com/vmware/carbon-black-adapter-for-harbor/cmd/harboradapter
