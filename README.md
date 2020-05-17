# gglmm-redis
## 依赖
+ github.com/gomodule/redigo/redis
## Factory
```golang
type Factory struct {
	pool *redis.Pool
}

func NewFactoryConfig(config ConfigRedis) *Factory
func NewFactory(network string, address string, maxActive int, maxIdle int, idleTimeout time.Duration) *Factory
func (factory *Factory) Close()

func (factory *Factory) NewCacher(expires int) *Cacher
func (factory *Factory) NewCounter(name string) *Counter
func (factory *Factory) NewToper(name string, limit int) *Toper
func (factory *Factory) NewHoter(name string) *Hoter
func (factory *Factory) NewMessageQueue(channel string) *MessageQueue
```