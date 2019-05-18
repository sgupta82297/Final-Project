package services

import (
	b64 "encoding/base64"
	"errors"
	"fmt"
	"log"
	"soilprotection-service/constants"
	"soilprotection-service/models"
	"time"
)

//SaveUser : ""
func (s *Service) SaveUser(collectionName string, user *models.User, dept models.IDepartment) (string, error) {
	suffix, prefix, digits, department, err := models.DeptConfigProviders(collectionName)
	if err != nil {
		return "", errors.New("Invalid department")
	}
	str := fmt.Sprintf("%v%v%v", "%s", digits, "%s")
	username := fmt.Sprintf(str, suffix, s.Daos.Next(collectionName), prefix)
	// username := s.Daos.Next(collectionName)
	user.UserName = username
	user.CreatedOn = time.Now()
	user.Department = department
	if department == models.DEPARTMENTFARMER {
		user.Status = constants.USERSTATUSACTIVE
	} else {
		user.Status = constants.USERSTATUSINIT
	}

	if user.Pass == "" {
		user.Password = constants.DEFAULTLOGINPASSWORD
	} else {
		user.Password = user.Pass
	}
	err = s.Daos.SaveUser(user)
	if err != nil {
		return "", errors.New("Error in saving user - " + err.Error())
	}
	dept.SetUserName(username)
	err = s.Daos.SaveDepartment(collectionName, dept)
	if err != nil {
		return "", errors.New("Error in saving dept - " + err.Error())
	}
	return username, nil
}

//GetAllUser : ""
func (s *Service) GetAllUser() ([]models.User, error) {

	usr, err := s.Daos.GetAllUser()

	return usr, err

}

//ForgetPassword : ""
func (s *Service) ForgetPassword(uniqueID string) error {
	user, err := s.Daos.GetuserWithUniqueID(uniqueID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New(constants.INTERNALSERVERERROR)
	}
	_, err = s.GenerateOTP(constants.OTPFORGETPASSWORD, user.Mobile, constants.PHONEOTPLENGTH, constants.OTPEXPIRY)
	if err != nil {
		return err
	}
	return nil
}

//ValidateForgetPassWordOTP : ""
func (s Service) ValidateForgetPassWordOTP(uniqueID, OTP string) (string, error) {
	user, err := s.Daos.GetuserWithUniqueID(uniqueID)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New(constants.INTERNALSERVERERROR)
	}
	err = s.ValidateOTP(constants.OTPFORGETPASSWORD, user.Mobile, OTP)
	if err != nil {
		return "", err
	}
	Token, err := s.GenerateOTP(constants.OTPTOKEN, user.Mobile, constants.TOKENOTPLENGTH, constants.OTPEXPIRY)
	if err != nil {
		return "", err
	}

	sEnc := b64.StdEncoding.EncodeToString([]byte(Token))
	fmt.Println(sEnc)
	return sEnc, nil

}

//ChangePassWithToken : ""
func (s Service) ChangePassWithToken(input *models.ChangeWithToken) error {
	user, err := s.Daos.GetuserWithUniqueID(input.UserName)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New(constants.INTERNALSERVERERROR)
	}
	fmt.Println("TOKEN INPUT ", input.Token)
	sEnc := input.Token
	sDec, _ := b64.StdEncoding.DecodeString(sEnc)
	token := string(sDec)
	fmt.Println(".....", sDec)
	fmt.Println(".....", token)
	err = s.ValidateOTP(constants.OTPTOKEN, user.Mobile, token)
	if err != nil {
		return err
	}
	err = s.Daos.ChangePassword(input.UserName, input.PassWord)
	if err != nil {
		return err
	}
	return nil
}

//GetuserWithUniqueID : ""
func (s Service) GetuserWithUniqueID(userName string) (*models.User, error) {
	user, err := s.Daos.GetuserWithUniqueID(userName)
	return user, err
}

//ProfileUpdate : ""
func (s Service) ProfileUpdate(user *models.User) error {

	count, err := s.Daos.CheckUniquness(user)
	if err != nil {
		return err
	}
	fmt.Println("COUNT", count)
	if count == 0 {
		err = s.Daos.UpdateUserWithUniqueID(user.UserName, user)
		if err != nil {
			return err
		}

	} else if count == 1 {
		return errors.New("Mobile Number or Email exists already")
	} else if count == 2 {

		return errors.New("Mobile Number or Email exists already")
	} else if count > 2 {
		log.Printf("Alarming Number of duplicate Entries Found in The Database ")
		return errors.New("Mobile Number or Email exists already")
	}
	return nil
}

//UserStatusBulkChange : " "
func (s *Service) UserStatusBulkChange(arr []string, stat string) error {

	var err error

	switch {
	case stat == "active":
		err = s.Daos.UserStatusBulkChange(arr, constants.USERSTATUSACTIVE)
	case stat == "deactive":
		err = s.Daos.UserStatusBulkChange(arr, constants.USERSTATUSDEACTIVE)
	case stat == "delete":
		err = s.Daos.UserStatusBulkChange(arr, constants.USERSTATUSDELETED)
	default:
		return errors.New("Enter Valid URL")
	}

	if err != nil {
		return err
	}
	return nil
}

//GetNearByUsers : ""
func (s *Service) GetNearByUsers(userType string, km float64, coordinates []float64) ([]models.User, error) {
	return s.Daos.GetNearByUsers(userType, km, coordinates)
}

// SearchUserByKeywords : ""
func (s *Service) SearchUserByKeywords(userSearch *models.UserSearch, projectedFields []string, pagination *models.Pagination) (users []models.User, err error) {
	return s.Daos.SearchUserByKeywords(userSearch, projectedFields, pagination)
}

//ActivteNewUsers : ""
func (s *Service) ActivteNewUsers(user models.User) error {
	user.Status = constants.USERSTATUSACTIVE
	return s.Daos.ActivteNewUsers(user)
}

//UpdateCurrentLocation : ""
func (s *Service) UpdateCurrentLocation(username string, loc []float64) error {
	var location models.Location
	location.Type = "point"
	location.Coordinates = loc
	return s.Daos.UpdateCurrentLocation(username, location)
}

//AddressConversion : ""
func (s *Service) AddressConversion(address *models.AddressV2) (*models.AddressV2, error) {
	if address == nil {
		return nil, errors.New("Address Missing!!")
	}
	if address.County != "" {
		cou, err := s.Daos.GetSingleGeo("Counties", "V1", "code:"+address.County)
		if err != nil {
			log.Println("county err==>", err.Error())
		}
		if err == nil {
			var county models.County
			s.Shared.BsonToType(cou, &county)
			// log.Println("county==>", county)
			address.CountyName = county.Name
		}
	}

	if address.State != "" {
		stat, err := s.Daos.GetSingleGeo("states", "V1", "code:"+address.State)
		if err != nil {
			log.Println("county err==>", err.Error())
		}
		if err == nil {
			var state models.State
			s.Shared.BsonToType(stat, &state)
			// log.Println("county==>", state)
			address.StateName = state.Name
		}
	}
	if address.ZipCode != "" {
		blk, err := s.Daos.GetSingleGeo("zipcode", "V1", "code:"+address.ZipCode)
		if err != nil {
			log.Println("state err==>", err.Error())
		}
		if err == nil {
			var zipcode models.ZipCode
			s.Shared.BsonToType(zip, &zipcode)
			// log.Println("zipcode==>", zipcode)
			address.ZipCode = zipcode.Name
		}
	}
	if address.City != "" {
		ci, err := s.Daos.GetSingleGeo("cities", "V1", "code:"+address.City)
		if err != nil {
			log.Println("cities err==>", err.Error())
		}
		if err == nil {
			var cities models.City
			s.Shared.BsonToType(ci, &cities)
			// log.Println("cities==>", cities)
			address.CityName = cities.Name
		}
	}
	if address.Street != "" {
		str, err := s.Daos.GetSingleGeo("streets", "S1", "code:"+address.Street)
		if err != nil {
			log.Println("str err==>", err.Error())
		}
		if err == nil {
			var street models.Street
			s.Shared.BsonToType(str, &street)
			// log.Println("street==>", street)
			address.Street = street.Name
		}
	}
	return address, nil
}
