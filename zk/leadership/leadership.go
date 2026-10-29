package leadership

import (
	"context"
	"time"

	"github.com/Comcast/go-leaderelection"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-zookeeper/zk"
)

// retryDelay bounds the leader-election retry loop. Without it a ZooKeeper
// outage would turn TakeLeader's `for {}` into a hot spin.
const retryDelay = 5 * time.Second

func takeLeader(zkConn *zk.Conn, leaderRootPath string, logger *log.Helper, onTakeLeadership func(ctx context.Context) error) {
	candidate, err := leaderelection.NewElection(zkConn, leaderRootPath, "electron")

	if err != nil {
		// Must not be Fatalf: a transient ZooKeeper i/o timeout used to call
		// os.Exit(1) here, which killed the whole process (indicator-calculator
		// restarted twice a day right after the CN open for exactly this), and a
		// restart also drops the kafka consumer group membership. The caller
		// retries, so return and let it do that.
		logger.Errorf("New election error: %v, will retry leader election in %v", err, retryDelay)
		return
	}

	go candidate.ElectLeader()

	var ctx context.Context
	var cancelFunc context.CancelFunc

	for {
		select {
		case status, ok := <-candidate.Status():
			if !ok {
				logger.Info("Channel closed, election is terminated! Will retry leader election.")
				candidate.Resign()
				if cancelFunc != nil {
					cancelFunc()
				}
				return
			}
			if status.Err != nil {
				logger.Infof("Received election status error: %v! Will retry leader election", status.Err)
				candidate.Resign()
				if cancelFunc != nil {
					cancelFunc()
				}
				return
			}

			logger.Infof("Candidate received status message: <%v>.", status)
			if status.Role == leaderelection.Leader {
				// doLeaderStuff(candidate, status, respCh, connFailCh, waitFor)
				logger.Info("Now this node is the leader")

				ctx, cancelFunc = context.WithCancel(context.Background())

				go onTakeLeadership(ctx)
			} else if cancelFunc != nil {
				logger.Info("Cancel leader runner")
				cancelFunc()
				cancelFunc = nil
			}
		}
	}
}

func TakeLeader(zkConn *zk.Conn, leaderRootPath string, logger *log.Helper, onTakeLeadership func(ctx context.Context) error) {
	if exists, err, _ := zkConn.Exists(leaderRootPath); !exists {
		// Create the election node in ZooKeeper. A failure here is the same
		// transient ZooKeeper outage as in takeLeader: log it and let the retry
		// loop below try again instead of exiting the process.
		if _, err := zkConn.Create(leaderRootPath, []byte(""), 0, zk.WorldACL(zk.PermAll)); err != nil {
			logger.Errorf("Error creating the election node %s: %v, will retry", leaderRootPath, err)
		}
	} else if err != nil {
		logger.Errorf("Error checking the election node %s: %v, will retry", leaderRootPath, err)
	}

	for {
		takeLeader(zkConn, leaderRootPath, logger, onTakeLeadership)
		time.Sleep(retryDelay)
	}
}
