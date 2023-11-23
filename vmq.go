package vmq

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
	"net"
)

// Session ...
type Session interface {
	Close() error
	request(qName string, cmd uint8, msg any) (io.Reader, error)
}

// Options ...
type Options struct {
	Addr     string
	UserID   string
	Password string
}

// NewSession ...
func NewSession(opt *Options) (Session, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", opt.Addr)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}

	return &session{
		conn: conn,
	}, nil
}

const (
	ping = iota
	createQueue
	deleteQueue
	publish
	consume
	delete
)

// session ...
type session struct {
	conn *net.TCPConn
}

// Close ...
func (s session) Close() error {
	return s.conn.Close()
}

// Ping ...
func Ping(s Session) error {
	r, err := s.request("", ping, nil)
	if err != nil {
		return err
	}

	res, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	if len(res) == 0 {
		return fmt.Errorf("ping failed")
	}

	return nil
}

// Publish ...
func Publish(s Session, queueName string, msg any) error {
	_, err := s.request(queueName, publish, msg)
	return err
}

// Consume ...
func Consume[T any](s Session, queueName string) (T, error) {
	r, err := s.request(queueName, consume, nil)
	if err != nil {
		var v T
		return v, err
	}

	buf, err := io.ReadAll(r)
	if err != nil {
		var v T
		return v, err
	}

	var v T
	if err := binary.Read(bytes.NewReader(buf), binary.BigEndian, &v); err != nil {
		var v T
		return v, err
	}

	return v, nil
}

// request ...
func (c session) request(qName string, cmd uint8, msg any) (io.Reader, error) {
	hf := headerField{
		Command: cmd,
	}
	copy(hf.QueueName[:], []rune(qName))

	hr, err := hf.encode()
	if err != nil {
		return nil, err
	}

	mr, err := encode(msg)
	if err != nil {
		return nil, err
	}

	req, err := io.ReadAll(io.MultiReader(hr, mr))
	if err != nil {
		return nil, err
	}

	if _, err := c.conn.Write(req); err != nil {
		return nil, err
	}

	return c.conn, nil
}

// headerField ...
type headerField struct {
	AccountID [32]rune
	QueueName [128]rune
	Command   uint8
}

// encode ...
func (hf headerField) encode() (io.Reader, error) {
	buf := new(bytes.Buffer)
	binary.Write(
		buf,
		binary.BigEndian,
		hf,
	)
	return buf, nil
}

// encode ...
func encode(data any) (io.Reader, error) {
	buf := new(bytes.Buffer)
	if data == nil {
		return buf, nil
	}

	enc := gob.NewEncoder(buf)
	if err := enc.Encode(data); err != nil {
		return nil, err
	}

	return buf, nil
}
