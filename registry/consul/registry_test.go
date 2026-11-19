package consul

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/registry"
	"github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
)

func tcpServer(t *testing.T, lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			return
		}
		fmt.Println("get tcp")
		conn.Close()
	}
}

func TestRegister(t *testing.T) {
	addr := fmt.Sprintf("%s:8081", getIntranetIP())
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		t.Errorf("listen tcp %s failed!", addr)
		t.Fail()
	}
	defer lis.Close()
	go tcpServer(t, lis)
	time.Sleep(time.Millisecond * 100)
	cli, err := api.NewClient(&api.Config{Address: "127.0.0.1:8500"})
	if err != nil {
		t.Fatalf("create consul client failed: %v", err)
	}
	r := New(cli)
	assert.Nil(t, err)
	version := strconv.FormatInt(time.Now().Unix(), 10)
	svc := &registry.ServiceInstance{
		ID:        "test2233",
		Name:      "test-provider",
		Version:   version,
		Metadata:  map[string]string{"app": "kratos"},
		Endpoints: []string{fmt.Sprintf("tcp://%s?isSecure=false", addr)},
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err = r.Deregister(ctx, svc)
	assert.Nil(t, err)
	err = r.Register(ctx, svc)
	assert.Nil(t, err)
	w, err := r.Watch(ctx, "test-provider")
	assert.Nil(t, err)

	services, err := w.Next()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(services))
	assert.EqualValues(t, "test2233", services[0].ID)
	assert.EqualValues(t, "test-provider", services[0].Name)
	assert.EqualValues(t, version, services[0].Version)
}

func getIntranetIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}

// fakeAgent records the payload the consul client posts to the agent, so the
// registration body can be asserted without a live consul.
func fakeAgent(t *testing.T) (*httptest.Server, func() (string, map[string]interface{})) {
	t.Helper()
	var (
		mu   sync.Mutex
		path string
		body map[string]interface{}
	)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		parsed := map[string]interface{}{}
		_ = json.Unmarshal(raw, &parsed)
		mu.Lock()
		path, body = r.URL.Path, parsed
		mu.Unlock()
		_, _ = w.Write([]byte("true"))
	}))
	t.Cleanup(srv.Close)
	return srv, func() (string, map[string]interface{}) {
		mu.Lock()
		defer mu.Unlock()
		return path, body
	}
}

func newTestRegistry(t *testing.T, srv *httptest.Server, opts ...Option) *Registry {
	t.Helper()
	cli, err := api.NewClient(&api.Config{Address: srv.URL})
	if err != nil {
		t.Fatal(err)
	}
	return New(cli, opts...)
}

func TestRegisterAdvertisesConfiguredServiceName(t *testing.T) {
	srv, payload := fakeAgent(t)
	r := newTestRegistry(t, srv, WithServiceName("bullgpt-quote-canary"))

	svc := &registry.ServiceInstance{
		ID:        "bullgpt-quote_10_0_0_1",
		Name:      "bullgpt-quote",
		Version:   "1.0",
		Endpoints: []string{"http://10.0.0.1:8080"},
	}
	if err := r.Register(context.Background(), svc); err != nil {
		t.Fatal(err)
	}

	path, body := payload()
	if path != "/v1/agent/service/register" {
		t.Errorf("posted to %q, want the agent service register endpoint", path)
	}
	if got := body["Name"]; got != "bullgpt-quote-canary" {
		t.Errorf("registered name = %v, want bullgpt-quote-canary", got)
	}
	if got := body["ID"]; got != "bullgpt-quote_10_0_0_1" {
		t.Errorf("registered id = %v, want the caller's id untouched", got)
	}
	if svc.Name != "bullgpt-quote" {
		t.Errorf("caller's instance was mutated to %q; Deregister keys on it", svc.Name)
	}
}

func TestRegisterKeepsNameWithoutOverride(t *testing.T) {
	srv, payload := fakeAgent(t)
	r := newTestRegistry(t, srv)

	svc := &registry.ServiceInstance{
		ID:        "hello_10_0_0_1",
		Name:      "hello-world",
		Version:   "1.0",
		Endpoints: []string{"http://10.0.0.1:8080"},
	}
	if err := r.Register(context.Background(), svc); err != nil {
		t.Fatal(err)
	}

	if _, body := payload(); body["Name"] != "hello-world" {
		t.Errorf("registered name = %v, want the caller's name when no override is set", body["Name"])
	}
}
