package main

import (
	"context"
	"crypto/x509"
	"log"
	"strings"

	pb "github.com/caos/zitadel/examples/client/zitadel/admin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const zitadelAPI = "api.zitadel.ch:433"

func main() {
	conn, err := grpc.Dial(zitadelAPI, grpc.WithTransportCredentials(cert()))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewAdminServiceClient(conn)
	res, err := client.Healthz(context.TODO(), nil)
	log.Println(res, err)
}

func cert() credentials.TransportCredentials {
	ca, err := x509.SystemCertPool()
	if err != nil {
		log.Println("unable to load cert pool")
	}
	if ca == nil {
		ca = x509.NewCertPool()
	}

	servernameWithoutPort := strings.Split(zitadelAPI, ":")[0]
	return credentials.NewClientTLSFromCert(ca, servernameWithoutPort)
}
