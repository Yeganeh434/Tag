package nats

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"service1/internal/adapters/databases/mysql"
	"service1/internal/application/usecases"
	"service1/internal/config"
	"service1/internal/domain/service"

	"github.com/nats-io/nats.go"
)

type Tag struct {
	ID    uint64 `json:"id"`
	Title string `json:"title"`
}

func Subscribe() {
	tagRepo := mysql.NewMySQLTagRepository(mysql.TagDB.DB)
	tagService := service.NewTagService(tagRepo)
	tagManagementUseCase := usecases.NewTagManagementUseCase(tagService)

	ctx := context.Background()
	ctx, span := config.Tracer.Start(ctx, "RegisterApprovedTag_handler")
	defer span.End()

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Printf("error connecting nats:%v", err)
		return
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Printf("error connecting jetStream:%v", err)
		return
	}

	_, err = js.Subscribe("events.EditTag", func(msg *nats.Msg) {
		var data Tag
		err := json.Unmarshal(msg.Data, &data)
		if err != nil {
			log.Printf("error unmarshaling:%v", err)
			return
		}
		tag, err := tagManagementUseCase.DeleteTag(data.ID, ctx)
		if err != nil {
			if errors.Is(err, service.ErrNoTagExistsWithThisID) {
				log.Printf("error:%v", err)
				return
			}
			log.Printf("error deleting tag:%v", err)
			return
		}
		log.Println("tag deleted successfully!") //////////////////////////////////////////////

		tag.Title = data.Title
		err = tagManagementUseCase.RegisterTag(tag, ctx)
		if err != nil {
			log.Printf("error registering tag:%v", err)
			return
		}
		log.Println("tag registered successfully!") //////////////////////////////////////////////

		err = msg.Ack()
		if err != nil {
			log.Printf("error acknowledging message: %v", err)
			return
		}
	}, nats.Durable("tag_proj_consumer"))
	if err != nil {
		log.Printf("error subscribing: %v", err)
		return
	}

	select {}
}
