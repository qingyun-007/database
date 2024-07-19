package helper

import (
	"crypto/md5"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"time"

	"gopkg.in/gomail.v2"

	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"

	"getcharzp.cn/define"
)

type UserClaims struct {
	Identity string `json:"identity"`
	Name     string `json:"name"`
	IsAdmin  int    `json:"is_admin"`
	jwt.StandardClaims
}

// GetMd5
// 生成 md5
func GetMd5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

var myKey = []byte("gin-gorm-oj-key")

// GenerateToken
// 生成 token
func GenerateToken(identity, name string, isAdmin int) (string, error) {
	UserClaim := &UserClaims{
		Identity:       identity,
		Name:           name,
		IsAdmin:        isAdmin,
		StandardClaims: jwt.StandardClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim)
	tokenString, err := token.SignedString(myKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// AnalyseToken
// 解析 token
func AnalyseToken(tokenString string) (*UserClaims, error) {
	userClaim := new(UserClaims)
	claims, err := jwt.ParseWithClaims(tokenString, userClaim, func(token *jwt.Token) (interface{}, error) {
		return myKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return nil, fmt.Errorf("analyse Token Error:%v", err)
	}
	return userClaim, nil
}

// SendCode
// 发送验证码
func SendCode(toUserEmail, code string) error {
	//e := email.NewEmail()
	//e.From = "Get <1070505345@qq.com>"
	//e.To = []string{toUserEmail}
	//e.Subject = "验证码已发送，请查收"
	//e.HTML = []byte("您的验证码：<b>" + code + "</b>")
	//return e.SendWithTLS("pop.qq.com:995",
	//	smtp.PlainAuth("", "1070505345@qq.com", "kutbkqideqskbdig", "smtp.qq.com"),
	//	//&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.163.com"})
	//	&tls.Config{InsecureSkipVerify: true})

	//host := "smtpdm.aliyun.com"
	host := "smtp.qq.com"
	port := 465
	userName := "1070505345@qq.com"
	password := "kutbkqideqskbdig"

	d := gomail.NewDialer(host, port, userName, password)
	m := gomail.NewMessage()
	m.SetHeader("From", "Get <1070505345@qq.com>") // 发件人
	// m.SetHeader("From", "alias"+"<"+userName+">") // 增加发件人别名

	m.SetHeader("To", toUserEmail) // 收件人，可以多个收件人，但必须使用相同的 SMTP 连接
	// m.SetHeader("To", "fanwd@fuya.live", "chenyan@fuya.live", "caizhaoliang@fuya.live", "hufeijie@fuya.live") // 收件人，可以多个收件人，但必须使用相同的 SMTP 连接
	// m.SetHeader("Cc", "fanwd@fuya.live", "chenyan@fuya.live", "caizhaoliang@fuya.live", "hufeijie@fuya.live")  // 抄送，可以多个
	// m.SetHeader("Bcc", "fanwd@fuya.live", "chenyan@fuya.live", "caizhaoliang@fuya.live", "hufeijie@fuya.live") // 暗送，可以多个
	m.SetHeader("Subject", "哈罗") // 邮件主题
	// text/html 的意思是将文件的 content-type 设置为 text/html 的形式，浏览器在获取到这种文件时会自动调用html的解析器对文件进行相应的处理。
	// 可以通过 text/html 处理文本格式进行特殊处理，如换行、缩进、加粗等等
	codeByte := []byte("您的验证码：<b>" + code + "</b>")
	codeStr := string(codeByte)
	m.SetBody("text/html", codeStr)

	// text/plain的意思是将文件设置为纯文本的形式，浏览器在获取到这种文件时并不会对其进行处理
	// m.SetBody("text/plain", "纯文本")
	// m.Attach("test.sh")   // 附件文件，可以是文件，照片，视频等等
	// m.Attach("lolcatVideo.mp4") // 视频
	// m.Attach("lolcat.jpg") // 照片
	// 关闭SSL协议认证
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
	return nil
}

// GetUUID
// 生成唯一码
func GetUUID() string {
	return uuid.NewV4().String()
}

// GetRand
// 生成验证码
func GetRand() string {
	rand.Seed(time.Now().UnixNano())
	s := ""
	for i := 0; i < 6; i++ {
		s += strconv.Itoa(rand.Intn(10))
	}
	return s
}

// CodeSave
// 保存代码
func CodeSave(code []byte) (string, error) {
	dirName := "code/" + GetUUID()
	path := dirName + "/main.go"
	err := os.Mkdir(dirName, 0777)
	if err != nil {
		return "", err
	}
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	f.Write(code)
	defer f.Close()
	return path, nil
}

// CheckGoCodeValid
// 检查golang代码的合法性
func CheckGoCodeValid(path string) (bool, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return false, err
	}
	code := string(b)
	for i := 0; i < len(code)-6; i++ {
		if code[i:i+6] == "import" {
			var flag byte
			for i = i + 7; i < len(code); i++ {
				if code[i] == ' ' {
					continue
				}
				flag = code[i]
				break
			}
			if flag == '(' {
				for i = i + 1; i < len(code); i++ {
					if code[i] == ')' {
						break
					}
					if code[i] == '"' {
						t := ""
						for i = i + 1; i < len(code); i++ {
							if code[i] == '"' {
								break
							}
							t += string(code[i])
						}
						if _, ok := define.ValidGolangPackageMap[t]; !ok {
							return false, nil
						}
					}
				}
			} else if flag == '"' {
				t := ""
				for i = i + 1; i < len(code); i++ {
					if code[i] == '"' {
						break
					}
					t += string(code[i])
				}
				if _, ok := define.ValidGolangPackageMap[t]; !ok {
					return false, nil
				}
			}
		}
	}
	return true, nil
}
