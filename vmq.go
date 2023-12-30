package vmq

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
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

	// --

	return nil
}

// CreateQueue ...
func CreateQueue(s Session, queueName string) error {
	_, err := s.request(queueName, createQueue, nil)
	return err
}

// ListQueue ...
func ListQueue(s Session) ([]string, error) {
	r, err := s.request("", listQueue, nil)
	if err != nil {
		return []string{}, err
	}

	buf := make([]byte, r.header.DataSize)
	if _, err := r.Read(buf); err != nil {
		return []string{}, err
	}

	type ListQueuesOutput struct {
		Queues []string `json:"queues"`
	}
	data, err := decodeJSON[ListQueuesOutput](buf)

	return data.Queues, err
}

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
func Consume[T any](s Session, queueName string) (Message[T], error) {
	r, err := s.request(queueName, consume, nil)
	if err != nil {
		return Message[T]{}, err
	}

	var messageID MessageID
	if _, err := r.Read(messageID[:]); err != nil {
		return Message[T]{}, err
	}

	buf := make([]byte, r.header.DataSize)
	if _, err := r.Read(buf); err != nil {
		return Message[T]{}, err
	}

	data, err := decodeGob[T](buf)
	if err != nil {
		return Message[T]{}, err
	}

	return Message[T]{ID: messageID, Data: data}, nil
}

// Delete ...
func Delete(s Session, queueName string, messageID MessageID) error {
	_, err := s.request(queueName, delete, messageID[:])
	return err
}

// MessageID ...
type MessageID [26]byte

// String ...
func (m MessageID) String() string {
	return string(m[:])
}

// Message ...
type Message[T any] struct {
	ID   MessageID
	Data T
}

// headerField ...
type headerField struct {
	SessionID SessionID
	Command   uint8
	QueueName [128]rune
	DataSize  uint64
}

// encode ...
func (hf headerField) encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	binary.Write(
		buf,
		binary.BigEndian,
		hf,
	)
	return buf.Bytes(), nil
}

// encode ...
func encode(data any) ([]byte, error) {
	switch typedData := data.(type) {
	case []byte:
		return typedData, nil
	}

	buf := new(bytes.Buffer)
	if data == nil {
		return []byte{}, nil
	}

	enc := gob.NewEncoder(buf)
	if err := enc.Encode(data); err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

func decodeGob[T any](data []byte) (T, error) {
	buf := bytes.NewBuffer(data)

	var v T
	if err := gob.NewDecoder(buf).Decode(&v); err != nil {
		return *new(T), err
	}

	return v, nil
}

func decodeJSON[T any](data []byte) (T, error) {
	buf := bytes.NewBuffer(data)

	var v T
	if err := json.NewDecoder(buf).Decode(&v); err != nil {
		return *new(T), err
	}

	return v, nil
}
