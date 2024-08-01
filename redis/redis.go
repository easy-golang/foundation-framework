package redis

import (
	"context"
	"github.com/easy-golang/foundation-framework/constant/symbol"
	"github.com/easy-golang/foundation-framework/err/comm"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
	"sync/atomic"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"time"
)

type Conf struct {
	redis.RedisConf
	LockPath string `json:",optional"`
}

type Client struct {
	*redis.Redis
	LockPath string
}

func NewClient(conf Conf, opts ...redis.Option) *Client {
	return &Client{
		Redis:    redis.MustNewRedis(conf.RedisConf, opts...),
		LockPath: conf.LockPath,
	}
}

func (c Client) NewRedisLock(keys ...string) *lock {
	key := strings.Join(keys, symbol.Colon)
	return &lock{
		RedisLock: redis.NewRedisLock(c.Redis, "lock:"+c.LockPath+":"+key),
	}
}

type lock struct {
	*redis.RedisLock
	//锁状态 0 未获取锁 1 获取锁
	state uint32
}

func (l *lock) Lock() bool {
	for {
		ok, err := l.Acquire()
		if err != nil {
			logx.Error(comm.Wrap(err))
			return false
		}
		if ok {
			l.SetExpire(60)
			// 锁重入的时候直接返回，不用执行锁续期
			if atomic.LoadUint32(&l.state) == 1 {
				return true
			}
			atomic.StoreUint32(&l.state, 1)
			// 锁续期
			go l.renewLock()

			return true
		}
		// 10 毫秒重试直到获取锁成功
		time.Sleep(time.Millisecond * 10)
	}
}

func (l *lock) LockWithContext(ctx context.Context) bool {
	_, ok := ctx.Deadline()
	if !ok {
		logx.Error(comm.New("context must with timeout"))
		return false
	}
	l.SetExpire(60)
	for {
		ok, err := l.Acquire()
		if err != nil {
			logx.Error(comm.Wrap(err))
			return false
		}
		if ok {
			// 锁重入的时候直接返回，不用执行锁续期
			if atomic.LoadUint32(&l.state) == 1 {
				return true
			}
			atomic.StoreUint32(&l.state, 1)
			// 锁续期
			go l.renewLock()
			return true
		} else {
			select {
			// 如果等待时间超时，或者退出则直接返回false获取锁失败
			case <-ctx.Done():
				return false
			}
		}
		// 200 毫秒重试直到获取锁成功
		time.Sleep(time.Millisecond * 200)
	}
}

func (l *lock) TryLock() bool {
	ok, err := l.Acquire()
	if err != nil {
		logx.Error(comm.Wrap(err))
		return false
	}
	if ok {
		l.SetExpire(60)
		// 锁重入的时候直接返回，不用执行锁续期
		if atomic.LoadUint32(&l.state) == 1 {
			return true
		}
		atomic.StoreUint32(&l.state, 1)
		// 锁续期
		go l.renewLock()
		return true
	}
	return false
}

func (l *lock) Unlock() bool {
	ok, err := l.Release()
	if err != nil {
		logx.Error(comm.Wrap(err))
		return false
	}
	if ok {
		atomic.StoreUint32(&l.state, 0)
		return true
	}
	logx.Error(comm.New("unlock failed"))
	return false
}

func (l *lock) LockFunc(fn func() error) error {
	b := l.Lock()
	defer l.Unlock()
	if b {
		return fn()
	}
	return comm.New("LockFunc failed")
}

func (l *lock) LockFuncWithContext(ctx context.Context, fn func() error) error {
	b := l.LockWithContext(ctx)
	defer l.Unlock()
	if b {
		return fn()
	}
	return comm.New("LockFunc failed")
}

func (l *lock) TryLockFunc(fn func() error) error {
	b := l.TryLock()
	defer l.Unlock()
	if b {
		return fn()
	}
	return comm.New("TryLockFunc failed")
}

func (l *lock) renewLock() {
	for {
		time.Sleep(time.Second * 50)
		if atomic.LoadUint32(&l.state) == 0 {
			return
		}
		// 没有释放进行锁续期
		ok, _ := l.Acquire()
		if !ok {
			return
		}
	}
}
