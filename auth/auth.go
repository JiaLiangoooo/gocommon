package auth

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.icsoc.net/cc/gocommon/config"
	"gitlab.icsoc.net/cc/gocommon/zredis"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	REDIS_KEY                  = "auth"
	VERIFY_URL_PATH            = "/oauth2/verify"
	USERINFO_URL_PATH          = "/oauth2/userinfo"
	AUTHORIZATION_HEADER       = "Authorization"
	AUTHORIZATION_VALUE_PREFIX = "Bearer "
)

//用户信息结构体
type UserInfo struct {
	VccId     string `json:"vcc_id"`
	VccCode   string `json:"vcc_code"`
	AgId      string `json:"ag_id"`
	GroupId   string `json:"group_id"`
	GroupName string `json:"group_name"`
	DeptId    string `json:"dept_id"`   // 部门id
	DeptName  string `json:"dept_name"` // 部门名称
	UserNum   string `json:"user_num"`
	UserName  string `json:"user_name"`
	RoleType  string `json:"role_type"`
	ClientId  string `json:"client_id"`
	LoginType string `json:"login_type"`
	Token     string `json:"token"`
	UserSub   string `json:"sub"`
	Expire    int64  `json:"expires"`
}

//实现BinaryMarshaler接口
func (m *UserInfo) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

func (m *UserInfo) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

//校验token
func CheckAuth(token string) (*UserInfo, error) {
	authConfig := config.GetAuthConfig()
	if authConfig == nil {
		return nil, fmt.Errorf("未配置认证host,请配置{auth.host}参数")
	}
	authHost := authConfig.Host
	return verifyOauth(token, authHost)
}

// 验证token的可用性
func verifyOauth(token string, authHost string) (user *UserInfo, err error) {
	redisCli := zredis.GetRedisClient(REDIS_KEY)
	user = checkRedis(token, redisCli)
	if user != nil {
		return user, nil
	}
	err, expire := verifyToken(token, authHost)
	if err != nil {
		return nil, err
	}
	user, err = getUserInfo(token, authHost)
	if user == nil || err != nil {
		return nil, fmt.Errorf("获取用户失败,authorization: %s, ", token)
	}

	user.Expire = expire
	exp := time.Duration(user.Expire-time.Now().Unix()) * time.Second
	e := redisCli.Set(token, user, exp).Err()
	if e != nil {
		return user, e
	}

	return user, nil
}

// 检查redis中是否存在token
func checkRedis(token string, redisCli zredis.Redis) (user *UserInfo) {
	// redis 去查找authentication 若存在就返回return
	if redisCli.Exists(token).Val() > 0 {
		val := redisCli.Get(token).Val()
		ui := &UserInfo{}
		if er := json.Unmarshal([]byte(val), ui); er != nil {
			fmt.Printf("get redis token:%s unmarshal error: %v", token, er)
			return nil
		}
		return ui
	} else {
		return nil
	}
}

// 检查Token是否有效
func verifyToken(token string, authHost string) (err error, expire int64) {
	client := &http.Client{}
	// manager token 与 ekt token 都能校验接口
	verifyUrl := authHost + VERIFY_URL_PATH

	req, e := http.NewRequest("GET", verifyUrl, nil)
	if e != nil {
		return e, 0
	}
	req.Header = getHeader(token)
	resp, e := client.Do(req)
	if e != nil {
		return e, 0
	}
	if resp.StatusCode != 200 {
		return errors.Errorf("无效的token:%s statuscode:%d", token, resp.StatusCode), 0
	}
	body, e := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if e != nil {
		return errors.Errorf("Read error: %s", body), 0
	}

	var dat map[string]interface{}
	if err := json.Unmarshal(body, &dat); err != nil {
		return fmt.Errorf("get http expire:%s unmarshal error: %v", token, err), 0
	}

	if tempExpire, ok := dat["expires"].(float64); ok {
		expire = int64(tempExpire)
	}
	return nil, expire
}

//获取用户信息
func getUserInfo(token string, authHost string) (info *UserInfo, e error) {
	client := &http.Client{}
	userInfoUrl := authHost + USERINFO_URL_PATH
	request, err := http.NewRequest("GET", userInfoUrl, nil)
	if err != nil {
		return nil, errors.Errorf("获取userInfo失败: %v", err)
	}
	request.Header = getHeader(token)

	resp, err := client.Do(request)
	if err != nil {
		return nil, errors.Errorf("request error: %v", e)
	}

	if resp.StatusCode != 200 {
		return nil, errors.Errorf("获取userInfo error: %v", resp.Status)
	}

	var user = &UserInfo{}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err := json.Unmarshal(body, user); err != nil {
		return nil, errors.Errorf("unmarshal body error %v", err)
	}
	user.Token = token
	return user, nil
}

//组装http请求header
func getHeader(authorization string) map[string][]string {
	header := map[string][]string{}
	aus := []string{AUTHORIZATION_VALUE_PREFIX + authorization}
	header[AUTHORIZATION_HEADER] = aus
	return header
}
