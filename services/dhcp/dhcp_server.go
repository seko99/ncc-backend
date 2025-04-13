package dhcp

import (
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	"code.evixo.ru/ncc/ncc-backend/pkg/events/interfaces"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"github.com/gogf/gf/net/gipv4"
	"net"
	"sync"
)

type Server struct {
	cfg    *config.Config
	log    logger.Logger
	events interfaces.Events

	Sock     []*net.UDPConn
	Stopping bool

	bindingMap BindingMap
	poolMap    PoolMap

	handlers map[byte]func(*Packet) ([]byte, error)

	serverLeases ServerLeaseMap

	Wg      sync.WaitGroup
	allocWg sync.WaitGroup
	startup bool
}

func (ths *Server) Serve(conn *net.UDPConn) error {
	for {
		buf := make([]byte, 1024)
		n, client, e := conn.ReadFromUDP(buf)
		if e != nil {
			// TODO: Silently ignore?
			continue
		}
		go func(n int, client *net.UDPAddr, buf []byte) {
			ths.log.Trace("Request from", client.IP.String())

			pkt, err := NewPacket(buf)
			if err != nil {
				ths.log.Error("Can't parse packet: %v", err)
				return
			}

			ths.log.Debug("packet: type=%d mac=%s ip=%s", pkt.Options.Opt53.Type, byte2mac(pkt.Packet.ClientMAC), gipv4.Long2ip(pkt.Packet.Ciaddr))

			handler, ok := ths.handlers[pkt.Options.Opt53.Type]
			if ok {
				var reply []byte

				reply, err = handler(pkt)
				if err != nil {
					ths.log.Error("Can't handle packet: %v", err)

					replyPkt := Packet{}
					reply, err = replyPkt.Nak()
					if err != nil {
						ths.log.Error("Can't assemble NAK packet: %w", err)
					}
				}

				_, err = conn.WriteTo(reply, client)
				if err != nil {
					ths.log.Error("Can't send reply: %v", err)
				}
			} else {
				ths.log.Error("No handler for packet type=%d", pkt.Packet.MsgType)
			}

		}(n, client, buf)
	}
}

func (ths *Server) Listen(addr string) (*net.UDPConn, error) {
	udpAddr, e := net.ResolveUDPAddr("udp", addr)
	if e != nil {
		return nil, e
	}
	return net.ListenUDP("udp", udpAddr)
}

func (ths *Server) ListenAndServe(addr string) error {

	conn, err := ths.Listen(addr)
	if err != nil {
		return err
	}
	ths.Sock = append(ths.Sock, conn)

	err = ths.Serve(conn)
	if err != nil {
		return err
	}

	return nil
}

func (ths *Server) Start() error {
	ths.handlers[MsgTypeDiscover] = ths.discoverHandler
	ths.handlers[MsgTypeRequest] = ths.requestHandler
	ths.handlers[MsgTypeInform] = ths.informHandler
	ths.handlers[MsgTypeDecline] = ths.declineHandler

	go ths.configUpdater()

	var err error
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err = ths.ListenAndServe(ths.cfg.Dhcp.Listen)
		if err != nil {
			ths.log.Error("Can't listen: %v", err)
		}
		wg.Done()
	}()

	wg.Wait()

	return err
}

func NewDhcpServer(cfg *config.Config, log logger.Logger, events interfaces.Events) *Server {

	r := &Server{
		cfg:      cfg,
		log:      log,
		events:   events,
		handlers: map[byte]func(*Packet) ([]byte, error){},
		serverLeases: ServerLeaseMap{
			m: map[uint32]ServerLease{},
		},
		Wg:      sync.WaitGroup{},
		allocWg: sync.WaitGroup{},
	}
	r.Wg.Add(1)
	r.startup = true
	return r
}
