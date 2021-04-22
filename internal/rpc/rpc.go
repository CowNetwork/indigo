package rpc

import (
	"github.com/cownetwork/indigo/internal/dao"
	"github.com/cownetwork/indigo/internal/model"
	pb "github.com/cownetwork/mooapis-go/cow/indigo/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"regexp"
)

var snakeCaseRegex = regexp.MustCompile("^[a-z]+(_[a-z]+)*$")

type IndigoServiceServer struct {
	pb.UnimplementedIndigoServiceServer
	Dao dao.DataAccessor
}

func ValidateRole(r *model.Role) error {
	n := r.Name
	t := r.Type
	if len(n) == 0 || len(t) == 0 {
		return status.Error(codes.InvalidArgument, "name or type can not be empty.")
	}
	if !snakeCaseRegex.MatchString(n) || !snakeCaseRegex.MatchString(t) {
		return status.Error(codes.InvalidArgument, "name and type must be snake_case.")
	}
	if len(r.Color) > 6 {
		return status.Error(codes.InvalidArgument, "role color length must be 6 or less.")
	}
	return nil
}
