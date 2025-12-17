package gateway

import "fmt"

func sscanf(s, f string, a ...any) (int,error){return fmt.Sscanf(s,f,a...)}
func sprintf(f string, a ...any) string {return fmt.Sprintf(f,a...)}
