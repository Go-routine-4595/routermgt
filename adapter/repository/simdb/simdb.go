package simdb

import (
	"github.com/Go-routine-4995/routermgt/domain"
	"sync"
)

type Simdb struct {
	tenantdbLock *sync.RWMutex
	tenantdb     map[string]map[string]domain.Router
}

func NewSimDB() *Simdb {
	return &Simdb{
		tenantdb:     make(map[string]map[string]domain.Router),
		tenantdbLock: &sync.RWMutex{},
	}
}

func (s *Simdb) GetRouter(router domain.Router, tenant string) (domain.Router, bool) {

	var (
		re domain.Router
		ok bool
	)

	s.tenantdbLock.RLock()
	defer s.tenantdbLock.RUnlock()

	re, ok = s.tenantdb[tenant][router.RouterSerial]

	return re, ok

}

// GetPaged return a pointer of a slice of routers, and the total number of page with the given limit.
func (s *Simdb) GetPaged(page domain.Pagination, tenant string) (*[]domain.Router, int) {
	var (
		re *[]domain.Router
		l  int
		p  int
		i  int
		r  int
	)
	// page.Page start at index 0 to ... ceil(l/page.Limit)
	s.tenantdbLock.RLock()
	defer s.tenantdbLock.RUnlock()

	l = len(s.tenantdb[tenant])
	p = l / page.Limit
	r = l % page.Limit
	if r != 0 {
		p++
	}

	re = new([]domain.Router)

	// we are out of range!
	if (page.Page * page.Limit) > l {
		return re, p - 1
	}

	// /!\ p and page.Page index are different p [1..n] page.Page [0..n-1] page.Page is 0 indexed
	// p indicates the number of page(s)
	// p == 1 && page.Limit > l we have only one page and the limit asked is bigger than the number of elements
	// p == page.Limit + 1 && page.Limit > l this is the last page and there is fewer elements than the limit asked
	if (p == 1 && page.Limit > l) || (p == page.Page+1 && page.Limit > l) {
		*re = make([]domain.Router, l)
	} else {
		*re = make([]domain.Router, page.Limit)
	}
	// /!\ we are using a map to simulate a DB,  Go spec says that the iteration order over maps is not specified.
	// That is to say, you should not expect the map keys to appear in any particular order.
	for _, v := range s.tenantdb[tenant] {
		if i >= (page.Page*page.Limit) && i < ((page.Page+1)*page.Limit) {
			(*re)[i-(page.Page*(page.Limit))] = v
		}
		i++
	}

	return re, p - 1
}

// Add a list of router and return a list of routers that are already in the DB
func (s *Simdb) Add(routers []domain.Router, tenant string) *[]domain.Router {

	var (
		re *[]domain.Router
		ok bool
	)

	s.tenantdbLock.Lock()
	defer s.tenantdbLock.Unlock()

	_, ok = s.tenantdb[tenant]
	if !ok {
		s.tenantdb[tenant] = make(map[string]domain.Router)
	}
	for _, v := range routers {
		_, ok = s.tenantdb[tenant][v.RouterSerial]
		if !ok {
			s.tenantdb[tenant][v.RouterSerial] = v
		} else {
			if re == nil {
				re = new([]domain.Router)
				*re = make([]domain.Router, 0)
			}
			*re = append(*re, v)
		}
	}
	return re
}

func (s *Simdb) Delete(routers []domain.Router, tenant string) {
	s.tenantdbLock.Lock()
	defer s.tenantdbLock.Unlock()

	for _, v := range routers {
		delete(s.tenantdb[tenant], v.RouterSerial)
	}
}
