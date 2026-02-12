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

// type declaration
type contextKey string
type UserMiddleware struct {
	UserStore store.UserStore
}


// ? - to avoid issues with dataType, that's why explictidly declared const for the type	
const UserContextKey = contextKey("user") 



//! functions having *req --> coz we are gonna modify it --> injects the context to the *req
func SetUser(r *http.Request,user *store.User) *http.Request {
	contxt := context.WithValue(r.Context(),UserContextKey,user)
	return r.WithContext(contxt)
}

//! get context attached user as client
func GetUser(r *http.Request) *store.User {
	user,ok := r.Context().Value(UserContextKey).(*store.User) //* pulling key from context from req and ensuring if it matces *store.User
	
	// if there is some sort of bad actor call
	if !ok {
		panic("missing user in request") //* as that context with key would been injected by above func that sets the context
	}

	return user
}


// mw fnc that do its and call the next func that proccesses rea --> this is like check fnc if passed calls the requested fnc
func (um *UserMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// processing req via this passed anony func
		w.Header().Add("Vary","Authorization")
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			r = SetUser(r,store.AnonymousUser)
			//! imp --> this { next.ServeHttp(passing rWriter,req)} calls the next function in the lineup chain
			next.ServeHTTP(w,r)
			return
		}

		headerParts := strings.Split(authHeader," ") // Bearer token --> need to learn more about this bearer thing
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			utils.WriteJson(w,http.StatusUnauthorized,utils.Envelope{"error":"invalid authorization header"})
			return 
		}

		// * we know other than that there will be a token there
		token := headerParts[1]
		user,err := um.UserStore.GetUserToken(tokens.ScopeAuth,token)
		if err != nil {
			utils.WriteJson(w,http.StatusUnauthorized,utils.Envelope{"error":"invalid token"})
			return
		}
		if user == nil {
			utils.WriteJson(w,http.StatusUnauthorized,utils.Envelope{"error":"token has been expired or invalid"})
			return
		}
		r = SetUser(r,user)
		next.ServeHTTP(w,r)
	})
}


func (um *UserMiddleware) RequireUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r)

		if user.IsAnonymousUser() {
			utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error": "you must be logged in to access this route"})
			return
		}
		// ? - if incoming user is not anony and has context attched (cause getUser returns user if r.Context has that key)
		next.ServeHTTP(w, r)
	})
}