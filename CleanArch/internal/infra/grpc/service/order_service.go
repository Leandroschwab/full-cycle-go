package service

import (
	"context"

	"github.com/devfullcycle/20-CleanArch/internal/infra/grpc/pb"
	"github.com/devfullcycle/20-CleanArch/internal/usecase"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase usecase.CreateOrderUseCase
	ListAllOrdersUseCase usecase.ListOrdersUseCase
}

func NewOrderService(createOrderUseCase usecase.CreateOrderUseCase, listAllOrdersUseCase usecase.ListOrdersUseCase) *OrderService {
	return &OrderService{
		CreateOrderUseCase: createOrderUseCase,
		ListAllOrdersUseCase: listAllOrdersUseCase,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	dto := usecase.OrderInputDTO{
		ID:    in.Id,
		Price: float64(in.Price),
		Tax:   float64(in.Tax),
	}
	output, err := s.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrderResponse{
		Id:         output.ID,
		Price:      float32(output.Price),
		Tax:        float32(output.Tax),
		FinalPrice: float32(output.FinalPrice),
	}, nil
}

func (s *OrderService) ListAll(ctx context.Context, in *pb.ListOrderRequest) (*pb.ListOrderResponse, error) {
	orders, err := s.ListAllOrdersUseCase.Execute()
	if err != nil {
		return nil, err
	}
	var pbOrders []*pb.Order
	for _, order := range orders {
		pbOrders = append(pbOrders, &pb.Order{
			Id:         order.ID,
			Price:      float32(order.Price),
			Tax:        float32(order.Tax),
			FinalPrice: float32(order.FinalPrice),
		})
	}	
	return &pb.ListOrderResponse{Orders: pbOrders}, nil	
}