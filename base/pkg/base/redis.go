package base

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
)

type RedisOptions struct {
	Username string
	Password string
	Host     string
	Port     int
	URL      string
}

func (o *RedisOptions) IsSet() bool {
	return o.URL != "" || o.Host != ""
}

func (o *RedisOptions) ConnectionString() string {
	if o.URL != "" {
		return o.URL
	}
	host := o.Host
	i := strings.Index(host, "://")
	var protocol string
	if i != -1 {
		// url has protocol
		protocol = host[:i+3]
		host = host[i+3:]
	} else {
		// for security reasons, default to encrypted
		protocol = "rediss://"
	}
	return fmt.Sprintf(
		"%s%s:%s@%s:%d",
		protocol,
		o.Username,
		o.Password,
		host,
		o.Port,
	)
}

func RedisEnv(o *RedisOptions, required bool) *RedisOptions {
	if o == nil {
		o = &RedisOptions{}
	}
	CheckEnv("REDIS_URL", &o.URL)
	CheckEnv("REDIS_HOST", &o.Host)
	CheckEnvInt("REDIS_PORT", &o.Port)
	CheckEnv("REDIS_USERNAME", &o.Username)
	CheckEnv("REDIS_PASSWORD", &o.Password)
	if required && o.URL == "" {
		if o.Host == "" {
			panic("missing --redis-host")
		}
		if o.Port == 0 {
			panic("missing --redis-port")
		}
	}
	return o
}

func AddRedisFlags(cmd *cobra.Command, o *RedisOptions) {
	cmd.PersistentFlags().StringVar(&o.URL, "redis-url", "", "redis url (full connection string)")
	cmd.PersistentFlags().StringVar(&o.Host, "redis-host", "", "redis host")
	cmd.PersistentFlags().IntVar(&o.Port, "redis-port", 6379, "redis port")
	cmd.PersistentFlags().StringVar(&o.Username, "redis-username", "", "redis username")
	cmd.PersistentFlags().StringVar(&o.Password, "redis-password", "", "redis password")
}

func ConnectRedis(
	ctx context.Context,
	o *RedisOptions,
) *redis.Client {
	if !o.IsSet() {
		panic("missing redis options")
	}
	url := o.ConnectionString()
	options, err := redis.ParseURL(url)
	if err != nil {
		panic(fmt.Errorf("failed to parse redis url '%s': %v", url, err))
	}
	redisClient := redis.NewClient(options)
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		panic(errors.Wrap(err, "failed to ping redis"))
	}
	DefaultLog.Debug("connected to redis", Elapsed(start))
	return redisClient
}
