package conf

import (
	"gopkg.in/yaml.v3"
	"os"
)

func FromYaml(file string) (Conf, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return Conf{}, err
	}

	var conf Conf
	if err := yaml.Unmarshal(data, &conf); err != nil {
		return Conf{}, err
	}

	return conf, nil
}
