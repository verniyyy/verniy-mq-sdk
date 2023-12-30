package test

import (
	"testing"

	vmq "github.com/verniyyy/verniy-mq-sdk"
)

func TestMQ(t *testing.T) {
	// example message type
	type User struct {
		ID   int
		Name string
	}

	sess, err := vmq.NewSession(SessionOptionsExample)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := sess.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	t.Run("CreateQueue", func(t *testing.T) {
		if err := vmq.CreateQueue(sess, "example queue"); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Publish", func(t *testing.T) {
		msg1 := User{ID: 1, Name: "Jhon"}
		if err := vmq.Publish(sess, "example queue", msg1); err != nil {
			t.Fatal(err)
		}
		msg2 := User{ID: 2, Name: "Harry"}
		if err := vmq.Publish(sess, "example queue", msg2); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Consume and Delete", func(t *testing.T) {
		msg1, err := vmq.Consume[User](sess, "example queue")
		if err != nil {
			t.Fatal(err)
		}

		msg2, err := vmq.Consume[User](sess, "example queue")
		if err != nil {
			t.Fatal(err)
		}

		if err := vmq.Delete(sess, "example queue", msg1.ID); err != nil {
			t.Fatal(err)
		}

		if err := vmq.Delete(sess, "example queue", msg2.ID); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("DeleteQueue", func(t *testing.T) {
		if err := vmq.DeleteQueue(sess, "example queue"); err != nil {
			t.Fatal(err)
		}
	})
}
