package exec

import (
	"fmt"
	"getcharzp.cn/bytes"
	"getcharzp.cn/models"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"strings"
)

type defaultModel struct {
	conn *sqlx.DB
}

type CodeExe struct {
	Case   models.TestCase
	Stderr interface{} `json:"stderr"`
	Stdout interface{} `json:"stdout"`
	Path   string      `json:"path"`
	m      defaultModel
}

type ExecSpace struct{}

func Command(language, exec string, codePath string) CodeExe {
	dsn := "root:www.163.com@tcp(127.0.0.1:3306)/sql_test?charset=utf8mb4&parseTime=True"
	db, _ := sqlx.Connect("mysql", dsn)
	return CodeExe{
		Stderr: bytes.Stderr,
		Stdout: bytes.Stdout,
		Path:   codePath,
		m: defaultModel{
			conn: db,
		},
	}
}

// 执行数据库
func (c CodeExe) Run() error {
	content, _ := ioutil.ReadFile(c.Path)
	query := fmt.Sprintf("%s", content)
	query = strings.ToTitle(query)
	if query != c.Case.Output {
		fmt.Println("错误")
		c.Stdout = "xxx"
	}
	//if strings.Contains(query, c.Case.Output) {
	//	fmt.Println("错误")
	//}
	//if strings.Contains(query, c.Case.Yaoqiuk) {
	//	fmt.Println("错误")
	//}
	//if strings.Contains(query, c.Case.Yaoqiuv) {
	//	fmt.Println("错误")
	//}
	//db, _ := sqlx.Open("mysql", "dev:123456@tcp(127.0.0.1:3306)/oms?charset=utf8mb4&parseTime=True")
	affect, err := c.m.conn.Exec(query)
	if err != nil {
		fmt.Println("exec:", err) //直接返回错误
		c.Stdout = "xxx"

	}
	rowaffect, _ := affect.RowsAffected()
	if rowaffect < 1 {
		fmt.Println("插入错误") // 直接返回错误
		c.Stdout = "xxx"

	}
	//querys := strings.Split(query, " ")
	//if querys[0] == "select" {
	//	// out = 输出的结果
	//}
	c.Stdout = query

	// out = ok

	//db, _ := sqlx.Open("mysql", "dev:123456@tcp(127.0.0.1:3306)/oms?charset=utf8mb4&parseTime=True")
	//var user []User
	//err := db.Select(&user, query)
	//if err != nil {
	//	fmt.Println("select err:", err)
	//}
	//fmt.Println("result", user)
	return nil
}
