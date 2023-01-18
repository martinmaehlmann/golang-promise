package promise

import "sync"

type Promise[T any] struct {
	fetchWG sync.WaitGroup
	writeWG sync.WaitGroup
	res     T
	err     error
}

// NewPromise creates a new Promise to return a value of T or an error. It is meant to be used with an outer
// sync.WaitGroup
func NewPromise[T any](fun func() (T, error)) *Promise[T] {
	promise := &Promise[T]{}
	promise.fetchWG.Add(1)
	promise.writeWG.Add(1)

	go func() {
		defer promise.fetchWG.Done()
		promise.res, promise.err = fun()
	}()

	return promise
}

// Then is a function that takes a result function res and an error function err. The result function res is used to
// assign the result of the computation to a variable to be used later. The err function is used to catch any ocurring
// error.
func (p *Promise[T]) Then(res func(T), err func(error)) *Promise[T] {
	go func() {
		defer p.writeWG.Done()
		p.fetchWG.Wait()

		if p.err != nil {
			err(p.err)

			return
		}

		res(p.res)
	}()

	return p
}

// Finally is used to do anyting after writing errors or values. It should be used to call wg.Done as seen in the above
// example.
func (p *Promise[T]) Finally(fun func()) {
	go func() {
		p.writeWG.Wait()
		fun()
	}()
}
