package commands

import (
	"context"
	"encoding/json"
	"errors"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	app_interface "seal/internal/app/interface"
	"seal/internal/domain/command"
	"seal/internal/repository/pg/query"
	"seal/pkg/app_error"
	"time"

	commands_v1 "gitlab.kvant.online/seal/grpc-contracts/pkg/commands/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var dlErrTxt = errors.New("Сервис команд временно недоступен, повторите оерацию позже.")

type grpcClient struct {
	ctx     context.Context
	addr    string
	timeout uint8
	logger  app_interface.Logger
}

func NewGRPC(ctx context.Context, addr string, timeout uint8, logger app_interface.Logger) grpcClient {
	return grpcClient{
		ctx,
		addr,
		timeout,
		logger,
	}
}

func (s grpcClient) Send(imei string, name string, params any, author string) (bool, error) {
	conn, err := s.getConnection()

	if err != nil {
		return false, err
	}

	defer conn.Close()

	ctx, cancel := context.WithTimeout(s.ctx, time.Second*time.Duration(s.timeout))
	defer cancel()

	c := commands_v1.NewCommandsServiceClient(conn)

	jsonParams, err := json.Marshal(params)
	if err != nil {
		s.logger.Error(err.Error())
		return false, err
	}

	rBody := &commands_v1.AddRequest{
		Imei:     imei,
		Protocol: 1,
		Name:     name,
		Params:   string(jsonParams),
		Author:   author,
	}

	r, err := c.Add(ctx, rBody)

	if err != nil {

		st, _ := status.FromError(err)
		s.logger.Error(st.Message())

		if st.Code() == codes.DeadlineExceeded || st.Code() == codes.Unavailable {
			return false, dlErrTxt
		}

		for _, detail := range st.Details() {

			switch t := detail.(type) {
			case *errdetails.BadRequest:

				validateErrors := map[string]string{}

				for _, violation := range t.GetFieldViolations() {
					validateErrors[violation.GetField()] = violation.GetDescription()
				}

				return false, app_error.ValidationError(validateErrors)
			}
		}

		return false, app_error.ValidationError(map[string]string{"command": st.Message()})
	}

	if r.Code > 0 {
		return false, app_error.ValidationError(r.Messages)
	}

	return true, nil
}

func (s grpcClient) List(imei string) (query.List[command.Command], error) {
	conn, err := s.getConnection()

	if err != nil {
		return query.List[command.Command]{}, err
	}

	defer conn.Close()

	ctx, cancel := context.WithTimeout(s.ctx, time.Second*time.Duration(s.timeout))
	defer cancel()

	c := commands_v1.NewCommandsServiceClient(conn)

	rBody := &commands_v1.ListRequest{
		Imei:    imei,
		Offset:  0,
		Limit:   10,
		NewOnly: false,
	}

	r, err := c.List(ctx, rBody)

	if err != nil {

		s.logger.Error(err.Error())
		st, _ := status.FromError(err)
		if st.Code() == codes.DeadlineExceeded || st.Code() == codes.Unavailable {
			return query.List[command.Command]{}, dlErrTxt
		}

		return query.List[command.Command]{}, err
	}

	cmdList := query.List[command.Command]{
		Data: make([]command.Command, 0),
	}

	for _, cmd := range r.GetCommands() {

		var tryDate, responseDate, abortDate *time.Time

		if cmd.TryDate != nil {
			td := cmd.GetTryDate().AsTime()
			tryDate = &td
		}
		if cmd.ResponseDate != nil {
			responseD := cmd.GetResponseDate().AsTime()
			responseDate = &responseD
		}
		if cmd.AbortDate != nil {
			abortD := cmd.GetAbortDate().AsTime()
			abortDate = &abortD
		}

		response := cmd.GetResponse()

		cmdList.Data = append(cmdList.Data, command.Command{
			Id:           cmd.GetId(),
			Serial:       cmd.GetImei(),
			Name:         cmd.GetName(),
			Params:       cmd.GetParams(),
			Author:       cmd.GetAuthor(),
			Dateon:       cmd.GetDateon().AsTime(),
			TryNumber:    cmd.TryNumber,
			TryDate:      tryDate,
			ResponseDate: responseDate,
			Response:     &response,
			RawRequest:   cmd.RawRequest,
			RawResponse:  cmd.RawResponse,
			AbortDate:    abortDate,
		})
	}

	return cmdList, nil
}

func (s grpcClient) getConnection() (*grpc.ClientConn, error) {
	ctx, cf1 := context.WithTimeout(s.ctx, time.Second*time.Duration(s.timeout))
	defer cf1()

	conn, err := grpc.DialContext(ctx, s.addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())

	if err != nil {
		s.logger.Error(err.Error())

		if errors.Is(err, context.DeadlineExceeded) {
			return nil, dlErrTxt
		}

		return nil, err
	}

	return conn, nil
}
