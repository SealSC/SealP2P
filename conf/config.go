package conf

type Config struct {
	ID         string
	ClientOnly bool
	PKFile     string

	ServerPort    int
	MulticastAddr string
	MulticastPort int
}

func (c Config) Clone(id string) *Config {
	return &Config{
		ID:            id,
		ClientOnly:    c.ClientOnly,
		PKFile:        c.PKFile,
		ServerPort:    c.ServerPort,
		MulticastAddr: c.MulticastAddr,
		MulticastPort: c.MulticastPort,
	}
}

var DefaultConfig = &Config{
	ID:            "",
	ClientOnly:    false,
	PKFile:        "SealP2PPK",
	ServerPort:    3333,
	MulticastAddr: "224.0.0.3",
	MulticastPort: 5678,
}
