package usecase

import (
	"context"
	"log"
	"os"
	"time"
)

type repository interface {
	CreateSegment(ctx context.Context, name string) error
	DeleteSegment(ctx context.Context, name string) error
	UnSubscribeFromSegment(ctx context.Context, name string) error
	SubscribeUserOnSegment(ctx context.Context, userId int, segmentName string) error
	SubscribeUserOnSegmentWithTimeout(ctx context.Context, userId int, segmentName string, timeout time.Time) error
	UnsubscribeUser(ctx context.Context, userId int, segmentName string) error
}

type avitoService struct {
	logger *log.Logger
	repository
}

func NewService(rep repository) *avitoService {
	return &avitoService{
		logger:     log.New(os.Stderr, "Servece: ", log.Lshortfile|log.LstdFlags),
		repository: rep,
	}
}

func (as *avitoService) CreateSegment(ctx context.Context, name string) error {
	return as.repository.CreateSegment(ctx, name)
}

func (as *avitoService) DeleteSegment(ctx context.Context, name string) error {
	err := as.repository.UnSubscribeFromSegment(ctx, name)
	if err != nil {
		return err
	}
	err = as.repository.DeleteSegment(ctx, name)
	if err != nil {
		return err
	}

	return nil
}

func (as *avitoService) SubscribeUser(ctx context.Context, userId int, segmentName string, timeout int) error {
	if timeout != 0 {
		// ticker
		ticker := time.NewTicker(time.Duration(timeout * 1000000000))
		go as.ttl(ctx, userId, segmentName, ticker)

		return as.repository.SubscribeUserOnSegmentWithTimeout(ctx,
			userId, segmentName, time.Now().Add(time.Duration(timeout*3600000000000)))

	} else {
		return as.repository.SubscribeUserOnSegment(ctx, userId, segmentName)
	}
}

func (as *avitoService) ttl(ctx context.Context, userId int, segmentName string, ticker *time.Ticker) {
	for {
		select {
		case <-ticker.C:
			ticker.Stop()
			err := as.repository.UnsubscribeUser(ctx, userId, segmentName)
			if err != nil {
				as.logger.Println("cant dell segment on time\n" + err.Error())
			}
			as.logger.Printf("User %d unsubscribe from %s", userId, segmentName)
			break
		}
	}
}
