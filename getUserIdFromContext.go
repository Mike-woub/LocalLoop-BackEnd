package main

import (
	"context"
	"fmt"
)

func getUserIDFromContext(ctx context.Context) (int32, error) {
	userID, ok := ctx.Value(userIDKey).(int32)
	if !ok {
		return 0, fmt.Errorf("user ID not found in context")
	}
	return userID, nil
}
