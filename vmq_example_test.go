package vmq

import "log"

// SessionOptionsExample ...
var SessionOptionsExample = &Options{
	Addr:     "localhost:9000",
	UserID:   "01HG17X22440GTQW3AS6WHCF0K",
	Password: "P@ssw0rd",
}

func ExampleSession() {
	sess, err := NewSession(SessionOptionsExample)
	if err != nil {
		return
	}
	defer func() {
		if err := sess.Close(); err != nil {
			panic(err)
		}
	}()

	//Output:
	//
}

func ExamplePing() {
	sess, err := NewSession(SessionOptionsExample)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		if err := sess.Close(); err != nil {
			panic(err)
		}
	}()

	if err := Ping(sess); err != nil {
		log.Println(err)
	}

	//Output:
	//
}

// ExampleCreateQueue ...
func ExampleCreateQueue() {
	sess, err := NewSession(SessionOptionsExample)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		if err := sess.Close(); err != nil {
			panic(err)
		}
	}()

	if err := CreateQueue(sess, "example queue"); err != nil {
		log.Println(err)
	}

	//Output:
	//
}

// ExampleDeleteQueue ...
func ExampleDeleteQueue() {
	sess, err := NewSession(SessionOptionsExample)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		if err := sess.Close(); err != nil {
			panic(err)
		}
	}()

	if err := DeleteQueue(sess, "example queue"); err != nil {
		log.Println(err)
	}

	//Output:
	//
}

// ExamplePublish ...
func ExamplePublish() {
	sess, err := NewSession(SessionOptionsExample)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		if err := sess.Close(); err != nil {
			panic(err)
		}
	}()

	msg := User{ID: 1, Name: "Jhon"}

	if err := Publish(sess, "example queue", msg); err != nil {
		log.Println(err)
	}

	//Output:
	//
}

type User struct {
	ID   int
	Name string
}
