package tcpClient

import (
	"bufio"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/mock"
)

var mockReaderObj mockReader

// implements the target interface
type mockReader struct {
	mock.Mock
}

func (m *mockReader) ReadString(delim byte) (string, error) {
	return "foo", nil
}

// ********** "Implemented" interface ***********
func (m *mockReader) Reset(r io.Reader) {

}

func (m *mockReader) Peek(n int) ([]byte, error) {

	return nil, nil
}

func (m *mockReader) Discard(n int) (discarded int, err error) {

	return 0, nil
}

func (m *mockReader) Read(p []byte) (n int, err error) {

	return 0, nil
}

func (m *mockReader) ReadByte() (byte, error) {

	return '\n', nil
}

func (m *mockReader) UnreadByte() error {

	return nil
}

func (m *mockReader) ReadRune() (r rune, size int, err error) {

	return '\n', 0, nil
}

func (m *mockReader) Buffered() int {

	return 0
}

func (m *mockReader) ReadSlice(delim byte) (line []byte, err error) {

	return nil, nil
}

func (m *mockReader) ReadLine() (line []byte, isPrefix bool, err error) {

	return nil, true, nil
}

func (m *mockReader) ReadBytes(delim byte) ([]byte, error) {

	return nil, nil
}

func (m *mockReader) WriteTo(w io.Writer) (n int64, err error) {

	return 0, nil
}

// *********************

// Attempting to discern best route to go for mocking this bad boy (reading process from stdin) out.

func mockCreateReader() *bufio.Reader {
	return &mockReaderObj // this should be our mock object with mocked return values
}

func TestCreateUser(t *testing.T) {
	expectedUserName := "Foobarbaz"
	mockReaderObj := new(mockReader)
	mockReaderObj.On("ReadString", '\n').Return(expectedUserName, nil)

	//Mock out the createReader func
	client := Client{createReader: mockCreateReader}

	user, err := client.CreateUser()
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	if user != expectedUserName {
		fmt.Println("Failed to create the user. Actual user name is not the expected one.")
		t.Fail()
	}

}

func TestConstructMessage(t *testing.T) {
	message := constructMessage("Foo", "Bar")
	if message.Body != "Bar" || message.Sender != "Foo" {
		fmt.Printf("Construction of message failed. Username: %s, Body: %s\n", message.Sender, message.Body)
		t.Fail()
	}
}
