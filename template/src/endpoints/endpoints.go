package endpoints

import (
	"github.com/go-kit/kit/endpoint"

	"<%= domainDir + _.folderName %>/src/service"
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
