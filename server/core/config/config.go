package config

import (
	"fmt"
	"os"
	"path"
	"regexp"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

// Unmarshal 함수는 설정 파일을 로드하고 구조체로 변환합니다.
// 동작 과정:
// 1. .env 파일을 로드하여 환경 변수를 설정합니다.
// 2. SERVICE_ENV 환경 변수를 기반으로 적절한 config.*.yaml 파일을 선택합니다.
// 3. YAML 파일 내용을 읽고 환경 변수(${VAR_NAME} 형식)를 실제 값으로 치환합니다.
// 4. YAML 데이터를 제공된 구조체 타입(T)으로 언마샬링합니다.
//
// 예시:
//
//	type Config struct {
//		Database struct {
//			Host string `yaml:"host"`
//			Port int    `yaml:"port"`
//		} `yaml:"database"`
//	}
//
//	cfg := config.Unmarshal(&Config{})
func Unmarshal[T any](out *T) error {
	_ = godotenv.Load()
	var fileName = configYaml()

	// YAML 파일 읽기
	yamlData, err := readFile(fileName)
	if err != nil {
		return fmt.Errorf("error reading YAML file: %v", err)
	}

	// 환경 변수 치환
	yamlData = replaceEnvVariables(yamlData)

	err = yaml.Unmarshal(yamlData, out)
	if err != nil {
		return fmt.Errorf("error unmarshalling YAML: %v", err)
	}

	return nil
}

// replaceEnvVariables 함수는 YAML 콘텐츠 내의 환경 변수 참조를 실제 값으로 대체합니다.
// 지원하는 형식:
// 1. ${VAR_NAME} - 환경 변수 VAR_NAME의 값으로 대체됩니다. 변수가 없으면 빈 문자열로 대체됩니다.
// 2. ${VAR_NAME:-default} - 환경 변수 VAR_NAME의 값으로 대체됩니다. 변수가 없으면 'default' 값으로 대체됩니다.
//
// 예시:
//   - ${DB_HOST} → 환경 변수 DB_HOST 값 (없으면 빈 문자열)
//   - ${DB_PORT:-5432} → 환경 변수 DB_PORT 값 (없으면 5432)
func replaceEnvVariables(yamlContent []byte) []byte {
	re := regexp.MustCompile(`\$\{(\w+)(?::-([^}]*))?}`)

	return re.ReplaceAllFunc(yamlContent, func(match []byte) []byte {
		matchStr := string(match)

		submatches := re.FindStringSubmatch(matchStr)
		if len(submatches) < 2 {
			return match
		}

		envVar := submatches[1]
		if val, exists := os.LookupEnv(envVar); exists {
			return []byte(val)
		}
		if len(submatches) >= 3 && submatches[2] != "" {
			return []byte(submatches[2])
		}

		return []byte{}
	})
}

func configYaml() string {
	var cfgFile string
	if env, exists := os.LookupEnv("SERVICE_ENV"); exists {
		cfgFile = "config." + env + ".yaml"
	}

	if !fileExists(cfgFile) {
		return "config.default.yaml"
	}

	return cfgFile
}

// #nosec G304 - filename is restricted to config.*.yaml pattern
func readFile(filename string) ([]byte, error) {
	if !regexp.MustCompile(`^config\.([a-zA-Z0-9]+)?\.yaml$`).MatchString(filename) {
		return nil, fmt.Errorf("invalid config file: %s", filename)
	}

	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("unable to get working directory: %v", err)
	}

	file := path.Join(dir, filename)

	if !fileExists(file) {
		return nil, fmt.Errorf("file not found: %s", filename)
	}

	return os.ReadFile(file)
}

func fileExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}
