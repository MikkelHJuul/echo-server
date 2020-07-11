DOCKER_REPO = mjuul/echo-server

GO_MATRIX_OS := linux
GO_MATRIX_ARCH := amd64

-include .makefiles/Makefile
-include .makefiles/pkg/go/v1/Makefile
-include .makefiles/pkg/docker/v1/Makefile

.makefiles/%:
	@curl -sfL https://makefiles.dev/v1 | bash /dev/stdin "$@"
