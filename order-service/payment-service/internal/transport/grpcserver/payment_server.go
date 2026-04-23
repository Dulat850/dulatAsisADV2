package grpcserver

import (
    "context"
    "log"

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
    log.Printf("Received payment request: OrderID=%s, Amount=%d", req.OrderId, req.Amount)

    payment, err := s.paymentUseCase.ProcessPayment(ctx, req.OrderId, req.Amount)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to process payment: %v", err)
    }

    return &pb.PaymentResponse{
        PaymentId:     payment.ID,
        OrderId:       payment.OrderID,
        TransactionId: payment.TransactionID,
        Status:        payment.Status,
    }, nil
}
