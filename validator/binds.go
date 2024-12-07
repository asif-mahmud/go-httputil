package validator

import (
	"context"
	"encoding/json"
	"io"
	"net/url"
)

// BindUrlValues binds url.Values into a struct instance.
// Additionally it runs the struct through mold transformer.
func BindUrlValues(ctx context.Context, v url.Values, s any) error {
	if err := formDecoder.Decode(s, v); err != nil {
		return err
	}

	return conform.Struct(ctx, v)
}

// BindJSON binds body into a struct instance.
// Additionally it runs the struct through mold transformer.
func BindJSON(ctx context.Context, body io.ReadCloser, s any) error {
	defer body.Close()
	if err := json.NewDecoder(body).Decode(s); err != nil {
		return err
	}

	return conform.Struct(ctx, s)
}
