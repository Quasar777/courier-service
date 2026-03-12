package courier

import (
	"github.com/Quasar777/courier-service/internal/handler/dto"
	uc "github.com/Quasar777/courier-service/internal/usecase/courier"
)

func toUseCaseCreateModel(r dto.CreateCourierRequest) uc.CreateCourierInput {
	res := uc.CreateCourierInput{
		Name: r.Name,
		Lastname: r.Lastname,
		Phone: r.Phone,
		Status: r.Status,
		TransportType: r.TransportType,
	}

	return res
}

func toUseCaseUpdateModel(r dto.UpdateCourierRequest) uc.UpdateCourierInput {
	res := uc.UpdateCourierInput{
		ID: r.Id,
		Name: r.Name,
		Lastname: r.Lastname,
		Phone: r.Phone,
		Status: r.Status,
		TransportType: r.TransportType,
	}

	return res
}