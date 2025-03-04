package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

func replaceConstantValue(content, constantName, newValue string) string {
	// 构建正则表达式，匹配常量定义
	regex := regexp.MustCompile(fmt.Sprintf(`(%s\s*=\s*`+"`)[^`]*(`)", constantName))
	// 替换为新值
	return regex.ReplaceAllString(content, fmt.Sprintf("${1}%s${2}", newValue))
}

func replaceConstantValueWithMarks(content, constantName, newValue string) string {
	// 构建正则表达式，匹配常量定义
	regex := regexp.MustCompile(fmt.Sprintf(`(%s\s*=\s*")[^"]*(")`, constantName))
	// 替换为新值
	return regex.ReplaceAllString(content, fmt.Sprintf("${1}%s${2}", newValue))
}

func WriteSettingsVal(path, key, val string) error {
	if content, err := ioutil.ReadFile(path); err != nil {
		return err
	} else {
		updatedContent := replaceConstantValueWithMarks(string(content), key, val)
		if err := ioutil.WriteFile(path, []byte(updatedContent), os.ModePerm); err != nil {
			return err
		}
		return nil
	}
}

func WriteSettingsRsa(path, prv, pub string) error {
	if content, err := ioutil.ReadFile(path); err != nil {
		return err
	} else {
		// 使用正则表达式替换内容
		updatedContent := replaceConstantValue(string(content), "RSA_PUBLIC_KEY", pub)
		updatedContent = replaceConstantValue(updatedContent, "RSA_PRIVATE_KEY", prv)
		// 将文件内容回写到文件中
		if err := ioutil.WriteFile(path, []byte(updatedContent), os.ModePerm); err != nil {
			return err
		}
		return nil
	}
}
