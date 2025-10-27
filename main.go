package main

import (
	"Gecko/internal/pkg/handler"
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"io/fs"
	"os"
)

func main() {
	// 输入项目名
	var projectName string
	_ = survey.AskOne(&survey.Input{
		Message: "请输入项目名:",
		Default: "my_test_project",
	}, &projectName)

	// create dir
	err := os.Mkdir(projectName, 0777)
	if err != nil && !os.IsExist(err) {
		return
	}

	// 选择模板
	templatesList := []string{"api-pytest", "pytest-request", "unittest-basic", "远程模板（GitHub URL / user/repo）"}
	var templateChoice string
	_ = survey.AskOne(&survey.Select{
		Message: "选择模板类型:",
		Options: templatesList,
	}, &templateChoice)

	err = os.Mkdir(projectName, 0777)
	if err != nil && !errors.Is(err, fs.ErrExist) {
		return
	}

	switch templateChoice {
	case "base":
		handler.WScript(projectName, "api-pytest")
	default:
		fmt.Println(fmt.Sprintf("暂时没有可用模板{%s}", templateChoice))
	}

}
