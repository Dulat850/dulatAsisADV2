package grpc

import (
    "context"
    "time"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    
    paymentPb "dulatAsisADV2/proto/gen/payment"
)

type PaymentServiceClient struct {
    client paymentPb.PaymentServiceClient
    conn   *grpc.ClientConn
}

func NewPaymentServiceClient(addr string) (*PaymentServiceClient, error) {
    conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        return nil, err
    }
    
    client := paymentPb.NewPaymentServiceClient(conn)
    
    return &PaymentServiceClient{
        client: client,
        conn:   conn,
    }, nil
}

func (c *PaymentServiceClient) ProcessPayment(ctx context.Context, orderID string, amount float64) (string, error) {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()
    
    req := &paymentPb.PaymentRequest{
        OrderId:  orderID,
        Amount:   amount,
        Currency: "USD",
    }
    
    resp, err := c.client.ProcessPayment(ctx, req)
    if err != nil {
        return "", err
    }
    
    return resp.PaymentId, nil
}

func (c *PaymentServiceClient) Close() error {
    return c.conn.Close()
}
