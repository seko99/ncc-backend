package dhcp

import "fmt"

func (ths *Server) discoverHandler(pkt *Packet) ([]byte, error) {
	replyPkt := Packet{}
	var reply []byte

	lease, err := ths.serverLeases.GetByPacket(pkt)
	if err != nil {
		newLease, err := ths.newLease(pkt)
		if err != nil {
			return nil, fmt.Errorf("can't create lease: %w", err)
		}

		reply, err = replyPkt.Offer(newLease)
		if err != nil {
			return nil, fmt.Errorf("can't assemble OFFER packet: %w", err)
		}
	} else {
		renewLease, err := ths.renewLease(*lease)

		reply, err = replyPkt.Offer(renewLease)
		if err != nil {
			return nil, fmt.Errorf("can't assemble OFFER packet: %w", err)
		}
	}

	return reply, nil
}

func (ths *Server) requestHandler(pkt *Packet) ([]byte, error) {
	replyPkt := Packet{}
	var reply []byte

	lease, err := ths.serverLeases.GetByPacket(pkt)
	if err != nil {
		reply, e := replyPkt.Nak()
		if e != nil {
			return nil, fmt.Errorf("can't assemble NAK packet: %w", e)
		}
		return reply, fmt.Errorf("can't process request: %w", err)
	}

	var acceptedLease *ServerLease

	if lease.Status == LeaseStatusAllocated {
		acceptedLease, err = ths.acceptLease(*lease)
		if err != nil {
			return nil, fmt.Errorf("can't accept lease: %w", err)
		}
	} else {
		acceptedLease, err = ths.renewLease(*lease)
		if err != nil {
			return nil, fmt.Errorf("can't renew lease: %w", err)
		}
	}

	reply, err = replyPkt.Ack(acceptedLease)
	if err != nil {
		return nil, fmt.Errorf("can't assemble ACK packet: %w", err)
	}

	return reply, nil
}

func (ths *Server) informHandler(pkt *Packet) ([]byte, error) {
	return ths.requestHandler(pkt)
}

func (ths *Server) declineHandler(pkt *Packet) ([]byte, error) {
	replyPkt := Packet{}
	var reply []byte

	lease, err := ths.serverLeases.GetByPacket(pkt)
	if err != nil {
		reply, e := replyPkt.Nak()
		if e != nil {
			return nil, fmt.Errorf("can't assemble NAK packet: %w", e)
		}
		return reply, fmt.Errorf("can't process request: %w", err)
	}

	err = ths.removeLease(*lease)
	if err != nil {
		reply, e := replyPkt.Nak()
		if e != nil {
			return nil, fmt.Errorf("can't assemble NAK packet: %w", e)
		}
		return reply, fmt.Errorf("can't process decline: %w", err)
	}

	reply, err = replyPkt.Ack(lease)
	if err != nil {
		return nil, fmt.Errorf("can't assemble ACK packet: %w", err)
	}

	return reply, nil
}
