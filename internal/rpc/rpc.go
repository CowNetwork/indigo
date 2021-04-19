package rpc

import (
	"github.com/cownetwork/indigo/internal/dao"
	pb "github.com/cownetwork/mooapis-go/cow/indigo/v1"
)

type IndigoServiceServer struct {
	pb.UnimplementedIndigoServiceServer
	Dao dao.DataAccessor
}
