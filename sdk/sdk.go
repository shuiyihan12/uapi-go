// Package sdk is the single import surface for using uapi-go as a Go
// library instead of deploying the REST gateway.
//
// New assembles a Client from functional options; the per-domain accessors
// lazily create and cache the typed SOAP services:
//
//	c, err := sdk.New(
//		sdk.WithEndpoint("https://apac.universal-api.travelport.com/B2BGateway/connect/uAPI"),
//		sdk.WithLogger(logging.Noop()),
//	)
//	defer c.Close()
//
//	hotel, err := c.Hotel()
//	out, err := hotel.SearchAvailability(ctx, req)
//
// Credentials are request-level, mirroring the gateway: pass the Travelport
// Authorization value and the region through pkg/requestctx on the context
// of each call (see ExampleNew).
//
// The SDK shares the gateway's version tag (v0.WSDL.PATCH): the module
// version and the container image version are always the same.
package sdk

import (
	"time"

	"github.com/shuiyihan12/uapi-go/pkg/client"
	"github.com/shuiyihan12/uapi-go/pkg/logging"
	"github.com/shuiyihan12/uapi-go/pkg/manager"
	airservice "github.com/shuiyihan12/uapi-go/pkg/services/air"
	gdsqueueservice "github.com/shuiyihan12/uapi-go/pkg/services/gdsQueue"
	hotelservice "github.com/shuiyihan12/uapi-go/pkg/services/hotel"
	railservice "github.com/shuiyihan12/uapi-go/pkg/services/rail"
	sharedbookingservice "github.com/shuiyihan12/uapi-go/pkg/services/sharedBooking"
	shareduprofileservice "github.com/shuiyihan12/uapi-go/pkg/services/sharedUprofile"
	systemservice "github.com/shuiyihan12/uapi-go/pkg/services/system"
	terminalservice "github.com/shuiyihan12/uapi-go/pkg/services/terminal"
	universalservice "github.com/shuiyihan12/uapi-go/pkg/services/universal"
	uprofileservice "github.com/shuiyihan12/uapi-go/pkg/services/uprofile"
	utilservice "github.com/shuiyihan12/uapi-go/pkg/services/util"
	vehicleservice "github.com/shuiyihan12/uapi-go/pkg/services/vehicle"
)

// Client is the SDK entry point holding the shared service manager. All
// domain accessors are lazy and cached; Close releases every created
// service and its connections.
type Client struct {
	mgr *manager.ServiceManager
}

// Option mutates the SDK configuration (seeded from
// manager.DefaultServiceConfig).
type Option func(*manager.ServiceConfig)

// New creates a Client from the default configuration plus the given
// options. Construction is offline; only actual SOAP calls touch the
// network.
func New(opts ...Option) (*Client, error) {
	cfg := manager.DefaultServiceConfig()
	for _, opt := range opts {
		opt(&cfg)
	}
	mgr, err := manager.NewServiceManager(cfg)
	if err != nil {
		return nil, err
	}
	return &Client{mgr: mgr}, nil
}

// WithEndpoint overrides the default endpoint prefix, used when a request
// carries no region.
func WithEndpoint(endpoint string) Option {
	return func(c *manager.ServiceConfig) { c.BaseEndpoint = endpoint }
}

// WithEnvironment sets "test" or "production"; it only affects the built-in
// logger shape when no logger is injected.
func WithEnvironment(env string) Option {
	return func(c *manager.ServiceConfig) {
		c.Environment = env
		c.IsDevelopment = env != "production"
	}
}

// WithTimeouts sets the connect, response-header and per-request total
// timeouts.
func WithTimeouts(connect, readHeader, total time.Duration) Option {
	return func(c *manager.ServiceConfig) {
		c.ConnectionTimeout = connect
		c.ReadTimeout = readHeader
		c.RequestTimeout = total
	}
}

// WithKeepAlivePool sizes the warm keep-alive connection pool; per-host is
// capped by global.
func WithKeepAlivePool(global, perHost int) Option {
	return func(c *manager.ServiceConfig) {
		c.MaxIdleConns = global
		c.MaxIdleConnsPerHost = perHost
	}
}

// WithLogLevel sets the built-in logger's level ("debug", "info", "warn",
// "error"); it only matters when no logger is injected via WithLogger.
func WithLogLevel(level string) Option {
	return func(c *manager.ServiceConfig) { c.LogLevel = level }
}

// WithLogger injects a logger (pkg/logging). nil keeps the built-in zap
// logger derived from the log level; logging.Noop() silences all output.
func WithLogger(l logging.Logger) Option {
	return func(c *manager.ServiceConfig) { c.Logger = l }
}

// WithMetrics injects a per-call observability hook (client.Metrics); nil
// keeps the no-op default.
func WithMetrics(m client.Metrics) Option {
	return func(c *manager.ServiceConfig) { c.Metrics = m }
}

// WithTLSSkipVerify skips upstream TLS certificate verification — only for
// private environments with self-signed certificates.
func WithTLSSkipVerify(skip bool) Option {
	return func(c *manager.ServiceConfig) { c.SkipTLSVerify = skip }
}

// Air returns the cached Air service (lazily created).
func (c *Client) Air() (*airservice.AirService, error) {
	return manager.Get[*airservice.AirService](c.mgr, "air")
}

// GdsQueue returns the cached GdsQueue service (lazily created).
func (c *Client) GdsQueue() (*gdsqueueservice.GdsQueueService, error) {
	return manager.Get[*gdsqueueservice.GdsQueueService](c.mgr, "gdsQueue")
}

// Hotel returns the cached Hotel service (lazily created).
func (c *Client) Hotel() (*hotelservice.HotelService, error) {
	return manager.Get[*hotelservice.HotelService](c.mgr, "hotel")
}

// Rail returns the cached Rail service (lazily created).
func (c *Client) Rail() (*railservice.RailService, error) {
	return manager.Get[*railservice.RailService](c.mgr, "rail")
}

// SharedBooking returns the cached SharedBooking service (lazily created).
func (c *Client) SharedBooking() (*sharedbookingservice.SharedBookingService, error) {
	return manager.Get[*sharedbookingservice.SharedBookingService](c.mgr, "sharedBooking")
}

// SharedUprofile returns the cached SharedUprofile service (lazily created).
func (c *Client) SharedUprofile() (*shareduprofileservice.SharedUprofileService, error) {
	return manager.Get[*shareduprofileservice.SharedUprofileService](c.mgr, "sharedUprofile")
}

// System returns the cached System service (lazily created).
func (c *Client) System() (*systemservice.SystemService, error) {
	return manager.Get[*systemservice.SystemService](c.mgr, "system")
}

// Terminal returns the cached Terminal service (lazily created).
func (c *Client) Terminal() (*terminalservice.TerminalService, error) {
	return manager.Get[*terminalservice.TerminalService](c.mgr, "terminal")
}

// Universal returns the cached UniversalRecord service (lazily created).
func (c *Client) Universal() (*universalservice.UniversalService, error) {
	return manager.Get[*universalservice.UniversalService](c.mgr, "universal")
}

// Uprofile returns the cached Uprofile service (lazily created).
func (c *Client) Uprofile() (*uprofileservice.UprofileService, error) {
	return manager.Get[*uprofileservice.UprofileService](c.mgr, "uprofile")
}

// Util returns the cached Util service (lazily created).
func (c *Client) Util() (*utilservice.UtilService, error) {
	return manager.Get[*utilservice.UtilService](c.mgr, "util")
}

// Vehicle returns the cached Vehicle service (lazily created).
func (c *Client) Vehicle() (*vehicleservice.VehicleService, error) {
	return manager.Get[*vehicleservice.VehicleService](c.mgr, "vehicle")
}

// Close releases every created service and its connections.
func (c *Client) Close() error {
	return c.mgr.Close()
}
