package parse

type Stat struct {
	Key string `toml:"key"`
	Val string `toml:"value"`
}

type Base struct {
	ID     string `toml:"id"`
	IconID string `toml:"icon_id"`
	Stats  []Stat `toml:"stats"`
}
