package interactions

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
)

type Param struct {
	Key   string
	Value string
}

func Encode(prefix string, params ...Param) (string, error) {
	if len(params) == 0 {
		return prefix, nil
	}

	var parts []string
	for _, param := range params {
		parts = append(parts, fmt.Sprintf("%s=%s", param.Key, param.Value))
	}

	encodedParams := strings.Join(parts, "&")

	encodedID := fmt.Sprintf("%s:%s", prefix, encodedParams)
	if len(encodedID) > 100 {
		return "", errors.New("character limit reached on custom id")
	}

	return encodedID, nil
}

func Decode(customID string) (string, map[string]string, error) {
	params := make(map[string]string)
	if !strings.Contains(customID, ":") {
		return customID, params, nil
	}

	parts := strings.Split(customID, ":")
	if len(parts) != 2 {
		return "", params, fmt.Errorf("malformed custom id: %q", customID)
	}

	prefix := parts[0]
	kvps := strings.Split(parts[1], "&")

	for _, kvp := range kvps {
		kvpParts := strings.Split(kvp, "=")
		if len(kvpParts) != 2 {
			slog.Warn("malformed kvp", "kvp", kvp, "customID", prefix)
			continue
		}

		params[kvpParts[0]] = kvpParts[1]
	}

	return prefix, params, nil
}
