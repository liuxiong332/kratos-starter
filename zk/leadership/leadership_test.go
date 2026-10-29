package leadership

import (
	"context"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-zookeeper/zk"
)

// A ZooKeeper outage must not take the process down. Before this was fixed
// takeLeader called logger.Fatalf on a NewElection error, so an unreachable
// ZooKeeper called os.Exit(1) — the service died and restarted instead of
// retrying, and the restart also dropped its kafka consumer group membership.
//
// If the Fatalf comes back, this test does not fail politely: the test binary
// exits with status 1.
func TestTakeLeaderDoesNotExitOnZooKeeperOutage(t *testing.T) {
	// Nothing listens on port 1: every zk operation returns ErrNoServer.
	conn, _, err := zk.Connect([]string{"127.0.0.1:1"}, 200*time.Millisecond)
	if err != nil {
		t.Fatalf("zk.Connect: %v", err)
	}
	defer conn.Close()

	logger := log.NewHelper(log.DefaultLogger)
	returned := make(chan struct{})

	go func() {
		takeLeader(conn, "/unittest-leader", logger, func(ctx context.Context) error { return nil })
		close(returned)
	}()

	select {
	case <-returned:
	case <-time.After(10 * time.Second):
		t.Fatal("takeLeader did not return: it must return so TakeLeader can back off and retry")
	}
}
