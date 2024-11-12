package config

import "github.com/spf13/viper"

type App struct {
	AppPort string `json:"app_port"`
	AppEnv  string `json:"app_env"`

	JwtSecretKey          string `json:"jwt_key"`
	JwtIssuer             string `json:"jwt_issuer"`
	JwtDurationAccessKey  int    `json:"jwt_duration_access_key"`
	JwtDurationRefreshKey int    `json:"jwt_duration_refresh_key"`
}

type PostgresDB struct {
	DBPort    string `json:"port"`
	DBHost    string `json:"host"`
	DBUser    string `json:"user"`
	DBPass    string `json:"pass"`
	DBName    string `json:"name"`
	DBMaxOpen int    `json:"max_open"`
	DBMaxIdle int    `json:"max_idle"`
}

type CloudflareR2 struct {
	BucketName string `json:"bucket_name"`
	ApiKey     string `json:"api_key"`
	ApiSecret  string `json:"api_secret"`
	Token      string `json:"token"`
	AccountID  string `json:"account_id"`
	PublicUrl  string `json:"public_url"`
}

type Config struct {
	App          App
	PostgresDB   PostgresDB
	CloudflareR2 CloudflareR2
}

func NewConfig() Config {
	return Config{
		App: App{
			AppPort:               viper.GetString("APP_PORT"),
			AppEnv:                viper.GetString("APP_ENV"),
			JwtSecretKey:          viper.GetString("JWT_SECRET_KEY"),
			JwtIssuer:             viper.GetString("JWT_ISSUER"),
			JwtDurationAccessKey:  viper.GetInt("JWT_DURATION_ACCESS_KEY"),
			JwtDurationRefreshKey: viper.GetInt("JWT_DURATION_REFRESH_KEY"),
		},
		PostgresDB: PostgresDB{
			DBPort:    viper.GetString("DATABASE_PORT"),
			DBHost:    viper.GetString("DATABASE_HOST"),
			DBUser:    viper.GetString("DATABASE_USER"),
			DBPass:    viper.GetString("DATABASE_PASSWORD"),
			DBName:    viper.GetString("DATABASE_NAME"),
			DBMaxOpen: 0,
			DBMaxIdle: 0,
		},
		CloudflareR2: CloudflareR2{
			BucketName: viper.GetString("CLOUDFLARE_R2_BUCKET_NAME"),
			ApiKey:     viper.GetString("CLOUDFLARE_R2_API_KEY"),
			ApiSecret:  viper.GetString("CLOUDFLARE_R2_API_SECRET"),
			Token:      viper.GetString("CLOUDFLARE_R2_TOKEN"),
			AccountID:  viper.GetString("CLOUDFLARE_R2_ACCOUNT_ID"),
			PublicUrl:  viper.GetString("CLOUDFLARE_R2_PUBLIC_URL"),
		},
	}
}
