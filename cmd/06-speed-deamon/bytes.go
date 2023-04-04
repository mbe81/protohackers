package main

import (
	"encoding/binary"
	"io"
)

func WriteByte(w io.Writer, b byte) error {
	if _, err := w.Write([]byte{b}); err != nil {
		return err
	}

	return nil
}

func WriteString(w io.Writer, str string) error {
	// first byte contains string length, following n bytes the actual string

	length := uint8(len(str))
	if _, err := w.Write([]byte{length}); err != nil {
		return err
	}
	if _, err := w.Write([]byte(str)); err != nil {
		return err
	}
	return nil
}

func WriteUint32(w io.Writer, i uint32) error {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, i)
	if _, err := w.Write(data); err != nil {
		return err
	}
	return nil
}

func WriteUint16(w io.Writer, i uint16) error {
	data := make([]byte, 2)
	binary.BigEndian.PutUint16(data, i)
	if _, err := w.Write(data); err != nil {
		return err
	}
	return nil
}

func ReadByte(r io.Reader) (byte, error) {
	data := make([]byte, 1)
	if _, err := io.ReadFull(r, data); err != nil {
		return 0, err
	}
	return data[0], nil
}

func ReadString(r io.Reader) (string, error) {
	// first byte contains string length, following n bytes the actual string
	length := make([]byte, 1)
	if _, err := io.ReadFull(r, length); err != nil {
		return "", err
	}

	data := make([]byte, length[0])
	if _, err := io.ReadFull(r, data); err != nil {
		return "", err
	}
	return string(data), nil
}

func ReadUInt32(r io.Reader) (uint32, error) {
	var err error

	data := make([]byte, 4)
	_, err = io.ReadFull(r, data)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(data), nil
}

func ReadUInt16(r io.Reader) (uint16, error) {
	var err error

	data := make([]byte, 2)
	_, err = io.ReadFull(r, data)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(data), nil
}

func ReadUInt8(r io.Reader) (uint8, error) {
	return ReadByte(r)
}
