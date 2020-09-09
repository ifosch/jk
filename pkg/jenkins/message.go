package jenkins

// Message is a response
type Message struct {
	Message string
	Error   bool
	Done    bool
}

func reply(msg string, error, done bool, ch chan Message) {
	ch <- Message{
		Message: msg,
		Error:   error,
		Done:    done,
	}
}
