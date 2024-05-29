package services

import (
	"context"
	"fmt"
	"log"
	"mup/clients"
	"mup/data"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MupService struct {
	repo   *data.MUPRepo
	logger *log.Logger
	ssoc   clients.SSOClient
	cc     clients.CourtClient
}

func NewMupService(r *data.MUPRepo, log *log.Logger, ssoc clients.SSOClient, cc clients.CourtClient) *MupService {
	return &MupService{repo: r, logger: log, ssoc: ssoc, cc: cc}
}

func (ms *MupService) CheckForPersonsDrivingBans(ctx context.Context, personID primitive.ObjectID) (data.DrivingBans, error) {
	return ms.repo.CheckForPersonsDrivingBans(ctx, personID)
}

func (ms *MupService) RetrieveRegisteredVehicles(ctx context.Context) (data.Vehicles, error) {
	return ms.repo.RetrieveRegisteredVehicles(ctx)
}

func (ms *MupService) SubmitRegistrationRequest(ctx context.Context, registration *data.Registration) error {
	registration.Approved = false
	registration.IssuedDate = time.Now()

	err := ms.repo.SubmitRegistrationRequest(ctx, registration)
	if err != nil {
		return err
	}

	err = ms.repo.SaveRegistrationIntoVehicle(ctx, registration)
	if err != nil {
		return err
	}

	return nil
}

func (ms *MupService) SubmitTrafficPermitRequest(ctx context.Context, trafficPermit *data.TrafficPermit, jmbg, tokenStr string) error {
	user, err := ms.ssoc.GetUserByJMBG(ctx, jmbg, tokenStr)
	if err != nil {
		return err
	}

	warrants, err := ms.cc.CheckForPersonsWarrant(ctx, trafficPermit.Person, tokenStr)
	if err != nil {
		return err
	}

	if len(warrants) != 0 {
		return fmt.Errorf("user is on warrant list")
	}

	trafficPermit.Person = user.JMBG
	trafficPermit.Approved = false
	trafficPermit.IssuedDate = time.Now()

	return ms.repo.SubmitTrafficPermitRequest(ctx, trafficPermit)
}

func (ms *MupService) SaveVehicle(ctx context.Context, vehicle *data.Vehicle) error {
	return ms.repo.SaveVehicle(ctx, vehicle)
}

func (ms *MupService) IssueDrivingBan(ctx context.Context, drivingBan *data.DrivingBan) error {
	return ms.repo.IssueDrivingBan(ctx, drivingBan)
}

func (ms *MupService) ApproveRegistration(ctx context.Context, registration data.Registration) error {
	expirationDate := time.Now().AddDate(1, 0, 0)
	registration.Approved = true
	registration.ExpirationDate = expirationDate

	registration.RegistrationNumber = clients.GenerateRegistration()
	platesNumber := clients.GeneratePlates()

	err := ms.repo.ApproveRegistration(ctx, registration)
	if err != nil {
		return err
	}

	plates := data.Plates{
		RegistrationNumber: registration.RegistrationNumber,
		PlatesNumber:       platesNumber,
		PlateType:          "plateType",
		VehicleID:          registration.VehicleID,
	}

	return ms.repo.IssuePlates(ctx, plates)
}

func (ms *MupService) ApproveTrafficPermitRequest(ctx context.Context, permitID primitive.ObjectID) error {
	return ms.repo.ApproveTrafficPermitRequest(ctx, permitID)
}

func (ms *MupService) SaveMup() error {
	err := ms.repo.SaveMup(context.Background())
	if err != nil {
		return err
	}
	return nil
}
