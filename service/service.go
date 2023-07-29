package service

import (
	"context"
	"github.com/Go-routine-4995/routermgt/domain"
)

type IRepository interface {
	Add(routes []domain.Router, tenant string) *[]domain.Router
	GetPaged(page domain.Pagination, tenant string) *[]domain.Router
	GetRouter(router domain.Router, tenant string) (domain.Router, bool)
	Delete(router []domain.Router, tenant string)
}

type IService interface {
}

type Service struct {
	rep IRepository
}

func NewService(r interface{}) IService {
	return &Service{
		rep: r.(IRepository),
	}
}

func (s *Service) AddRouters(ctx context.Context, routers []domain.Router, tenant string) *[]domain.Router {
	return s.rep.Add(routers, tenant)
}

func (s *Service) GetPagedRouters(ctx context.Context, page domain.Pagination, tenant string) *[]domain.Router {
	return s.rep.GetPaged(page, tenant)
}

func (s *Service) DeleteRouters(ctx context.Context, routers []domain.Router, tenant string) {
	s.rep.Delete(routers, tenant)
}

func (s *Service) GetRouter(ctx context.Context, router domain.Router, tenant string) *domain.Router {
	var (
		re     domain.Router
		status bool
	)
	re, status = s.rep.GetRouter(router, tenant)
	if status {
		return &re
	} else {
		return nil
	}

}
