package orders

import (
	"github.com/nats-io/stan.go"
)

type ConfigNast struct {
	ClusterID   string
	ClientID    string
	Subject     string
	QGroup      string
	DurableName string
}

func Connect(cfg ConfigNast, cb stan.MsgHandler) error {
	sc, err := stan.Connect(cfg.ClusterID, cfg.ClientID)
	if err != nil {
		return err
	}

	_, err = sc.QueueSubscribe(cfg.Subject, cfg.QGroup, cb, stan.DurableName(cfg.DurableName))
	if err != nil {
		return err
	}

	return nil
}
