package queueCollection

import (
	"reflect"
	"testing"
)

func Test_queue_set(t *testing.T) {
	type args[T any] struct {
		elem T
	}
	type testCase[T any] struct {
		name        string
		elemInQueue []T
		args        args[T]
		wantElem    []T
		wantErr     bool
	}
	tests := []testCase[int]{
		{
			name:        "",
			elemInQueue: []int{},
			args: args[int]{
				elem: 1,
			},
			wantElem: []int{1},
			wantErr:  false,
		},
		{
			name:        "",
			elemInQueue: []int{1, 2, 3},
			args: args[int]{
				elem: 4,
			},
			wantElem: []int{1, 2, 3, 4},
			wantErr:  false,
		},
		{
			name:        "",
			elemInQueue: make([]int, queueMaxLen),
			args: args[int]{
				elem: 21,
			},
			wantElem: make([]int, queueMaxLen),
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			que := newQueue[int]()
			for _, elem := range tt.elemInQueue {
				que.queue <- elem
				que.len++
			}
			err := que.set(tt.args.elem)
			if (err != nil) != tt.wantErr {
				t.Errorf("set() error = %v, wantErr %v", err, tt.wantErr)
			}

			var factElem []int
			for que.len != 0 {
				que.len--

				factElem = append(factElem, <-que.queue)
			}
			if !reflect.DeepEqual(tt.wantElem, factElem) {
				t.Errorf("set() elem = %v, wantElem = %v", factElem, tt.wantElem)
			}
		})
	}
}

func Test_queue_get(t *testing.T) {
	type testCase[T any] struct {
		name        string
		elemInQueue []T
		wantElem    T
		wantErr     bool
	}
	tests := []testCase[int]{
		{
			name:        "",
			elemInQueue: []int{1},
			wantElem:    1,
			wantErr:     false,
		},
		{
			name:        "",
			elemInQueue: []int{},
			wantElem:    0,
			wantErr:     true,
		},
		{
			name:        "",
			elemInQueue: []int{12133, 2, 3, 4, 5},
			wantElem:    12133,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			que := newQueue[int]()
			for _, elem := range tt.elemInQueue {
				que.queue <- elem
				que.len++
			}
			gotElem, err := que.get()
			if (err != nil) != tt.wantErr {
				t.Errorf("get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotElem, tt.wantElem) {
				t.Errorf("get() gotElem = %v, want %v", gotElem, tt.wantElem)
			}
		})
	}
}
