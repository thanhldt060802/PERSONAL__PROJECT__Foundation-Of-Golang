package auth

import (
	"context"
	"fmt"
)

type SimpleAuthMiddleware struct {
}

func NewSimpleAuthMiddleware() IAuthMiddleware {
	return &SimpleAuthMiddleware{}
}

func (mdw *SimpleAuthMiddleware) AuthMiddleware(ctx context.Context) (string, error) {
	token, ok := ctx.Value("auth-header").(string)
	if !ok {
		return "", nil
	}

	result := fmt.Sprintf("Token %v is accepted", token)

	return result, nil
}
