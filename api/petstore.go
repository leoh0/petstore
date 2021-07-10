//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --package=api --generate server -o server.gen.go spec.yaml
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --package=api --generate types -o type.gen.go spec.yaml
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --package=api --generate spec -o spec.gen.go spec.yaml
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --package=api --generate client -o client.gen.go spec.yaml

package api
