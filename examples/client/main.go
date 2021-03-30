package main

import (
	"context"
	"crypto/x509"
	"log"
	"strings"

	// the generated zitadel files for management api
	pb "github.com/caos/zitadel/examples/client/zitadel/management"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

//zitadelAPI is the default zitadel api
const zitadelAPI = "api.zitadel.ch:443"

func main() {
	conn, err := grpc.Dial(zitadelAPI, grpc.WithTransportCredentials(cert()))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewManagementServiceClient(conn)

	//call ZITADEL. the response has no payload so we ignore the res
	// the call was successful if no error responded
	_, err = client.Healthz(context.TODO(), &empty.Empty{})
	if err != nil {
		log.Fatalln("call failed: ", err)
	}
	log.Println("call was successful")
}

//cert load default cert pool for tls
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
