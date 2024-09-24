package apis

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"regexp"

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
	files, err := ioutil.ReadDir(api.getHooksDir())
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
	content, err := ioutil.ReadFile(filepath.Join(api.getHooksDir(), filename + ".js"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "文件未找到")
	}

	return c.JSON(http.StatusOK, string(content))
}

func (api *hookEditApi) writeHook(c echo.Context) error {
	filename := c.PathParam("filename")
	var content struct {
		Content string `json:"content"`
	}
	if err := c.Bind(&content); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "参数错误")
	}

	err := ioutil.WriteFile(filepath.Join(api.getHooksDir(), filename + ".js"), []byte(content.Content), 0644)
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

	err := os.Rename(filepath.Join(api.getHooksDir(), content.OldName + ".js"), filepath.Join(api.getHooksDir(), content.NewName + ".js"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "重命名失败")
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "重命名成功"})
}

func (api *hookEditApi) createHook(c echo.Context) error {
	var file struct {
		Name    string `json:"filename"`
	}
	if err := c.Bind(&file); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "参数错误")
	}

	if api.isInvalidFilename(file.Name) {
		return echo.NewHTTPError(http.StatusInternalServerError, "无效的文件名")
	}

	err := ioutil.WriteFile(filepath.Join(api.getHooksDir(), file.Name + ".js"), []byte(""), 0644)
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

	err := os.Remove(filepath.Join(api.getHooksDir(), filename + ".js"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "删除文件失败")
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "删除文件成功"})
}
