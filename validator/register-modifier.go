package validator

import "github.com/go-playground/mold/v4"

// RegisterModifier registers a new modifier or replaces an existing one with tag.
// The modifier can be used using the tag inside struct
// tag mod.
func RegisterModifier(tag string, fn mold.Func) {
	conform.Register(tag, fn)
}
