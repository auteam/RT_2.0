package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"

	//"github.com/auth0/go-jwt-middleware"
	"../model"
	"../store"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const (
	ctxKeyUser ctxKey = iota
	ctxKeyRequestID
)

// Claims ...
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var (
	errIncorrectEmailOrPassword = errors.New("Incorrect email or password")
	errNotAuthenticated         = errors.New("Not authenticated")
	errIncorrectEmail           = errors.New("Email not found")
	errIncorrectPassword        = errors.New("Incorrect password")
	errEmailIsUsed              = errors.New("Email is already used")
	errNoData                   = errors.New("No Data")
	errPermission               = errors.New("Permission Denied")
	errNoChamp                  = errors.New("You are not registered for more than one championship, or the administrator has not confirmed that you are a participant. If you are a member and see this message, please contact technical support")
)

type ctxKey int8

type server struct {
	router           *mux.Router
	logger           *logrus.Logger
	store            *store.Store
	secretKey        string
	secretKeyRefresh string

	email string
	champ string
}

// Start ...
func Start() error {
	Config := store.Configurate()
	db, err := store.NewDB(Config.DatabaseURL)
	if err != nil {
		return err
	}

	defer db.Close()
	store := store.New(db)

	srv := newServer(store)

	handler := cors.Default().Handler(srv)

	return http.ListenAndServe(Config.BindAddr, handler)
}

func newServer(store *store.Store) *server {
	s := &server{
		router: mux.NewRouter(),
		logger: logrus.New(),
		store:  store,

		secretKey:        "Bbibaboba",
		secretKeyRefresh: "abobabibB",
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))

	s.router.HandleFunc("/registration", s.handleUsersCreate()).Methods("POST")
	s.router.HandleFunc("/authorization", s.handleTokenCreate()).Methods("POST")

	s.router.HandleFunc("/main", s.JWTMiddleware(s.MainCheck())).Methods("POST")
	s.router.HandleFunc("/admin", s.JWTMiddleware(s.Admin())).Methods("POST")

	s.router.HandleFunc("/topology", s.JWTMiddleware(s.Topology())).Methods("POST")
	s.router.HandleFunc("/topologyVNC", s.JWTMiddleware(s.VNCTopology())).Methods("POST")

	s.router.HandleFunc("/topology/create", s.AdminMiddleware(s.CreateTopology())).Methods("POST")
	s.router.HandleFunc("/topology/save", s.AdminMiddleware(s.SaveTopology())).Methods("POST")
	s.router.HandleFunc("/topology/remove", s.AdminMiddleware(s.RemoveTopology())).Methods("POST")
	s.router.HandleFunc("/topology/clone", s.AdminMiddleware(s.CloneTopology())).Methods("POST")
	s.router.HandleFunc("/topology/get", s.AdminMiddleware(s.GetTopology())).Methods("POST")
	s.router.HandleFunc("/topology/link", s.AdminMiddleware(s.TopologyLink())).Methods("POST")

	s.router.HandleFunc("/champ/create", s.AdminMiddleware(s.CreateChamp())).Methods("POST")
	s.router.HandleFunc("/champ/remove", s.AdminMiddleware(s.DeleteChamp())).Methods("POST")
	s.router.HandleFunc("/champ/get", s.AdminMiddleware(s.AllChamp())).Methods("GET")

	s.router.HandleFunc("/module/create", s.AdminMiddleware(s.CreateModule())).Methods("POST")
	s.router.HandleFunc("/module/remove", s.AdminMiddleware(s.DeleteModule())).Methods("POST")
	s.router.HandleFunc("/module/get", s.AdminMiddleware(s.GetModule())).Methods("GET")

	s.router.HandleFunc("/stand/create", s.AdminMiddleware(s.CreateStand())).Methods("POST")
	s.router.HandleFunc("/stand/remove", s.AdminMiddleware(s.RemoveStand())).Methods("POST")
	s.router.HandleFunc("/stand/update", s.AdminMiddleware(s.UpdateStand())).Methods("POST")
	s.router.HandleFunc("/stand/allupdate", s.AdminMiddleware(s.AllUpdateStand())).Methods("POST")
	s.router.HandleFunc("/stand/get", s.AdminMiddleware(s.AllStand())).Methods("POST")

	s.router.HandleFunc("/device/ticket", s.JWTMiddleware(s.GetTicket())).Methods("POST")
	s.router.HandleFunc("/device/clear", s.JWTMiddleware(s.Snapshot())).Methods("POST")

	s.router.HandleFunc("/admin/settime", s.AdminMiddleware(s.SetTime())).Methods("POST")

	s.router.HandleFunc("/admin/userfromcsv", s.AdminMiddleware(s.UserFromCSV())).Methods("POST")
	s.router.HandleFunc("/admin/userfromcsv/create", s.AdminMiddleware(s.UserFromCSVCreate())).Methods("POST")

	s.router.HandleFunc("/admin/addtochampcsv", s.AdminMiddleware(s.AddToChampCSV())).Methods("POST")
	s.router.HandleFunc("/admin/addtochampcsv/create", s.AdminMiddleware(s.AddToChampCSVCreate())).Methods("POST")

	s.router.HandleFunc("/admin/standfromcsv", s.AdminMiddleware(s.StandFromCsv())).Methods("POST")
	s.router.HandleFunc("/admin/standfromcsv/create", s.AdminMiddleware(s.StandFromCsvCreate())).Methods("POST")

	s.router.HandleFunc("/admin/alluser", s.AdminMiddleware(s.AllUser())).Methods("GET")
	s.router.HandleFunc("/admin/resetpass", s.AdminMiddleware(s.ResetPass())).Methods("POST")
	s.router.HandleFunc("/admin/addtochamp", s.AdminMiddleware(s.AddToChamp())).Methods("POST")
	s.router.HandleFunc("/admin/addtomodule", s.AdminMiddleware(s.AddToModule())).Methods("POST")
	s.router.HandleFunc("/admin/changename", s.AdminMiddleware(s.ChangeName())).Methods("POST")

	
	s.router.HandleFunc("/admin/trystate", s.AdminMiddleware(s.ResetTryState())).Methods("POST")
}

/* Post Init HTTP*/
func (s *server) MainCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := s.store.User().FindByEmail(s.email)
		if err != nil {

		}
		if u.Role != "admin" {
			w.Header().Set("Content-Type", "application/json")
			s.Main().ServeHTTP(w, r)
			return
		}
		s.AllChamp().ServeHTTP(w, r)
	}
}

func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(ctxKeyRequestID),
		})
		logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		var level logrus.Level
		switch {
		case rw.code >= 500:
			level = logrus.ErrorLevel
		case rw.code >= 400:
			level = logrus.WarnLevel
		default:
			level = logrus.InfoLevel
		}
		logger.Logf(
			level,
			"completed with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start),
		)
	})
}

/*Create User ... REGISTRATION*/
func (s *server) handleUsersCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := &model.User{
			Email:    req.Email,
			Password: req.Password,
			Name:     req.Name,
		}

		if len(req.Email) == 0 {
			s.error(w, r, http.StatusBadRequest, errors.New("bad email"))
			return
		}
		if len(req.Password) == 0 {
			s.error(w, r, http.StatusBadRequest, errors.New("bad password"))
			return
		}
		if len(req.Name) == 0 {
			s.error(w, r, http.StatusBadRequest, errors.New("bad name"))
			return
		}


		if err := s.store.User().Create(u); err != nil {
			s.error(w, r, http.StatusTeapot, err)
			return
		}

		u.Sanitize()

		//u, tokenString, tokenStringRefresh, errT := s.myjwt(w, r, req.Email, req.Password)
		//if errT != 0 {
		//	return
		//}
		//
		//expirationTimeRefresh := time.Now().Add(5 * 60 * 24 * time.Minute)
		//
		//http.SetCookie(w, &http.Cookie{
		//	Name:     "tokenAccess",
		//	Value:    tokenString,
		//	Expires:  expirationTimeRefresh,
		//	HttpOnly: false,
		//	Path:     "/",
		//})
		//http.SetCookie(w, &http.Cookie{
		//	Name:     "tokenRefresh",
		//	Value:    tokenStringRefresh,
		//	Expires:  expirationTimeRefresh,
		//	HttpOnly: false,
		//	Path:     "/",
		//})

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"status\":\"OK\",\"role\":\"" + u.Role + "\"}"))
		w.WriteHeader(http.StatusCreated)
	}
}

/* Tokens */
//Refresh
func (s *server) Refresh(w http.ResponseWriter, r *http.Request) {
	//c, err := r.Cookie("tokenRefresh")

	//tknStr := c.Value
	claims := &Claims{}
	// tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
	// 	return s.secretKeyRefresh, nil
	// })

	// if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	//return
	// }

	expirationTime := time.Now().Add(1 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "tokenAccess",
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: false,
		Path:     "/",
	})

	expirationTimeRefresh := time.Now().Add(5 * 24 * 60 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	tokenRefresh := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = tokenRefresh.SignedString(s.secretKeyRefresh)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "tokenRefresh",
		Value:    tokenString,
		Expires:  expirationTimeRefresh,
		HttpOnly: false,
		Path:     "/",
	})
	w.Write([]byte("{\"status\":\"OK\"}"))
	w.WriteHeader(http.StatusOK)
}

//Create Token
func (s *server) handleTokenCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"Name"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "10.11.4.18")
		// w.Header().Set("Access-Control-Allow-Headers", "*")
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		err := s.store.User().CheckTryState(req.Email)
		if err != nil {
			s.error(w, r, http.StatusTeapot, err)
			return
		}


		u, tokenString, tokenStringRefresh, errT := s.myjwt(w, r, req.Email, req.Password)
		if errT != 0 {
			return
		}

		expirationTimeRefresh := time.Now().Add(5 * 60 * 24 * time.Minute)

		http.SetCookie(w, &http.Cookie{
			Name:     "tokenAccess",
			Value:    tokenString,
			Expires:  expirationTimeRefresh,
			HttpOnly: false,
			Path:     "/",
		})
		http.SetCookie(w, &http.Cookie{
			Name:     "tokenRefresh",
			Value:    tokenStringRefresh,
			Expires:  expirationTimeRefresh,
			HttpOnly: false,
			Path:     "/",
		})

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"status\":\"OK\", \"role\":\"" + u.Role + "\"}"))
		w.WriteHeader(http.StatusOK)
	}
}

/*JWT Validate Token JWTMiddleware*/
func (s *server) JWTMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("tokenRefresh")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tknStr := c.Value
		claims := &Claims{}

		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(s.secretKeyRefresh), nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		//
		c, err = r.Cookie("tokenAccess")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tknStr = c.Value
		claims = &Claims{}

		tknRefresh, errR := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(s.secretKey), nil
		})
		if errR != nil {
			if errR == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		if !tkn.Valid {
			if !tknRefresh.Valid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.Write([]byte(fmt.Sprintf("REFRESH \n")))
			s.Refresh(w, r)
		}

		s.email = claims.Username
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	}
}

/*JWT Validate Token AdminMiddleware*/
func (s *server) AdminMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("tokenRefresh")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tknStr := c.Value
		claims := &Claims{}

		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(s.secretKeyRefresh), nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		//
		c, err = r.Cookie("tokenAccess")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tknStr = c.Value
		claims = &Claims{}

		tknRefresh, errR := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(s.secretKey), nil
		})
		if errR != nil {
			if errR == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		if !tkn.Valid {
			if !tknRefresh.Valid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.Write([]byte(fmt.Sprintf("REFRESH \n")))
			s.Refresh(w, r)
		}

		u, err := s.store.User().FindByEmail(claims.Username)
		if u.Role != "admin" {
			w.Header().Set("Content-Type", "application/json")
			s.error(w, r, http.StatusForbidden, errPermission)
			return
		}
		s.email = claims.Username
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	}
}

//Welcome Test
func (s *server) Welcome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//fmt.Printf("%v", claims.Username)
		//w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf("Welcome %v!", s.email)))
		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"status": strconv.Itoa(code), "error": err.Error()})
}

func (s *server) ValidateSQL(t string) string {
	return ""
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
