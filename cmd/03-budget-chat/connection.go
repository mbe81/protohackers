package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"regexp"
)

type Connection struct {
	UserName string
	conn     net.Conn
	ch       chan Message
}

func NewConnection(conn net.Conn) (Connection, error) {
	c := Connection{conn: conn}

	err := c.writeLine("Welcome to budget chat! What shall I call you?")
	if err != nil {
		return c, err
	}

	userName, err := c.readLine()
	if err != nil {
		return c, err
	}

	isValidUsername := len(userName) > 0 && regexp.MustCompile("^[a-zA-Z0-9]*$").MatchString(userName)
	if !isValidUsername {
		return c, errors.New("invalid username")
	}

	c.UserName = userName
	c.ch = make(chan Message)

	return c, nil
}

func (c *Connection) Run(ch chan Message) {
	go c.StartReceiver()
	c.StartSession(ch)
}

func (c *Connection) StartReceiver() {
	for {
		err := c.writeMessage(<-c.ch)
		if err != nil {
			return
		}
	}
}

func (c *Connection) StartSession(ch chan Message) {
	for {
		msg, err := c.readLine()
		if err != nil {
			return
		}

		ch <- NewMessage(msg, c.UserName)
	}
}

func (c *Connection) readLine() (string, error) {
	scanner := bufio.NewScanner(c.conn)

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

func (c *Connection) writeMessage(msg Message) error {
	var err error
	if msg.System == false {
		_, err = fmt.Fprintf(c.conn, "[%s] %s\n", msg.UserName, msg.Text)
	} else {
		_, err = fmt.Fprintf(c.conn, "* %s\n", msg.Text)
	}
	return err
}

func (c *Connection) writeLine(line string) error {
	_, err := c.conn.Write([]byte(line + "\n"))
	return err
}
