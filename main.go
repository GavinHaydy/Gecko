package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"embed"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

// =========================
// 内置模板 embed
// =========================

//go:embed templates/*
var templatesFS embed.FS

// =========================
// 文件复制与模板替换
// =========================
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
		return copyFileWithTemplate(path, target, data)
	})
}

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
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}
	return os.WriteFile(dst, buf.Bytes(), 0644)
}

func copyEmbedDir(fsys fs.FS, src, dst string, data map[string]string) error {
	entries, err := fs.ReadDir(fsys, src)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		// Windows 文件系统用 filepath.Join 生成 dstPath
		dstPath := filepath.Join(dst, entry.Name())

		// embed 用 Unix 风格路径
		srcPath := path.Join(src, entry.Name())

		if entry.IsDir() {
			os.MkdirAll(dstPath, os.ModePerm)
			if err := copyEmbedDir(fsys, srcPath, dstPath, data); err != nil {
				return err
			}
		} else {
			content, err := fs.ReadFile(fsys, srcPath)
			if err != nil {
				return err
			}
			tmpl, err := template.New(entry.Name()).Parse(string(content))
			if err != nil {
				return err
			}
			var buf bytes.Buffer
			if err := tmpl.Execute(&buf, data); err != nil {
				return err
			}
			if err := os.WriteFile(dstPath, buf.Bytes(), 0644); err != nil {
				return err
			}
		}
	}
	return nil
}

// =========================
// 远程模板下载
// =========================
func gitClone(repo, dst string) error {
	cmd := exec.Command("git", "clone", "--depth", "1", repo, dst)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return os.RemoveAll(filepath.Join(dst, ".git"))
}

func downloadZip(repo, dst string) error {
	url := fmt.Sprintf("https://github.com/%s/archive/refs/heads/main.zip", repo)
	fmt.Println("🌐 下载 ZIP 模板:", url)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下载失败: %s", resp.Status)
	}

	zipPath := filepath.Join(os.TempDir(), "template.zip")
	file, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := io.Copy(file, resp.Body); err != nil {
		return err
	}

	return unzip(zipPath, dst)
}

func unzip(src, dst string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()
	for _, f := range r.File {
		parts := strings.SplitN(f.Name, "/", 2)
		if len(parts) < 2 {
			continue
		}
		fpath := filepath.Join(dst, parts[1])
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}
		out, err := os.Create(fpath)
		if err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		_, err = io.Copy(out, rc)
		out.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// =========================
// 判断本地模板是否存在
// =========================
func localTemplateExists(name string) bool {
	path := filepath.Join("templates", name)
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// =========================
// 主程序
// =========================
func main() {
	fmt.Println("🐍 Gecko Python Test Project Generator")
	fmt.Println("--------------------------------------")

	// 输入项目名
	var projectName string
	_ = survey.AskOne(&survey.Input{
		Message: "请输入项目名:",
		Default: "my_test_project",
	}, &projectName)

	// 选择模板
	templates := []string{"pytest-request", "unittest-basic", "远程模板（GitHub URL / user/repo）"}
	var templateChoice string
	_ = survey.AskOne(&survey.Select{
		Message: "选择模板类型:",
		Options: templates,
	}, &templateChoice)

	projectDir := filepath.Join(".", projectName)
	if _, err := os.Stat(projectDir); !os.IsNotExist(err) {
		fmt.Println("❌ 目录已存在:", projectDir)
		os.Exit(1)
	}

	data := map[string]string{"ProjectName": projectName}

	// 远程模板
	if strings.Contains(templateChoice, "远程模板") {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("请输入模板地址 (Git URL 或 GitHub user/repo): ")
		repo, _ := reader.ReadString('\n')
		repo = strings.TrimSpace(repo)

		fmt.Println("🔄 获取模板中...")

		if strings.HasPrefix(repo, "http") {
			if err := gitClone(repo, projectDir); err != nil {
				fmt.Println("❌ 克隆失败:", err)
				os.Exit(1)
			}
		} else {
			if err := downloadZip(repo, projectDir); err != nil {
				fmt.Println("❌ 下载模板失败:", err)
				os.Exit(1)
			}
		}

		fmt.Println("✅ 模板已下载到:", projectDir)
		return
	}

	// 本地或内置模板
	if localTemplateExists(templateChoice) {
		fmt.Println("📁 使用本地模板:", templateChoice)
		if err := copyDir(filepath.Join("templates", templateChoice), projectDir, data); err != nil {
			fmt.Println("❌ 复制模板失败:", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("📦 使用内置模板:", templateChoice)
		if err := copyEmbedDir(templatesFS, "templates/"+templateChoice, projectDir, data); err != nil {
			fmt.Println("❌ 内置模板复制失败:", err)
			os.Exit(1)
		}
	}

	fmt.Println("✅ 项目已生成在:", projectDir)
	fmt.Println()
	fmt.Println("👉 下一步:")
	fmt.Println("   cd", projectName)
	fmt.Println("   pip install -r requirements.txt")
	fmt.Println("   pytest -v")
}
