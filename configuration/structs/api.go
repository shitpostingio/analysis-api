package structs

// Config represents the analysis-api configuration.
type Config struct {
	Database            DatabaseConfiguration
	Redis               RedisConfiguration
	NSFWEndpoint        string
	FingerprintEndpoint string
	GibberishEndpoint   string
	Testing             bool `type:"optional"`
	MaxDownloadSize 	int64
}
