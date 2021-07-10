# heavily modified petstore

![](https://source.unsplash.com/collection/983219/1600x900)

초보자를 위한 OpenAPI 3 기반의 golang API 서버를 개발 및 배포 하기 위한 예제

kubernetes 기반으로 bazel 형상관리와 개발을 위한 tilt 이용

## 사전 준비

mac 기반으로만 작성했습니다. `brew`, `docker` 등은 기본으로 설치되있다고 가정 함

* [asdf](https://github.com/asdf-vm/asdf)

    ``` sh
    brew install asdf
    ```

* [bazel](https://bazel.build/)

    ``` sh
    brew install bazel
    ```

* [kind](https://kind.sigs.k8s.io/)

    ``` sh
    brew install kind
    ```

* [tilt](https://tilt.dev/)

    ``` sh
    brew install tilt
    ```

## 단계 설명

### golang 기반의 OpenAPI server

#### 1. go project setup

asdf 를 이용해서 golang version을 지정하고 필요한 버전을 인스톨

``` sh
$ cat << 'EOF' > .tool-versions
golang 1.16.5
EOF

$ asdf install
Platform 'darwin' supported!
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100  124M  100  124M    0     0  24.7M      0  0:00:05  0:00:05 --:--:-- 26.6M
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100    64  100    64    0     0    204      0 --:--:-- --:--:-- --:--:--   203
verifying checksum
/Users/al/.asdf/downloads/golang/1.16.5/archive.tar.gz: OK
checksum verified
```

`go mod init` 으로 `go.mod` 를 생성하고 go module을 사용할 준비를 함

``` sh
$ go mod init github.com/leoh0/petstore
go: creating new go.mod: module github.com/leoh0/petstore

$ cat go.mod
module github.com/leoh0/petstore

go 1.16
```

#### 2. .gitignore 생성

`go` 와 향후에 사용할 `bazel`을 위한 .gitignore 파일 생성

``` sh
$ curl -o .gitignore https://www.toptal.com/developers/gitignore/api/go,bazel
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   763  100   763    0     0   1099      0 --:--:-- --:--:-- --:--:--  1097
```

#### 3. openapi 서버를 위해 spec 파일 생성

[OpenAPI 3.0 예제](https://github.com/OAI/OpenAPI-Specification/blob/main/examples/v3.0/petstore-expanded.yaml)를 다운로드

``` sh
$ mkdir api

$ curl -o api/spec.yaml https://raw.githubusercontent.com/OAI/OpenAPI-Specification/main/examples/v3.0/petstore-expanded.yaml
...
```

아니면 직접 편집해서 사용. [OpenAPI 3.1.0](https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.1.0.md) 스펙 참고

``` sh
$ open https://editor.swagger.io/?url=https://raw.githubusercontent.com/OAI/OpenAPI-Specification/main/examples/v3.0/petstore-expanded.yaml
```

편집한 spec은 `api/spec.yaml` 에 위치시킴

#### 4. spec 파일 기준으로 code generation 준비

generation 을 위한 파일 생성

``` sh
cat << 'EOF' > api/petstore.go
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --package=api --generate server -o server.gen.go spec.yaml
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --package=api --generate types -o type.gen.go spec.yaml
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --package=api --generate spec -o spec.gen.go spec.yaml
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --package=api --generate client -o client.gen.go spec.yaml

package api
EOF
```

#### 5. spec 파일 기준으로 code generation

oapi-codegen 으로 코드 생성

``` sh
$ go get github.com/deepmap/oapi-codegen/pkg/codegen
go get: added github.com/deepmap/oapi-codegen v1.8.1

$ go generate -v ./...
api/petstore.go

$ tree -a -I .git
.
├── .gitignore
├── .tool-versions
├── api
│   ├── client.gen.go
│   ├── petstore.go
│   ├── server.gen.go
│   ├── spec.gen.go
│   ├── spec.yaml
│   └── type.gen.go
├── go.mod
└── go.sum

1 directory, 10 files
```

generation 된 코드들에서 필요한 의존성을 업데이트 함

``` sh
$ go mod tidy
```

#### 6. 실제 비지니스 로직 작성

이제 실제 해당 스펙에 맞는 비지니스 로직을 작성해야 함

여기에서는 간단히 [oapi-codegen petstore sample code](https://github.com/deepmap/oapi-codegen/tree/master/examples/petstore-expanded/echo) 로 예제 코드 생성

``` sh
$ tree -a -I .git
.
├── .gitignore
├── .tool-versions
├── api
│   ├── client.gen.go
│   ├── petstore.go # 예제 비지니스 로직 코드 추가 우측 코드 변경 https://github.com/deepmap/oapi-codegen/blob/master/examples/petstore-expanded/echo/api/petstore.go
│   ├── server.gen.go
│   ├── spec.gen.go
│   ├── spec.yaml
│   └── type.gen.go
├── cmd
│   ├── petstore.go # 예제 비지니스 로직 코드 추가 우측 코드 변경 https://github.com/deepmap/oapi-codegen/blob/master/examples/petstore-expanded/echo/petstore.go
│   └── petstore_test.go # 예제 비지니스 로직 코드 추가 우측 코드 변경 https://github.com/deepmap/oapi-codegen/blob/master/examples/petstore-expanded/echo/petstore_test.go
├── go.mod
└── go.sum

2 directories, 12 files
```

관련 코드 추가 뒤 또 필요한 의존성 업데이트

``` sh
$ go mod tidy
```

여기까지 진행했으면 문제 없이 테스트 진행 가능

``` sh
$ go test -v ./...
?   	github.com/leoh0/petstore/api	[no test files]
=== RUN   TestPetStore
{"time":"2021-07-10T20:24:20.124991+09:00","id":"","remote_ip":"192.0.2.1","host":"example.com","method":"POST","uri":"/pets","user_agent":"","status":201,"error":"","latency":31234,"latency_human":"31.234µs","bytes_in":0,"bytes_out":44}
{"time":"2021-07-10T20:24:20.125145+09:00","id":"","remote_ip":"192.0.2.1","host":"example.com","method":"GET","uri":"/pets/1000","user_agent":"","status":200,"error":"","latency":11200,"latency_human":"11.2µs","bytes_in":0,"bytes_out":44}
{"time":"2021-07-10T20:24:20.125187+09:00","id":"","remote_ip":"192.0.2.1","host":"example.com","method":"GET","uri":"/pets/27179095781","user_agent":"","status":404,"error":"","latency":12001,"latency_human":"12.001µs","bytes_in":0,"bytes_out":64}
{"time":"2021-07-10T20:24:20.125222+09:00","id":"","remote_ip":"192.0.2.1","host":"example.com","method":"POST","uri":"/pets","user_agent":"","status":201,"error":"","latency":4290,"latency_human":"4.29µs","bytes_in":0,"bytes_out":44}
{"time":"2021-07-10T20:24:20.125249+09:00","id":"","remote_ip":"192.0.2.1","host":"example.com","method":"GET","uri":"/pets","user_agent":"","status":200,"error":"","latency":9116,"latency_human":"9.116µs","bytes_in":0,"bytes_out":90}
{"time":"2021-07-10T20:24:20.125282+09:00","id":"","remote_ip":"192.0.2.1","host":"example.com","method":"GET","uri":"/pets?tags=TagOfFido","user_agent":"","status":200,"error":"","latency":5566,"latency_human":"5.566µs","bytes_in":0,"bytes_out":46}
{"time":"2021-07-10T20:24:20.125306+09:00","id":"","remote_ip":"192.0.2.1","host":"example.com","method":"GET","uri":"/pets?tags=NotExists","user_agent":"","status":200,"error":"","latency":6023,"latency_human":"6.023µs","bytes_in":0,"bytes_out":5}
{"time":"2021-07-10T20:24:20.12534+09:00","id":"","remote_ip":"192.0.2.1","host":"example.com","method":"DELETE","uri":"/pets/7","user_agent":"","status":404,"error":"","latency":9123,"latency_human":"9.123µs","bytes_in":0,"bytes_out":54}
{"time":"2021-07-10T20:24:20.12536+09:00","id":"","remote_ip":"192.0.2.1","host":"example.com","method":"DELETE","uri":"/pets/1000","user_agent":"","status":204,"error":"","latency":1173,"latency_human":"1.173µs","bytes_in":0,"bytes_out":0}
{"time":"2021-07-10T20:24:20.125374+09:00","id":"","remote_ip":"192.0.2.1","host":"example.com","method":"DELETE","uri":"/pets/1001","user_agent":"","status":204,"error":"","latency":797,"latency_human":"797ns","bytes_in":0,"bytes_out":0}
{"time":"2021-07-10T20:24:20.125389+09:00","id":"","remote_ip":"192.0.2.1","host":"example.com","method":"GET","uri":"/pets","user_agent":"","status":200,"error":"","latency":2276,"latency_human":"2.276µs","bytes_in":0,"bytes_out":5}
--- PASS: TestPetStore (0.00s)
PASS
ok  	github.com/leoh0/petstore/cmd	(cached)
```

또한 아래와 같이 실행 가능

``` sh
$ go run -v cmd/petstore.go
command-line-arguments

   ____    __
  / __/___/ /  ___
 / _// __/ _ \/ _ \
/___/\__/_//_/\___/ v4.2.1
High performance, minimalist Go web framework
https://echo.labstack.com
____________________________________O/_______
                                    O\
⇨ http server started on [::]:8080
```

### bazel 을 이용한 선언적인 형상 관리

#### 7. 기본적인 bazel 환경을 setup

우선 아래와 같이 4개의 파일을 추가함

``` sh
$ tree -a -I .git
.
├── .bazelignore # bazel 사용시 무시하고 사용할 파일들
├── .bazelversion # bazel 버전을 특정 지어서 사용
├── .gitignore
├── .tool-versions
├── BUILD # bazel에서 build 시 사용할 target을 위한 파일
├── WORKSPACE # bazel 의 사용을 위한 repository와 dependency등을 정의하는 파일
├── api
│   ├── client.gen.go
│   ├── petstore.go
│   ├── server.gen.go
│   ├── spec.gen.go
│   ├── spec.yaml
│   └── type.gen.go
├── cmd
│   ├── petstore.go
│   └── petstore_test.go
├── go.mod
└── go.sum

2 directories, 16 files
```

`.bazelignore` 과 `bazelversion` 의 내용은 아래와 같음

``` sh
$ cat << 'EOF' > .bazelignore
.git
EOF

$ cat << 'EOF' > .bazelversion
4.1.0
EOF
```

bazel이 `go`와 `gazelle(Bazel build file generator)` 를 사용하기 위한 선언

``` sh
$ cat << 'EOF' > WORKSPACE
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

# https://github.com/bazelbuild/rules_go#initial-project-setup
http_archive(
    name = "io_bazel_rules_go",
    sha256 = "69de5c704a05ff37862f7e0f5534d4f479418afc21806c887db544a316f3cb6b",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.27.0/rules_go-v0.27.0.tar.gz",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.27.0/rules_go-v0.27.0.tar.gz",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains(version = "1.16")

# https://github.com/bazelbuild/bazel-gazelle#running-gazelle-with-bazel
http_archive(
    name = "bazel_gazelle",
    sha256 = "62ca106be173579c0a167deb23358fdfe71ffa1e4cfdddf5582af26520f1c66f",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.23.0/bazel-gazelle-v0.23.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.23.0/bazel-gazelle-v0.23.0.tar.gz",
    ],
)

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()
EOF
```

아래와 같이 `gazelle` target을 사용하기 위해 prefix 를 go mod 의 값과 맞춰 놓음 `# gazelle:prefix github.com/leoh0/petstore`

``` sh
cat << 'EOF' > BUILD
load("@bazel_gazelle//:def.bzl", "gazelle")

# gazelle:prefix github.com/leoh0/petstore
gazelle(name = "gazelle")
EOF
```

gazelle이 정상적으로 셋업 되었다면 아래와 같이 `query`를 통해 `gazelle` target이 보임

``` sh
$ bazel query '...'
Starting local Bazel server and connecting to it...
//:gazelle
//:gazelle-runner
Loading: 1 packages loaded
```

`gazelle` 을 통해 go.mod 파일 기준으로 bazel go repository를 업데이트 함

``` sh
$ bazel run //:gazelle -- update-repos -from_file=go.mod -to_macro repos.bzl%go_repositories
INFO: SHA256 (https://golang.org/dl/?mode=json&include=all) = 6efc06a1bd0a710df5cbaa2fd314f9a3f702f7d9cd59ee2bd53c2a02aa8c4475
INFO: Analyzed target //:gazelle (69 packages loaded, 7310 targets configured).
INFO: Found 1 target...
Target //:gazelle up-to-date:
  bazel-bin/gazelle-runner.bash
  bazel-bin/gazelle
INFO: Elapsed time: 52.997s, Critical Path: 15.08s
INFO: 54 processes: 13 internal, 41 darwin-sandbox.
INFO: Build completed successfully, 54 total actions
INFO: Running command line: bazel-bin/gazelle update-repos '-from_file=go.mod' -to_macro repos.b
INFO: Build completed successfully, 54 total actions
```

이후에 `repos.bzl` 파일에 go repository 들이 정리된 것을 확인 할 수 있음

``` sh
$ tree -a -I ".git|bazel-*"
.
├── .bazelignore
├── .bazelversion
├── .gitignore
├── .tool-versions
├── BUILD
├── WORKSPACE
├── api
│   ├── client.gen.go
│   ├── petstore.go
│   ├── server.gen.go
│   ├── spec.gen.go
│   ├── spec.yaml
│   └── type.gen.go
├── cmd
│   ├── petstore.go
│   └── petstore_test.go
├── go.mod
├── go.sum
└── repos.bzl

2 directories, 17 files
```

#### 8. Makefile 추가와 gazelle을 이용한 기본 build 파일 생성

여태까지의 업데이트를 한번에 하기 위한 Makefile 추가

``` sh
cat << 'EOF' > Makefile
.PHONY: update
update:
	go get github.com/deepmap/oapi-codegen/pkg/codegen
	go generate -v ./...
	go mod tidy
	bazel run //:gazelle
	bazel run //:gazelle -- update-repos -from_file=go.mod -to_macro repos.bzl%go_repositories
EOF
```

이를 실행하여 업데이트

``` sh
$ make
go get github.com/deepmap/oapi-codegen/pkg/codegen
go generate -v ./...
api/client.gen.go
api/petstore.go
api/server.gen.go
api/spec.gen.go
api/type.gen.go
cmd/petstore.go
cmd/petstore_test.go
go mod tidy
warning: ignoring symlink /Users/al/github.com/leoh0/petstore/bazel-bin
warning: ignoring symlink /Users/al/github.com/leoh0/petstore/bazel-out
warning: ignoring symlink /Users/al/github.com/leoh0/petstore/bazel-petstore
warning: ignoring symlink /Users/al/github.com/leoh0/petstore/bazel-testlogs
bazel run //:gazelle
INFO: Analyzed target //:gazelle (12 packages loaded, 101 targets configured).
INFO: Found 1 target...
Target //:gazelle up-to-date:
  bazel-bin/gazelle-runner.bash
  bazel-bin/gazelle
INFO: Elapsed time: 2.281s, Critical Path: 0.01s
INFO: 1 process: 1 internal.
INFO: Build completed successfully, 1 total action
INFO: Build completed successfully, 1 total action
bazel run //:gazelle -- update-repos -from_file=go.mod -to_macro repos.bzl%go_repositories
INFO: Analyzed target //:gazelle (1 packages loaded, 2 targets configured).
INFO: Found 1 target...
Target //:gazelle up-to-date:
  bazel-bin/gazelle-runner.bash
  bazel-bin/gazelle
INFO: Elapsed time: 0.295s, Critical Path: 0.00s
INFO: 1 process: 1 internal.
INFO: Build completed successfully, 1 total action
INFO: Running command line: bazel-bin/gazelle update-repos '-from_file=go.mod' -to_macro repos.b
INFO: Build completed successfully, 1 total action
```

`gazelle` 을 통해서 build 파일이 생성됨을 확인

``` sh
$ cat api/BUILD.bazel
load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "api",
    srcs = [
        "client.gen.go",
        "petstore.go",
        "server.gen.go",
        "spec.gen.go",
        "type.gen.go",
    ],
    importpath = "github.com/leoh0/petstore/api",
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_deepmap_oapi_codegen//pkg/runtime",
        "@com_github_getkin_kin_openapi//openapi3",
        "@com_github_labstack_echo_v4//:echo",
    ],
)

$ cat cmd/BUILD.bazel
load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_test")

go_library(
    name = "cmd_lib",
    srcs = ["petstore.go"],
    importpath = "github.com/leoh0/petstore/cmd",
    visibility = ["//visibility:private"],
    deps = [
        "//api",
        "@com_github_deepmap_oapi_codegen//pkg/middleware",
        "@com_github_labstack_echo_v4//:echo",
        "@com_github_labstack_echo_v4//middleware",
    ],
)

go_binary(
    name = "cmd",
    embed = [":cmd_lib"],
    visibility = ["//visibility:public"],
)

go_test(
    name = "cmd_test",
    srcs = ["petstore_test.go"],
    embed = [":cmd_lib"],
    deps = [
        "//api",
        "@com_github_deepmap_oapi_codegen//pkg/middleware",
        "@com_github_deepmap_oapi_codegen//pkg/testutil",
        "@com_github_labstack_echo_v4//:echo",
        "@com_github_labstack_echo_v4//middleware",
        "@com_github_stretchr_testify//assert",
        "@com_github_stretchr_testify//require",
    ],
)
```

이후 bazel query를 통해 사용할 수 있는 `target`이 업데이트 된것을 확인 할 수 있음

``` sh
$ bazel query '...'
//cmd:cmd_test
//cmd:cmd
//cmd:cmd_lib
//api:api
//:gazelle
//:gazelle-runner
Loading: 0 packages loaded
```

이후 아래와 같이 테스트를 실행가능 만약 test에 관련된 부분중 변경된 부분이 없으면 아래처럼 `cached` 되어 테스트를 진행 하지 않을 수 있음

``` sh
$ bazel test //...
INFO: Analyzed 4 targets (0 packages loaded, 0 targets configured).
INFO: Found 3 targets and 1 test target...
INFO: Elapsed time: 0.158s, Critical Path: 0.00s
INFO: 1 process: 1 internal.
INFO: Build completed successfully, 1 total action
//cmd:cmd_test                                                  (cached) PASSED in 1.0s

Executed 0 out of 1 test: 1 test passes.
There were tests whose specified size is too big. Use the --test_verbose_timeout_warnings comman
INFO: Build completed successfully, 1 total action
```

그리고 아래와 같이 빌드된 바이너리를 바로 실행 할 수 있음

``` sh
$ bazel run //cmd:cmd
INFO: Analyzed target //cmd:cmd (0 packages loaded, 0 targets configured).
INFO: Found 1 target...
Target //cmd:cmd up-to-date:
  bazel-bin/cmd/cmd_/cmd
INFO: Elapsed time: 0.194s, Critical Path: 0.00s
INFO: 1 process: 1 internal.
INFO: Build completed successfully, 1 total action
INFO: Build completed successfully, 1 total action

   ____    __
  / __/___/ /  ___
 / _// __/ _ \/ _ \
/___/\__/_//_/\___/ v4.2.1
High performance, minimalist Go web framework
https://echo.labstack.com
____________________________________O/_______
                                    O\
⇨ http server started on [::]:8080
```

#### 9. 컨테이너 이미지 빌드용 타겟 추가

컨테이너 빌드를 위한 종속성을 `WORKSPACE` 에 추가함

``` sh
cat << 'EOF' >> WORKSPACE

http_archive(
    name = "io_bazel_rules_docker",
    sha256 = "59d5b42ac315e7eadffa944e86e90c2990110a1c8075f1cd145f487e999d22b3",
    strip_prefix = "rules_docker-0.17.0",
    urls = ["https://github.com/bazelbuild/rules_docker/releases/download/v0.17.0/rules_docker-v0.17.0.tar.gz"],
)

load(
    "@io_bazel_rules_docker//repositories:repositories.bzl",
    container_repositories = "repositories",
)

container_repositories()

load(
    "@io_bazel_rules_docker//go:image.bzl",
    _go_image_repos = "repositories",
)

_go_image_repos()
EOF
```

golang 기반의 이미지를 빌드하기 위한 target 추가

``` sh
cat << 'EOF' >> cmd/BUILD.bazel

load("@io_bazel_rules_docker//go:image.bzl", "go_image")

go_image(
    name = "image",
    embed = [":cmd_lib"],
    goarch = "amd64",
    goos = "linux",
    pure = "on",
)
EOF
```

이후 아래와 같이 이미지 빌드 및 실행 가능

``` sh
$ bazel run //cmd:image
DEBUG: /private/var/tmp/_bazel_al/5e5e74e2f5250506af20feb95a2222e2/external/bazel_gazelle/internal/go_repository.bzl:189:18: com_github_google_go_containerregistry: gazelle: finding module path for import github.com/vdemeester/k8s-pkg-credentialprovider: go: downloading github.com/vdemeester/k8s-pkg-credentialprovider v1.21.0
go get: github.com/vdemeester/k8s-pkg-credentialprovider@v1.21.0-1 updating to
	github.com/vdemeester/k8s-pkg-credentialprovider@v1.21.0 requires
	k8s.io/kubelet@v0.0.0: reading k8s.io/kubelet/go.mod at revision v0.0.0: unknown revision v0.0.0
gazelle: finding module path for import github.com/vdemeester/k8s-pkg-credentialprovider: go: downloading github.com/vdemeester/k8s-pkg-credentialprovider v1.21.0
go get: github.com/vdemeester/k8s-pkg-credentialprovider@v1.21.0-1 updating to
	github.com/vdemeester/k8s-pkg-credentialprovider@v1.21.0 requires
	k8s.io/kubelet@v0.0.0: reading k8s.io/kubelet/go.mod at revision v0.0.0: unknown revision v0.0.0
gazelle: finding module path for import github.com/vdemeester/k8s-pkg-credentialprovider: go: downloading github.com/vdemeester/k8s-pkg-credentialprovider v1.21.0
go get: github.com/vdemeester/k8s-pkg-credentialprovider@v1.21.0-1 updating to
	github.com/vdemeester/k8s-pkg-credentialprovider@v1.21.0 requires
	k8s.io/kubelet@v0.0.0: reading k8s.io/kubelet/go.mod at revision v0.0.0: unknown revision v0.0.0
gazelle: finding module path for import github.com/vdemeester/k8s-pkg-credentialprovider: go: downloading github.com/vdemeester/k8s-pkg-credentialprovider v1.21.0
go get: github.com/vdemeester/k8s-pkg-credentialprovider@v1.21.0-1 updating to
	github.com/vdemeester/k8s-pkg-credentialprovider@v1.21.0 requires
	k8s.io/kubelet@v0.0.0: reading k8s.io/kubelet/go.mod at revision v0.0.0: unknown revision v0.0.0
gazelle: finding module path for import github.com/vdemeester/k8s-pkg-credentialprovider: go: downloading github.com/vdemeester/k8s-pkg-credentialprovider v1.21.0
go get: github.com/vdemeester/k8s-pkg-credentialprovider@v1.21.0-1 updating to
	github.com/vdemeester/k8s-pkg-credentialprovider@v1.21.0 requires
	k8s.io/kubelet@v0.0.0: reading k8s.io/kubelet/go.mod at revision v0.0.0: unknown revision v0.0.0
INFO: Analyzed target //cmd:image (0 packages loaded, 0 targets configured).
INFO: Found 1 target...
Target //cmd:image up-to-date:
  bazel-bin/cmd/image-layer.tar
INFO: Elapsed time: 0.281s, Critical Path: 0.00s
INFO: 1 process: 1 internal.
INFO: Build completed successfully, 1 total action
INFO: Build completed successfully, 1 total action
6757c1b19b49: Loading layer  11.87MB/11.87MB
Loaded image ID: sha256:075cd58fb3a965374e9625e4eb3485f601f28ac45876b6294663bbb3572e22d1
Tagging 075cd58fb3a965374e9625e4eb3485f601f28ac45876b6294663bbb3572e22d1 as bazel/cmd:image

   ____    __
  / __/___/ /  ___
 / _// __/ _ \/ _ \
/___/\__/_//_/\___/ v4.2.1
High performance, minimalist Go web framework
https://echo.labstack.com
____________________________________O/_______
                                    O\
⇨ http server started on [::]:8080
```

`bazel run //cmd:image -- --norun` 과 같이 `--norun` 커맨드를 추가하면 실행하지 않고 이미지만 빌드 가능

#### 10. k8s 사용 준비

k8s 사용을 위한 종속성을 `WORKSPACE`에 추가

``` sh
cat << 'EOF' >> WORKSPACE

http_archive(
    name = "io_bazel_rules_k8s",
    strip_prefix = "rules_k8s-0.6",
    urls = ["https://github.com/bazelbuild/rules_k8s/archive/v0.6.tar.gz"],
    sha256 = "51f0977294699cd547e139ceff2396c32588575588678d2054da167691a227ef",
)

load("@io_bazel_rules_k8s//k8s:k8s.bzl", "k8s_repositories")

k8s_repositories()

load("@io_bazel_rules_k8s//k8s:k8s_go_deps.bzl", k8s_go_deps = "deps")

k8s_go_deps()
EOF
```

기본적으로 사용할 deployment를 생성하고 k8s yaml을 생성하기 위한 target을 추가

``` sh
mkdir deployments

cat << 'EOF' > deployments/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: petstore
  labels:
    app: petstore
spec:
  selector:
    matchLabels:
      app: petstore
  template:
    metadata:
      labels:
        app: petstore
    spec:
      containers:
      - name: petstore
        image: leoh0/petstore-image
        ports:
        - containerPort: 8080
EOF

cat << 'EOF' > deployments/BUILD.bazel
load("@io_bazel_rules_k8s//k8s:object.bzl", "k8s_object")

k8s_object(
    name = "k8s",
    kind = "deployment",
    template = ":deployment.yaml",
    images = {
        "leoh0/petstore-image": "//cmd:image",
    },
)
EOF
```

그리고 아래와 같이 기존에 생성되던 이미지를 deployments 디렉토리에서 사용할 수 있도록 public으로 visibility 설정

``` sh
$ git diff cmd/BUILD.bazel
diff --git a/cmd/BUILD.bazel b/cmd/BUILD.bazel
index d1d5dcc..f7b5322 100644
--- a/cmd/BUILD.bazel
+++ b/cmd/BUILD.bazel
@@ -42,4 +42,5 @@ go_image(
     goarch = "amd64",
     goos = "linux",
     pure = "on",
+    visibility = ["//visibility:public"],
 )
```

잘 설정이 되었다면 아래와 같이 bazel로 build 해서 push한 이미지의 sha 값으로 yaml이 업데이트 됨을 알 수 있음

``` sh
$ bazel run //deployments:k8s
INFO: Analyzed target //deployments:k8s (0 packages loaded, 0 targets configured).
INFO: Found 1 target...
Target //deployments:k8s up-to-date:
  bazel-bin/deployments/k8s
INFO: Elapsed time: 0.257s, Critical Path: 0.00s
INFO: 1 process: 1 internal.
INFO: Build completed successfully, 1 total action
INFO: Build completed successfully, 1 total action
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: petstore
  name: petstore
spec:
  selector:
    matchLabels:
      app: petstore
  template:
    metadata:
      labels:
        app: petstore
    spec:
      containers:
      - image: index.docker.io/leoh0/petstore-image@sha256:02366c5a2d3c97fc5de2f382dae8e883eefc405cdedb5b24e8a4c641c6ac0638
        name: petstore
        ports:
        - containerPort: 8080
```

이 결과값을 그대로 apply 하는식으로도 적용 가능하고 cluster 정보등을 k8s target에 넣어서 bazel 만으로도 배포할 수 있음

``` sh
$ bazel run //deployments:k8s | kubectl apply -f -
```

#### 11. human readable image bundle 추가

bazel을 이용하면 deterministic 하게 이미지의 sha 값으로 배포하게 되지만 사람이 찾기 어려운 sha 값만이 존재하여 이를 인간이 구분하기 위한 tag 값을 부여하는게 편할 수 있음

우선 bazel에서 사용할 자동생성시킬 tag 값을 아래와 같은 스크립트 형태로 제작

``` sh
$ mkdir hack

$ cat << 'EOD' > hack/print-workspace-status.sh
#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

git_commit="$(git describe --tags --always --dirty)"
build_date="$(date -u '+%Y%m%d')"
docker_tag="v${build_date}-${git_commit}"

cat <<EOF
DOCKER_TAG ${docker_tag}
EOF
EOD

$ chmod +x hack/print-workspace-status.sh
```

그리고 해당 스크립트를 `.bazelrc` 를 통하여 커맨드마다 실행 할 수 있도록 설정

``` sh
cat << 'EOF' > .bazelrc
run --workspace_status_command=./hack/print-workspace-status.sh
build --workspace_status_command=./hack/print-workspace-status.sh
EOF
```

그리고 아래와 같은 target 을 추가해서 이미지 이름,tag 변경 및 push를 만약 원한다면 할 수 있도록 함

``` sh
cat << 'EOF' > cmd/BUILD.bazel
load("@io_bazel_rules_docker//container:bundle.bzl", "container_bundle")
load("@io_bazel_rules_docker//contrib:push-all.bzl", "container_push")

container_bundle(
    name = "bundle",
    images = {
        "docker.io/leoh0/petstore" + ":{DOCKER_TAG}": ":image",
        "docker.io/leoh0/petstore" + ":latest": ":image",
    }
)

container_push(
    name = "push",
    format = "Docker",
    bundle = ":bundle",
)
EOF
```

이후 아래와 같은 커맨드로 이미지를 원하는 이름으로도 build & push 가능

``` sh
$ bazel run //cmd:push
INFO: Analyzed target //cmd:push (1 packages loaded, 6 targets configured).
INFO: Found 1 target...
Target //cmd:push up-to-date:
  bazel-bin/cmd/push
INFO: Elapsed time: 0.346s, Critical Path: 0.03s
INFO: 3 processes: 3 internal.
INFO: Build completed successfully, 3 total actions
INFO: Build completed successfully, 3 total actions
2021/07/10 23:09:57 Destination docker.io/leoh0/petstore:{DOCKER_TAG} was resolved to docker.io/leoh0/petstore:v20210710-5e3adab-dirty after stamping.
2021/07/10 23:10:03 Successfully pushed Docker image to docker.io/leoh0/petstore:v20210710-5e3adab-dirty
2021/07/10 23:10:04 Successfully pushed Docker image to docker.io/leoh0/petstore:latest
```

### tilt 를 이용한 간편한 개발환경 셋업

우선 tilt 사용을 위한 kind 와 registry를 생성

``` sh
$ ctlptl create cluster kind --registry=ctlptl-registry
Creating cluster "kind" ...
 ✓ Ensuring node image (kindest/node:v1.21.1) 🖼
 ✓ Preparing nodes 📦
 ✓ Writing configuration 📜
 ✓ Starting control-plane 🕹️
 ✓ Installing CNI 🔌
 ✓ Installing StorageClass 💾
Set kubectl context to "kind-kind"
You can now use your cluster with:

kubectl cluster-info --context kind-kind

Have a question, bug, or feature request? Let us know! https://kind.sigs.k8s.io/#community 🙂
Switched to context "kind-kind".
 🔌 Connected cluster kind-kind to registry ctlptl-registry at localhost:53217
 👐 Push images to the cluster like 'docker push localhost:53217/alpine'
cluster.ctlptl.dev/kind-kind created
```

#### 12. Tiltfile 추가 및 k8s-dev target 추가

tilt 가 사용할 dev 용 target을 추가. 이미지를 tilt가 바꿔서 적용 할 수 있도록 함

``` sh
cat << 'EOF' >> deployments/BUILD.bazel

k8s_object(
    name = "k8s-dev",
    kind = "deployment",
    template = ":deployment.yaml",
)
EOF
```

tilt 사용을 위한 `Tiltfile` 추가

`bazel run //deployments:k8s-dev` 기반의 yaml에서 `bazel run {image_target} -- --norun` 로 빌드되는 이미지를 스트링 변환으로 적용하여 실행하게 됨

``` sh
cat << 'EOF' > Tiltfile
# -*- mode: Python -*-

# Use Bazel to generate the Kubernetes YAML
watch_file('./deployments/deployment.yaml')
k8s_yaml(local('bazel run //deployments:k8s-dev'))

# Use Bazel to build the image

# The go_image BUILD rule
image_target='//cmd:image'

# Where go_image puts the image in Docker (bazel/path/to/target:name)
bazel_image='bazel/cmd:image'

custom_build(
  ref='leoh0/petstore-image',
  command=(
    'bazel run {image_target} -- --norun && ' +
    'docker tag {bazel_image} $EXPECTED_REF').format(image_target=image_target, bazel_image=bazel_image),
  deps=['./api', './cmd', './deployments'],
)

k8s_resource('petstore', port_forwards=8080)
EOF
```

이후 tilt 를 실행해서 실행

``` sh
$ tilt up
Tilt started on http://localhost:10350/
v0.21.2, built 2021-07-06

(space) to open the browser
(s) to stream logs (--stream=true)
(t) to open legacy terminal mode (--legacy=true)
(ctrl-c) to exit
Tilt started on http://localhost:10350/
v0.21.2, built 2021-07-06

Initial Build • (Tiltfile)
Beginning Tiltfile execution
local: bazel run //deployments:k8s-dev
creating uisession: uisessions.tilt.dev "Tiltfile" already exists
 → Loading:
 → Loading: 0 packages loaded
 → Analyzing: target //deployments:k8s-dev (0 packages loaded, 0 targets configured)
 → INFO: Analyzed target //deployments:k8s-dev (0 packages loaded, 0 targets configured).
 → INFO: Found 1 target...
 → [0 / 2] [Prepa] BazelWorkspaceStatusAction stable-status.txt
 → Target //deployments:k8s-dev up-to-date:
 →   bazel-bin/deployments/k8s-dev
 → INFO: Elapsed time: 0.274s, Critical Path: 0.00s
 → INFO: 1 process: 1 internal.
 → INFO: Build completed successfully, 1 total action
 → INFO: Running command line: bazel-bin/deployments/k8s-dev
 → INFO: Build completed successfully, 1 total action
 → apiVersion: apps/v1
 → kind: Deployment
 → metadata:
 →   labels:
 →     app: petstore
 →   name: petstore
 → spec:
 →   selector:
 →     matchLabels:
 →       app: petstore
 →   template:
 →     metadata:
 →       labels:
 →         app: petstore
 →     spec:
 →       containers:
 →       - image: leoh0/petstore-image
 →         name: petstore
 →         ports:
 →         - containerPort: 8080
 →
Auto-detected local registry from environment: {localhost:53217  }
Successfully loaded Tiltfile (340.615117ms)
     petstore │
     petstore │ Initial Build • petstore
     petstore │ STEP 1/3 — Building Custom Build: [leoh0/petstore-image]
     petstore │ Custom Build: Injecting Environment Variables
     petstore │   EXPECTED_REF=localhost:53217/leoh0_petstore-image:tilt-build-1625925285
     petstore │ Running custom build cmd "bazel run //cmd:image -- --norun && docker tag bazel/cmd:image $EXPECTED_REF"
     petstore │ Loading:
     petstore │ Loading: 0 packages loaded
     petstore │ Analyzing: target //cmd:image (0 packages loaded, 0 targets configured)
     petstore │ INFO: Analyzed target //cmd:image (0 packages loaded, 0 targets configured).
     petstore │ INFO: Found 1 target...
     petstore │ [0 / 1] [Prepa] BazelWorkspaceStatusAction stable-status.txt
     petstore │ Target //cmd:image up-to-date:
     petstore │   bazel-bin/cmd/image-layer.tar
     petstore │ INFO: Elapsed time: 0.250s, Critical Path: 0.00s
     petstore │ INFO: 1 process: 1 internal.
     petstore │ INFO: Build completed successfully, 1 total action
     petstore │ INFO: Running command line: bazel-bin/cmd/image.executable --norun
     petstore │ INFO: Build completed successfully, 1 total action
     petstore │ Loaded image ID: sha256:075cd58fb3a965374e9625e4eb3485f601f28ac45876b6294663bbb3572e22d1
     petstore │ Tagging 075cd58fb3a965374e9625e4eb3485f601f28ac45876b6294663bbb3572e22d1 as bazel/cmd:image
     petstore │
     petstore │ STEP 2/3 — Pushing localhost:53217/leoh0_petstore-image:tilt-075cd58fb3a96537
     petstore │      Pushing with Docker client
     petstore │      Authenticating to image repo: localhost:53217
     petstore │      Sending image data
     petstore │      417cb9b79ade: Layer already exists
     petstore │      6757c1b19b49: Pushing  135.2kB/11.86MB
     petstore │      6757c1b19b49: Pushed
     petstore │
     petstore │ STEP 3/3 — Deploying
     petstore │      Injecting images into Kubernetes YAML
     petstore │      Applying via kubectl:
     petstore │      → petstore:deployment
     petstore │
     petstore │      Step 1 - 1.45s (Building Custom Build: [leoh0/petstore-image])
     petstore │      Step 2 - 0.83s (Pushing localhost:53217/leoh0_petstore-image:tilt-075cd58fb3a96537)
     petstore │      Step 3 - 0.04s (Deploying)
     petstore │      DONE IN: 2.32s
     petstore │
     petstore │
     petstore │ Tracking new pod rollout (petstore-6777d68c7f-tjsps):
     petstore │      ┊ Scheduled       - <1s
     petstore │      ┊ Initialized     - (…) Pending
     petstore │      ┊ Ready           - (…) Pending
     petstore │ [K8s EVENT: Pod petstore-6777d68c7f-tjsps (ns: default)] Pulling image "localhost:53217/leoh0_petstore-image:tilt-075cd58fb3a96537"
     petstore │ [K8s EVENT: Pod petstore-6777d68c7f-tjsps (ns: default)] Successfully pulled image "localhost:53217/leoh0_petstore-image:tilt-075cd58fb3a96537" in 202.1675ms
     petstore │
     petstore │    ____    __
     petstore │   / __/___/ /  ___
     petstore │  / _// __/ _ \/ _ \
     petstore │ /___/\__/_//_/\___/ v4.2.1
     petstore │ High performance, minimalist Go web framework
     petstore │ https://echo.labstack.com
     petstore │ ____________________________________O/_______
     petstore │                                     O\
     petstore │ ⇨ http server started on [::]:8080
```

이후 파일들을 수정시 tilt 가 잘 업데이트 하는지를 확인

## 참고

[https://github.com/leoh0/petstore](https://github.com/leoh0/petstore) 레포의 git commit 단위로 위 설명과 일치하게 작업해 둔 결과물로 참고 가능
