package future

import "context"

// Fn is a function that returns a value of type T and an error.
type Fn[T any] func(context.Context) (T, error)

// New returns a function that returns a value of type T and an error.
func New[T any](ctx context.Context, f Fn[T]) Fn[T] {
	var result T
	var err error
	c := make(chan struct{}, 1)
	go func() {
		defer close(c)
		result, err = f(ctx)
	}()
	return func(ctx context.Context) (T, error) {
		select {
		case <-ctx.Done():
			return result, ctx.Err()
		case <-c:
			return result, err
		}
	}
}
