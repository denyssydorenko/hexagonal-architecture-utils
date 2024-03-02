package domain

const (
	XTRACEID = "x-trace-id"
)

func ToPointer[T any](t T) *T {
	return &t
}
