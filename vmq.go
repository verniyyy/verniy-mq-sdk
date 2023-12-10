package vmq

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"io"
	"log"
)

const (
	_ = iota
	quit
	ping
	createQueue
	listQueue
	deleteQueue
	publish
	consume
	delete
)

// Ping ...
func Ping(s Session) error {
	r, err := s.request("", ping, nil)
	if err != nil {
		return err
	}
	buf := make([]byte, 64)
	if _, err := r.Read(buf); err != nil {
		log.Println(err)
		return err
	}
	log.Printf("%v\n", string(buf))

	// --

	return nil
}

// CreateQueue ...
func CreateQueue(s Session, queueName string) error {
	_, err := s.request(queueName, createQueue, nil)
	return err
}

// ListQueue ...
// func ListQueue(s Session, queueName string) ([]string, error) {
// 	r, err := s.request(queueName, listQueue, nil)
// 	if err != nil {
// 		return []string{}, err
// 	}

// 	buf, err := io.ReadAll(r)
// 	if err != nil {
// 		return []string{}, err
// 	}

// 	return []string{}, err
// }

// DeleteQueue ...
func DeleteQueue(s Session, queueName string) error {
	_, err := s.request(queueName, deleteQueue, nil)
	return err
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

// headerField ...
type headerField struct {
	SessionID SessionID
	Command   uint8
	QueueName [128]rune
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
