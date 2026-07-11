package apiErrors

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// RenderJSON returns a JSON catalog mapping every error's code to its message.
// Input order is preserved and a blank line is inserted whenever the code's
// "century" (code/100) changes between consecutive entries. Returns an error
// if any code is duplicated, malformed, or paired with a nil pointer.
func RenderJSON(errs []*APIError) ([]byte, error) {
	seen := make(map[string]string, len(errs))
	for _, e := range errs {
		if e == nil {
			return nil, fmt.Errorf("nil entry in errors slice")
		}
		if prev, dup := seen[e.Code]; dup {
			return nil, fmt.Errorf("duplicate code %s: %q vs %q", e.Code, prev, e.Message)
		}
		seen[e.Code] = e.Message
	}

	var buf bytes.Buffer
	buf.WriteString("{\n")
	prevGroup := -1
	for i, e := range errs {
		num, err := strconv.Atoi(strings.TrimPrefix(e.Code, "SE-"))
		if err != nil {
			return nil, fmt.Errorf("invalid code %q: %w", e.Code, err)
		}
		group := num / 100
		if i > 0 && group != prevGroup {
			buf.WriteString("\n")
		}
		prevGroup = group

		fmt.Fprintf(&buf, "  %q: %s", e.Code, strconv.Quote(e.Message))
		if i < len(errs)-1 {
			buf.WriteByte(',')
		}
		buf.WriteByte('\n')
	}
	buf.WriteString("}\n")
	return buf.Bytes(), nil
}
