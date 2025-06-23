package errors

import (
	"context"
	"io"

	"github.com/chise0904/golang_template/pkg/grpc/interceptor"
	"github.com/chise0904/golang_template/pkg/trace"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// UnaryServerErrorInterceptor recode error and convert error to grpc error
func UnaryServerErrorInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			logFields := map[string]interface{}{}
			logFields["requestID"] = trace.XRequestIDFromContextForGRPC(ctx)
			//logFields["time"] = trace.SubTimeFromContextForGRPC(ctx)
			logFields["input"] = req
			logFields["output"] = resp
			causeErr := errors.Cause(err)
			oErr, ok := causeErr.(*_error)
			if !ok || causeErr == nil {
				oErr = &_error{grpcCode: codes.Internal, category: CategoryInternalServiceError, code: "500001", message: "internal error", top: new(int)}

			}
			stake, ok := err.(interface {
				StackTrace() errors.StackTrace
			})

			logger := log.With().Fields(logFields).Logger()
			if oErr.category >= CategoryInternalServiceError {
				if !ok {
					logger.Error().Msgf("msg: %+v", err)
				} else {
					logger.Error().Msgf("msg: %s %+v", err.Error(), stake.StackTrace())
				}
			} else {
				if !ok {
					logger.Debug().Msgf("msg: %+v", err)
				} else {
					logger.Debug().Msgf("msg: %s %+v", err.Error(), stake.StackTrace())
				}
			}
			e := ConvertErrorToGrpcErr(err)
			if e != nil {
				return nil, e
			}

			return resp, err
		}
		return resp, err
	}
}

func StreamServerErrorInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		_ss := interceptor.WrapServerStream(ss)
		err := handler(srv, _ss)
		if err != nil {
			logFields := map[string]interface{}{}
			logFields["requestID"] = trace.XRequestIDFromContextForGRPC(_ss.WrappedContext)
			//logFields["time"] = trace.SubTimeFromContextForGRPC(_ss.WrappedContext)

			causeErr := errors.Cause(err)
			oErr, ok := causeErr.(*_error)
			if !ok || causeErr == nil {
				oErr = &_error{grpcCode: codes.Internal, category: CategoryInternalServiceError, code: "500001", message: "internal error", top: new(int)}

			}
			stake, ok := err.(interface {
				StackTrace() errors.StackTrace
			})

			logger := log.With().Fields(logFields).Logger()
			if oErr.category >= CategoryInternalServiceError {
				if !ok {
					logger.Error().Msgf("msg: %+v", err)
				} else {
					logger.Error().Msgf("msg: %s %+v", err.Error(), stake.StackTrace())
				}
			} else {
				if !ok {
					logger.Debug().Msgf("msg: %+v", err)
				} else {
					logger.Debug().Msgf("msg: %s %+v", err.Error(), stake.StackTrace())
				}
			}
			return ConvertErrorToGrpcErr(err)
		}
		return err
	}
}

// UnaryClientInterceptor catch grpc error to errors.err
func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return ConvertGrpcErrToHttpErr(invoker(ctx, method, req, reply, cc, opts...))
	}
}

type wrappedStream struct {
	grpc.ClientStream
}

func StreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		c, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			return nil, ConvertGrpcErrToHttpErr(err)
		}

		return &wrappedStream{c}, nil
	}
}

func (c *wrappedStream) RecvMsg(m any) error {
	err := c.ClientStream.RecvMsg(m)
	if err != nil {
		if err == io.EOF {
			return err
		}
		return ConvertGrpcErrToHttpErr(err)
	}
	return nil
}

func (c *wrappedStream) CloseSend() error {
	err := c.ClientStream.CloseSend()
	if err != nil {
		return ConvertGrpcErrToHttpErr(err)
	}
	return nil
}
