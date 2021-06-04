package conf

type LogTransferCfg struct {
	KafkaCfg `ini:"kafka"`
	EsCfg `ini:"es"`
	EtcdCfg `ini:"etcd"`
}

type KafkaCfg struct {
	Brokers string `ini:"brokers"`
	Topics string `ini:"topics"`
	Group string `ini:"group"`
	Oldest bool `ini:"oldest"`
	Assignor string `ini:"assignor"`

}

type EsCfg struct {
	Address string `ini:"address"`
	ChanSize int `ini:"chansize"`
}

type EtcdCfg struct {
	Address string `ini:"address"`
	Timeout int `ini:"timeout"`
}