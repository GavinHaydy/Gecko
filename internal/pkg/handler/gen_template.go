package handler

import (
	"Gecko/internal/pkg/dal/rao"
	"Gecko/templates"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func WScript(path string, templateName string) {
	var templateList []rao.Template
	var fullPath string
	switch templateName {
	case "api-pytest":
		templateList = templates.ApiPytest()
	}

	for _, t := range templateList {
		fullPath = fmt.Sprintf("%s/%s", path, t.FileName)

		// 获取文件所在目录
		dir := filepath.Dir(fullPath)

		// 递归创建目录（如果已存在不会报错）
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatalf("创建目录失败: %v", err)
		}

		create, fullErr := os.Create(fullPath)
		if fullErr != nil {
			return
		}

		err = create.Close()
		if err != nil {
			return
		}

		err = os.WriteFile(fullPath, []byte(t.Content), 0644)
		if err != nil {
			return
		}
	}

}
