package transport

import (
	"context"
	"fmt"

	"github.com/ghettovoice/gosip/core"
	"github.com/ghettovoice/gosip/log"
)

type tlsProtocol struct {
	protocol
	listeners   ListenerPool
	connections ConnectionPool
	conns       chan Connection
}

func NewTlsProtocol(ctx context.Context, output chan<- *IncomingMessage, errs chan<- error) Protocol {
	tls := new(tlsProtocol)
	tls.network = "tls"
	tls.reliable = true
	tls.streamed = true
	tls.conns = make(chan Connection)
	tls.listeners = NewListenerPool(ctx, tls.conns, errs)
	tls.connections = NewConnectionPool(ctx, output, errs)
	tls.SetLog(log.StandardLogger())
	// start up pools
	go tls.listeners.Manage()
	go tls.connections.Manage()
	go tls.manage(ctx)

	return tls
}

func (tls *tlsProtocol) SetLog(logger log.Logger) {
	tls.protocol.SetLog(logger)
	tls.listeners.SetLog(tls.Log())
	tls.connections.SetLog(tls.Log())
}

// piping new connections to connection pool for serving
func (tls *tlsProtocol) manage(ctx context.Context) {
	defer func() {
		tls.Log().Debugf("stop %s managing", tls)
		tls.dispose()
	}()
	tls.Log().Debugf("start %s managing", tls)

	for {
		select {
		case <-ctx.Done():
			return
		case conn := <-tls.conns:
			if err := tls.connections.Add(conn.RemoteAddr(), conn, socketTtl); err != nil {
				// TODO should it be passed up to UA?
				tls.Log().Errorf("%s failed to add new %s to %s: %s", tls, conn, tls.connections, err)
				continue
			}
		}
	}
}

func (tls *tlsProtocol) dispose() {
	tls.Log().Debugf("dispose %s", tls)
	close(tls.conns)
}

func (tls *tlsProtocol) Listen(target *Target) error {
	return fmt.Errorf("not implemented method in %s", tls)
	//target = FillTargetHostAndPort(tls.Network(), target)
	//network := strings.ToLower(tls.Network())
	//// resolve local TCP endpoint
	//laddr, err := tls.resolveTarget(target)
	//if err != nil {
	//	return err
	//}
	//// create listener
	//listener, err := net.ListenTCP(network, laddr)
	//if err != nil {
	//	return &ProtocolError{
	//		fmt.Errorf("failed to listen address %s: %s", laddr, err),
	//		fmt.Sprintf("create %s listener", tls.Network()),
	//		tls,
	//	}
	//}
	//// index listeners by local address
	//tls.listeners.Add(listener.Addr(), listener)
	//
	//return err // should be nil here
}

func (tls *tlsProtocol) Send(target *Target, msg core.Message) error {
	return fmt.Errorf("not implemented method in %s", tls)
	//target = FillTargetHostAndPort(tls.Network(), target)
	//
	//tls.Log().Infof("sending message '%s' to %s", msg.Short(), target.Addr())
	//tls.Log().Debugf("sending message '%s' to %s:\r\n%s", msg.Short(), target.Addr(), msg)
	//
	//// validate remote address
	//if target.Host == "" || target.Host == DefaultHost {
	//	return &ProtocolError{
	//		fmt.Errorf("invalid remote host resolved %s", target.Host),
	//		"resolve destination address",
	//		tls,
	//	}
	//}
	//// resolve remote address
	//raddr, err := tls.resolveTarget(target)
	//if err != nil {
	//	return err
	//}
	//// find or create connection
	//conn, err := tls.getOrCreateConnection(raddr)
	//if err != nil {
	//	return err
	//}
	//// send message
	//_, err = conn.Write([]byte(msg.String()))
	//
	//return err
}

//func (tls *tlsProtocol) resolveTarget(target *Target) (*net.TCPAddr, error) {
//	addr := target.Addr()
//	network := strings.ToLower(tls.String())
//	// resolve remote address
//	raddr, err := net.ResolveTCPAddr(network, addr)
//	if err != nil {
//		return nil, &ProtocolError{
//			fmt.Errorf("failed to resolve address %s: %s", addr, err),
//			fmt.Sprintf("resolve %s address", addr),
//			tls,
//		}
//	}
//
//	return raddr, nil
//}
//
//func (tls *tlsProtocol) getOrCreateConnection(raddr *net.TCPAddr) (Connection, error) {
//	network := strings.ToLower(tls.String())
//	laddr := &net.TCPAddr{
//		IP:   net.IP(DefaultHost),
//		Port: int(DefaultUdpPort),
//		Zone: "",
//	}
//
//	conn, ok := tls.connections.Get(raddr)
//	if !ok {
//		tls.Log().Debugf("connection for address %s not found; create a new one", raddr)
//		tcpConn, err := net.DialTCP(network, laddr, raddr)
//		if err != nil {
//			return nil, &ProtocolError{
//				fmt.Errorf("failed to create connection to remote address %s: %s", raddr, err),
//				fmt.Sprintf("create %s connection", tls.Network()),
//				tls,
//			}
//		}
//
//		conn = NewConnection(tcpConn)
//		conn.SetLog(tls.Log())
//		tls.connections.Add(conn.RemoteAddr(), conn, socketTtl)
//	}
//
//	return conn, nil
//}