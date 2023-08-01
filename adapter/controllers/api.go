package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Go-routine-4995/routermgt/domain"
	"github.com/nats-io/nats.go"
)

const (
	messageGet = iota
	messageGetPaged
	messageCreate
	messageDelete

	queue = "worker_group_router"
)

type message struct {
	Mtype int    `json:"mtype"`
	Data  []byte `json:"Data"`
}

type IService interface {
	AddRouters(ctx context.Context, routers []domain.Router, tenant string) *[]domain.Router
	GetPagedRouters(ctx context.Context, page domain.Pagination, tenant string) (*[]domain.Router, int)
	GetRouter(ctx context.Context, router domain.Router, tenant string) *domain.Router
	DeleteRouters(ctx context.Context, routers []domain.Router, tenant string)
}

type ApiServer struct {
	ctx       context.Context
	urlBroker string
	subject   string
	con       *nats.Conn
	next      IService
}

func NewApiService(svc interface{}, u string, s string) *ApiServer {

	c, err := connect(u)
	if err != nil {
		fmt.Println("Broker connection error: ", err)
	}
	return &ApiServer{
		ctx:       context.Background(),
		urlBroker: u,
		con:       c,
		subject:   s,
		next:      svc.(IService),
	}
}

func connect(u string) (*nats.Conn, error) {
	nc, err := nats.Connect(u)
	return nc, err
}

func (a *ApiServer) AddRouters(routers []domain.Router, tenant string) *[]domain.Router {
	return a.next.AddRouters(a.ctx, routers, tenant)
}

func (a *ApiServer) GetRouters(routers domain.Router, tenant string) *domain.Router {
	return a.next.GetRouter(a.ctx, routers, tenant)
}

func (a *ApiServer) GetPagedRouters(page domain.Pagination, tenant string) (*[]domain.Router, int) {
	return a.next.GetPagedRouters(a.ctx, page, tenant)
}

func (a *ApiServer) DeleteRouters(routers []domain.Router, tenant string) {
	a.next.DeleteRouters(a.ctx, routers, tenant)
}

func (a *ApiServer) Start() {
	fmt.Println(" subscribing to: ", a.subject)
	_, err := a.con.QueueSubscribe(a.subject, queue, func(msg *nats.Msg) {
		var (
			err error
			b   []byte
			m   message
		)
		//fmt.Println("Create: ", string(msg.Data))
		err = json.Unmarshal(msg.Data, &m)
		if err != nil {
			fmt.Println("error unmarshalling: ", err)
			//return
		}
		// Call the right action here

		switch m.Mtype {
		case messageCreate:
			b, err = a.createCB(m.Data, "test")
		case messageGet:
			b, err = a.getCB(m.Data, "test")
		case messageGetPaged:
			b, err = a.getPagedCB(m.Data, "test")
		case messageDelete:
			b, err = a.deleteCB(m.Data, "test")
		}

		err = msg.Respond(b)
		a.con.Flush()

	})

	//TODO manage reconnection strategy
	if err != nil {
		fmt.Println("error while subscribing to subject", err)
		return
	}

	for {
	}
}

func (a *ApiServer) createCB(in []byte, tenant string) ([]byte, error) {
	var (
		routers []domain.Router
		ret     *[]domain.Router
		out     []byte
		err     error
	)
	err = json.Unmarshal(in, &routers)
	if err != nil {
		fmt.Println("error unmarshalling: ", err)
		return out, err
	}
	ret = a.AddRouters(routers, tenant)
	if ret != nil {
		out, err = json.Marshal(ret)
		if err != nil {
			fmt.Println("err unmarshalling answer: ", err)
		}
		return out, nil
	}
	return []byte("sucess!"), nil
}

func (a *ApiServer) getCB(in []byte, tenant string) ([]byte, error) {
	var (
		router domain.Router
		ret    *domain.Router
		out    []byte
		err    error
	)
	err = json.Unmarshal(in, &router)
	if err != nil {
		fmt.Println("error unmarshalling: ", err)
		return out, err
	}

	ret = a.GetRouters(router, tenant)

	if ret != nil {
		out, err = json.Marshal(ret)
		if err != nil {
			fmt.Println("err unmarshalling answer: ", err)
			return out, err
		}
		return out, nil
	}
	return []byte(""), nil
}

func (a *ApiServer) getPagedCB(in []byte, tenant string) ([]byte, error) {
	var (
		page     domain.Pagination
		out      []byte
		err      error
		response struct {
			Last    int              `json:"last"`
			Routers *[]domain.Router `json:"routers"`
		}
	)
	err = json.Unmarshal(in, &page)
	if err != nil {
		fmt.Println("error unmarshalling: ", err)
		return out, err
	}

	response.Routers, response.Last = a.GetPagedRouters(page, tenant)

	out, err = json.Marshal(response)
	if err != nil {
		fmt.Println("err unmarshalling answer: ", err)
		return out, err
	}

	return out, nil

}

/*
	{
	  "_metadata":
	  {
	      "page": 5,
	      "per_page": 20,
	      "page_count": 20,
	      "total_count": 521,
	      "Links": [
	        {"self": "/products?page=5&per_page=20"},
	        {"first": "/products?page=0&per_page=20"},
	        {"previous": "/products?page=4&per_page=20"},
	        {"next": "/products?page=6&per_page=20"},
	        {"last": "/products?page=26&per_page=20"},
	      ]
	  },
	  "records": [
	    {
	      "id": 1,
	      "name": "Widget #1",
	      "uri": "/products/1"
	    },
	    {
	      "id": 2,
	      "name": "Widget #2",
	      "uri": "/products/2"
	    },
	    {
	      "id": 3,
	      "name": "Widget #3",
	      "uri": "/products/3"
	    }
	  ]
	}
*/
func (a *ApiServer) deleteCB(in []byte, tenant string) ([]byte, error) {
	var (
		routers []domain.Router
		out     []byte
		err     error
	)
	err = json.Unmarshal(in, &routers)
	if err != nil {
		fmt.Println("error unmarshalling: ", err)
		return out, err
	}
	a.DeleteRouters(routers, tenant)

	return []byte("sucess!"), nil
}
