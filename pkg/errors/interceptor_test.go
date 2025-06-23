package errors

import (
	"context"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/interop/grpc_testing"
	"google.golang.org/grpc/test/bufconn"
)

type testSVC struct {
	grpc_testing.UnimplementedTestServiceServer
	err error
}

func (t *testSVC) EmptyCall(context.Context, *grpc_testing.Empty) (*grpc_testing.Empty, error) {
	return nil, t.err
}

func TestErrorInterceptor(t *testing.T) {

	acceptErr := NewError(ErrorTooManyRequest, "test error, don't care")
	listener := bufconn.Listen(1024 * 1024)

	svc := &testSVC{
		err: acceptErr,
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(UnaryServerErrorInterceptor()))

	grpc_testing.RegisterTestServiceServer(s, svc)
	go func() {
		if err := s.Serve(listener); err != nil {
			t.Errorf("Server exited with error: %v", err)
		}
	}()
	defer s.Stop()

	// client
	conn, err := grpc.DialContext(context.Background(),
		"bufnet",
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return listener.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(UnaryClientInterceptor()),
	)
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := grpc_testing.NewTestServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = client.EmptyCall(ctx, &grpc_testing.Empty{})

	if !Is(err, acceptErr) {
		t.Errorf("want err: %v\r\nbut got err: %v", acceptErr, err)
	}

}
