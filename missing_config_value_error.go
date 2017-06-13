package main

import "fmt"

type MissingConfigValueError struct {
	ConfigPath string
	ValueName  string
	OptionName string
	KeyName    string
}

func (err MissingConfigValueError) Error() string {
	return NewError(
		fmt.Errorf("%s is not specified.", err.ValueName),
		`Either specify --%s command line option or set "%s" `+
			"option in the configuration file:\n\n%s:\n\t%s",
		err.OptionName,
		err.KeyName,
		err.ConfigPath,
		fmt.Sprintf(`%s: "PUT_VALUE_HERE"`, err.KeyName),
	).Error()
}
