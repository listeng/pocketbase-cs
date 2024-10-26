package apis

import (
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/core"
)

// bindHookEditApi registers the hook editing api endpoints.
func bindHookEditApi(app core.App, rg *echo.Group) {
	api := hookEditApi{app: app}

	subGroup := rg.Group("/hooks", ActivityLogger(app), RequireAdminAuth())
	subGroup.GET("", api.listHooks)
	subGroup.GET("/:filename", api.readHook)
	subGroup.POST("/:filename", api.writeHook)
	subGroup.PUT("", api.renameHook)
	subGroup.POST("", api.createHook)
	subGroup.DELETE("/:filename", api.deleteHook)
	subGroup.PUT("/restart", api.restartBackend)
}

type hookEditApi struct {
	app core.App
}

func (api *hookEditApi) isInvalidFilename(filename string) bool {
	ext := filepath.Ext(filename)

	// 只允许 .pb 扩展名或无扩展名
	if ext != "" && ext != ".pb" {
		return true
	}

	// 获取不带扩展名的基本文件名
	baseName := strings.TrimSuffix(filename, ext)

	// 验证文件名
	if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(baseName) {
		return true
	}

	// 清理并验证路径
	cleanFilename := filepath.Clean(filename)
	if strings.Contains(cleanFilename, "..") {
		return true
	}

	fullPath := filepath.Join(api.getHooksDir(), cleanFilename+".js")
	if !strings.HasPrefix(fullPath, api.getHooksDir()) {
		return true
	}

	return false
}

// getHooksDir returns the path to the hooks directory
func (api *hookEditApi) getHooksDir() string {
	return filepath.Join(api.app.DataDir(), "../pb_hooks")
}

func (api *hookEditApi) listHooks(c echo.Context) error {
	files, err := os.ReadDir(api.getHooksDir())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Hooks 目录不存在")
	}

	var fileNames []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".js") {
			fileNames = append(fileNames, strings.TrimSuffix(file.Name(), ".js"))
		}
	}

	return c.JSON(http.StatusOK, fileNames)
}

func (api *hookEditApi) readHook(c echo.Context) error {
	filename := c.PathParam("filename")
	content, err := os.ReadFile(filepath.Join(api.getHooksDir(), filename+".js"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "文件未找到")
	}

	return c.JSON(http.StatusOK, string(content))
}

func decryptString(encryptedString string) (string, error) {
	// 解密函数
	decrypt := func(r rune) rune {
		switch {
		case r >= 'A' && r <= 'Z':
			return 'A' + (r-'A'+23)%26 // 26 - 3 = 23
		case r >= 'a' && r <= 'z':
			return 'a' + (r-'a'+23)%26 // 26 - 3 = 23
		case r >= '0' && r <= '9':
			return '0' + (r-'0'+5)%10 // 对于数字，+5 和 -5 在模10的情况下是等价的
		default:
			return r
		}
	}

	// 对加密的字符串进行解密
	decrypted := strings.Map(decrypt, encryptedString)

	// 解码Base64
	decodedBytes, err := base64.StdEncoding.DecodeString(decrypted)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %v", err)
	}

	return string(decodedBytes), nil
}

func (api *hookEditApi) writeHook(c echo.Context) error {
	filename := c.PathParam("filename")
	var content struct {
		Content string `json:"content"`
	}
	if err := c.Bind(&content); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "参数错误")
	}

	decryptedContent, err := decryptString(content.Content)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "解码失败")
	}

	// URL解码
	decodedString, err := url.QueryUnescape(decryptedContent)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "URL解码失败")
	}

	api.app.Logger().Debug(
		"write to file",
		slog.String("name", "hook write"),
		slog.String("data", decodedString),
		slog.String("base64", content.Content),
	)

	err = os.WriteFile(filepath.Join(api.getHooksDir(), filename+".js"), []byte(decodedString), 0644)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "保存文件失败")
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "保存文件成功"})
}

func (api *hookEditApi) renameHook(c echo.Context) error {
	var content struct {
		OldName string `json:"oldname"`
		NewName string `json:"newname"`
	}
	if err := c.Bind(&content); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "参数错误")
	}

	if api.isInvalidFilename(content.OldName) {
		return echo.NewHTTPError(http.StatusInternalServerError, "无效的原文件名")
	}

	if api.isInvalidFilename(content.NewName) {
		return echo.NewHTTPError(http.StatusInternalServerError, "无效的新文件名")
	}

	err := os.Rename(filepath.Join(api.getHooksDir(), content.OldName+".js"), filepath.Join(api.getHooksDir(), content.NewName+".js"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "重命名失败")
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "重命名成功"})
}

func (api *hookEditApi) createHook(c echo.Context) error {
	var file struct {
		Name string `json:"filename"`
	}
	if err := c.Bind(&file); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "参数错误")
	}

	if api.isInvalidFilename(file.Name) {
		return echo.NewHTTPError(http.StatusInternalServerError, "无效的文件名")
	}

	err := os.WriteFile(filepath.Join(api.getHooksDir(), file.Name+".js"), []byte(""), 0644)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "创建文件失败")
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "创建文件成功"})
}

func (api *hookEditApi) deleteHook(c echo.Context) error {
	filename := c.PathParam("filename")

	if api.isInvalidFilename(filename) {
		return echo.NewHTTPError(http.StatusInternalServerError, "无效的文件名")
	}

	err := os.Remove(filepath.Join(api.getHooksDir(), filename+".js"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "删除文件失败")
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "删除文件成功"})
}

func (api *hookEditApi) restartBackend(c echo.Context) error {
	if runtime.GOOS != "windows" {
		return echo.NewHTTPError(http.StatusInternalServerError, "只有Windows平台需要重启，其他平台不需要")
	}

	go func() {
		// 等待一小段时间，确保 HTTP 响应已经发送
		time.Sleep(100 * time.Millisecond)

		// 获取当前可执行文件的路径
		executable, err := os.Executable()
		if err != nil {
			api.app.Logger().Error(
				"Error getting executable path:" + err.Error(),
			)
			return
		}

		exePath := filepath.Dir(executable)
		batPath := filepath.Join(exePath, "restart.bat")

		_, fileexterr := os.Stat(batPath)
		if !os.IsNotExist(fileexterr) {

			// 创建 restart.bat 文件，包含 2 秒延迟
			batContent := fmt.Sprintf(`@echo off
cd /d "%s"
echo Waiting for 2 seconds before restarting...
timeout /t 2 /nobreak > NUL
start "" "%s" serve --http 0.0.0.0:8090
exit
`, exePath, filepath.Base(executable))

			err = os.WriteFile(batPath, []byte(batContent), 0644)
			if err != nil {
				api.app.Logger().Error(
					"Error creating restart.bat:" + err.Error(),
				)
				return
			}
		}

		// 启动 restart.bat
		cmd := exec.Command("cmd", "/C", "start", batPath)
		cmd.Dir = exePath
		err = cmd.Start()
		if err != nil {
			api.app.Logger().Error(
				"Error starting restart.bat:" + err.Error(),
			)
			return
		}

		// 退出当前进程
		os.Exit(0)
	}()

	return c.JSON(http.StatusOK, map[string]string{"message": "重启中……"})
}
