package client

import (
	"io"
	"net"
	"os/exec"
	"time"

	"github.com/gw31415/texc/proto"
	"google.golang.org/grpc"
)

type file_conn struct {
	wc io.WriteCloser
}
func (file_conn *file_conn)Read() error {
	return 
}
func (fileconn *file_conn)Close() error {
	return fileconn.wc.Close()
}
func (wc *file_conn) Write(p []byte) (n int, err error) {
	return wc.wc.Write(p)
}
func (fc *file_conn)LocalAddr() net.Addr {
	return nil
}
func NewFileConn(wc io.WriteCloser) *file_conn {
	return &file_conn{
		wc: wc,
	}
}

func main() {
	dialer := func(ex string, t time.Duration) (net.Conn, error) {
		cmd := exec.Command(ex)
		file, err := cmd.StdinPipe()
		return NewFileConn(file), nil

	}
	grpc.Dial("", grpc.WithDialer(dialer))
	proto.NewTexcServiceClient
}
