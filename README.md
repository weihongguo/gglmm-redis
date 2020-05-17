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
## 缓存 -- 实现了gglmm的Cacher接口
```golang
func NewCacherConfig(config ConfigCacher) *Cacher
func NewCacher(network string, address string, maxActive int, maxIdle int, idleTimeout time.Duration, expires int) *Cacher
func NewCacherPool(pool *redis.Pool, expires int) *Cacher
func (cacher *Cacher) Close()
```
## 记数器
```golang
func NewCounterConfig(config ConfigCounter, name string) *Counter
func NewCounter(network string, address string, maxActive int, maxIdle int, idleTimeout time.Duration, name string) *Counter
func NewCounterPool(pool *redis.Pool, name string) *Counter
func (counter *Counter) Close()
```