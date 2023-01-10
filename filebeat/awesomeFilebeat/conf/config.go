package conf

type Config struct {
	Etcd Etcd `ini:"etcd"`
}

type Etcd struct {
	Address string `ini:"address"`
	Key     string `ini:"collect_log_key"`
	Timeout int    `ini:"timeout"`
}
