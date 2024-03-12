package utils

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"net"
)

const MsgHeadLen = 4

// This is to send message between HTTP and TCP
// write the length of message as first and then append message
func SendMsg(conn net.Conn, req interface{}) error {
	data, err := json.Marshal(req)
	if err != nil {
		return errors.New("JSON Serialization error" + err.Error())
	}
	lengthBytes := make([]byte, MsgHeadLen)
	binary.BigEndian.PutUint32(lengthBytes, uint32(len(data)))
	if _, err := conn.Write(append(lengthBytes, data...)); err != nil {
		return errors.New("Write request data failed: " + err.Error())
	}
	return nil
}

// This is to Receive message between HTTP and TCP
// receive the length of message as first and then receive message
func ReadMsg(conn net.Conn, resp interface{}) error {
	lengthBytes := make([]byte, MsgHeadLen)
	if n, err := io.ReadAtLeast(conn, lengthBytes, MsgHeadLen); err != nil || n != MsgHeadLen {
		if err == io.EOF {
			return err
		}
		return errors.New("Read response length failed: " + err.Error())
	}
	length := int(binary.BigEndian.Uint32(lengthBytes))
	data := make([]byte, length)
	if n, err := io.ReadAtLeast(conn, data, length); err != nil || n != length {
		if err == io.EOF {
			return err
		}
		return errors.New("Read response data failed: " + err.Error())
	}
	if err := json.Unmarshal(data, resp); err != nil {
		return errors.New("JSON Serialization error" + err.Error())
	}
	return nil

}
