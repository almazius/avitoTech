package handlers

import (
	"avitoTech/config"
	models "avitoTech/internal"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
)

type services interface {
	CreateSegment(ctx context.Context, name string) error
	DeleteSegment(ctx context.Context, name string) error
	SubscribeUser(ctx context.Context, userId int, segmentName string, timeout int) error
}

func NewFiberServer(config *config.Config, serv services) *FiberServer {
	var fiberServer FiberServer
	fiberServer.app = fiber.New(fiber.Config{
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	})

	fiberServer.services = serv

	fiberServer.logger = log.New(os.Stderr, "Fiber: ", log.Lshortfile|log.LstdFlags)

	// create handlers for segments
	segment := fiberServer.app.Group("segment")
	segment.Post("create", fiberServer.createSegment)
	segment.Delete("delete", fiberServer.dellSegment)
	segment.Post("addUser", fiberServer.addUserOnSegment)
	segment.Put("addPercent", fiberServer.addPercentUsers)

	// create handlers for users
	users := fiberServer.app.Group("users")
	users.Get("getSegments/...", fiberServer.getUsersSegments)

	return &fiberServer
}

func (fs *FiberServer) StartServer(config *config.Config) error {
	err := fs.app.Listen(config.ServerPort)
	if err != nil {
		fs.logger.Print(err)
	}
	return err
}

type FiberServer struct {
	app    *fiber.App
	logger *log.Logger
	services
}

func (fs *FiberServer) createSegment(ctx *fiber.Ctx) error {
	ctx.Status(201)

	object := ctx.Body()
	segment := models.Segment{}
	err := json.Unmarshal(object, &segment)
	if err != nil {
		fs.logger.Print(err)
		respond := models.Respond{
			Status:  500,
			Message: "cant unmarshal json",
			Err:     err,
		}
		ctx.Status(respond.Status)
		return ctx.SendString(respond.Error())
	}

	err = ValidateStruct(segment)
	if err != nil {
		fs.logger.Print(err)
		respond := models.Respond{
			Status:  400,
			Message: "Uncorrected body",
			Err:     err,
		}
		ctx.Status(respond.Status)
		return ctx.SendString(respond.Error())
	}

	err = fs.services.CreateSegment(context.Background(), segment.NameSegment)
	if err != nil {
		ctx.Status(err.(*models.Respond).Status)
		return ctx.SendString(err.Error())
	}
	return nil
}

func (fs *FiberServer) dellSegment(ctx *fiber.Ctx) error {
	ctx.Status(204)

	segmentName := ctx.Query("segmentName", "none")
	if segmentName == "none" {
		respond := models.Respond{
			Status:  400,
			Message: "need fill segmentName query",
			Err:     errors.New("uncorrected segment name"),
		}
		ctx.Status(respond.Status)
		return ctx.SendString(respond.Error())
	}
	err := fs.services.DeleteSegment(context.Background(), segmentName)
	if err != nil {
		ctx.Status(err.(*models.Respond).Status)
		return ctx.SendString(err.Error())
	}
	return nil

}

func (fs *FiberServer) addUserOnSegment(ctx *fiber.Ctx) error {
	request := models.SubscribeRequest{}
	err := ctx.BodyParser(&request)
	if err != nil {
		fs.logger.Print(err)
		respond := models.Respond{
			Status:  400,
			Message: "",
			Err:     errors.New("uncorrected body"),
		}
		ctx.Status(respond.Status)
		return ctx.SendString(respond.Error())
	}

	err = ValidateStruct(request)
	if err != nil {
		fs.logger.Print(err)
		respond := models.Respond{
			Status:  400,
			Message: "",
			Err:     errors.New("uncorrected body"),
		}
		ctx.Status(respond.Status)
		return ctx.SendString(respond.Error())
	}

	err = fs.services.SubscribeUser(context.Background(),
		request.UserId, request.SegmentName, request.TimeoutHours)
	if err != nil {
		ctx.Status(err.(*models.Respond).Status)
		return ctx.SendString(err.Error())
	}
	return nil
}

func (fs *FiberServer) addPercentUsers(ctx *fiber.Ctx) error {
	return nil
}

func (fs *FiberServer) getUsersSegments(ctx *fiber.Ctx) error {
	return nil
}

func ValidateStruct(data interface{}) error {
	validate := validator.New()

	err := validate.Struct(data)
	if err != nil {
		return err
	}
	return nil
}
