package repository

import (
	"ecommerce/internal/domain"
	"errors"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository interface {
	CreateUser(u domain.User) (domain.User, error)
	FindUser(email string) (domain.User, error)
	FindUserById(id uint) (domain.User, error)
	UpdateUser(id uint, u domain.User) (domain.User, error)
	CreateBankAccount(e domain.BankAccount) error

	//cart
	FindCartItems(uId uint) ([]domain.Cart, error)
	FindCartItem(uId, pId uint) (domain.Cart, error)
	CreateCart(c domain.Cart) error
	UpdateCart(c domain.Cart) error
	DeleteCartById(id uint) error
	DeleteCartItems(uId uint) error

	// Order
	CreateOrder(o domain.Order) error
	FindOrders(uid uint) ([]domain.Order, error)
	FindOrderById(id uint) (domain.Order, error)
	FindUserOrderById(id uint, uId uint) (domain.Order, error)

	// Profile
	CreateAddress(e domain.Address) error
	UpdateAddress(e domain.Address) (*domain.Address, error)
}

type userRepository struct {
	db *gorm.DB
}

// FindUserOrderById implements UserRepository.
func (r *userRepository) FindUserOrderById(id uint, uId uint) (domain.Order, error) {
	var order domain.Order
	err := r.db.Where("order_ref=? AND user_id=?", id, uId).First(&order).Error
	if err != nil {
		log.Printf("find order by user id error %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return order, domain.ErrorOrderNotFound
		}
		return order, errors.New("order not found")
	}
	return order, nil
}

// FindOrderById implements UserRepository.
func (r *userRepository) FindOrderById(id uint) (domain.Order, error) {
	var order domain.Order
	err := r.db.Preload("Itens").First(&order, id).Error
	if err != nil {
		log.Printf("find order by id error %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return order, domain.ErrorOrderNotFound
		}
		return order, errors.New("order not found")
	}
	return order, nil
}

// FindOrders implements UserRepository.
func (r *userRepository) FindOrders(uid uint) ([]domain.Order, error) {

	var orders []domain.Order
	err := r.db.Preload("Itens").Where("user_id=?", uid).Find(&orders).Error
	if err != nil {
		log.Printf("find order by user id error %v", err)
		return orders, errors.New("failed to fetch orders")
	}
	return orders, nil

}

// CreateOrder implements UserRepository.
func (r *userRepository) CreateOrder(o domain.Order) error {
	err := r.db.Create(&o).Error
	if err != nil {
		log.Printf("create order error : %v", err)
		return errors.New("error creating order")
	}
	return nil
}

// CreateAddress implements UserRepository.
func (r *userRepository) CreateAddress(e domain.Address) error {
	err := r.db.Create(&e).Error
	if err != nil {
		log.Printf("create profile address error : %v", err)
		return errors.New("error creating profile address")
	}
	return nil
}

// UpdateAddress implements UserRepository.
func (r *userRepository) UpdateAddress(e domain.Address) (*domain.Address, error) {
	err := r.db.Where("user_id=?", e.UserId).Updates(&e).Error
	if err != nil {
		log.Printf("update profile address error : %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrorUserNotFound
		}
		return nil, errors.New("error updating profile address")
	}
	return &e, nil
}

// DeleteCartItems implements UserRepository.
func (r *userRepository) DeleteCartItems(uId uint) error {
	err := r.db.Where("user_id=?", uId).Delete(&domain.Cart{}).Error
	if err != nil {
		log.Printf("delete cart items error : %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrorCartItemNotFound
		}
		return errors.New("error deleting cart items")
	}
	return nil
}

// CreateCart implements UserRepository.
func (r *userRepository) CreateCart(c domain.Cart) error {
	err := r.db.Create(&c).Error
	if err != nil {
		log.Printf("create cart items error : %v", err)
		return errors.New("error creating cart items")
	}
	return nil
}

// DeleteCartById implements UserRepository.
func (r *userRepository) DeleteCartById(id uint) error {
	err := r.db.Delete(&domain.Cart{}, id).Error
	if err != nil {
		log.Printf("delete cart error : %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrorCartItemNotFound
		}
		return errors.New("error deleting cart")
	}
	return nil
}

// FindCartItem implements UserRepository.
func (r *userRepository) FindCartItem(uId uint, pId uint) (domain.Cart, error) {
	var cartItem domain.Cart
	err := r.db.Where("user_id=? AND product_id=?", uId, pId).First(&cartItem).Error
	if err != nil {
		log.Printf("error finding cart item for user %d and product %d : %v", uId, pId, err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return cartItem, domain.ErrorUserProductCartNotFound
		}
		return cartItem, errors.New("error finding cart item for user and product")
	}
	return cartItem, nil
}

// FindCartItems implements UserRepository.
func (r *userRepository) FindCartItems(uId uint) ([]domain.Cart, error) {
	var carts []domain.Cart
	err := r.db.Where("user_id=?", uId).Find(&carts).Error
	if err != nil {
		log.Printf("error finding cart item for user %d : %v", uId, err)
		return carts, errors.New("error finding cart items")
	}
	return carts, nil
}

// UpdateCart implements UserRepository.
func (r *userRepository) UpdateCart(c domain.Cart) error {
	var cart domain.Cart
	err := r.db.Model(&cart).Clauses(clause.Returning{}).Where("id=?", c.ID).Updates(c).Error
	if err != nil {
		log.Printf("error updating cart item : %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrorCartItemNotFound
		}
		return errors.New("error updating cart")
	}
	return nil
}

// CreateBankAccount implements UserRepository.
func (r userRepository) CreateBankAccount(e domain.BankAccount) error {
	err := r.db.Create(&e).Error
	if err != nil {
		log.Printf("error creating bank account info : %v", err)
		return errors.New("error creating bank account")
	}
	return nil
}

func (r userRepository) CreateUser(usr domain.User) (domain.User, error) {

	err := r.db.Create(&usr).Error
	if err != nil {
		log.Printf("create user error %v", err)
		return domain.User{}, errors.New("failed to create user")
	}

	return usr, nil
}

func (r userRepository) FindUser(email string) (domain.User, error) {

	var user domain.User

	err := r.db.Preload("Address").First(&user, "email=?", email).Error
	if err != nil {
		log.Printf("find user by email error %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, errors.New("user doesn't exists")
		}
		return domain.User{}, errors.New("user doesn't exists")
	}

	return user, nil
}

func (r userRepository) FindUserById(id uint) (domain.User, error) {
	var user domain.User

	err := r.db.Preload("Address").
		Preload("Cart").
		Preload("Orders").
		First(&user, id).Error
	if err != nil {
		log.Printf("find user by id error %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, domain.ErrorUserNotFound
		}
		return domain.User{}, errors.New("user doesn't exists")
	}

	return user, nil
}

func (r userRepository) UpdateUser(id uint, u domain.User) (domain.User, error) {

	var user domain.User

	err := r.db.Model(&user).Clauses(clause.Returning{}).Where("id=?", id).Updates(u).Error
	if err != nil {
		log.Printf("error on update %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, domain.ErrorUserNotFound
		}
		return domain.User{}, errors.New("failed update user")
	}

	return user, nil
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}
