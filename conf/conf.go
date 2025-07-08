package conf

type Conf struct {
	Keys       Keys   `json:"keys" yaml:"keys"`
	PlayListId string `json:"playlist_id" yaml:"playlist_id"`
	ChunkSize  int    `json:"chunk_size" yaml:"chunk_size"`
	Aws        Aws    `json:"aws" yaml:"aws"`
}

type Keys struct {
	Youtube string `json:"youtube" yaml:"youtube"`
}

type Aws struct {
	AccessKey string `json:"access_key_id" yaml:"access_key_id"`
	SecretKey string `json:"secret_access_key" yaml:"secret_access_key"`
	S3        S3     `json:"s3" yaml:"s3"`
}

type S3 struct {
	Bucket string `json:"bucket" yaml:"bucket"`
}
