package metrics

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

// Very small Prometheus text exporter without extra deps.

type Counters struct {
	TickNanosTotal atomic.Int64
	TickCount      atomic.Int64

	Entities atomic.Int64
	Players  atomic.Int64

	RepBytes atomic.Int64
}

func (c *Counters) ObserveTick(d time.Duration) {
	c.TickNanosTotal.Add(d.Nanoseconds())
	c.TickCount.Add(1)
}

func (c *Counters) AddRepBytes(n int) {
	c.RepBytes.Add(int64(n))
}

func (c *Counters) Serve(addr string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4")
		ticks := c.TickCount.Load()
		total := c.TickNanosTotal.Load()
		avg := int64(0)
		if ticks > 0 {
			avg = total / ticks
		}
		fmt.Fprintf(w, "zone_tick_count %d
", ticks)
		fmt.Fprintf(w, "zone_tick_avg_nanos %d
", avg)
		fmt.Fprintf(w, "zone_entities %d
", c.Entities.Load())
		fmt.Fprintf(w, "zone_players %d
", c.Players.Load())
		fmt.Fprintf(w, "zone_rep_bytes_total %d
", c.RepBytes.Load())
	})
	srv := &http.Server{Addr: addr, Handler: mux}
	go func() { _ = srv.ListenAndServe() }()
	return srv
}
