package services

import (
	"context"
	"fmt"
	"log"
	"mup/clients"
	"mup/data"
	"mup/utils"
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

func (ms *MupService) CheckForPersonsDrivingBans(ctx context.Context, jmbg string) (data.DrivingBans, error) {
	return ms.repo.CheckForPersonsDrivingBans(ctx, jmbg)
}

func (ms *MupService) GetPersonsRegistrations(ctx context.Context, jmbg string) (data.Registrations, error) {
	return ms.repo.GetPersonsRegistrations(ctx, jmbg)
}

func (ms *MupService) GetUserDrivingPermit(ctx context.Context, jmbg string) (data.TrafficPermits, error) {
	return ms.repo.GetUserDrivingPermit(ctx, jmbg)
}

func (ms *MupService) GetUserDrivingPermitDetails(ctx context.Context, jmbg, tokenStr string) (data.DrivingPermitDetailsList, error) {
	drivingPermits, err := ms.repo.GetUserDrivingPermits(ctx, jmbg)
	if err != nil {
		return nil, err
	}

	var drivingPermitDetailsList data.DrivingPermitDetailsList
	for _, permit := range drivingPermits {
		user, err := ms.ssoc.GetUserByJMBG(ctx, permit.Person, tokenStr)
		if err != nil {
			return nil, err
		}

		drivingPermitDetails := data.DrivingPermitDetails{
			ID:             permit.ID,
			Number:         permit.Number,
			IssuedDate:     permit.IssuedDate,
			ExpirationDate: permit.ExpirationDate,
			Approved:       permit.Approved,
			Person:         permit.Person,
			FirstName:      user.FirstName,
			LastName:       user.LastName,
		}

		drivingPermitDetailsList = append(drivingPermitDetailsList, drivingPermitDetails)
	}

	return drivingPermitDetailsList, nil
}

func (ms *MupService) GetPendingRegistrationRequests(ctx context.Context, tokenStr string) (data.RegistrationDetailsList, error) {
	pendingRequests, err := ms.repo.GetPendingRegistrationRequests(ctx)
	if err != nil {
		return nil, err
	}

	var registrationDetailsList data.RegistrationDetailsList
	for _, reg := range pendingRequests {
		user, err := ms.ssoc.GetUserByJMBG(ctx, reg.Owner, tokenStr)
		if err != nil {
			return nil, err
		}

		vehicle, err := ms.repo.GetVehicleByID(ctx, reg.VehicleID)
		if err != nil {
			return nil, err
		}

		registrationDetails := data.RegistrationDetails{
			RegistrationNumber: reg.RegistrationNumber,
			IssuedDate:         reg.IssuedDate,
			ExpirationDate:     reg.ExpirationDate,
			VehicleID:          reg.VehicleID,
			Owner:              reg.Owner,
			Plates:             reg.Plates,
			Approved:           reg.Approved,
			FirstName:          user.FirstName,
			LastName:           user.LastName,
			VehicleBrand:       vehicle.Brand,
			VehicleModel:       vehicle.Model,
		}

		registrationDetailsList = append(registrationDetailsList, registrationDetails)
	}

	return registrationDetailsList, nil
}

func (ms *MupService) GetPendingTrafficPermitRequests(ctx context.Context, tokenStr string) (data.TrafficPermitDetailsList, error) {
	pendingRequests, err := ms.repo.GetPendingTrafficPermitRequests(ctx)
	if err != nil {
		return nil, err
	}

	var trafficPermitDetailsList data.TrafficPermitDetailsList
	for _, permit := range pendingRequests {
		user, err := ms.ssoc.GetUserByJMBG(ctx, permit.Person, tokenStr)
		if err != nil {
			return nil, err
		}

		trafficPermitDetails := data.TrafficPermitDetails{
			ID:             permit.ID,
			Number:         permit.Number,
			IssuedDate:     permit.IssuedDate,
			ExpirationDate: permit.ExpirationDate,
			Approved:       permit.Approved,
			Person:         permit.Person,
			FirstName:      user.FirstName,
			LastName:       user.LastName,
		}

		trafficPermitDetailsList = append(trafficPermitDetailsList, trafficPermitDetails)
	}

	return trafficPermitDetailsList, nil
}

func (ms *MupService) RetrieveRegisteredVehicles(ctx context.Context) (data.Vehicles, error) {
	return ms.repo.RetrieveRegisteredVehicles(ctx)
}

func (ms *MupService) SubmitRegistrationRequest(ctx context.Context, registration *data.Registration) error {
	registration.Approved = false
	registration.IssuedDate = time.Now()
	registration.ExpirationDate = registration.IssuedDate
	registration.RegistrationNumber = utils.GenerateRegistration()
	registration.Plates = ""

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

	trafficPermit.ID = primitive.NewObjectID()
	trafficPermit.Person = user.JMBG
	trafficPermit.Approved = false
	trafficPermit.IssuedDate = time.Now()
	trafficPermit.Number = utils.GenerateRegistration()

	return ms.repo.SubmitTrafficPermitRequest(ctx, trafficPermit)
}

func (ms *MupService) GetPersonsVehicles(ctx context.Context, jmbg string) ([]data.Vehicle, error) {
	return ms.repo.GetPersonsVehicles(ctx, jmbg)
}

func (ms *MupService) SaveVehicle(ctx context.Context, vehicle *data.Vehicle) error {
	vehicle.Registration = ""
	vehicle.Plates = ""
	vehicle.ID = primitive.NewObjectID()
	return ms.repo.SaveVehicle(ctx, vehicle)
}

func (ms *MupService) IssueDrivingBan(ctx context.Context, drivingBan *data.DrivingBan) error {
	return ms.repo.IssueDrivingBan(ctx, drivingBan)
}

func (ms *MupService) ApproveRegistration(ctx context.Context, registration data.Registration) error {
	expirationDate := time.Now().AddDate(5, 0, 0)
	registration.Approved = true
	registration.ExpirationDate = expirationDate

	platesNumber := utils.GeneratePlates()

	plates := data.Plates{
		RegistrationNumber: registration.RegistrationNumber,
		PlatesNumber:       platesNumber,
		PlateType:          "vehicle plates",
		VehicleID:          registration.VehicleID,
		Owner:              registration.Owner,
	}

	registration.Plates = platesNumber

	err := ms.repo.ApproveRegistration(ctx, registration)
	if err != nil {
		return err
	}

	return ms.repo.IssuePlates(ctx, plates)
}

func (ms *MupService) DeletePendingRegistration(ctx context.Context, registrationNumber string) error {
	return ms.repo.DeletePendingRegistration(ctx, registrationNumber)
}

func (ms *MupService) DeletePendingTrafficPermit(ctx context.Context, permitID primitive.ObjectID) error {
	return ms.repo.DeletePendingTrafficPermit(ctx, permitID)
}

func (ms *MupService) GetVehiclesDTOByJMBG(ctx context.Context, jmbg string) (data.VehiclesDTO, error) {
	vehicles, err := ms.repo.GetPersonsVehicles(ctx, jmbg)
	if err != nil {
		return nil, err
	}

	var vehicleDTOs data.VehiclesDTO
	for _, vehicle := range vehicles {
		plates, err := ms.repo.GetPlatesByVehicleID(ctx, vehicle.ID)
		if err != nil {
			return nil, err
		}

		registration, err := ms.repo.GetRegistrationByVehicleID(ctx, vehicle.ID)
		if err != nil {
			return nil, err
		}

		vehicleDTO := data.VehicleDTO{
			ID:           vehicle.ID,
			Brand:        vehicle.Brand,
			Model:        vehicle.Model,
			Year:         vehicle.Year,
			Registration: registration,
			Plates:       plates,
			Owner:        vehicle.Owner,
		}
		vehicleDTOs = append(vehicleDTOs, vehicleDTO)
	}

	return vehicleDTOs, nil
}

func (ms *MupService) ApproveTrafficPermitRequest(ctx context.Context, permitID primitive.ObjectID) error {
	return ms.repo.ApproveTrafficPermitRequest(ctx, permitID)
}

func (ms *MupService) GetRegistrationByPlate(ctx context.Context, plate string) (data.Registration, error) {
	return ms.repo.GetRegistrationByPlate(ctx, plate)
}
func (ms *MupService) GetDrivingBan(ctx context.Context, jmbg string) (data.DrivingBan, error) {
	return ms.repo.GetDrivingBan(ctx, jmbg)
}

func (ms *MupService) GetDrivingPermitByJMBG(ctx context.Context, jmbg string) (data.TrafficPermit, error) {
	return ms.repo.GetDrivingPermitByJMBG(ctx, jmbg)
}

func (ms *MupService) SaveMup() error {
	err := ms.repo.SaveMup(context.Background())
	if err != nil {
		return err
	}
	return nil
}
