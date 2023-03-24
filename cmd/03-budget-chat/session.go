package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"regexp"
)

type Session struct {
	UserName string
	conn     net.Conn
	channel  chan Message
}

func NewSession(conn net.Conn) (Session, error) {
	session := Session{conn: conn}

	err := session.writeLine("Welcome to budget chat! What shall I call you?")
	if err != nil {
		return session, err
	}

	userName, err := session.readLine()
	if err != nil {
		return session, err
	}

	isValidUsername := len(userName) > 0 && regexp.MustCompile("^[a-zA-Z0-9]*$").MatchString(userName)
	if !isValidUsername {
		return session, errors.New("invalid username")
	}

	session.UserName = userName
	session.channel = make(chan Message)

	return session, nil
}

func (s *Session) Run(c chan Message) {
	go s.StartReceiver()
	s.StartSession(c)
}

func (s *Session) StartReceiver() {
	for {
		err := s.writeMessage(<-s.channel)
		if err != nil {
			return
		}
	}
}

func (s *Session) StartSession(c chan Message) {
	for {
		msg, err := s.readLine()
		if err != nil {
			return
		}

		c <- NewMessage(msg, s.UserName)
	}
}

func (s *Session) readLine() (string, error) {
	scanner := bufio.NewScanner(s.conn)

	if scanner.Scan() {
		return string(scanner.Bytes()), nil
	} else {
		err := scanner.Err()
		if err != nil {
			return "", err
		}
		return "", io.EOF
	}
}

func (s *Session) writeMessage(msg Message) error {
	var err error
	if msg.System == false {
		_, err = fmt.Fprintf(s.conn, "[%s] %s\n", msg.UserName, msg.Text)
	} else {
		_, err = fmt.Fprintf(s.conn, "* %s\n", msg.Text)
	}
	return err
}

func (s *Session) writeLine(line string) error {
	_, err := s.conn.Write([]byte(line + "\n"))
	return err
}
