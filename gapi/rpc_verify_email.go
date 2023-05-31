package gapi

import (
	"context"

	db "github.com/mativm02/bank_system/db/sqlc"
	"github.com/mativm02/bank_system/pb"
	"github.com/mativm02/bank_system/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	violations := validateVerifyEmailRequest(req)
	if len(violations) > 0 {
		return nil, invalidArgumentError(violations)
	}

	result, err := server.store.VerifyEmailTx(ctx, db.VerifyEmailTxParams{
		EmailID:    req.EmailId,
		SecretCode: req.SecretCode,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot verify email: %v", err)
	}

	rsp := &pb.VerifyEmailResponse{
		IsVerified: result.User.IsEmailVerified,
	}

	return rsp, nil

}

func validateVerifyEmailRequest(req *pb.VerifyEmailRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateEmailID(req.EmailId); err != nil {
		violations = append(violations, &errdetails.BadRequest_FieldViolation{
			Field:       "email_id",
			Description: err.Error(),
		})
	}

	if err := val.ValidateSecretCode(req.SecretCode); err != nil {
		violations = append(violations, &errdetails.BadRequest_FieldViolation{
			Field:       "secret_code",
			Description: err.Error(),
		})
	}

	return violations
}
