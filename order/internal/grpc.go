package internal

import (
	"context"
	"log"

	mapset "github.com/deckarep/golang-set/v2"
	account "github.com/kushaljangra/e-commerce/account/client"
	"github.com/kushaljangra/e-commerce/order/models"
	"github.com/kushaljangra/e-commerce/order/proto/pb"
	product "github.com/kushaljangra/e-commerce/product/client"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type grpcServer struct {
	pb.UnimplementedOrderServiceServer
	service       Service
	accountClient *account.Client
	productClient *product.Client
}

func (server *grpcServer) PostOrder(ctx context.Context, request *pb.PostOrderRequest) (*pb.PostOrderResponse, error) {
	_, err := server.accountClient.GetAccount(ctx, request.AccountId)
	if err != nil {
		log.Println("Error getting account", err)
		return nil, err
	}
	var productIDs []string
	for _, p := range request.Products {
		productIDs = append(productIDs, p.Id)
	}
	orderedProducts, err := server.productClient.GetProducts(ctx, 0, 0, productIDs, "")
	if err != nil {
		log.Println("Error getting ordered products", err)
		return nil, err
	}

	var products []*models.OrderedProduct
	totalPrice := 0.0

	for _, p := range orderedProducts {
		productObj := &models.OrderedProduct{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    0,
		}
		for _, requestProduct := range request.Products {
			if requestProduct.Id == p.ID {
				productObj.Quantity = requestProduct.Quantity
				break
			}
		}

		if productObj.Quantity != 0 {
			products = append(products, productObj)
			totalPrice += productObj.Price * float64(productObj.Quantity)
		}
	}

	postOrder, err := server.service.PostOrder(ctx, request.AccountId, totalPrice, products)
	if err != nil {
		log.Println("Error posting postOrder", err)
		return nil, err
	}

	orderProto := &pb.Order{
		Id:         uint64(postOrder.ID),
		AccountId:  postOrder.AccountID,
		TotalPrice: postOrder.TotalPrice,
		Products:   []*pb.ProductInfo{},
	}
	orderProto.CreatedAt, _ = postOrder.CreatedAt.MarshalBinary()
	for _, p := range postOrder.Products {
		orderProto.Products = append(orderProto.Products, &pb.ProductInfo{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    p.Quantity,
		})
	}

	return &pb.PostOrderResponse{
		Order: orderProto,
	}, nil
}

func (server *grpcServer) GetOrdersForAccount(ctx context.Context, request *wrapperspb.UInt64Value) (*pb.GetOrdersForAccountResponse, error) {
	accountOrders, err := server.service.GetOrdersForAccount(ctx, request.Value)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	productIDsSet := mapset.NewSet[string]()
	for _, o := range accountOrders {
		for _, p := range o.Products {
			productIDsSet.Add(p.ID)
		}
	}

	productIDs := productIDsSet.ToSlice()

	products, err := server.productClient.GetProducts(ctx, 0, 0, productIDs, "")
	if err != nil {
		log.Println("Error getting account products: ", err)
		return nil, err
	}

	var orders []*pb.Order
	for _, order := range accountOrders {
		encodedOrder := &pb.Order{
			AccountId:  order.AccountID,
			Id:         uint64(order.ID),
			TotalPrice: order.TotalPrice,
			Products:   []*pb.ProductInfo{},
		}
		encodedOrder.CreatedAt, _ = order.CreatedAt.MarshalBinary()

		for _, orderedProduct := range order.Products {
			for _, prod := range products {
				if prod.ID == orderedProduct.ID {
					orderedProduct.Name = prod.Name
					orderedProduct.Description = prod.Description
					orderedProduct.Price = prod.Price
					break
				}
			}

			encodedOrder.Products = append(encodedOrder.Products, &pb.ProductInfo{
				Id:          orderedProduct.ID,
				Name:        orderedProduct.Name,
				Description: orderedProduct.Description,
				Price:       orderedProduct.Price,
				Quantity:    orderedProduct.Quantity,
			})
		}

		orders = append(orders, encodedOrder)
	}
	return &pb.GetOrdersForAccountResponse{Orders: orders}, nil
}

func (server *grpcServer) UpdateOrderStatus(ctx context.Context, request *pb.UpdateOrderStatusRequest) (*emptypb.Empty, error) {
	err := server.service.UpdateOrderPaymentStatus(ctx, request.OrderId, request.Status)

	if err != nil {
		log.Println("Error updating order status", err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
