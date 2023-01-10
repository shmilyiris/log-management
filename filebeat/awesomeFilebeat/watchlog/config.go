package watchlog

type Configs struct {
	Filebeat Filebeat `yaml:"filebeat"`
	Output   Output   `yaml:"output"`
}

type Filebeat struct {
	Inputs []Inputs `yaml:"inputs"`
}

type Inputs struct {
	Type    string   `yaml:"type"`
	Enabled bool     `yaml:"enabled"`
	Paths   []string `yaml:"paths"`
	Fields  Fields   `yaml:"fields"`
}

type Fields struct {
	LogTopics string `yaml:"log_topics"`
	LogTypes  string `yaml:"log_types"`
}

type Config struct {
	Modules Modules `yaml:"modules"`
}

type Modules struct {
	Enabled bool   `yaml:"enabled"`
	Path    string `yaml:"path"`
}

type Output struct {
	Kafka Kafka `yaml:"kafka"`
}

type Kafka struct {
	Hosts []string `yaml:"hosts"`
	Topic string   `yaml:"topic"`
	Codec Codec    `yaml:"codec"`
}

type Codec struct {
	JSON JSON `yaml:"jSON"`
}

type JSON struct {
	Pretty bool `yaml:"pretty"`
}
