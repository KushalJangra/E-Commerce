package internal

import (
	"fmt"
	"net"

	account "github.com/kushaljangra/e-commerce/account/client"
	"github.com/kushaljangra/e-commerce/order/proto/pb"
	product "github.com/kushaljangra/e-commerce/product/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func ListenGRPC(service Service, accountURL string, productURL string, port int) error {
	accountClient, err := account.NewClient(accountURL)
	if err != nil {
		return err
	}

	productClient, err := product.NewClient(productURL)
	if err != nil {
		accountClient.Close()
		return err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		accountClient.Close()
		productClient.Close()
		return err
	}

	serv := grpc.NewServer()
	pb.RegisterOrderServiceServer(serv, &grpcServer{
		pb.UnimplementedOrderServiceServer{},
		service,
		accountClient,
		productClient,
	})
	reflection.Register(serv)

	return serv.Serve(lis)
}
