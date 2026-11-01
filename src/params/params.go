package params

import (
	"errors"
	"regexp"
)

type WorkMode string



// ValidationErrors
var (
	NameValidationError = errors.New("Name validation failed")
	URLValidationError = errors.New("URL validation failed")
	WorkModeValidationError = errors.New("Restricted value to param: work_mode, appropriate values: testing, production")
	FilenameValidationError = errors.New("Restricted name of file")
	PathValidationError = errors.New("Path validation failed")
)


// WorkMode values
const (
	WorkModeTesting WorkMode = "testing"
	WorkModeProduction WorkMode = "production"
)


// Regexes to validate fields
var (
	nameRegex = regexp.MustCompile(`^[a-z0-9][a-z0-9+.-]+$`)
	urlRegex = regexp.MustCompile(`^https?://[^/\s]+/ubuntu/dists/[a-z][a-z0-9.-]*/(main|universe|restricted|multiverse)/binary-[a-z0-9-]+/Packages(\.gz)?$`)
	filenameRegex = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*\.png$`)
	pathRegex = regexp.MustCompile(`^(/[^/]+)+/?$`)
)

type Params struct {
	Name string 
	URL string
	Mode WorkMode
	Filename string
}

func validateName(name string) (string, error) {
	if nameRegex.MatchString(name) {
		return name, nil
	}
	return "", NameValidationError
}

func validateURL(value string, mode WorkMode) (string, error) {
	switch mode {

	case WorkModeProduction:
		if urlRegex.MatchString(value) {
			return value, nil
		}
		return "", URLValidationError

	case WorkModeTesting:
		if pathRegex.MatchString(value) {
			return value, nil
		}
		return "", PathValidationError

	default:
		return "", WorkModeValidationError
	}
}
func validateMode(mode string) (WorkMode, error) {
	if mode == string(WorkModeTesting) {
		return WorkModeTesting, nil
	}
	if mode == string(WorkModeProduction) {
		return WorkModeProduction, nil
	}
	return "", WorkModeValidationError
}

func validateFilename(filename string) (string, error) {
	if filenameRegex.MatchString(filename) {
		return filename, nil
	}
	return "", FilenameValidationError
}


func NewParams(name string, url string, mode string, filename string) (*Params, error) {
	name, err := validateName(name)
	if err != nil {
		return &Params{}, err
	}

	workMode, err := validateMode(mode)
	if err != nil {
		return nil, err
	}
	
	url, err = validateURL(url, workMode)
	if err != nil {
		return nil, err
	}
	
	filename, err = validateFilename(filename) 
	if err != nil {
		return nil, err
	}

	return &Params{
		Name: name,
		URL: url,
		Mode: workMode,
		Filename: filename,
	}, nil
}