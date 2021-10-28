package conf

type Config struct {
	ID         string
	ClientOnly bool
	PKFile     string

	ServerPort    int
	MulticastAddr string
	MulticastPort int
}

var DefaultConfig = &Config{
	ID:            "",
	ClientOnly:    false,
	PKFile:        "SealP2PPK",
	ServerPort:    3333,
	MulticastAddr: "224.0.0.3",
	MulticastPort: 5678,
}
