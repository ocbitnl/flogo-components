/*
* Copyright © 2017. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	ftrigger "github.com/TIBCOSoftware/flogo-lib/core/trigger"
)

var githubRawContent = "https://raw.githubusercontent.com"

//GetGithubResource used to get github files present in given path
func GetGithubResource(gitHubPath string, resourceFile string) ([]byte, error) {
	ref := ReplaceNth(gitHubPath+"/"+resourceFile, "/", "/master/", 3)
	remoteFile := strings.Replace(ref, "github.com", githubRawContent, 1)
	response, rerr := http.Get(remoteFile)
	if rerr != nil {
		return nil, rerr
	}
	responseData, resperr := ioutil.ReadAll(response.Body)
	if resperr != nil {
		return nil, resperr
	}

	return responseData, nil
}

//GetTriggerMetadata returns trigger.json for supplied trigger github path
func GetTriggerMetadata(gitHubPath string) (*ftrigger.Metadata, error) {
	goPathVendor := filepath.Join(os.Getenv("GOPATH"), "vendor")
	triggerMetadata := &ftrigger.Metadata{}
	if _, err := os.Stat(filepath.Join(goPathVendor, gitHubPath, Gateway_Trigger_Metadata_JSON_Name)); os.IsNotExist(err) {
		if _, err := os.Stat(filepath.Join(goPathVendor, gitHubPath)); os.IsNotExist(err) {
			os.MkdirAll(filepath.Join(goPathVendor, gitHubPath), os.ModePerm)
		}
		data, err := GetGithubResource(gitHubPath, Gateway_Trigger_Metadata_JSON_Name)
		if err != nil {
			return nil, err
		}
		json.Unmarshal(data, triggerMetadata)

		os.Create(filepath.Join(goPathVendor, gitHubPath, Gateway_Trigger_Metadata_JSON_Name))
		err = ioutil.WriteFile(filepath.Join(goPathVendor, gitHubPath, Gateway_Trigger_Metadata_JSON_Name), data, os.ModePerm)
		if err != nil {
			return nil, err
		}
	} else {
		data, err := ioutil.ReadFile(filepath.Join(goPathVendor, gitHubPath, Gateway_Trigger_Metadata_JSON_Name))
		if err != nil {
			return nil, err
		}
		json.Unmarshal(data, triggerMetadata)
	}

	return triggerMetadata, nil
}

func IsValidTriggerSetting(metadata *ftrigger.Metadata, property string) bool {
	settings := metadata.Settings
	for key := range settings {
		if key == property {
			return true
		}
	}

	return false
}

func IsValidTriggerHandlerSetting(metadata *ftrigger.Metadata, property string) bool {
	settings := metadata.Handler.Settings

	for _, element := range settings {
		if element.Name() == property {
			return true
		}
	}

	return false
}

func ValidateTriggerConfigExpr(expression *string) (bool, *string) {
	if expression == nil {
		return false, nil
	}

	exprValue := *expression
	if strings.HasPrefix(exprValue, Gateway_Trigger_Config_Prefix) && strings.HasSuffix(exprValue, Gateway_Trigger_Config_Suffix) {
		//get name of the config
		str := exprValue[len(Gateway_Trigger_Config_Prefix) : len(exprValue)-1]
		return true, &str
	} else {
		return false, &exprValue
	}
}

func CheckTriggerOptimization(triggerSettings map[string]interface{}) bool {
	if val, ok := triggerSettings[Gateway_Trigger_Optimize_Property]; ok {
		optimize, err := strconv.ParseBool(val.(string))
		if err != nil {
			//check if its a boolean
			optimize, found := val.(bool)
			if !found {
				return found
			}
			return optimize
		}
		return optimize
	} else {
		return Gateway_Trigger_Optimize_Property_Default
	}
}

func validateEnvPropertySettingExpr(expression *string) (bool, *string) {
	if expression == nil {
		return false, nil
	}

	exprValue := *expression
	if strings.HasPrefix(exprValue, Gateway_Trigger_Setting_Env_Prefix) && strings.HasSuffix(exprValue, Gateway_Trigger_Setting_Env_Suffix) {
		//get name of the property
		str := exprValue[len(Gateway_Trigger_Setting_Env_Prefix) : len(exprValue)-1]
		return true, &str
	}
	return false, &exprValue
}

// ResolveEnvironmentProperties resolves environment properties mentioned in the settings map.
func ResolveEnvironmentProperties(settings map[string]interface{}) error {
	for k, v := range settings {
		value := v.(string)
		valid, propertyName := validateEnvPropertySettingExpr(&value)
		if !valid {
			continue
		}
		//lets get the env property value
		propertyNameStr := *propertyName
		propertyValue, found := os.LookupEnv(propertyNameStr)
		if !found {
			return fmt.Errorf("environment property [%v] is not set", propertyNameStr)
		}
		settings[k] = propertyValue
	}
	return nil
}

// ReplaceNth Replaces the nth occurrence of old in s by new.
func ReplaceNth(s, old, new string, n int) string {
	i := 0
	for m := 1; m <= n; m++ {
		x := strings.Index(s[i:], old)
		if x < 0 {
			break
		}
		i += x
		if m == n {
			return s[:i] + new + s[i+len(old):]
		}
		i += len(old)
	}
	return s
}