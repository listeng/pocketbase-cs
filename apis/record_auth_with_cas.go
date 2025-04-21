package apis

import (
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/security"
)

func validateCASTicket(c *core.Collection, e *core.RequestEvent) (string, string, string, error) {
	if !c.CASAuth.Enabled {
		return "", "", "", e.ForbiddenError("CAS 登录未启用.", nil)
	}

	casTicket := e.Request.URL.Query().Get("ticket")
	serviceURL := e.Request.URL.Query().Get("service")

	if casTicket == "" {
		return "", "", "", e.UnauthorizedError("没有票据", nil)
	}

	validateURL := fmt.Sprintf("%s?ticket=%s&service=%s", c.CASAuth.ValidateUrl, url.QueryEscape(casTicket), url.QueryEscape(serviceURL))
	resp, err := http.Get(validateURL)
	if err != nil {
		return "", "", "", e.InternalServerError("CAS登录失败", nil)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", e.InternalServerError("读取CAS服务器响应失败", nil)
	}

	if !strings.Contains(string(body), "cas:authenticationSuccess") {
		return "", "", "", e.UnauthorizedError("CAS拒绝登录", nil)
	}

	re := regexp.MustCompile(`<cas:user>(.*?)</cas:user>`)
	userMatch := re.FindStringSubmatch(string(body))
	if len(userMatch) < 2 {
		return "", "", "", e.UnauthorizedError("CAS没有返回用户名（username）", nil)
	}
	username := userMatch[1]

	re = regexp.MustCompile(`<cas:utype>(.*?)</cas:utype>`)
	userTypeMatch := re.FindStringSubmatch(string(body))
	if len(userTypeMatch) < 2 {
		return "", "", "", e.UnauthorizedError("CAS没有返回用户类型（utype）", nil)
	}
	userType := userTypeMatch[1]

	re = regexp.MustCompile(`<cas:roles>(.*?)</cas:roles>`)
	userRolesMatch := re.FindStringSubmatch(string(body))
	// Extract roles but skip if no match found
	roles := ""
	if len(userRolesMatch) >= 2 {
		roles = userRolesMatch[1]

		// URL decode the roles string
		decodedRoles, err := url.QueryUnescape(roles)
		if err == nil {
			roles = decodedRoles
		}
	}

	e.App.Logger().Info("cas login", slog.String("user", username), slog.String("utype", userType), slog.String("roles", roles))

	return username, userType, roles, nil
}

func recordAuthWithCAS(e *core.RequestEvent) error {
	collection, _ := findAuthCollection(e)
	if collection == nil {
		return NewNotFoundError("缺少数据集上下文", nil)
	}

	username, userType, userRoles, err := validateCASTicket(collection, e)
	if err != nil {
		return err
	}

	if collection.Name == core.CollectionNameSuperusers &&
		userType != "admin" && userType != "super" {
		return e.UnauthorizedError(username+" 不是管理员用户", nil)
	}

	user, err := e.App.FindAuthRecordByEmail(collection, username+"@"+collection.CASAuth.Realm)
	if err == sql.ErrNoRows {
		if collection.CASAuth.CreateNewUser {
			record := core.NewRecord(collection)
			record.Set("email", username+"@"+collection.CASAuth.Realm)
			record.SetTokenKey(security.RandomString(50))

			if collection.Fields.GetByName("username") != nil {
				record.Set("username", username)
			}

			if collection.Fields.GetByName("name") != nil {
				record.Set("name", username)
			}

			if collection.Fields.GetByName("verified") != nil {
				record.Set("verified", true)
			}

			if userType == "admin" || userType == "super" {
				if collection.Fields.GetByName("is_admin") != nil {
					user.Set("is_admin", true)
				}
			}

			if collection.CASAuth.AdminRole != "" && userRoles != "" {
				roles := strings.Split(userRoles, ",")
				for _, role := range roles {
					if role == collection.CASAuth.AdminRole {
						if collection.Fields.GetByName("is_admin") != nil {
							user.Set("is_admin", true)
						}
						break
					}
				}
			}

			if err := e.App.Save(record); err != nil {
				return e.UnauthorizedError("CAS登录失败，用户无效[新用户创建失败]", nil)
			}
		} else {
			return e.UnauthorizedError("CAS登录失败，用户无效", nil)
		}
	} else if err != nil {
		return e.UnauthorizedError("CAS登录失败，用户无效[未知错误]", nil)
	}

	e.App.Logger().Info("cas login as "+userType, slog.String("user", username))

	if user != nil {
		return RecordAuthResponse(e, user, core.MFAMethodOTP, nil)
	}

	return e.UnauthorizedError("CAS登录失败，无效的登录信息", nil)
}
