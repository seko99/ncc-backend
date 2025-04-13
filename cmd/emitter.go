package cmd

import (
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	"code.evixo.ru/ncc/ncc-backend/pkg/events"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger/zero"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"code.evixo.ru/ncc/ncc-backend/services/radius"
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"time"
)

const (
	emitterPrefix = "/v1"
)

var emitterCmd = &cobra.Command{
	Use:   "emitter",
	Short: "Event emitter",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.NewConfig()
		if err != nil {
			panic(err)
		}
		log := zero.NewLogger()

		log.Info("Starting Event emitter")

		storage := psqlstorage.NewStorage(cfg, log)
		err = storage.Connect()
		if err != nil {
			panic(fmt.Sprintf("can't connect to storage: %v", err))
		}

		id := uuid.NewString()
		log.Info("Radius serverId: %s", id)
		radiusEvents, err := events.NewEvents(cfg, log, id, radius.Queue)

		for i := 0; i < 5; i++ {
			radiusEvents.PublishRequest(events.Event{
				Type: radius.ConfigRequest,
			}, radius.ConfigResponse, func(event events.Event, params ...interface{}) {
				fmt.Printf("i=%d\n", i)
			})
		}

		/*		radiusEvents.SubscribeOnRequest(radius.ReloadRequest, radius.ReloadResponse, func(event events.Event) map[string]interface{} {
					shared.Log.Infof("ReloadRequest")
					return map[string]interface{}{
						"success": true,
					}
				})
		*/
		go func() {

			/*			err = radiusEvents.PublishRequest(events.Event{
							Type: radius.ConfigRequest,
						}, radius.ConfigResponse, func(event events.Event) {
							shared.Log.Infof("ConfigResponse: %d", len(event.Payload))
						})
						if err != nil {
							shared.Log.Errorf("Request error: %+v", err)
						}

						err = radiusEvents.PublishRequest(events.Event{
							Type: radius.SessionsRequest,
						}, radius.SessionsResponse, func(event events.Event) {
							var sessions []models.SessionData

							b, _ := json.Marshal(event.Payload["sessions"])
							json.Unmarshal(b, &sessions)

							shared.Log.Infof("SessionsResponse: %d sessions", len(sessions))
						})
						if err != nil {
							shared.Log.Errorf("Request error: %+v", err)
						}
			*/
			time.Sleep(3 * time.Second)
		}()

		for {
			time.Sleep(cfg.Watcher.Delay)

			/*			err = radiusEvents.PublishEvent(events.Event{
							Type:    radius.InterimUpdate,
							Payload: map[string]interface{}{},
						})
						if err != nil {
							shared.Log.Errorf("Can't publish event: %v", err)
						}

						err = radiusEvents.PublishRequest(events.Event{
							Type: radius.AccessRequestType,
							Payload: map[string]interface{}{
								"request": radius.AccessRequestEvent{
									Login: "andersondawn",
								},
							},
						}, radius.AccessResponseType, func(event events.Event) {
							var response radius.AccessResponseEvent

							b, _ := json.Marshal(event.Payload["response"])
							json.Unmarshal(b, &response)

							shared.Log.Infof("AccessResponse: %v", response)
						})
						if err != nil {
							shared.Log.Errorf("Request error: %+v", err)
						}
			*/
		}
	},
}
