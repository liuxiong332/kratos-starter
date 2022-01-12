package leadership

import (
	"context"

	"github.com/Comcast/go-leaderelection"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-zookeeper/zk"
)

func takeLeader(zkConn *zk.Conn, leaderRootPath string, logger *log.Helper, onTakeLeadership func(ctx context.Context) error) {
	candidate, err := leaderelection.NewElection(zkConn, leaderRootPath, "electron")

	if err != nil {
		logger.Fatalf("New election error: %v", err)
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
	if exists, _, _ := zkConn.Exists(leaderRootPath); !exists {
		// Create the election node in ZooKeeper
		_, err := zkConn.Create(leaderRootPath, []byte(""), 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			logger.Fatalf("Error creating the election node: %v", err)
		}
	}

	for {
		takeLeader(zkConn, leaderRootPath, logger, onTakeLeadership)
	}
}
