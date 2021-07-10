package client

import (
	"context"

	"github.com/rkrmr33/onka/pkg/proto/v1alpha1"
	"github.com/rkrmr33/onka/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	defaultDaemonAddr = "localhost:6543"
)

type ClientOptions struct {
	Addr      string
	TLSCert   string
	TLSKey    string
	TLSCA     string
	TLSVerify bool
}

func AddFlags(cmd *cobra.Command) ClientOptions {
	o := ClientOptions{}

	fs := &pflag.FlagSet{}
	fs.StringVar(&o.Addr, "host", defaultDaemonAddr, "Address of onkad instance to talk to")
	fs.StringVar(&o.TLSCert, "tlscert", "", "Path to TLS certificate file")
	fs.StringVar(&o.TLSKey, "tlskey", "", "Path to TLS key file")
	fs.StringVar(&o.TLSCA, "tlscacert", "", "Path to CA certificate file")
	fs.BoolVar(&o.TLSVerify, "tlsverify", true, "Use TLS and verify the remote certificate")

	util.Must(viper.BindPFlags(fs))

	bindEnv()

	return o
}

func FromEnv() ClientOptions {
	o := ClientOptions{}
	bindEnv()
	return o
}

func New(ctx context.Context) (v1alpha1.DaemonServiceClient, error) {
	var (
		host    = viper.GetString("host")
		tlsCert = viper.GetString("tlscert")
		tlsKey  = viper.GetString("tlskey")
		//tlsCACert   = viper.GetString("tlscacert")
		tlsCertPath = viper.GetString("tlscertpath")
		tlsVerify   = viper.GetBool("tlsverify")
	)

	return newClient(ctx, host, tlsCert, tlsKey, tlsCertPath, tlsVerify)
}

func (o ClientOptions) New(ctx context.Context) (v1alpha1.DaemonServiceClient, error) {
	return newClient(
		ctx,
		o.Addr,
		o.TLSCert,
		o.TLSKey,
		"",
		o.TLSVerify,
	)
}

func newClient(ctx context.Context, host, tlsCert, tlsKey, tlsCertPath string, tlsVerify bool) (v1alpha1.DaemonServiceClient, error) {
	connOptions := []grpc.DialOption{}

	if !tlsVerify {
		connOptions = append(connOptions, grpc.WithInsecure())
	} else {
		if tlsCertPath != "" {
			if tlsCert == "" {
				tlsCert = tlsCertPath + "cert.pem"
			}
			if tlsKey == "" {
				tlsKey = tlsCertPath + "key.pem"
			}
		}

		if tlsCert != "" && tlsKey != "" {
			certBund, err := credentials.NewServerTLSFromFile(tlsCert, tlsKey)
			if err != nil {
				return nil, err
			}

			connOptions = append(connOptions,
				grpc.WithTransportCredentials(certBund),
			)
		} else {
			connOptions = append(connOptions, grpc.WithInsecure())
		}
	}

	grpcCon, err := grpc.DialContext(
		ctx,
		host,
		connOptions...,
	)
	if err != nil {
		return nil, err
	}

	return v1alpha1.NewDaemonServiceClient(grpcCon), nil
}

func bindEnv() {
	util.Must(viper.BindEnv("host", "ONKA_HOST"))
	util.Must(viper.BindEnv("tlsverify", "ONKA_TLS_VERIFY"))
	util.Must(viper.BindEnv("tlscertpath", "ONKA_CERT_PATH"))

	viper.SetDefault("host", defaultDaemonAddr)
}
