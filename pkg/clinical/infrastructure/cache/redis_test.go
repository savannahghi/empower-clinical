package cache_test

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
	"github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/cache"
	fakeCache "github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/cache/mock"
)

type mockHandler struct {
	cache *fakeCache.MockCacheService
}

func TestCacheStore_Get(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name  string
		setup func(mh *mockHandler) args
		want  *redis.StringCmd
	}{
		{
			name: "Happy Case: Successfully retrieve data from cache",
			setup: func(mh *mockHandler) args {
				mh.cache.EXPECT().Get(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, key string) *redis.StringCmd {
						return redis.NewStringCmd(context.Background())
					})
				return args{
					ctx: ctx,
					key: gofakeit.UUID(),
				}
			},
		},
		{
			name: "Sad Case: return nil",
			setup: func(mh *mockHandler) args {
				mh.cache.EXPECT().Get(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, key string) *redis.StringCmd {
						return redis.NewStringCmd(context.Background())
					})

				return args{
					ctx: ctx,
					key: gofakeit.UUID(),
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := fakeCache.NewMockCacheService(t)
			c := cache.NewCacheStore(mock)

			args := tt.setup(&mockHandler{cache: mock})

			got := c.Get(args.ctx, args.key)
			if got.Err() != nil {
				t.Errorf("CacheStore.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCacheStore_Set(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx        context.Context
		key        string
		value      interface{}
		expiration time.Duration
	}
	tests := []struct {
		name  string
		setup func(mh *mockHandler) args
		want  *redis.StatusCmd
	}{
		{
			name: "Happy Case: Successfully set data in the cache",
			setup: func(mh *mockHandler) args {
				mh.cache.EXPECT().Set(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
						return redis.NewStatusCmd(context.Background())
					})

				return args{
					ctx:        ctx,
					key:        gofakeit.UUID(),
					value:      gofakeit.Name(),
					expiration: time.Hour,
				}
			},
		},
		{
			name: "Sad Case: Fail to pass a key value",
			setup: func(mh *mockHandler) args {
				mh.cache.EXPECT().Set(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
						return redis.NewStatusCmd(context.Background())
					})

				return args{
					ctx:        ctx,
					value:      gofakeit.Name(),
					expiration: time.Hour,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := fakeCache.NewMockCacheService(t)
			c := cache.NewCacheStore(mock)

			args := tt.setup(&mockHandler{cache: mock})

			got := c.Set(args.ctx, args.key, args.value, args.expiration)
			if got.Err() != nil {
				t.Errorf("CacheStore.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
