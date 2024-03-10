package config

import "os"

type Config struct {
	DBPath             string
	CloudflareAPIToken string
	CloudflareZoneID   string
	CloudflareRuleID   string
	MandrillAPIKey     string
	SenderEmail        string
	RealIPHeader       string
}

func New() *Config {
	return &Config{
		DBPath:             os.Getenv("DATABASE_PATH"),
		CloudflareAPIToken: os.Getenv("CLOUDFLARE_API_TOKEN"),
		CloudflareZoneID:   os.Getenv("CLOUDFLARE_ZONE_ID"),
		CloudflareRuleID:   os.Getenv("CLOUDFLARE_RULE_ID"),
		MandrillAPIKey:     os.Getenv("MANDRILL_API_KEY"),
		SenderEmail:        os.Getenv("EMAIL_FROM"),
		RealIPHeader:       os.Getenv("REAL_IP_HEADER"),
	}
}
