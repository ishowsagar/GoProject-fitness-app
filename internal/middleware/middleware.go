package middleware

// importing packages
import (
	"context"
	"fem/internal/store"
	"fem/internal/tokens"
	"fem/internal/utils"
	"net/http"
	"strings"
)

//! type declaration
type contextKey string //* custom type for context keys to avoid collisions
type UserMiddleware struct {
	UserStore store.UserStore //* needed to fetch user from token
}


//! UserContextKey --> unique key for storing user in request context
//? using custom type prevents accidental key conflicts with other middleware
const UserContextKey = contextKey("user") 



//! SetUser --> injects user into request context for downstream handlers
//! Takes pointer because we modify the request by adding context
func SetUser(r *http.Request,user *store.User) *http.Request {
	//* create new context with user value attached
	contxt := context.WithValue(r.Context(),UserContextKey,user)
	return r.WithContext(contxt) //* return modified request
}

//! GetUser --> extracts user from request context
//! Panics if user not found (should never happen after Authenticate middleware)
func GetUser(r *http.Request) *store.User {
	//* extract value from context and type assert to *store.User
	user,ok := r.Context().Value(UserContextKey).(*store.User)
	
	//? if type assertion fails, middleware wasn't run properly
	if !ok {
		panic("missing user in request") //* critical error - indicates middleware chain broken
	}

	return user
}


//! Authenticate --> middleware that validates Bearer token from Authorization header
//! Sets user in context (either authenticated user or AnonymousUser)
func (um *UserMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//* tell caches that response varies by Authorization header
		w.Header().Add("Vary","Authorization")
		//* extract Authorization header from request
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			//* no token provided, set as anonymous user
			r = SetUser(r,store.AnonymousUser)
			//* call next handler in chain
			next.ServeHTTP(w,r)
			return
		}

		//* split "Bearer TOKEN" into ["Bearer", "TOKEN"]
		headerParts := strings.Split(authHeader," ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			//? header format wrong (should be: "Bearer <token>")
			utils.WriteJson(w,http.StatusUnauthorized,utils.Envelope{"error":"invalid authorization header"})
			return 
		}

		//* extract token string (second part after "Bearer ")
		token := headerParts[1]
		//* lookup user by token hash in database
		user,err := um.UserStore.GetUserToken(tokens.ScopeAuth,token)
		if err != nil {
			//? database error or token not found
			utils.WriteJson(w,http.StatusUnauthorized,utils.Envelope{"error":"invalid token"})
			return
		}
		if user == nil {
			//? token expired or doesn't exist
			utils.WriteJson(w,http.StatusUnauthorized,utils.Envelope{"error":"token has been expired or invalid"})
			return
		}
		//* valid token! attach authenticated user to request context
		r = SetUser(r,user)
		next.ServeHTTP(w,r) //* call next handler with authenticated user
	})
}


//! RequireUser --> ensures user is authenticated (not anonymous)
//! Must be used after Authenticate middleware
func (um *UserMiddleware) RequireUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//* get user from context (set by Authenticate middleware)
		user := GetUser(r)

		if user.IsAnonymousUser() {
			//? user didn't provide valid token
			utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error": "you must be logged in to access this route"})
			return
		}
		//* user is authenticated, proceed to handler
		next.ServeHTTP(w, r)
	})
}