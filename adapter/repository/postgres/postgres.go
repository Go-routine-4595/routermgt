package postgres

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/Go-routine-4995/routermgt/domain"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Rule struct {
	Id       int    `pg:",pk"`
	RuleID   string `json:"action" pg:"type:uuid""`
	RuleByte []byte `json:"condition"`
}

type Profile struct {
	domain.Router
	Id    int      `pg:",pk"`
	Rules []string `json:"rules"`
}

type Postgres struct {
	db         *pg.DB
	Address    string
	User       string
	Password   string
	Database   string
	ClientCert string
	ClientKey  string
	ServerCert string
	wg         *sync.WaitGroup
}

func NewPostgres(address string, user string, password string, database string, clCert string, clKey string, serCert string, wg *sync.WaitGroup) *Postgres {

	var (
		conf *tls.Config
		db   *pg.DB
		ctx  context.Context
		err  error
	)

	conf = ConfTLS(clCert, clKey, serCert)
	db = pg.Connect(&pg.Options{
		Addr:      address,
		User:      user,
		Password:  password,
		Database:  database,
		TLSConfig: conf,
	})

	ctx = context.Background()
	err = db.Ping(ctx)

	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	// trap SIGINT / SIGTERM to exit cleanly
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Shutting down DB...")
		_ = db.Close()
		fmt.Println("BD connection closed gracefully")
		wg.Done()
	}()

	return &Postgres{
		db:         db,
		Address:    address,
		User:       user,
		Password:   password,
		Database:   database,
		ClientCert: clCert,
		ClientKey:  clKey,
		ServerCert: serCert,
		wg:         wg,
	}
}

// Add a list of router and return a list of routers that are already in the DB
func (p *Postgres) Add(routes []domain.Router, tenant string) *[]domain.Router {
	var (
		err        error
		r          domain.Router
		res        orm.Result
		resRouters *[]domain.Router
	)

	for _, v := range routes {
		r = v
		res, err = p.db.Model(&r).
			OnConflict("DO NOTHING").
			Insert()
		if err != nil {
			fmt.Println(err)
		}
		if res.RowsAffected() <= 0 {
			if resRouters == nil {
				resRouters = new([]domain.Router)
				*resRouters = make([]domain.Router, 0)
			}
			*resRouters = append(*resRouters, v)
			fmt.Println("row already existing")
		}
	}

	return resRouters
}

// GetPaged return a pointer of a slice of routers, and the total number of page with the given limit.
func (p *Postgres) GetPaged(page domain.Pagination, tenant string) (*[]domain.Router, int) {
	var (
		routers   *[]domain.Router
		err       error
		count     int
		ps        int
		r         int
		fetchSize int
	)

	routers = new([]domain.Router)
	*routers = make([]domain.Router, 0)

	count, err = p.db.Model((*domain.Router)(nil)).Count()

	ps = count / page.Limit
	r = count % page.Limit
	if r != 0 {
		ps++
	}

	// we are out of range!
	if (page.Page * page.Limit) > count {
		return routers, ps - 1
	}

	// /!\ p and page.Page index are different p [1..n] page.Page [0..n-1] page.Page is 0 indexed
	// p indicates the number of page(s)
	// p == 1 && page.Limit > l we have only one page and the limit asked is bigger than the number of elements
	// p == page.Limit + 1 && page.Limit > l this is the last page and there is fewer elements than the limit asked
	if (ps == 1 && page.Limit > count) || (ps == page.Page+1 && page.Limit > count) {
		*routers = make([]domain.Router, count)
		fetchSize = count
	} else {
		*routers = make([]domain.Router, page.Limit)
		fetchSize = page.Limit
	}

	// we need to find the right page in the DB to do so we are fetching one row (sorted by router_serial) the last row of a page
	// if the last row of the page is smaller that the previous router_serial then we are in the right offset/page
	err = p.db.Model(routers).Limit(fetchSize).Offset(ps - 1).Select()
	if err != nil {
		fmt.Println(err)
	}
	//for i, k := range receiver {
	//	(*routers)[i] = k.Router
	//}

	return routers, ps - 1
}

func (p *Postgres) GetRouter(router domain.Router, tenant string) (domain.Router, bool) {
	var (
		res domain.Router
		err error
	)

	fmt.Println("GetRouter")
	err = p.db.Model(&res).Where("router_serial = ?", res.RouterSerial).Limit(1).Select()
	fmt.Printf("GetRouter returned: %+v \n", res)
	if err != nil {
		fmt.Println(err)
		return res, false
	}
	return res, true
}

func (p *Postgres) Delete(routers []domain.Router, tenant string) {
	var (
		router domain.Router
		err    error
	)

	for _, k := range routers {
		_, err = p.db.Model(&router).Where("router_serial =?", k.RouterSerial).Delete()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func ConfTLS(clientCert string, clientKey string, serverCert string) *tls.Config {
	cert, err := tls.LoadX509KeyPair(clientCert, clientKey)
	if err != nil {
		log.Println("failed to load client certificate: %v", err)
	}

	CACert, err := os.ReadFile(serverCert)
	if err != nil {
		log.Println("failed to load server certificate: %v", err)
	}

	CACertPool := x509.NewCertPool()
	CACertPool.AppendCertsFromPEM(CACert)

	return &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            CACertPool,
		InsecureSkipVerify: true,
	}
}
