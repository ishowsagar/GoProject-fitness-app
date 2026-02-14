package api

import (
	"encoding/json"
	"fem/internal/store"
	"fem/internal/tokens"
	"fem/internal/utils"
	"log"
	"net/http"
	"time"
)

type TokenHandler struct {
	tokenStore store.TokenStore //* for creating/storing tokens
	userStore store.UserStore //* for validating user credentials
	logger *log.Logger //* for error logging
}

//! createTokenRequest --> login credentials from client
type createTokenRequest struct {
	Username string `json:"username"` //* user's login name
	Password string `json:"password"` //* plaintext password to verify
}

//! NewTokenHandler --> constructor for token handler
func NewTokenHandler(tokenStore store.TokenStore,userStore store.UserStore,logger *log.Logger) *TokenHandler {
	return &TokenHandler{
		tokenStore: tokenStore,
		userStore: userStore,
		logger: logger,
	}
}

// func NewTokenHandler(tokenstore store.TokenStore,userStore store.UserStore,logger *log.Logger) *TokenHandler{
// 	// return instance of TokenHandler type struct
// 	return &TokenHandler{
// 		tokenStore: tokenstore,
// 		userStore: userStore,
// 		logger : logger,

// 	}
// }

//! HandleCreateToken --> POST /tokens/authentication (login endpoint)
//! Validates credentials and returns authentication token
func (h *TokenHandler) HandleCreateToken(w http.ResponseWriter,req *http.Request)  {
	var tokenRequestingUser createTokenRequest //* holds username and password from client
	//* decode JSON body into struct
	err := json.NewDecoder(req.Body).Decode(&tokenRequestingUser)
	if err!= nil {
		h.logger.Printf("ERROR : createTokenRequest %v", err)
		utils.WriteJson(w,http.StatusBadRequest,utils.Envelope{"error":"invalid payload request"})
	}
	//* lookup user by username in database
	user, err := h.userStore.GetUserByUsername(tokenRequestingUser.Username)
	if err != nil || user == nil {
		h.logger.Printf("ERROR: GetUserByUsername: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	//! comparing plaintext password with hashed password using bcrypt
	passwordsDoMatch, err := user.PasswordHash.Matches(tokenRequestingUser.Password)
	if err != nil {
		h.logger.Printf("ERORR: PasswordHash.Mathes %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if !passwordsDoMatch {
		//? wrong password
		utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid credentials"})
		return
	}

	//* credentials valid! generate new authentication token (expire in 24 hours)
	token, err := h.tokenStore.CreateNewToken(user.ID, 24*time.Hour, tokens.ScopeAuth)
	if err != nil {
		h.logger.Printf("ERORR: Creating Token %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return

	}

	//* return token to client (they'll use this in Authorization header for protected routes)
	utils.WriteJson(w, http.StatusCreated, utils.Envelope{"auth_token": token})
}