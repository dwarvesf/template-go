package endpoints

import (
	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/template-go/template/service"
)

type Endpoints struct {
	Add endpoint.Endpoint
}

// MakeServerEndpoints returns an Endpoints struct
func MakeServerEndpoints(s service.Service) Endpoints {
	return Endpoints{
		Add: MakeAddEndpoint(s),
	}
}
