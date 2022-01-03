package logger

import "context"

type Options struct {
	metadata  map[string]interface{}
	err       error
	fmtArgs   []interface{}
	namespace string
	context   context.Context
}

var defaultOptions = &Options{
	metadata:  nil,
	err:       nil,
	fmtArgs:   nil,
	namespace: "",
	context:   nil,
}

type SetOptionFn = func(*Options)

func evaluateOptions(args []interface{}) *Options {
	optCopy := &Options{}
	*optCopy = *defaultOptions
	for _, v := range args {
		fn, ok := v.(SetOptionFn)
		if !ok {
			// Skipping
			continue
		}
		fn(optCopy)
	}
	return optCopy
}

func AddMetadata(key string, val interface{}) SetOptionFn {
	return func(o *Options) {
		if o.metadata == nil {
			o.metadata = make(map[string]interface{})
		}
		o.metadata[key] = val
	}
}

func Metadata(m map[string]interface{}) SetOptionFn {
	return func(o *Options) {
		o.metadata = m
	}
}

func Format(args ...interface{}) SetOptionFn {
	return func(o *Options) {
		o.fmtArgs = args
	}
}

func Error(err error) SetOptionFn {
	return func(o *Options) {
		o.err = err
	}
}

func WithNamespace(n string) SetOptionFn {
	return func(o *Options) {
		o.namespace = n
	}
}

func Context(ctx context.Context) SetOptionFn {
	return func(o *Options) {
		o.context = ctx
	}
}
