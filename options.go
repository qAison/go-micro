package micro

import (
	"time"

	"github.com/micro/cli"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/selector"
	"github.com/micro/go-micro/server"
	"github.com/micro/go-micro/transport"

	"golang.org/x/net/context"
)

type Options struct {
	Broker    broker.Broker
	Cmd       cmd.Cmd
	Client    client.Client
	Server    server.Server
	Registry  registry.Registry
	Transport transport.Transport

	// Register loop interval
	RegisterInterval time.Duration

	// Before and After funcs
	BeforeStart []func() error
	AfterStop   []func() error

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

func newOptions(opts ...Option) Options {
	opt := Options{
		Broker:    broker.DefaultBroker,
		Cmd:       cmd.DefaultCmd,
		Client:    client.DefaultClient,
		Server:    server.DefaultServer,
		Registry:  registry.DefaultRegistry,
		Transport: transport.DefaultTransport,
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

func Broker(b broker.Broker) Option {
	return func(o *Options) {
		o.Broker = b
		// Update Client and Server
		o.Client.Init(client.Broker(b))
		o.Server.Init(server.Broker(b))
	}
}

func Cmd(c cmd.Cmd) Option {
	return func(o *Options) {
		o.Cmd = c
	}
}

func Client(c client.Client) Option {
	return func(o *Options) {
		o.Client = c
	}
}

func Server(s server.Server) Option {
	return func(o *Options) {
		o.Server = s
	}
}

// Registry sets the registry for the service
// and the underlying components
func Registry(r registry.Registry) Option {
	return func(o *Options) {
		o.Registry = r
		// Update Client and Server
		o.Client.Init(client.Registry(r))
		o.Server.Init(server.Registry(r))
		// Update Selector
		o.Client.Options().Selector.Init(selector.Registry(r))
	}
}

// Transport sets the transport for the service
// and the underlying components
func Transport(t transport.Transport) Option {
	return func(o *Options) {
		o.Transport = t
		// Update Client and Server
		o.Client.Init(client.Transport(t))
		o.Server.Init(server.Transport(t))
	}
}

// Convenience options

// Name of the service
func Name(n string) Option {
	return func(o *Options) {
		o.Server.Init(server.Name(n))
	}
}

// Version of the service
func Version(v string) Option {
	return func(o *Options) {
		o.Server.Init(server.Version(v))
	}
}

// Command line Description
func Description(d string) Option {
	return func(o *Options) {
		o.Cmd.Init(cmd.Description(d))
	}
}

// Metadata associated with the service
func Metadata(md map[string]string) Option {
	return func(o *Options) {
		o.Server.Init(server.Metadata(md))
	}
}

func Flags(flags ...cli.Flag) Option {
	return func(o *Options) {
		o.Cmd.App().Flags = append(o.Cmd.App().Flags, flags...)
	}
}

func Action(a func(*cli.Context)) Option {
	return func(o *Options) {
		o.Cmd.App().Action = a
	}
}

func RegisterTTL(t time.Duration) Option {
	return func(o *Options) {
		o.Server.Init(server.RegisterTTL(t))
	}
}

func RegisterInterval(t time.Duration) Option {
	return func(o *Options) {
		o.RegisterInterval = t
	}
}

// WrapClient is a convenience method for wrapping a Client with
// some middleware component. A list of wrappers can be provided.
func WrapClient(w ...client.Wrapper) Option {
	return func(o *Options) {
		// apply in reverse
		for i := len(w); i > 0; i-- {
			o.Client = w[i-1](o.Client)
		}
	}
}

// WrapHandler adds a handler Wrapper to a list of options passed into the server
func WrapHandler(w ...server.HandlerWrapper) Option {
	return func(o *Options) {
		var wrappers []server.Option

		for _, wrap := range w {
			wrappers = append(wrappers, server.WrapHandler(wrap))
		}

		// Init once
		o.Server.Init(wrappers...)
	}
}

// WrapSubscriber adds a subscriber Wrapper to a list of options passed into the server
func WrapSubscriber(w ...server.SubscriberWrapper) Option {
	return func(o *Options) {
		var wrappers []server.Option

		for _, wrap := range w {
			wrappers = append(wrappers, server.WrapSubscriber(wrap))
		}

		// Init once
		o.Server.Init(wrappers...)
	}
}

// Before and Afters

func BeforeStart(fn func() error) Option {
	return func(o *Options) {
		o.BeforeStart = append(o.BeforeStart, fn)
	}
}

func AfterStop(fn func() error) Option {
	return func(o *Options) {
		o.AfterStop = append(o.AfterStop, fn)
	}
}
