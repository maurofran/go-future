package future_test

import (
	"condominio102/pkg/future"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFuture(t *testing.T) {
	t.Run("ReturnWithResults", func(t *testing.T) {
		ctx := context.Background()
		fixture := future.New(ctx, func(ctx context.Context) (int, error) {
			return 1, nil
		})
		result, err := fixture(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 1, result)
	})

	t.Run("ReturnWithError", func(t *testing.T) {
		ctx := context.Background()
		expectedErr := errors.New("error")
		fixture := future.New(ctx, func(ctx context.Context) (int, error) {
			return 0, expectedErr
		})
		_, err := fixture(ctx)
		assert.ErrorIs(t, err, expectedErr)
	})

	t.Run("ContextCanceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		fixture := future.New(ctx, func(ctx context.Context) (int, error) {
			return 0, nil
		})
		cancel()
		_, err := fixture(ctx)
		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("ContextDeadlineExceeded", func(t *testing.T) {
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(100*time.Millisecond))
		fixture := future.New(ctx, func(ctx context.Context) (int, error) {
			time.Sleep(200 * time.Millisecond)
			return 0, nil
		})
		_, err := fixture(ctx)
		assert.ErrorIs(t, err, context.DeadlineExceeded)
		cancel()
	})
}
