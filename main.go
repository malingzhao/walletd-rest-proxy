package main

import (
	"flag"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/theQRL/walletd-rest-proxy/generated"
	"github.com/theQRL/walletd-rest-proxy/qrlwallet"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net/http"
)

func run(walletServiceEndPoint string, serverIPPort string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := qrlwallet.RegisterWalletAPIHandlerFromEndpoint(ctx, mux, walletServiceEndPoint, opts)
	if err != nil {
		return err
	}
	err = generated.RegisterPublicAPIHandlerFromEndpoint(ctx, mux, walletServiceEndPoint, opts)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}

	return http.ListenAndServe(serverIPPort, mux)
}

func main() {

	walletServiceEndPoint := flag.String("walletServiceEndpoint",
		"127.0.0.1:19010",
		"endpoint of WalletAPIService")

	serverIPPort := flag.String("serverIPPort",
		"127.0.0.1:5359",
		"IP and Port at which REST proxy will be listening")

	flag.Parse()
	defer glog.Flush()

	if err := run(*walletServiceEndPoint, *serverIPPort); err != nil {
		glog.Fatal(err)
	}
}
