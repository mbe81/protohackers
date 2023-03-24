package main

type Message struct {
	Text     string
	UserName string
	System   bool
}

func NewMessage(text string, userName string) Message {
	return Message{Text: text, UserName: userName, System: false}
}

func NewEvent(text string, userName string) Message {
	return Message{Text: text, UserName: userName, System: true}
}
