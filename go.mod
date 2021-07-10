module github.com/rkrmr33/onka

go 1.16

require (
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/Masterminds/sprig v2.22.0+incompatible
	github.com/banzaicloud/logrus-runtime-formatter v0.0.0-20190729070250-5ae5475bae5e
	github.com/containerd/containerd v1.5.0 // indirect
	github.com/docker/docker v20.10.6+incompatible
	github.com/docker/docker/pkg/jsonlog v0.0.0-00010101000000-000000000000 // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/machine v0.16.2
	github.com/fatih/color v1.12.0
	github.com/ghodss/yaml v1.0.0
	github.com/gogo/protobuf v1.3.2
	github.com/golang/protobuf v1.4.3
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.1-0.20190118093823-f849b5445de4
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/intel-go/cpuid v0.0.0-20200819041909-2aa72927c3e2 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/moby/moby v20.10.6+incompatible
	github.com/moby/term v0.0.0-20201216013528-df9cb8a40635 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/samalba/dockerclient v0.0.0-20160531175551-a30362618471 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	google.golang.org/grpc v1.37.1
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c
)

replace github.com/docker/docker/pkg/jsonlog => ./assets/jsonlog
