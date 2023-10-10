//go:build k8s

package config

var Config = config{
	DB: DBConfig{
		DSN: "root:root@tcp(10.102.156.176:3308)/mysql",
	},
	Redis: RedisConfig{
		Addr: "mybook-live-redis:6380",
	},
}
