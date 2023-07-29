/**
 * @file    logging.go
 * @author  Christophe Buffard
 * @date    2023-06-17
 * @brief   DNSHelper Service for CoreDNS.
 *
 * License under GNU GENERAL PUBLIC LICENSE Version 3, 29 June 2007
 * service Logging
 */

package logging

//make init proto update tidy

import (
	"context"
	"fmt"
	"github.com/Go-routine-4995/routermgt/domain"
	"github.com/rs/zerolog"
	"os"
	"time"
)

type IService interface {
	AddRouters(ctx context.Context, routers []domain.Router, tenant string) *[]domain.Router
	GetPagedRouters(ctx context.Context, page domain.Pagination, tenant string) *[]domain.Router
	GetRouter(ctx context.Context, router domain.Router, tenant string) *domain.Router
	DeleteRouters(ctx context.Context, routers []domain.Router, tenant string)
}

type LoggingService struct {
	next IService
	log  zerolog.Logger
}

func NewLoggingService(n interface{}) IService {

	return &LoggingService{
		next: n.(IService),
		log:  zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).Level(zerolog.InfoLevel).With().Timestamp().Logger(),
	}
}

func (s *LoggingService) AddRouters(ctx context.Context, r []domain.Router, tenant string) (rep *[]domain.Router) {

	defer func(start time.Time) {
		var str string
		var sreq string
		if rep != nil {
			str = fmt.Sprintf("%v", rep)
		}
		if len(r) < 5 {
			sreq = fmt.Sprintf("%v", r)
		} else {
			sreq = fmt.Sprintf("request too large: %d routers being created", len(r))
		}
		s.log.Info().
			Str("method", "AddRouters").
			Str("request", sreq).
			Str("response", str).
			Str("tenant", tenant).
			Dur("took", time.Since(start)).Send()
	}(time.Now())

	return s.next.AddRouters(ctx, r, tenant)
}

func (s *LoggingService) DeleteRouters(ctx context.Context, r []domain.Router, tenant string) {

	defer func(start time.Time) {
		var sreq string
		if len(r) < 5 {
			sreq = fmt.Sprintf("%v", r)
		} else {
			sreq = fmt.Sprintf("request too large: %d routers being created", len(r))
		}
		s.log.Info().
			Str("method", "DeleteRouters").
			Str("request", sreq).
			Str("tenant", tenant).
			Dur("took", time.Since(start)).Send()
	}(time.Now())

	s.next.DeleteRouters(ctx, r, tenant)
}

func (s *LoggingService) GetPagedRouters(ctx context.Context, page domain.Pagination, tenant string) (rep *[]domain.Router) {

	defer func(start time.Time) {
		var str string
		var sreq string
		if rep != nil {
			if len(*rep) < 5 {
				str = fmt.Sprintf("%v", *rep)
			} else {
				str = fmt.Sprintf("%v", *rep)
			}
		}

		sreq = fmt.Sprintf("%v", page)

		s.log.Info().
			Str("method", "GetPagedRouters").
			Str("request", sreq).
			Str("response", str).
			Str("tenant", tenant).
			Dur("took", time.Since(start)).Send()
	}(time.Now())

	return s.next.GetPagedRouters(ctx, page, tenant)
}

func (s *LoggingService) GetRouter(ctx context.Context, r domain.Router, tenant string) (rep *domain.Router) {

	defer func(start time.Time) {
		var str string
		var sreq string

		str = fmt.Sprintf("%v", rep)

		sreq = fmt.Sprintf("%v", r)

		s.log.Info().
			Str("method", "GetRouters").
			Str("request", sreq).
			Str("response", str).
			Str("tenant", tenant).
			Dur("took", time.Since(start)).Send()
	}(time.Now())

	return s.next.GetRouter(ctx, r, tenant)
}
