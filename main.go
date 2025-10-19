package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

// copyDir 复制整个目录（支持模板变量替换）
func copyDir(src, dst string, data map[string]string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, _ := filepath.Rel(src, path)
		target := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(target, os.ModePerm)
		}
		// 复制文件并执行模板替换
		return copyFileWithTemplate(path, target, data)
	})
}

// copyFileWithTemplate 支持变量替换
func copyFileWithTemplate(src, dst string, data map[string]string) error {
	content, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	tmpl, err := template.New(filepath.Base(src)).Parse(string(content))
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, buf.Bytes(), 0644)
}

// gitClone 拉取远程模板（可选）
func gitClone(repo, dst string) error {
	cmd := exec.Command("git", "clone", "--depth", "1", repo, dst)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	fmt.Println("🐍 Python Test Project Generator (Go Edition)")
	fmt.Println("--------------------------------------------")

	// 1️⃣ 输入项目名
	var projectName string
	prompt := &survey.Input{
		Message: "请输入项目名:",
		Default: "my_test_project",
	}
	err := survey.AskOne(prompt, &projectName)
	if err != nil {
		return
	}

	// 2️⃣ 选择模板
	templates := []string{"pytest-api", "unittest-basic", "远程模板（GitHub URL）"}
	var templateChoice string
	err = survey.AskOne(&survey.Select{
		Message: "选择模板类型:",
		Options: templates,
	}, &templateChoice)
	if err != nil {
		return
	}

	projectDir := filepath.Join(".", projectName)
	if _, err := os.Stat(projectDir); !os.IsNotExist(err) {
		fmt.Println("❌ 目录已存在:", projectDir)
		os.Exit(1)
	}

	// 3️⃣ 本地模板或远程仓库
	if templateChoice == "远程模板（GitHub URL）" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("请输入 Git 仓库地址 (例如: https://github.com/user/python-template.git): ")
		repo, _ := reader.ReadString('\n')
		repo = strings.TrimSpace(repo)

		fmt.Println("🔄 正在从远程仓库克隆模板...")
		if err := gitClone(repo, projectDir); err != nil {
			fmt.Println("❌ 克隆失败:", err)
			os.Exit(1)
		}
		fmt.Println("✅ 已创建项目:", projectDir)
		return
	}

	// 4️⃣ 本地模板路径
	src := filepath.Join("templates", templateChoice)
	if _, err := os.Stat(src); os.IsNotExist(err) {
		fmt.Println("❌ 模板不存在:", src)
		os.Exit(1)
	}

	// 5️⃣ 执行模板复制
	data := map[string]string{
		"ProjectName": projectName,
	}
	if err := copyDir(src, projectDir, data); err != nil {
		fmt.Println("❌ 复制模板失败:", err)
		os.Exit(1)
	}

	fmt.Println("✅ 项目已生成在:", projectDir)
	fmt.Println()
	fmt.Println("👉 下一步:")
	fmt.Println("   cd", projectName)
	fmt.Println("   pip install -r requirements.txt")
	fmt.Println("   pytest -v")
}
