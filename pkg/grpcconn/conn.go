package grpcconn

type GrpcConnecter interface {
	Connect() error
	Request() (interface{}, error)
	Close() error
}
