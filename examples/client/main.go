package main

import (
	"context"
	"crypto/x509"
	"log"
	"strings"

	pb "github.com/caos/zitadel/examples/client/zitadel/management"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const zitadelAPI = "api.zitadel.ch:443"

func main() {
	conn, err := grpc.Dial(zitadelAPI, grpc.WithTransportCredentials(cert()))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewManagementServiceClient(conn)
	_, err = client.Healthz(context.TODO(), &empty.Empty{})
	if err != nil {
		log.Fatalln("call failed: ", err)
	}
	log.Println("call was successful")
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
