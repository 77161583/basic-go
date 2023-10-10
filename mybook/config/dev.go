//go:build !k8s

// Package config 本地启动的就是去启动docker-compose.yaml这个文件，去里面看配置
package config

var Config = config{
	DB: DBConfig{
		DSN: "root:root@tcp(localhost:13317)/webook",
	},
	Redis: RedisConfig{
		Addr: "localhost:6379",
	},
}
