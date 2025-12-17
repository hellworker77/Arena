package gateway

import "fmt"

func sscanf(s, format string, a ...any) (int, error) { return fmt.Sscanf(s, format, a...) }
func sprintf(format string, a ...any) string { return fmt.Sprintf(format, a...) }
