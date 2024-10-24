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

type CustomError struct{
	Code int `json:"code"`
	Error string `json:"error"`
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

	// js, err := nc.JetStream()
	// if err != nil {
	// 	log.Printf("error connecting jetStream:%v", err)
	// 	return
	// }

	_,err = nc.Subscribe("EditTag", func(msg *nats.Msg) {  //jetstream: "events.EditTag"
		var customError CustomError
		var data Tag
		err := json.Unmarshal(msg.Data, &data)
		if err != nil {
			customError.Code=400
			response,errState:=json.Marshal(customError)
			if errState!=nil{
				log.Printf("error marshaling:%v",errState)
				return
			}
			errState=msg.Respond(response)
			if errState!=nil{
				log.Printf("error sending response:%v",errState)
				return
			}
			log.Printf("error unmarshaling:%v", err)
			return
		}
		tag, err := tagManagementUseCase.DeleteTag(data.ID, ctx)
		if err != nil {
			if errors.Is(err, service.ErrNoTagExistsWithThisID) {
				customError.Error="no tag exist with this ID"
			}
			customError.Code=400
			response,errState:=json.Marshal(customError)
			if errState!=nil{
				log.Printf("error marshaling:%v",errState)
				return
			}
			errState=msg.Respond(response) 
			if errState!=nil{
				log.Printf("error sending response:%v",errState)
				return
			}
			log.Printf("error deleting tag:%v", err)
			return
		}

		tag.Title = data.Title
		err = tagManagementUseCase.RegisterTag(tag, ctx)
		if err != nil {
			customError.Code=400
			response,errState:=json.Marshal(customError)
			if errState!=nil{
				log.Printf("error marshaling:%v",errState)
				return
			}
			errState=msg.Respond(response)
			if errState!=nil{
				log.Printf("error sending response:%v",errState)
				return
			}
			log.Printf("error registering tag:%v", err)
			return
		}

		customError.Code=200
		response,errState:=json.Marshal(customError)
		if errState!=nil{
			log.Printf("error marshaling:%v",errState)
			return
		}
		errState=msg.Respond(response)
		if errState!=nil{
			log.Printf("error sending response:%v",errState)
			return
		}

		//jetstream:
		// err = msg.Ack()
		// if err != nil {
		// 	log.Printf("error acknowledging message: %v", err)
		// 	return
		// }
	}) //jetstream: nats.Durable("tag_proj_consumer")
	if err != nil {
		log.Printf("error subscribing: %v", err)
		return
	}

	select {}
}
