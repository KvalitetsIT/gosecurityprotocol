package securityprotocol

import "time"
import "go.uber.org/zap"


type Maintainable interface {
	MaintainCache() error
}



func StartMaintenance(m Maintainable, d time.Duration, logger *zap.SugaredLogger) {

	ticker := time.NewTicker(d)
	for _ = range ticker.C {
		logger.Info("Running cache maintenance")
                err := m.MaintainCache()
                if (err != nil) {
                	logger.Errorf("Error running cache maintenance: %s", err.Error())
                }
	}
}
