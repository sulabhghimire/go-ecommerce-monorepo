package service

import (
	"ecommerce/config"
	"ecommerce/internal/domain"
	"ecommerce/internal/dto"
	"ecommerce/internal/helper"
	"ecommerce/internal/repository"
	"ecommerce/pkg/notification"
	"errors"
	"fmt"
	"log"
	"time"
)

type UserService struct {
	Repo   repository.UserRepository
	PRepo  repository.ProductRepository
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

	if err != nil {
		return "", err
	}

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

func (s UserService) CreateProfile(id uint, input dto.ProfileInput) error {

	_, err := s.Repo.UpdateUser(id, domain.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
	})
	if err != nil {
		return err
	}

	address := domain.Address{
		AddressLine1: input.AddressInput.AddressLine1,
		AddressLine2: input.AddressInput.AddressLine2,
		City:         input.AddressInput.City,
		Country:      input.AddressInput.Country,
		PostCode:     input.AddressInput.PostCode,
		UserId:       id,
	}

	err = s.Repo.CreateAddress(address)
	if err != nil {
		return err
	}

	return nil

}

func (s UserService) GetProfile(id uint) (*domain.User, error) {

	user, err := s.Repo.FindUserById(id)
	if err != nil {
		return nil, err
	}

	return &user, nil

}

func (s UserService) UpdateProfile(id uint, input dto.ProfileInput) (*domain.User, error) {

	user, err := s.Repo.FindUserById(id)
	if err != nil {
		return nil, err
	}

	if input.FirstName != "" {
		user.FirstName = input.FirstName
	}
	if input.LastName != "" {
		user.LastName = input.LastName
	}

	updatedUser, err := s.Repo.UpdateUser(id, user)
	if err != nil {
		return nil, err
	}

	address := domain.Address{
		AddressLine1: input.AddressInput.AddressLine1,
		AddressLine2: input.AddressInput.AddressLine2,
		City:         input.AddressInput.City,
		Country:      input.AddressInput.Country,
		PostCode:     input.AddressInput.PostCode,
		UserId:       id,
	}
	updatedAddress, err := s.Repo.UpdateAddress(address)
	if err != nil {
		return nil, err
	}

	updatedUser.Address = *updatedAddress

	return &updatedUser, nil

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
		fmt.Printf("Error occurred while creating bank record for user %v", err)
		return "", errors.New("error creating bank info")
	}

	return token, nil

}

func (s UserService) FindCart(id uint) ([]domain.Cart, float64, error) {

	cartItems, err := s.Repo.FindCartItems(id)
	if err != nil {
		return nil, 0, err
	}

	var totalAmount float64
	for _, item := range cartItems {
		totalAmount += item.Price * float64(item.Qty)
	}

	return cartItems, totalAmount, nil

}

func (s UserService) CreateCart(input dto.CreateCartRequest, u domain.User) ([]domain.Cart, error) {

	if input.ProductId == 0 {
		return nil, errors.New("please provide a valid product id")
	}

	cart, _ := s.Repo.FindCartItem(u.ID, input.ProductId)
	if cart.ID > 0 {

		if input.Quantity < 1 {

			err := s.Repo.DeleteCartById(cart.ID)
			if err != nil {
				log.Panicln("Error on deleting cart item", err)
				return nil, errors.New("error on deleting cart item")
			}
		} else {
			cart.Qty = input.Quantity
			err := s.Repo.UpdateCart(cart)
			if err != nil {
				log.Panicln("Error on updating cart item", err)
				return nil, errors.New("error on updating cart item")
			}
		}

	} else {

		product, err := s.PRepo.GetProductById(input.ProductId)
		if err != nil {
			return nil, domain.ProductNotFound
		}

		err = s.Repo.CreateCart(domain.Cart{
			ProductId: product.ID,
			UserId:    u.ID,
			Name:      product.Name,
			Qty:       input.Quantity,
			Price:     product.Price,
			SellerId:  product.UserId,
			ImageUrl:  product.ImageUrl,
		})

		if err != nil {
			return nil, errors.New("error creating cart itrm")
		}

	}

	return s.Repo.FindCartItems(u.ID)

}

func (s UserService) CreateOrder(u domain.User) (string, error) {

	// get cart items
	cartItems, amount, err := s.FindCart(u.ID)
	if err != nil {
		return "", err
	}
	if len(cartItems) == 0 {
		return "", domain.ErrorCartItemNotFound
	}

	// process payment
	paymentId := "PAY12345678"
	txId := "TX12345678"
	orderRef, _ := helper.RandomString(8)

	// create order with generated order ref
	var orderItems []domain.OrderItem

	for _, item := range cartItems {
		orderItems = append(orderItems, domain.OrderItem{
			ProductId: item.ProductId,
			UserId:    item.UserId,
			Name:      item.Name,
			ImageUrl:  item.ImageUrl,
			SellerId:  item.SellerId,
			Price:     item.Price,
			Qty:       item.Qty,
		})
	}

	order := domain.Order{
		UserId:        u.ID,
		Amount:        amount,
		TransactionId: txId,
		OrderRef:      orderRef,
		PaymentId:     paymentId,
		Items:         orderItems,
	}

	err = s.Repo.CreateOrder(order)
	if err != nil {
		return "", err
	}

	// send order confirmation email to user

	// remove cart items
	err = s.Repo.DeleteCartItems(u.ID)
	if err != nil {
		return "", err
	}

	// return order number

	return orderRef, nil

}

func (s UserService) GetOrders(uId uint) ([]domain.Order, error) {

	return s.Repo.FindOrders(uId)

}

func (s UserService) GetOrderById(id string, uId uint) (domain.Order, error) {

	return s.Repo.FindUserOrderById(id, uId)

}
