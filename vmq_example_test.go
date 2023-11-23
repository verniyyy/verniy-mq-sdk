package vmq

import "log"

func ExampleSession() {
	const addr = "localhost:9000"
	sess, err := NewSession(&Options{
		Addr:     addr,
		UserID:   "",
		Password: "",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		log.Println("Close session")
		if err := sess.Close(); err != nil {
			panic(err)
		}
	}()

	if err := Ping(sess); err != nil {
		log.Fatal(err)
	}

	//Output:
	//
}

type User struct {
	ID   int
	Name string
}
