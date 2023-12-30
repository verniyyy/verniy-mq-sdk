package vmq

// QueueSessionConfig ...
type QueueSessionConfig struct {
	Addr      string
	QueueName string
	UserID    string
	Password  string
}

// QueueSession ...
type QueueSession[MESSAGE_TYPE any] interface {
	ID() string
	Close() error
	Ping() error
	Consume() (Message[MESSAGE_TYPE], error)
	Delete(messageID MessageID) error
	Publish(msg MESSAGE_TYPE) error
}

// NewQueueSession ...
func NewQueueSession[MESSAGE_TYPE any](cfg *QueueSessionConfig) (QueueSession[MESSAGE_TYPE], error) {
	sess, err := NewSession(&Config{
		Addr:     cfg.Addr,
		UserID:   cfg.UserID,
		Password: cfg.Password,
	})
	if err != nil {
		return nil, err
	}

	return queueSession[MESSAGE_TYPE]{
		session:   sess.(*session),
		queueName: cfg.QueueName,
	}, nil
}

// queueSession ...
type queueSession[T any] struct {
	*session

	queueName string
}

// Consume ...
func (s queueSession[T]) Consume() (Message[T], error) {
	return Consume[T](s.session, s.queueName)
}

// Delete ...
func (s queueSession[T]) Delete(messageID MessageID) error {
	return Delete(s.session, s.queueName, messageID)
}

// Publish ...
func (s queueSession[T]) Publish(msg T) error {
	return Publish[T](s.session, s.queueName, msg)
}
