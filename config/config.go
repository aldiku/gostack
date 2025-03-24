package config

import (
	"os"
	"regexp"
	"strconv"
)

type Config struct {
	AppName          string
	AppKey           string
	BaseUrl          string
	FrontEndUrl      string
	Environtment     string
	DatabaseUsername string
	DatabasePassword string
	DatabaseHost     string
	DatabasePort     string
	DatabaseName     string
	DatabaseURL      string
	PathDB           string
	PathDBPol        string
	CacheURL         string
	CachePassword    string
	LoggerLevel      string
	ContextTimeout   int
	JWTSecretKey     string
	MailMailer       string
	MailHost         string
	MailPort         int
	MailUsername     string
	MailPassword     string
	MailEncryption   string
	MidtransId       string
	MidtransClient   string
	MidtransServer   string
}

func LoadConfig() (config *Config) {
	appName := os.Getenv("APP_NAME")
	appKey := os.Getenv("APP_KEY")
	baseurl := os.Getenv("BASE_URL")
	frontendurl := os.Getenv("FRONT_END_URL")
	environment := os.Getenv("ENVIRONMENT")
	databaseUsername := os.Getenv("DATABASE_USERNAME")
	databasePassword := os.Getenv("DATABASE_PASSWORD")
	databaseHost := os.Getenv("DATABASE_HOST")
	databasePort := os.Getenv("DATABASE_PORT")
	databaseName := os.Getenv("DATABASE_NAME")
	databaseURL := os.Getenv("DATABASE_URL")
	PathDB := os.Getenv("PATH_DB")
	PathDBPol := os.Getenv("PATH_DB_POL")
	cacheURL := os.Getenv("CACHE_URL")
	cachePassword := os.Getenv("CACHE_PASSWORD")
	loggerLevel := os.Getenv("LOGGER_LEVEL")
	contextTimeout, _ := strconv.Atoi(os.Getenv("CONTEXT_TIMEOUT"))
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	mailMailer := os.Getenv("MAIL_MAILER")
	mailHost := os.Getenv("MAIL_HOST")
	mailPort, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	mailUsername := os.Getenv("MAIL_USERNAME")
	mailPassword := os.Getenv("MAIL_PASSWORD")
	mailEncryption := os.Getenv("MAIL_ENCRYPTION")
	midtransiD := os.Getenv("MIDTRANS_ID")
	midtransClient := os.Getenv("MIDTRANS_CLIENT")
	midtransServer := os.Getenv("MIDTRANS_SERVER")

	return &Config{
		AppName:          appName,
		AppKey:           appKey,
		BaseUrl:          baseurl,
		FrontEndUrl:      frontendurl,
		Environtment:     environment,
		DatabaseUsername: databaseUsername,
		DatabasePassword: databasePassword,
		DatabaseHost:     databaseHost,
		DatabasePort:     databasePort,
		DatabaseName:     databaseName,
		DatabaseURL:      databaseURL,
		PathDB:           PathDB,
		PathDBPol:        PathDBPol,
		CacheURL:         cacheURL,
		CachePassword:    cachePassword,
		LoggerLevel:      loggerLevel,
		ContextTimeout:   contextTimeout,
		JWTSecretKey:     jwtSecretKey,
		MailMailer:       mailMailer,
		MailHost:         mailHost,
		MailUsername:     mailUsername,
		MailPassword:     mailPassword,
		MailPort:         mailPort,
		MailEncryption:   mailEncryption,
		MidtransId:       midtransiD,
		MidtransClient:   midtransClient,
		MidtransServer:   midtransServer,
	}
}

func RootPath() string {
	projectDirName := os.Getenv("DIR_NAME")
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))
	return string(rootPath)
}

func PathDb() string {
	return os.Getenv("DATABASE_PATH")
}
