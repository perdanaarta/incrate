package config

type Server struct {
	Host string `yaml:"host"`
	Port uint16 `yaml:"port"`
}

type Logging struct {
}

type Artifact struct {
	Storage string `yaml:"storage"`
}

type Config struct {
	Server   Server   `yaml:"server"`
	Logging  Logging  `yaml:"log"`
	Artifact Artifact `yaml:"artifact"`
}
