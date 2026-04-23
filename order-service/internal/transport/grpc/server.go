package grpc

import (
    "context"
    "net"

    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
    
    "dulatAsisADV2/order-service/internal/domain"
    "dulatAsisADV2/order-service/internal/usecase"
    orderPb "dulatAsisADV2/proto/gen/order"
)

type OrderServer struct {
    orderPb.UnimplementedOrderServiceServer
    orderUseCase *usecase.OrderUseCase
}

func NewOrderGRPCServer(orderUseCase *usecase.OrderUseCase) *OrderServer {
    return &OrderServer{
        orderUseCase: orderUseCase,
    }
}

func (s *OrderServer) CreateOrder(ctx context.Context, req *orderPb.CreateOrderRequest) (*orderPb.CreateOrderResponse, error) {
    // Convert proto items to domain items
    var items []*domain.OrderItem
    for _, item := range req.Items {
        items = append(items, &domain.OrderItem{
            ProductID: item.ProductId,
            Quantity:  item.Quantity,
            Price:     item.Price,
        })
    }
    
    order, err := s.orderUseCase.CreateOrder(ctx, req.UserId, items, req.TotalAmount)
    if err != nil {
        return nil, err
    }
    
    return &orderPb.CreateOrderResponse{
        OrderId: order.ID,
        Status:  order.Status,
    }, nil
}

func (s *OrderServer) GetOrder(ctx context.Context, req *orderPb.GetOrderRequest) (*orderPb.Order, error) {
    order, err := s.orderUseCase.GetOrder(ctx, req.OrderId)
    if err != nil {
        return nil, err
    }
    
    // Convert domain items to proto items
    var protoItems []*orderPb.OrderItem
    for _, item := range order.Items {
        protoItems = append(protoItems, &orderPb.OrderItem{
            ProductId: item.ProductID,
            Quantity:  item.Quantity,
            Price:     item.Price,
        })
    }
    
    return &orderPb.Order{
        Id:          order.ID,
        UserId:      order.UserID,
        Items:       protoItems,
        TotalAmount: order.TotalAmount,
        Status:      order.Status,
        CreatedAt:   order.CreatedAt.String(),
    }, nil
}

func (s *OrderServer) UpdateOrderStatus(ctx context.Context, req *orderPb.UpdateOrderStatusRequest) (*orderPb.UpdateOrderStatusResponse, error) {
    err := s.orderUseCase.UpdateOrderStatus(ctx, req.OrderId, req.Status)
    if err != nil {
        return &orderPb.UpdateOrderStatusResponse{
            Success: false,
            Message: err.Error(),
        }, nil
    }
    
    return &orderPb.UpdateOrderStatusResponse{
        Success: true,
        Message: "Order status updated",
    }, nil
}

func StartGRPCServer(port string, orderUseCase *usecase.OrderUseCase) error {
    lis, err := net.Listen("tcp", ":"+port)
    if err != nil {
        return err
    }
    
    s := grpc.NewServer()
    orderServer := NewOrderGRPCServer(orderUseCase)
    orderPb.RegisterOrderServiceServer(s, orderServer)
    reflection.Register(s)
    
    return s.Serve(lis)
}
