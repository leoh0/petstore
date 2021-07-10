.PHONY: update
update:
	go get github.com/deepmap/oapi-codegen/pkg/codegen
	go generate -v ./...
	go mod tidy
	bazel run //:gazelle
	bazel run //:gazelle -- update-repos -from_file=go.mod -to_macro repos.bzl%go_repositories
