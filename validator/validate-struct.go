package validator

import "context"

// ValidateStruct validates the incoming structure.
func ValidateStruct(ctx context.Context, s any) error {
	return validate.StructCtx(ctx, s)
}
