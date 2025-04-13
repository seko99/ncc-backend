package radius

import (
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	"code.evixo.ru/ncc/ncc-backend/pkg/events/interfaces/mocks"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger/zero"
	"github.com/golang/mock/gomock"
	"testing"
	"time"
)

func TestRadiusServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Radius: config.RadiusConfig{
			Auth: config.RadiusServerConfig{
				Listen: "0.0.0.0:1812",
			},
			Acct: config.RadiusServerConfig{
				Listen: "0.0.0.0:1812",
			},
			Update: time.Minute,
			Secret: "secret",
			Watcher: config.RadiusWatcherConfig{
				Start:   time.Second,
				Stop:    time.Second,
				Interim: 5 * time.Second,
			},
		},
	}

	log := zero.NewLogger()

	events := mocks.NewMockEvents(ctrl)
	events.EXPECT().PublishRequest(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil)

	srv := NewRadiusServer(cfg, log, events)
	srv.UpdateConfig()
}
