package gateway

import (
	"fmt"
	"strconv"
	"strings"
)

type ZoneFlags map[uint32]string

func (z ZoneFlags) String() string {
	var s []string
	for id, addr := range z {
		s = append(s, fmt.Sprintf("%d=%s", id, addr))
	}
	return strings.Join(s, ",")
}

func (z ZoneFlags) Set(v string) error {
	parts := strings.Split(v, "=")
	if len(parts) != 2 {
		return fmt.Errorf("zone flag must be <id>=<addr>")
	}
	idU, err := strconv.ParseUint(parts[0], 10, 32)
	if err != nil || idU == 0 {
		return fmt.Errorf("bad zone id")
	}
	if parts[1] == "" {
		return fmt.Errorf("missing addr")
	}
	if z == nil {
		return fmt.Errorf("internal: ZoneFlags nil")
	}
	z[uint32(idU)] = parts[1]
	return nil
}
