package service

import (
	"ecommerce/config"
	"ecommerce/internal/domain"
	"ecommerce/internal/dto"
	"ecommerce/internal/helper"
	"ecommerce/internal/repository"
	"ecommerce/notification"
	"errors"
	"fmt"
	"time"
)

type UserService struct {
	Repo   repository.UserRepository
	Auth   helper.Auth
	Config config.AppConfig
}

func (s UserService) SignUp(input dto.UserSignUp) (string, error) {

	hPassword, err := s.Auth.CreateHashedPassword(input.Password)
	if err != nil {
		return "", err
	}

	user, err := s.Repo.CreateUser(domain.User{
		Email:    input.Email,
		Password: hPassword,
		Phone:    input.Phone,
	})

	return s.Auth.GenerateToken(user.ID, user.Email, user.UserType)

}

func (s UserService) findUserByEmail(email string) (*domain.User, error) {

	user, err := s.Repo.FindUser(email)

	return &user, err

}

func (s UserService) Login(email string, password string) (string, error) {

	user, err := s.findUserByEmail(email)
	if err != nil {
		return "", errors.New("user with given credentials doesn't exists")
	}

	err = s.Auth.VerifyPassword(password, user.Password)
	if err != nil {
		return "", err
	}

	return s.Auth.GenerateToken(user.ID, user.Email, user.UserType)

}

func (s UserService) isVerifiedUser(id uint) bool {

	currentUser, err := s.Repo.FindUserById(id)
	return err == nil && currentUser.Verified
}

func (s UserService) GetVerificationCode(u domain.User) error {

	if s.isVerifiedUser(u.ID) {
		fmt.Printf("user already verified \n")
		return errors.New("user already verified")
	}

	code, err := s.Auth.GenerateVerificationCode()
	if err != nil {
		fmt.Printf("error generating verification code %v\n", err)
		return errors.New("verification code generation failed")
	}

	user := domain.User{
		Expiry: time.Now().Add(30 * time.Minute),
		Code:   code,
	}

	_, err = s.Repo.UpdateUser(u.ID, user)
	if err != nil {
		fmt.Printf("error updating verification code to user: %v\n", err)
		return errors.New("verification code generation failed")
	}

	user, _ = s.Repo.FindUserById(u.ID)

	// Send SMS
	message := fmt.Sprintf("Your code for account verification is %v and is valid for next 30 minutes.", code)

	notificationClient := notification.NewNotificationClient(s.Config)
	err = notificationClient.SendSMS(user.Phone, message)
	if err != nil {
		fmt.Printf("error sending verification code to userId : %d, %v\n", user.ID, err)
	}

	return nil

}

func (s UserService) VerifyCode(id uint, code int) error {

	if s.isVerifiedUser(id) {
		fmt.Printf("user already verified \n")
		return errors.New("user already verified")
	}

	user, err := s.Repo.FindUserById(id)
	if err != nil {
		return errors.New("user not found")
	}

	if user.Code != code {
		return errors.New("verification code doesn't match")
	}

	if !time.Now().Before(user.Expiry) {
		return errors.New("verification code expired")
	}

	updateUser := domain.User{}
	updateUser.Verified = true

	_, err = s.Repo.UpdateUser(user.ID, updateUser)
	if err != nil {
		fmt.Printf("error updating verified status: %v\n", err)
		return errors.New("unable to verify user")
	}

	return nil

}

func (s UserService) CreateProfile(id uint, input any) error {

	return nil

}

func (s UserService) GetProfile(id uint) (*domain.User, error) {

	return nil, nil

}

func (s UserService) UpdateProfile(id uint, input any) error {

	return nil

}

func (s UserService) BecomeSeller(id uint, input dto.SellerInput) (string, error) {

	user, _ := s.Repo.FindUserById(id)

	if user.UserType == domain.SELLER {
		fmt.Println("The user is already a seller")
		return "", errors.New("you are already a seller")
	}

	updatedUser := domain.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Phone:     input.PhoneNumber,
		UserType:  domain.SELLER,
	}

	seller, err := s.Repo.UpdateUser(id, updatedUser)
	if err != nil {
		return "", err
	}

	token, err := s.Auth.GenerateToken(id, user.Email, seller.UserType)
	if err != nil {
		return "", err
	}

	account := domain.BankAccount{
		BankAccount: input.BankAccountNumber,
		SwiftCode:   input.SwiftCode,
		PaymentType: input.PaymentType,
		UserId:      id,
	}

	err = s.Repo.CreateBankAccount(account)
	if err != nil {
		fmt.Println("Error occurred while creating bank record for user %v", err)
		return "", errors.New("error creating bank info")
	}

	return token, nil

}

func (s UserService) FindCart(id uint) ([]interface{}, error) {

	return nil, nil

}

func (s UserService) CreateCart(input any, u domain.User) ([]interface{}, error) {

	return nil, nil

}

func (s UserService) CreateOrder(u domain.User) (int, error) {

	return 0, nil

}

func (s UserService) GetOrders(u domain.User) ([]interface{}, error) {

	return nil, nil

}

func (s UserService) GetOrderById(id uint, uId uint) (interface{}, error) {

	return nil, nil

}
