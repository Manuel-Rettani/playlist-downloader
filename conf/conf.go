package conf

type Conf struct {
	YoutubeKey     string `json:"youtube_key" yaml:"youtube_key"`
	PlayListId     string `json:"playlist_id" yaml:"playlist_id"`
	ChunkSize      int    `json:"chunk_size" yaml:"chunk_size"`
	Aws            Aws    `json:"aws" yaml:"aws"`
	MaxRetries     int    `json:"max_retries" yaml:"max_retries"`
	Email          Email  `json:"email" yaml:"email"`
	RequesterEmail string `json:"requester_email" yaml:"requester_email"`
}

type Aws struct {
	Region    string `yaml:"region" json:"region"`
	AccessKey string `json:"access_key_id" yaml:"access_key_id"`
	SecretKey string `json:"secret_access_key" yaml:"secret_access_key"`
	S3        S3     `json:"s3" yaml:"s3"`
}

type S3 struct {
	Bucket string `json:"bucket" yaml:"bucket"`
}

type Email struct {
	Email       string `json:"email" yaml:"email"`
	AppPassword string `json:"app_password" yaml:"app_password"`
	SmtpServer  string `json:"smtp_server" yaml:"smtp_server"`
	SmtpPort    int    `json:"smtp_port" yaml:"smtp_port"`
}
