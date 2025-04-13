package radius

import (
	"bytes"
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	"code.evixo.ru/ncc/ncc-backend/pkg/events"
	"code.evixo.ru/ncc/ncc-backend/pkg/events/interfaces"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net"
	"net/http"
	"regexp"
	"sync"
	"time"
)

type CustomersMap struct {
	sync.RWMutex
	m map[string]models.CustomerData
}

type AccessItem struct {
	State      int    `json:"state"`
	Limits     Limits `json:"limits"`
	HasLease   bool   `json:"has_lease"`
	CustomerId string `json:"customer_id"`
	NasId      string `json:"nas_id"`
	ServiceId  string `json:"service_id"`
}

type AccessMap struct {
	sync.RWMutex
	m map[string]AccessItem
}

type RadiusServer struct {
	cfg       *config.Config
	log       logger.Logger
	Handler   *Handler
	Sock      []*net.UDPConn
	Stopping  bool
	ids       map[uint8]bool
	handlers  map[string]func(io.Writer, *Packet)
	events    interfaces.Events
	accessMap AccessMap
	nasMap    NasMap
	Wg        sync.WaitGroup
	startup   bool
	strRegexp *regexp.Regexp
}

func (s *RadiusServer) ListenAndServe(addr string) {

	conn, e := s.Listen(addr)
	if e != nil {
		panic(e)
	}
	s.Sock = append(s.Sock, conn)
	if e := s.Serve(conn); e != nil {
		if s.Stopping {
			// Ignore close errors
			return
		}
		panic(e)
	}
}

func (s *RadiusServer) HandleFunc(code PacketCode, statusType int, handler func(io.Writer, *Packet)) {
	key := fmt.Sprintf("%d-%d", code, statusType)
	if _, inuse := s.handlers[key]; inuse {
		panic(fmt.Errorf("DevErr: HandleFunc-add for already assigned code=%d", code))
	}
	s.handlers[key] = handler
}

func (s *RadiusServer) Listen(addr string) (*net.UDPConn, error) {
	udpAddr, e := net.ResolveUDPAddr("udp", addr)
	if e != nil {
		return nil, e
	}
	return net.ListenUDP("udp", udpAddr)
}

func (s *RadiusServer) checkNas(ip string) (string, error) {
	//todo: logics
	s.nasMap.RLock()
	defer s.nasMap.RUnlock()

	for _, n := range s.nasMap.m {
		if n.Ip == ip {
			return n.Secret, nil
		}
	}

	return "", fmt.Errorf("unknown NAS: %s", ip)
}

func (s *RadiusServer) getNasAttrs(ip string) ([]NasAttr, error) {
	s.nasMap.RLock()
	nas, ok := s.nasMap.m[ip]
	s.nasMap.RUnlock()
	if !ok {
		return nil, fmt.Errorf("NAS not found: %s", ip)
	}

	nasAttrs := []NasAttr{}

	for _, a := range nas.NasType.NasAttributes {
		nasAttrs = append(nasAttrs, NasAttr{
			Attr:   a.Attribute.Name,
			Val:    a.Val,
			Code:   uint8(a.Attribute.Code),
			Vendor: uint32(a.Attribute.Vendor.Code),
		})
	}

	return nasAttrs, nil
}

func (s *RadiusServer) Serve(conn *net.UDPConn) error {
	for {
		buf := make([]byte, 1024)
		n, client, e := conn.ReadFromUDP(buf)
		if e != nil {
			// TODO: Silently ignore?
			continue
		}
		go func(n int, client *net.UDPAddr, buf []byte) {
			s.log.Trace("Request from", client.IP.String())

			secret, err := s.checkNas(client.IP.String())
			if err != nil {
				s.log.Error("Request dropped: %v", err)
				return
			}

			s.log.Trace("Packet from allowed IP=", client.IP.String(), " secret=", secret)

			s.log.Trace("raw.recv:", buf[:n])
			p, e := decode(buf, n, secret)
			if e != nil {
				// TODO: Silently ignore decode?
				//return e
				return
			}
			if !validate(p) {
				// TODO: Silently ignore invalidate package?
				//return fmt.Errorf("Invalid MessageAuthenticator")
				s.log.Warn("Invalid MessageAuthenticator")
				return
			}

			/*			processed, ok := s.ids[p.Identifier]
						if ok && processed {
							s.log.Warn("Duplicate packet id=%d", p.Identifier)
							return
						}

						s.ids[p.Identifier] = true
			*/
			statusType := uint32(0)
			if p.HasAttr(AcctStatusType) {
				attr := p.Attr(AcctStatusType)
				statusType = binary.BigEndian.Uint32(attr)
			}

			key := fmt.Sprintf("%d-%d", p.Code, statusType)
			handle, ok := s.handlers[key]
			if ok {
				readBuf := new(bytes.Buffer)
				handle(readBuf, p)
				if len(readBuf.Bytes()) != 0 {
					// Only send a packet if we got anything
					s.log.Debug("raw.send: %+v", readBuf.Bytes())
					if _, e := conn.WriteTo(readBuf.Bytes(), client); e != nil {
						// TODO: ignore clients that gone away?
						panic(e)
					}
				}
				readBuf.Reset()
			} else {
				s.log.Debug("Drop packet with code=%d, statusType=%d", p.Code, statusType)
			}

		}(n, client, buf)
	}
}

func (s *RadiusServer) getAccess(login string) *AccessItem {
	defer s.accessMap.RUnlock()
	s.accessMap.RLock()
	access, ok := s.accessMap.m[login]
	if !ok {
		return nil
	}
	return &access
}

func (s *RadiusServer) configUpdater() {
	for {
		s.UpdateConfig()
		time.Sleep(s.cfg.Radius.Update)
	}
}

func (s *RadiusServer) UpdateConfig() {
	err := s.events.PublishRequest(events.Event{
		Type: ConfigRequest,
	}, ConfigResponse, func(event events.Event, params ...interface{}) {
		b, _ := json.Marshal(event.Payload["access"])
		s.accessMap.Lock()
		s.accessMap.m = map[string]AccessItem{}
		err := json.Unmarshal(b, &s.accessMap.m)
		if err != nil {
			s.log.Error("Can't unmarshal access map: %v", err)
		}
		mapLen := len(s.accessMap.m)
		s.accessMap.Unlock()
		s.log.Info("Received access map: %d", mapLen)

		b, _ = json.Marshal(event.Payload["nases"])
		s.nasMap.Lock()
		s.nasMap.m = map[string]models.NasData{}
		err = json.Unmarshal(b, &s.nasMap.m)
		if err != nil {
			s.log.Error("Can't unmarshal NAS map: %v", err)
		}
		mapLen = len(s.nasMap.m)
		s.nasMap.Unlock()
		s.log.Info("Received NAS map: %d", mapLen)

		if s.startup {
			if mapLen > 0 {
				s.startup = false
				s.Wg.Done()
			}
		}
	})
	if err != nil {
		s.log.Error("can't get config: %v", err)
	}
}

func (s *RadiusServer) Start() error {
	s.HandleFunc(AccessRequest, 0, s.Auth)
	s.HandleFunc(AccountingRequest, 1, s.AcctBegin)
	s.HandleFunc(AccountingRequest, 3, s.AcctUpdate)
	s.HandleFunc(AccountingRequest, 2, s.AcctStop)
	s.HandleFunc(AccountingRequest, 7, s.AcctOn)
	s.HandleFunc(AccountingRequest, 8, s.AcctOff)

	go s.configUpdater()

	go s.ListenAndServe(s.cfg.Radius.Auth.Listen)
	go s.ListenAndServe(s.cfg.Radius.Acct.Listen)

	s.log.Info("RADIUS server started")

	if s.cfg.Radius.Control.Enabled {
		g := gin.Default()

		srv := http.Server{
			Addr:    s.cfg.Radius.Control.Listen,
			Handler: g,
		}

		g.POST("/ctrl/stop", func(context *gin.Context) {
			err := srv.Shutdown(context)
			if err != nil {
				s.log.Error("%v", err)
			}
			return
		})

		err := srv.ListenAndServe()
		if err != nil {
			return err
		}
	}

	select {}
}

func NewRadiusServer(cfg *config.Config, log logger.Logger, events interfaces.Events) *RadiusServer {
	r := &RadiusServer{
		cfg:       cfg,
		log:       log,
		Handler:   NewHanlder(events),
		ids:       map[uint8]bool{},
		handlers:  make(map[string]func(io.Writer, *Packet)),
		events:    events,
		Wg:        sync.WaitGroup{},
		strRegexp: regexp.MustCompile("^\\d+$"),
	}
	r.Wg.Add(1)
	r.startup = true
	return r
}
