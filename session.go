package vmq

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
	"strings"

	"github.com/oklog/ulid/v2"
)

// Config ...
type Config struct {
	Addr     string
	UserID   string
	Password string
}

// Session ...
type Session interface {
	ID() string
	Close() error
	Ping() error
	CreateQueue(queueName string) error
	ListQueue() ([]string, error)
	DeleteQueue(queueName string) error
	Consume(queueName string) (Message[any], error)
	Delete(queueName string, messageID MessageID) error
	Publish(queueName string, msg any) error
	request(qName string, cmd uint8, msg any) (*response, error)
}

// NewSession ...
func NewSession(cfg *Config) (Session, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", cfg.Addr)
	if err != nil {
		return nil, err
	}

	c, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}

	s := &session{
		conn: c,
	}
	af := authField{}
	copy(af.AccountID[:], []rune(cfg.UserID))
	copy(af.Password[:], []rune(cfg.Password))

	ab, err := af.encode()
	if err != nil {
		return nil, err
	}
	if _, err := s.conn.Write(ab); err != nil {
		return nil, err
	}

	if err := binary.Read(s.conn, binary.BigEndian, &s.id); err != nil {
		if err := s.conn.Close(); err != nil {
			log.Fatal(err)
		}
		return nil, errors.New("authentication failed1")
	}
	if _, err := ulid.Parse(s.ID()); err != nil {
		if err := s.conn.Close(); err != nil {
			log.Fatal(err)
		}
		return nil, errors.New("authentication failed2")
	}

	return s, nil
}

// session ...
type session struct {
	id   SessionID
	conn *net.TCPConn
}

// SessionID ...
type SessionID [32]rune

// ID ...
func (s session) ID() string {
	return strings.Trim(string(s.id[:]), "\x00")
}

// Close ...
func (s session) Close() error {
	return s.conn.Close()
}

// Ping ...
func (s session) Ping() error {
	return Ping(s)
}

// CreateQueue ...
func (s session) CreateQueue(queueName string) error {
	return CreateQueue(s, queueName)
}

// ListQueue ...
func (s session) ListQueue() ([]string, error) {
	return ListQueue(s)
}

// DeleteQueue ...
func (s session) DeleteQueue(queueName string) error {
	return DeleteQueue(s, queueName)
}

// Consume ...
func (s session) Consume(queueName string) (Message[any], error) {
	return Consume[any](s, queueName)
}

// Delete ...
func (s session) Delete(queueName string, messageID MessageID) error {
	return Delete(s, queueName, messageID)
}

// Publish
func (s session) Publish(queueName string, msg any) error {
	return Publish(s, queueName, msg)
}

// request ...
func (c session) request(qName string, cmd uint8, data any) (*response, error) {
	r := bufio.NewReader(c.conn)
	w := bufio.NewWriter(c.conn)

	hf := headerField{
		SessionID: c.id,
		Command:   cmd,
	}
	copy(hf.QueueName[:], []rune(qName))

	mr, err := encode(data)
	if err != nil {
		return nil, err
	}

	hf.DataSize = uint64(len(mr))
	hr, err := hf.encode()
	if err != nil {
		return nil, err
	}

	req := append(hr, mr...)
	if _, err := w.Write(req); err != nil {
		return nil, err
	}
	w.Flush()

	resHeaderField, err := read[resHeaderField](r, resHeaderFieldSize)
	if err != nil {
		return nil, err
	}
	if resHeaderField.Result != ResOK {
		errBuf := make([]byte, resHeaderField.DataSize)
		if _, err := r.Read(errBuf); err != nil {
			return nil, err
		}
		return nil, errors.New(string(errBuf))
	}

	return &response{header: resHeaderField, Reader: r}, nil
}

// read ...
func read[T any](r io.Reader, bufSize uint64) (T, error) {
	var v T
	buf := make([]byte, bufSize)
	_, err := r.Read(buf)
	if err != nil {
		return v, err
	}

	if err := binary.Read(bytes.NewReader(buf), binary.BigEndian, &v); err != nil {
		return v, err
	}

	return v, nil
}

// authField ...
type authField struct {
	AccountID [32]rune
	Password  [64]rune
}

// encode ...
func (af authField) encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	binary.Write(
		buf,
		binary.BigEndian,
		af,
	)
	return buf.Bytes(), nil
}

const (
	_ uint8 = iota
	ResOK
	ResError
)

const resHeaderFieldSize = 1 + 8

// resHeaderField ...
type resHeaderField struct {
	Result   uint8
	DataSize uint64
}

// response
type response struct {
	header resHeaderField
	io.Reader
}
