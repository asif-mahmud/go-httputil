package validator

import "github.com/go-playground/mold/v4"

// RegisterScrubber registers a new scrubber or replaces existing one with tag.
func RegisterScrubber(tag string, fn mold.Func) {
	scrub.Register(tag, fn)
}
