package workflows

import (
	"backend-sample/common"
	"backend-sample/database"

	"github.com/google/uuid"
)

type UserWorkflowService struct {
	repository database.UsersRepository
}

type UsersWorkflow interface {
	Create(UserRequest UserRequest) (*UserResponse, bool)
	Update(UserRequest UserRequest) (*UserResponse, bool)
	Delete(id string) bool
	GetUsers(id, name string) *[]UserResponse
}

type UserRequest struct {
	Id, Name, Email, Password string
}

type UserResponse struct {
	Id                    uuid.UUID
	Name, Email, Password string
}

func (w *UserWorkflowService) Create(req UserRequest) (*UserResponse, *common.BackendError) {
	if !common.StringMinMaxLength(req.Email, 1, 100) {
		return nil, common.NewBackendError(400, "Workflows.CreateUser.1", "invalid name", nil)
	}
	if !common.IsValidEmail(req.Email) {
		return nil, common.NewBackendError(400, "Workflows.CreateUser.2", "invalid email", nil)
	}
	if !common.StringMinMaxLength(req.Name, 1, 100) {
		return nil, common.NewBackendError(400, "Workflows.CreateUser.3", "invalid name", nil)
	}
	if !common.StringMinMaxLength(req.Password, 1, 100) {
		return nil, common.NewBackendError(400, "Workflows.CreateUser.4", "invalid password", nil)
	}

	user, err := w.repository.CreateUser(&database.UserEntity{Name: req.Name, Email: req.Email, Password: req.Password})

	if err != nil {
		return nil, err
	}

	return &UserResponse{Id: user.Id, Name: user.Name, Email: user.Email, Password: user.Password}, nil
}

func (w *UserWorkflowService) Update(req UserRequest) (*UserResponse, *common.BackendError) {
	if !common.IsValidUuid(req.Id) {
		return nil, common.NewBackendError(400, "Workflows.UpdateUser.1", "invalid name", nil)
	}
	uuid := uuid.MustParse(req.Id)
	user, err := w.repository.GetUserById(uuid)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}
	if !common.StringMinMaxLength(req.Email, 1, 100) {
		return nil, common.NewBackendError(400, "Workflows.UpdateUser.3", "invalid name", nil)
	}
	if !common.IsValidEmail(req.Email) {
		return nil, common.NewBackendError(400, "Workflows.UpdateUser.4", "invalid email", nil)
	}
	if !common.StringMinMaxLength(req.Name, 1, 100) {
		return nil, common.NewBackendError(400, "Workflows.UpdateUser.5", "invalid name", nil)
	}
	if !common.StringMinMaxLength(req.Password, 1, 100) {
		return nil, common.NewBackendError(400, "Workflows.UpdateUser.6", "invalid password", nil)
	}

	user.Email = req.Email
	user.Name = req.Name
	user.Password = req.Password
	err = w.repository.UpdateUser(*user)

	if err != nil {
		return nil, err
	}

	return &UserResponse{Id: user.Id, Name: user.Name, Email: user.Email, Password: user.Password}, nil
}

func (w *UserWorkflowService) Delete(id string) *common.BackendError {
	if common.IsValidUuid(id) {
		return common.NewBackendError(400, "Workflows.DeleteUser.1", "invalid uuid", nil)
	}
	value := uuid.MustParse(id)

	err := w.repository.DeleteUser(value)

	if err != nil {
		return err
	}

	return nil
}

func (w *UserWorkflowService) GetUsers(id, name string) (*[]UserResponse, *common.BackendError) {
	if id != "" && len(id) > 0 {
		user, err := w.getUserById(id)
		if err != nil {
			return nil, err
		}
		return &[]UserResponse{*user}, nil
	}

	if name != "" && len(name) > 0 {
		users, err := w.getUserByName(name)
		if err != nil {
			return nil, err
		}
		return users, nil
	}

	return nil, nil
}

func (w *UserWorkflowService) getUserById(id string) (*UserResponse, *common.BackendError) {
	if !common.IsValidUuid(id) {
		return nil, common.NewBackendError(400, "Workflows.getUserById.1", "invalid id %s", nil, id)
	}

	uuid := uuid.MustParse(id)

	user, err := w.repository.GetUserById(uuid)

	if err != nil {
		return nil, err
	}

	return parseEntityToResponse(*user), nil
}

func (w *UserWorkflowService) getUserByName(name string) (*[]UserResponse, *common.BackendError) {
	if !common.StringMinMaxLength(name, 1, 100) {
		return nil, common.NewBackendError(400, "Workflows.getUserByName.1", "invalid name", nil)
	}

	users, err := w.repository.GetUsersByName(name, false)

	if err != nil {
		return nil, err
	}

	return parseEntityListToResponse(*users), nil
}

func parseEntityToResponse(user database.UserEntity) *UserResponse {
	return &UserResponse{Id: user.Id, Name: user.Name, Email: user.Email, Password: user.Password}
}

func parseEntityListToResponse(users []database.UserEntity) *[]UserResponse {
	var response = make([]UserResponse, len(users))
	for i, u := range users {
		response[i] = *parseEntityToResponse(u)
	}
	return &response
}
