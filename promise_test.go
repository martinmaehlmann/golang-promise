package promise

import (
	"errors"
	"reflect"
	"sync"
	"testing"
)

func TestOnePromise(t *testing.T) {
	type testPromise struct {
		valueCallCount uint
		errorCallCount uint
		value          string
		err            error
	}
	type args struct {
		testPromise *testPromise
	}
	type want struct {
		wantedString         string
		wantedErr            error
		wantedCallCountValue uint
		wantedCallCountErr   uint
	}
	type testCase[T any] struct {
		name string
		args args
		want want
	}
	tests := []testCase[testPromise]{
		{
			name: "One Promise returns correctly",
			args: args{
				testPromise: &testPromise{
					valueCallCount: 0,
					errorCallCount: 0,
					value:          "test",
					err:            nil,
				},
			},
			want: want{
				wantedString:         "test",
				wantedErr:            nil,
				wantedCallCountValue: 1,
				wantedCallCountErr:   0,
			},
		},
		{
			name: "One Promise errors correctly",
			args: args{
				testPromise: &testPromise{
					valueCallCount: 0,
					errorCallCount: 0,
					value:          "test",
					err:            errors.New("test error"),
				},
			},
			want: want{
				wantedString:         "",
				wantedErr:            errors.New("test error"),
				wantedCallCountValue: 1,
				wantedCallCountErr:   1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result string
			var errr error

			promise := tt.args.testPromise

			wg := new(sync.WaitGroup)
			wg.Add(1)

			NewPromise[string](func() (string, error) {
				promise.valueCallCount++

				if promise.err != nil {
					return "", promise.err
				}

				return promise.value, nil
			}).Then(func(str string) {
				result = str
			}, func(e error) {
				promise.errorCallCount++
				errr = e
			}).Finally(func() {
				wg.Done()
			})

			wg.Wait()

			if !reflect.DeepEqual(result, tt.want.wantedString) {
				t.Errorf("result = %s, want %s", result, tt.want.wantedString)
			}

			if !reflect.DeepEqual(errr, tt.want.wantedErr) {
				t.Errorf("error = %s, want %s", errr, tt.want.wantedErr)
			}

			if promise.valueCallCount != tt.want.wantedCallCountValue {
				t.Errorf("value call count wrong = %d, want %d", promise.valueCallCount, tt.want.wantedCallCountValue)
			}

			if promise.errorCallCount != tt.want.wantedCallCountErr {
				t.Errorf("error call count wrong= %d, want %d", promise.errorCallCount, tt.want.wantedCallCountErr)
			}
		})
	}
}
