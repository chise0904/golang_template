package errors

type baseError func() *_error

// BadRequest
var ErrorInvalidInput baseError = func() *_error {
	return errorInvalidInput
}

// Unauthorized
var ErrorUnauthorized baseError = func() *_error {
	return errorUnauthorized
}
var ErrorPasswordNotCorrect baseError = func() *_error {
	return errorPasswordNotCorrect
}

// ResourceNotFound
var ErrorResourceNotFound baseError = func() *_error {
	return errorResourceNotFound
}
var ErrorPageNotFound baseError = func() *_error {
	return errorPageNotFound
}

// Conflict
var ErrorConflict baseError = func() *_error {
	return errorConflict
}

// 429
var ErrorTooManyRequest baseError = func() *_error {
	return errorTooManyRequest
}

// Forbidden
var ErrorNotAllow baseError = func() *_error {
	return errorNotAllow
}
var ErrorExceedStockQuantity baseError = func() *_error {
	return errorExceedStockQuantity
}
var ErrorExceededMaximumPurchaseQuantity baseError = func() *_error {
	return errorExceededMaximumPurchaseQuantity
}

var ErrorInternalError baseError = func() *_error {
	return errorInternalError
}
