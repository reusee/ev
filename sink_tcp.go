package ev

import (
	"context"
	"net"
	"time"

	"github.com/reusee/sb"
)

type TCPClient struct {
	ch  chan PutOp
	ctx context.Context
}

var _ Sink = new(TCPClient)

func NewTCPClient(ctx context.Context, remote string) (*TCPClient, error) {
	dialer := &net.Dialer{
		Timeout: time.Second * 32,
	}
	dialCtx, cancel := context.WithTimeout(ctx, time.Second*32)
	defer cancel()
	conn, err := dialer.DialContext(dialCtx, "net", remote)
	if err != nil {
		return nil, we(err)
	}

	ch := make(chan PutOp)

	go func() {
		for {
			select {

			case op := <-ch:
				err := sb.Copy(
					sb.Marshal(op.Ev),
					sb.Encode(conn),
				)
				select {
				case op.Err <- we(err):
				case <-ctx.Done():
					return
				}

			case <-ctx.Done():
				conn.Close()
				return
			}
		}
	}()

	return &TCPClient{
		ch:  ch,
		ctx: ctx,
	}, nil
}

func (t *TCPClient) Put(ev *Ev) error {
	op := NewPutOp(ev)
	select {
	case t.ch <- op:
	case <-t.ctx.Done():
		return nil
	}
	select {
	case err := <-op.Err:
		return we(err)
	case <-t.ctx.Done():
		return nil
	}
}

func NewTCPServer(ctx context.Context, addr string, upstream Sink, errCh chan error) error {
	config := net.ListenConfig{}
	ln, err := config.Listen(ctx, "net", addr)
	if err != nil {
		return we(err)
	}

	go func() {
		<-ctx.Done()
		ln.Close()
	}()

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				select {
				case errCh <- we(err):
				case <-ctx.Done():
				default:
				}
				return
			}

			go func() {
				defer conn.Close()

				for {
					var ev *Ev
					if err := sb.Copy(
						sb.Decode(conn),
						sb.Unmarshal(&ev),
					); err != nil {
						select {
						case errCh <- we(err):
						case <-ctx.Done():
						default:
						}
						return
					}

					if err := upstream.Put(ev); err != nil {
						select {
						case errCh <- we(err):
						case <-ctx.Done():
						default:
						}
						return
					}

				}
			}()

		}
	}()

	return nil
}
