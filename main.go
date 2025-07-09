package main

// maxprocs used as it's cgroup aware.  remove after runtime updates.
// https://github.com/uber-go/automaxprocs/issues/96

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"github.com/tristanfisher/patchpanel"
	"io"
	"log"
	"os"
	"reflect"
	"seedyfv2/header"
	"time"
)

type Config struct {
	FilePath             string        `default:"./input/private.xlsx"`
	MaxRuntime           time.Duration `default:"60m"`
	MaxFileSizeMegabytes int           `default:"1000"`
}

func ParseConfig(configPath string, configStruct Config) (*Config, error) {
	viperConf := viper.New()
	patch := patchpanel.NewPatchPanel(patchpanel.TokenSeparator, patchpanel.KeyValueSeparator)

	// get defaults off struct
	confType := reflect.TypeOf(configStruct)
	for i := 0; i < confType.NumField(); i++ {
		fieldVal, err := patch.GetDefault(confType.Field(i).Name, confType, []string{})
		if err != nil {
			return &Config{}, err
		}
		viperConf.SetDefault(confType.Field(i).Name, fieldVal)
	}

	// check configuration file
	if configPath != "" {
		viperConf.SetConfigFile(configPath)
		err := viperConf.ReadInConfig()
		if err != nil {
			var configFileNotFoundError viper.ConfigFileNotFoundError
			if errors.As(err, &configFileNotFoundError) {
				return nil, fmt.Errorf("file not found: %s", err)
			}
			return &Config{}, err
		}
	}

	viperConf.AutomaticEnv()

	err := viperConf.Unmarshal(&configStruct)
	if err != nil {
		return &Config{}, err
	}

	return &configStruct, nil
}

// LoadFile takes a max size to allow returning an error before potentially OOMing
func LoadFile(filePath string, maxSizeBytes int64) (io.Reader, error) {
	// note that we do not defer the file handle close
	// in order to be explicit
	f, err := os.Open(filePath)
	fClose := func(f *os.File, e error) error {
		fCloseErr := f.Close()
		if fCloseErr != nil {
			e = errors.Join(e, fCloseErr)
		}
		return e
	}

	if f == nil {
		return nil, err
	}

	lRead := io.LimitReader(f, maxSizeBytes)
	buf := bytes.NewBuffer([]byte{})
	bytesWritten, err := io.Copy(buf, lRead)
	if err != nil {
		err = fClose(f, err)
		return nil, err
	}
	if bytesWritten == 0 {
		return nil, errors.New("zero length file")
	}

	// check that we read the entire file
	stat, err := f.Stat()
	if err != nil {
		err = fClose(f, err)
		return nil, err
	}

	if bytesWritten < stat.Size() {
		return nil, errors.New("file larger than limit")
	}

	// with our file read in, close the file handle
	err = f.Close()
	if err != nil {
		return nil, err
	}

	return buf, nil

}

func main() {

	logInfo := log.New(os.Stdout, "", log.LstdFlags)
	_ = logInfo

	logErr := log.New(os.Stderr, "", log.LstdFlags)

	valuesFile := patchpanel.GetFileEnvOrPath(patchpanel.ENV_CONFIG_FILE, patchpanel.FLAG_CONFIG_FILE)

	conf, err := ParseConfig(valuesFile, Config{})
	if err != nil {
		logErr.Println("error parsing configuration: %s", err.Error())
		os.Exit(1)
	}

	inputFile, err := LoadFile(conf.FilePath, int64(conf.MaxFileSizeMegabytes*1024))
	if err != nil {
		logErr.Println(err)
		os.Exit(1)
	}

	fHeader, err := header.GetHeader(inputFile)
	if err != nil {
		logErr.Println(err)
		os.Exit(1)
	}
	_ = fHeader

	fmt.Println(fHeader)

}
