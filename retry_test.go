package retry

import (
	"context"
	"testing"
	"time"
	"unsafe"

	errnov1 "github.com/rejchev/errno"
)

type Data struct {
	Hello [64]byte
}

func DataFrom(v unsafe.Pointer) *Data {
	return (*Data)(v)
}

func (x *Data) Pointer() unsafe.Pointer {
	return unsafe.Pointer(x)
}

func TestThatItWork(t *testing.T) {
	i := 0
	data := Data{}

	errno := Try(LinearDelay(time.Microsecond), 3)(context.Background(), data.Pointer(), func(ctx context.Context, buff unsafe.Pointer) errnov1.Code {
		if i == 2 {
			copy(DataFrom(buff).Hello[:], []byte("Hello, World"))
			return errnov1.OK
		}

		i++
		return errnov1.ECALL
	})

	if errnov1.FAIL(errno) {
		t.Errorf("expected %d, but got %d", errnov1.OK, errno)
	}
}
