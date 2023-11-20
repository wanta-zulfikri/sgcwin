package errorr

type BadRequest struct {
	Err string
}

func (e BadRequest) Error() string {
	return e.Err
}

func NewBad(err string) BadRequest {
	return BadRequest{err}
}

type InternalServer struct {
	Err string
}

func (i InternalServer) Error() string {
	return i.Err
}

func NewInternal(err string) InternalServer {
	return InternalServer{err}
}
