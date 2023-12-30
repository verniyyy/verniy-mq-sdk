package test

import (
	"fmt"
	"log"

	vmq "github.com/verniyyy/verniy-mq-sdk"
)

type User struct {
	ID   int
	Name string
}

func ExampleSession() {
	sess, err := vmq.NewSession(SessionConfigExample)
	if err != nil {
		return
	}
	defer func() {
		if err := sess.Close(); err != nil {
			panic(err)
		}
	}()

	// Output:
	//
}

func ExamplePing() {
	sess, err := vmq.NewSession(SessionConfigExample)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		if err := sess.Close(); err != nil {
			panic(err)
		}
	}()

	if err := vmq.Ping(sess); err != nil {
		log.Println(err)
	}

	// Output:
	//
}

// ExampleCreateQueue ...
func ExampleCreateQueue() {
	sess, err := vmq.NewSession(SessionConfigExample)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		if err := sess.Close(); err != nil {
			panic(err)
		}
	}()

	if err := vmq.CreateQueue(sess, "example queue"); err != nil {
		log.Println(err)
	}

	// Output:
	//
}

// ExampleListQueue ...
func ExampleListQueue() {
	sess, err := vmq.NewSession(SessionConfigExample)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		if err := sess.Close(); err != nil {
			panic(err)
		}
	}()

	if err := vmq.CreateQueue(sess, "example queue1"); err != nil {
		log.Println(err)
	}
	if err := vmq.CreateQueue(sess, "example queue2"); err != nil {
		log.Println(err)
	}
	if err := vmq.CreateQueue(sess, "example queue3"); err != nil {
		log.Println(err)
	}
	qs, err := vmq.ListQueue(sess)
	if err != nil {
		log.Println(err)
	}
	for i, q := range qs {
		fmt.Printf("%d: %v\n", i, q)
	}

	// Output:
	// 0: example queue1
	// 1: example queue2
	// 2: example queue3
}

// ExampleDeleteQueue ...
func ExampleDeleteQueue() {
	sess, err := vmq.NewSession(SessionConfigExample)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		if err := sess.Close(); err != nil {
			panic(err)
		}
	}()

	if err := vmq.DeleteQueue(sess, "example queue"); err != nil {
		log.Println(err)
	}

	// Output:
	//
}

// ExamplePublish ...
func ExamplePublish() {
	sess, err := vmq.NewSession(SessionConfigExample)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		if err := sess.Close(); err != nil {
			panic(err)
		}
	}()

	msg := User{ID: 2, Name: "Jhon"}

	if err := vmq.Publish(sess, "example queue", msg); err != nil {
		log.Println(err)
	}

	// Output:
	//
}

// ExampleConsume ...
func ExampleConsume() {
	sess, err := vmq.NewSession(SessionConfigExample)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		if err := sess.Close(); err != nil {
			panic(err)
		}
	}()

	msg, err := vmq.Consume[User](sess, "example queue")
	if err != nil {
		log.Println(err)
	}

	fmt.Printf("msg: %+v\n", msg)

	// Output:
	// msg: {ID:01HJXMBDXK97A8M0V8Y05DG3C1 Data:{ID:1 Name:Jhon}}
}

// ExampleDelete ...
func ExampleDelete() {
	sess, err := vmq.NewSession(SessionConfigExample)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		if err := sess.Close(); err != nil {
			panic(err)
		}
	}()

	if err := vmq.Delete(sess, "example queue", vmq.MessageID([]byte("01HJXVRZ9QJSCAHEKFNAZHTFMS"))); err != nil {
		log.Println(err)
	}

	// Output:
	//
}
