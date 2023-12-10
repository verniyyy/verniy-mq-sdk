package vmq

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"github.com/oklog/ulid/v2"
)

// Session ...
type Session interface {
	ID() string
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

	c, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}

	s := &session{
		conn: c,
	}
	af := authField{}
	copy(af.AccountID[:], []rune(opt.UserID))
	copy(af.Password[:], []rune(opt.Password))

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

// request ...
func (c session) request(qName string, cmd uint8, msg any) (io.Reader, error) {
	r := bufio.NewReader(c.conn)
	w := bufio.NewWriter(c.conn)

	hf := headerField{
		SessionID: c.id,
		Command:   cmd,
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

	if _, err := w.Write(req); err != nil {
		return nil, err
	}
	w.Flush()

	log.Println("requested message")

	resHeaderField, err := read[resHeaderField](r, resHeaderFieldSize)
	if err != nil {
		return nil, err
	}
	log.Printf("resHeaderField.Result: %v\n", resHeaderField.Result)
	if resHeaderField.Result != ResOK {
		return nil, errors.New("response error")
	}

	return r, nil
}

// read ...
func read[T any](r io.Reader, bufSize int) (T, error) {
	var v T
	buf := make([]byte, bufSize)
	received, err := r.Read(buf)
	if err != nil {
		return v, err
	}

	if err := binary.Read(bytes.NewReader(buf), binary.BigEndian, &v); err != nil {
		log.Printf("received: %v\n", received)
		return v, err
	}

	return v, nil
}

// readAll ...
func readAll(r io.Reader) ([]byte, error) {
	// make a temporary bytes var to read from the connection
	tmp := make([]byte, 1024)
	// make 0 length data bytes (since we'll be appending)
	data := make([]byte, 0)
	// keep track of full length read
	length := 0

	// loop through the connection stream, appending tmp to data
	for {
		// read to the tmp var
		n, err := r.Read(tmp)
		if err != nil {
			// log if not normal error
			if err != io.EOF {
				fmt.Printf("Read error - %s\n", err)
			}
			break
		}

		// append read data to full data
		data = append(data, tmp[:n]...)

		// update total read var
		length += n
	}

	return data, nil
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

const resHeaderFieldSize = 1 + 4

// resHeaderField ...
type resHeaderField struct {
	Result   uint8
	DataSize uint32
}
