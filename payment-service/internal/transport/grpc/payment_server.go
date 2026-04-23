package grpc

import (
    "context"
    "log"
    "net"

    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
    
    "dulatAsisADV2/payment-service/internal/usecase"
    pb "dulatAsisADV2/proto/gen/payment"
    
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

type PaymentGRPCServer struct {
    pb.UnimplementedPaymentServiceServer
    paymentUseCase *usecase.PaymentUseCase
}

func NewPaymentGRPCServer(paymentUseCase *usecase.PaymentUseCase) *PaymentGRPCServer {
    return &PaymentGRPCServer{
        paymentUseCase: paymentUseCase,
    }
}

func (s *PaymentGRPCServer) ProcessPayment(ctx context.Context, req *pb.PaymentRequest) (*pb.PaymentResponse, error) {
    log.Printf("Received payment request: OrderID=%s, Amount=%f", req.OrderId, req.Amount)

    payment, err := s.paymentUseCase.ProcessPayment(ctx, req.OrderId, req.Amount)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to process payment: %v", err)
    }

    return &pb.PaymentResponse{
        PaymentId:     payment.ID,
        OrderId:       payment.OrderID,
        TransactionId: payment.TransactionID,
        Status:        payment.Status,
        CreatedAt:     payment.CreatedAt.String(),
    }, nil
}

func (s *PaymentGRPCServer) GetPaymentStatus(ctx context.Context, req *pb.PaymentStatusRequest) (*pb.PaymentStatusResponse, error) {
    payment, err := s.paymentUseCase.GetPaymentStatus(ctx, req.PaymentId)
    if err != nil {
        return nil, status.Errorf(codes.NotFound, "payment not found: %v", err)
    }
    
    return &pb.PaymentStatusResponse{
        PaymentId: payment.ID,
        Status:    payment.Status,
        Amount:    payment.Amount,
    }, nil
}

func (s *PaymentGRPCServer) StreamPayments(req *pb.StreamPaymentsRequest, stream pb.PaymentService_StreamPaymentsServer) error {
    log.Printf("StreamPayments called for user: %s", req.UserId)
    return nil
}

func (s *PaymentGRPCServer) ListPayments(ctx context.Context, req *pb.ListPaymentsRequest) (*pb.ListPaymentsResponse, error) {
    log.Printf("ListPayments called with status filter: %s", req.Status)
    
    payments, err := s.paymentUseCase.ListByStatus(ctx, req.Status)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to list payments: %v", err)
    }
    
    var pbPayments []*pb.PaymentResponse
    for _, payment := range payments {
        pbPayments = append(pbPayments, &pb.PaymentResponse{
            PaymentId:     payment.ID,
            OrderId:       payment.OrderID,
            TransactionId: payment.TransactionID,
            Status:        payment.Status,
            CreatedAt:     payment.CreatedAt.String(),
        })
    }
    
    return &pb.ListPaymentsResponse{
        Payments: pbPayments,
    }, nil
}

func StartGRPCServer(port string, paymentUseCase *usecase.PaymentUseCase) error {
    lis, err := net.Listen("tcp", ":"+port)
    if err != nil {
        return err
    }
    
    s := grpc.NewServer()
    paymentServer := NewPaymentGRPCServer(paymentUseCase)
    pb.RegisterPaymentServiceServer(s, paymentServer)
    reflection.Register(s) // Включаем reflection API
    
    log.Printf("Payment Service gRPC server listening on port %s", port)
    return s.Serve(lis)
}
