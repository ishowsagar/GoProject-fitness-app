package api

import (
	"encoding/json"
	"errors"
	"fem/internal/store"
	"fem/internal/utils"
	"log"
	"net/http"
	"regexp"
)

//  types declaration
type registerUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Bio      string `json:"bio"` 
}

type UserHandler struct {
	userStore store.UserStore
	logger *log.Logger
}

func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	// ! return instace of struct --> so this typw now could be accessible
	// * All those methods who use this as injected agrs , will be accessible through the type dot notation
	return &UserHandler{
		userStore: userStore,
		logger: logger,
	}
}

func (h *UserHandler) validateUserRegisterRequest (regUser *registerUserRequest) error {
	// ? - server side validation for incoming user data from signups
	if regUser.Username == "" {
	return	errors.New("Username is required")
	}

	if len(regUser.Username) > 50 {
	return	errors.New("Username cannot be greater than 50 characters")
	}

	if regUser.Email == "" {
	return	errors.New("Email is required")
	}
	
	if regUser.Password == "" {
		return errors.New("password is required")
	}

	//! regex tests
	emailRegexPattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegexPattern.MatchString(regUser.Email)  {
		// passing in email which to test 
		return errors.New("Invalid email format")
	} 

	return nil

}

func (h *UserHandler) HandleRegisterUser(w http.ResponseWriter, req *http.Request) {
	var r registerUserRequest

	err:= json.NewDecoder(req.Body).Decode(&r)
	if err!= nil {
		h.logger.Printf("Error : decoding Register request : %V ",err)
		utils.WriteJson(w,http.StatusBadRequest,utils.Envelope{"error":"invalid request payload"})
		return
	}

	err = h.validateUserRegisterRequest(&r)
	if err!= nil {
		// h.logger.Printf("Error : decoding Register request : %V ",err)
		utils.WriteJson(w,http.StatusBadRequest,utils.Envelope{"error":err.Error()})
		return
	}

	user := &store.User{
		Username: r.Username,
		Email: r.Email,
	}
	if r.Bio == "" {
		user.Bio = r.Bio
	}

	// dealing with client's pass with our Password validation func that we have imported
	err = user.PasswordHash.Set(r.Password)
	if err != nil {
		h.logger.Printf("ERROR : hashing password %v ",err)
		utils.WriteJson(w,http.StatusInternalServerError,utils.Envelope{"error":"internal server error"})
		return
	}

	err = h.userStore.CreateUser(user)
	if err != nil {
		h.logger.Printf("ERROR : registering user %v ",err)
		utils.WriteJson(w,http.StatusInternalServerError,utils.Envelope{"error":"internal server error"})
		return
	}

	// * successfully created the user
		utils.WriteJson(w,http.StatusCreated,utils.Envelope{"user":user })

}