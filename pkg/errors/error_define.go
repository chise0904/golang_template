package errors

import "google.golang.org/grpc/codes"

var (
	errorInvalidInput = &_error{category: CategoryBadRequest, grpcCode: codes.InvalidArgument, code: "400001", message: "invalid input", top: new(int)}

	// Unauthorized
	errorUnauthorized       = &_error{category: CategoryUnauthorized, grpcCode: codes.Unauthenticated, code: "401001", message: "unauthorized", top: new(int)}
	errorPasswordNotCorrect = &_error{category: CategoryUnauthorized, grpcCode: codes.Unauthenticated, code: "401002", message: "input password not correct", top: new(int)}

	// ResourceNotFound
	errorResourceNotFound = &_error{category: CategoryResourceNotFound, grpcCode: codes.NotFound, code: "404001", message: "resource not found", top: new(int)}
	errorPageNotFound     = &_error{category: CategoryResourceNotFound, grpcCode: codes.NotFound, code: "404001", message: "page not found", top: new(int)}

	// Conflict
	errorConflict = &_error{category: CategoryConflict, grpcCode: codes.AlreadyExists, code: "409001", message: "the request conflict", top: new(int)}

	// 429
	errorTooManyRequest = &_error{category: CategoryTooManyRequests, grpcCode: codes.ResourceExhausted, code: "429001", message: "too many request", top: new(int)}

	// Forbidden
	errorNotAllow                        = &_error{category: CategoryRequestForbidden, grpcCode: codes.PermissionDenied, code: "403001", message: "forbidden", top: new(int)}
	errorExceedStockQuantity             = &_error{category: CategoryRequestForbidden, grpcCode: codes.PermissionDenied, code: "403002", message: "forbidden", top: new(int)}
	errorExceededMaximumPurchaseQuantity = &_error{category: CategoryRequestForbidden, grpcCode: codes.PermissionDenied, code: "403003", message: "forbidden", top: new(int)}

	errorInternalError = &_error{category: CategoryInternalServiceError, grpcCode: codes.Internal, code: "500001", message: "internal error", top: new(int)}
)
