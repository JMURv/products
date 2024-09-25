package helpers

import (
	"bytes"
	"fmt"
	"github.com/JMURv/par-pro/products/internal/cache/redis"
	ctrl "github.com/JMURv/par-pro/products/internal/ctrl"
	handler "github.com/JMURv/par-pro/products/internal/hdl/http"
	"github.com/JMURv/par-pro/products/internal/repo/db"
	cfg "github.com/JMURv/par-pro/products/pkg/config"
	"github.com/goccy/go-json"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

var Conf *cfg.Config

func init() {
	Conf = cfg.MustLoad("../local.config.yaml")
}

func CreateUser(router *mux.Router) (user map[string]any, access string) {
	userData := map[string]string{
		"name":     "John Doe",
		"email":    "john@example.com",
		"password": "secret1234",
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	for key, val := range userData {
		_ = writer.WriteField(key, val)
	}

	err := writer.Close()

	rr := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/api/users", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(rr, req)

	var r map[string]any
	data, err := io.ReadAll(rr.Body)
	err = json.Unmarshal(data, &r)
	if err != nil {
		panic(err)
	}

	user, ok := r["data"].(map[string]any)["user"].(map[string]any)
	if !ok {
		panic(err)
	}

	cookies := rr.Result().Cookies()
	for _, v := range cookies {
		switch v.Name {
		case "access":
			access = v.Value
		}
	}

	return user, access
}

func CleanDB(t *testing.T) {
	conn, err := gorm.Open(
		postgres.Open(
			fmt.Sprintf(
				"postgres://%s:%s@%s:%v/%s",
				Conf.DB.User,
				Conf.DB.Password,
				Conf.DB.Host,
				Conf.DB.Port,
				Conf.DB.Database+"_test",
			),
		), &gorm.Config{},
	)
	if err != nil {
		t.Log(err)
	}

	sqlDB, err := conn.DB()
	if err != nil {
		t.Fatal(err)
	}
	defer sqlDB.Close()

	var tables []string
	if err := conn.Raw("SELECT tablename FROM pg_tables WHERE schemaname = 'public'").Scan(&tables).Error; err != nil {
		t.Fatal(err)
	}

	for _, table := range tables {
		if err := conn.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)).Error; err != nil {
			t.Fatal(err)
		}
	}

	t.Log("Database cleaned")
}

func SetupRouter() (router *mux.Router, cache ctrl.CacheRepo) {
	cache = redis.New(Conf.Redis)
	repo := db.New(&cfg.DBConfig{
		Host:     Conf.DB.Host,
		Port:     Conf.DB.Port,
		User:     Conf.DB.User,
		Password: Conf.DB.Password,
		Database: Conf.DB.Database + "_test",
	})

	svc := ctrl.New(repo, cache)
	h := handler.New(svc)

	router = mux.NewRouter()
	handler.RegisterAuthRoutes(router, h)
	handler.RegisterUserRoutes(router, h)

	return router, cache
}

func SendHttpRequest(t *testing.T, router *mux.Router, access string, method string, url string, body any) (map[string]any, *httptest.ResponseRecorder) {
	var err error
	var req *http.Request
	if body != nil {
		req, err = http.NewRequest(method, url, body.(*bytes.Buffer))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if !assert.NoError(t, err) {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	if access != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", access))
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return UnmarshallResponse(t, rr.Body), rr
}

func UnmarshallResponse(t *testing.T, rrBody *bytes.Buffer) (r map[string]any) {
	data, err := io.ReadAll(rrBody)
	if !assert.NoError(t, err) {
		t.Log(fmt.Sprintf("Error while reading: %v", err))
		t.Fatal(err)
	}

	if err = json.Unmarshal(data, &r); !assert.NoError(t, err) {
		t.Log(fmt.Sprintf("Error while unmarshalling: %v", err))
		t.Fatal(err)
	}
	return r
}
