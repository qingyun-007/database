package service

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"getcharzp.cn/define"
	"getcharzp.cn/helper"
	"getcharzp.cn/models"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/tiger1103/gfast/v3/library/libUtils"
)

var memInfo = CpuInfoInput{
	MemInfo: make([]CpuInfoSubInput, 0, 200),
}

type CpuInfoInput struct {
	MemInfo []CpuInfoSubInput
}

type CpuInfoSubInput struct {
	MemUsed float64 `json:"memUsed"`
	Time    string  `json:"time"`
}

type GetCpuInfoRes struct {
	CpuNum   interface{} `json:"cpuNum" description:"获取系统的Cpu数量" example:"1"`
	CpuUsed  interface{} `json:"cpuUsed" description:"获取系统的Cpu占用率" example:"1.00"`
	CpuAvg5  interface{} `json:"cpuAvg5" description:"获取系统的Cpu负载" example:"1.00"`
	CpuAvg15 interface{} `json:"cpuAvg15" description:"获取系统的Cpu数量" example:"1.00"`
	MemTotal interface{} `json:"memTotal" description:"获取系统的内存总数" example:"16"`
	GoTotal  interface{} `json:"goTotal"`
	MemUsed  interface{} `json:"memUsed" description:"获取系统的内存使用率" example:"1.00"`

	GoUsed          interface{} `json:"goUsed"`
	MemFree         interface{} `json:"memFree" description:"获取系统的内存的空闲" example:"1"`
	GoFree          interface{} `json:"goFree"`
	MemUsage        interface{} `json:"memUsage"`
	GoUsage         interface{} `json:"goUsage"`
	SysComputerName interface{} `json:"sysComputerName" description:"获取计算机名" example:"SERVER-13ABD"`
	SysOsName       interface{} `json:"sysOsName" description:"获取系统的名称" example:"windows"`
	SysComputerIp   interface{} `json:"sysComputerIp" description:"获取系统的IP" example:"192.168.1.1"`
	SysOsArch       interface{} `json:"sysOsArch"`
	GoName          interface{} `json:"goName"`
	GoVersion       interface{} `json:"goVersion" description:"获取Go语言版本" example:"1.12.0"`
	GoHome          interface{} `json:"goHome" description:"获取Go语言的目录" example:"/mnt/os"`
	GoUserDir       interface{} `json:"goUserDir" description:"获取系统的Cpu数量" example:"1"`
	DiskList        interface{} `json:"diskList" description:"获取磁盘列表" example:"1"`
}

func GetCpuInfo(c *gin.Context) {
	cpuNum := runtime.NumCPU() //核心数
	var cpuUsed float64 = 0    //用户使用率
	var cpuAvg5 float64 = 0    //CPU负载5
	var cpuAvg15 float64 = 0   //当前空闲率

	cpuInfo, err := cpu.Percent(time.Duration(time.Second), true)
	a := float64(0)
	if err == nil {
		for _, v := range cpuInfo {
			a += v
		}
		cpuUsed, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", a/16), 64)
	} else {
		fmt.Println("err:", err)
	}

	loadInfo, err := load.Avg()
	if err == nil {
		cpuAvg5, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", loadInfo.Load5), 64)
		cpuAvg15, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", loadInfo.Load5), 64)
	}

	var memTotal uint64 = 0  //总内存
	var memUsed uint64 = 0   //已用内存
	var memFree uint64 = 0   //剩余内存
	var memUsage float64 = 0 //使用率

	v, err := mem.VirtualMemory()
	if err == nil {
		memTotal = v.Total
		memUsed = v.Used
		memFree = memTotal - memUsed
		memUsage, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", v.UsedPercent), 64)
		if len(memInfo.MemInfo) > 200 {
			memInfo.MemInfo = memInfo.MemInfo[1:]
		}

		memUsedData := CpuInfoSubInput{
			MemUsed: memUsage,
			Time:    time.Now().Format("15:04:05"),
		}
		memInfo.MemInfo = append(memInfo.MemInfo, memUsedData)
	}
	fmt.Println("err:", err)
	var goTotal uint64 = 0  //go分配的总内存数
	var goUsed uint64 = 0   //go使用的内存数
	var goFree uint64 = 0   //go剩余的内存数
	var goUsage float64 = 0 //使用率

	var gomem runtime.MemStats
	runtime.ReadMemStats(&gomem)
	goUsed = gomem.Sys
	goUsage = gconv.Float64(fmt.Sprintf("%.2f", gconv.Float64(goUsed)/gconv.Float64(memTotal)*100))
	sysComputerIp := "" //服务器IP

	ip, err := libUtils.GetLocalIP()
	if err == nil {
		sysComputerIp = ip
	}
	fmt.Println("err:", err)

	sysComputerName := "" //服务器名称
	sysOsName := ""       //操作系统
	sysOsArch := ""       //系统架构

	sysInfo, err := host.Info()
	fmt.Println("err:", err)

	if err == nil {
		sysComputerName = sysInfo.Hostname
		sysOsName = sysInfo.OS
		sysOsArch = sysInfo.KernelArch
	}

	goName := "GoLang"             //语言环境
	goVersion := runtime.Version() //版本
	gtime.Date()
	//goStartTime := c.startTime //启动时间

	//goRunTime := gtime.Now().Timestamp() - c.startTime.Timestamp() //运行时长（秒）
	goHome := runtime.GOROOT() //安装路径
	goUserDir := ""            //项目路径

	curDir, err := os.Getwd()
	if err != nil {
		fmt.Println("err:", err)
	}

	if err == nil {
		goUserDir = curDir
	}

	//服务器磁盘信息
	diskList := make([]disk.UsageStat, 0)
	diskInfo, err := disk.Partitions(true) //所有分区
	if err == nil {
		for _, p := range diskInfo {
			diskDetail, err := disk.Usage(p.Mountpoint)
			if err == nil {
				diskDetail.UsedPercent, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", diskDetail.UsedPercent), 64)
				diskList = append(diskList, *diskDetail)
			}
		}
	}
	if err != nil {
		fmt.Println("err:", err)
	}
	res := GetCpuInfoRes{
		CpuNum:          cpuNum,
		CpuUsed:         cpuUsed,
		CpuAvg5:         gconv.String(cpuAvg5),
		CpuAvg15:        gconv.String(cpuAvg15),
		MemTotal:        memTotal,
		GoTotal:         goTotal,
		MemUsed:         memUsed,
		GoUsed:          goUsed,
		MemFree:         memFree,
		GoFree:          goFree,
		MemUsage:        memInfo,
		GoUsage:         goUsage,
		SysComputerName: sysComputerName,
		SysOsName:       sysOsName,
		SysComputerIp:   sysComputerIp,
		SysOsArch:       sysOsArch,
		GoName:          goName,
		GoVersion:       goVersion,
		//"goStartTime":     goStartTime,
		//"goRunTime":       goRunTime,
		GoHome:    goHome,
		GoUserDir: goUserDir,
		DiskList:  diskList,
	}
	if err != nil {
		fmt.Println("err:", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": res,
	})
}

// GetProblemList
// @Tags 公共方法
// @Summary 问题列表
// @Param page query int false "page"
// @Param size query int false "size"
// @Param keyword query string false "keyword"
// @Param category_identity query string false "category_identity"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /problem-list [get]
func GetProblemList(c *gin.Context) {
	size, _ := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Println("GetProblemList Page strconv Error:", err)
		return
	}

	page = (page - 1) * size

	var count int64
	keyword := c.Query("keyword")
	categoryIdentity := c.Query("category_identity")

	list := make([]*models.ProblemBasic, 0)

	// 查一下，这个关键词以及这个类型的题目有多少个
	err = models.GetProblemList(keyword, categoryIdentity).Distinct("`problem_basic`.`id`").Count(&count).Error
	if err != nil {
		log.Println("GetProblemList Count Error:", err)
		return
	}

	// 查询数据库题目
	err = models.GetProblemList(keyword, categoryIdentity).Offset(page).Limit(size).Find(&list).Error
	if err != nil {
		log.Println("Get Problem List Error:", err)
		return
	}

	// 响应
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"list":  list,
			"count": count,
		},
	})
}

// GetProblemDetail
// @Tags 公共方法
// @Summary 问题详情
// @Param identity query string false "problem identity"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /problem-detail [get]
func GetProblemDetail(c *gin.Context) {
	identity := c.Query("identity")
	if identity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "问题唯一标识不能为空",
		})
		return
	}
	data := new(models.ProblemBasic)
	err := models.DB.Where("identity = ?", identity).
		Preload("ProblemCategories").Preload("ProblemCategories.CategoryBasic").
		First(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "问题不存在",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get ProblemDetail Error:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": data,
	})
}

// ProblemCreate
// @Tags 管理员私有方法
// @Summary 问题创建
// @Accept json
// @Param authorization header string true "authorization"
// @Param data body define.ProblemBasic true "ProblemBasic"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/problem-create [post]
func ProblemCreate(c *gin.Context) {
	in := new(define.ProblemBasic)
	err := c.ShouldBindJSON(in)
	if err != nil {
		log.Println("[JsonBind Error] : ", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数错误",
		})
		return
	}

	if in.Title == "" || in.Content == "" || len(in.ProblemCategories) == 0 || len(in.TestCases) == 0 || in.MaxRuntime == 0 || in.MaxMem == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不能为空",
		})
		return
	}
	identity := helper.GetUUID()
	data := &models.ProblemBasic{
		Identity:   identity,
		Title:      in.Title,
		Content:    in.Content,
		MaxRuntime: in.MaxRuntime,
		MaxMem:     in.MaxMem,
		CreatedAt:  models.MyTime(time.Now()),
		UpdatedAt:  models.MyTime(time.Now()),
	}
	// 处理分类
	categoryBasics := make([]*models.ProblemCategory, 0)
	for _, id := range in.ProblemCategories {
		categoryBasics = append(categoryBasics, &models.ProblemCategory{
			ProblemId:  data.ID,
			CategoryId: uint(id),
			CreatedAt:  models.MyTime(time.Now()),
			UpdatedAt:  models.MyTime(time.Now()),
		})
	}
	data.ProblemCategories = categoryBasics
	// 处理测试用例
	testCaseBasics := make([]*models.TestCase, 0)
	for _, v := range in.TestCases {
		// 举个例子 {"input":"1 2\n","output":"3\n"}
		testCaseBasic := &models.TestCase{
			Identity:        helper.GetUUID(),
			ProblemIdentity: identity,
			Input:           v.Input,
			Output:          v.Output,
			CreatedAt:       models.MyTime(time.Now()),
			UpdatedAt:       models.MyTime(time.Now()),
		}
		testCaseBasics = append(testCaseBasics, testCaseBasic)
	}
	data.TestCases = testCaseBasics

	// 创建问题
	err = models.DB.Create(data).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Problem Create Error:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"identity": data.Identity,
		},
	})
}

// ProblemModify
// @Tags 管理员私有方法
// @Summary 问题修改
// @Param authorization header string true "authorization"
// @Param data body define.ProblemBasic true "ProblemBasic"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/problem-modify [put]
func ProblemModify(c *gin.Context) {
	in := new(define.ProblemBasic)
	err := c.ShouldBindJSON(in)
	if err != nil {
		log.Println("[JsonBind Error] : ", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数错误",
		})
		return
	}
	if in.Identity == "" || in.Title == "" || in.Content == "" || len(in.ProblemCategories) == 0 || len(in.TestCases) == 0 || in.MaxRuntime == 0 || in.MaxMem == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不能为空",
		})
		return
	}

	if err := models.DB.Transaction(func(tx *gorm.DB) error {
		// 问题基础信息保存 problem_basic
		problemBasic := &models.ProblemBasic{
			Identity:   in.Identity,
			Title:      in.Title,
			Content:    in.Content,
			MaxRuntime: in.MaxRuntime,
			MaxMem:     in.MaxMem,
			UpdatedAt:  models.MyTime(time.Now()),
		}
		err := tx.Where("identity = ?", in.Identity).Updates(problemBasic).Error
		if err != nil {
			return err
		}
		// 查询问题详情
		err = tx.Where("identity = ?", in.Identity).Find(problemBasic).Error
		if err != nil {
			return err
		}

		// 关联问题分类的更新
		// 1、删除已存在的关联关系
		err = tx.Where("problem_id = ?", problemBasic.ID).Delete(new(models.ProblemCategory)).Error
		if err != nil {
			return err
		}
		// 2、新增新的关联关系
		pcs := make([]*models.ProblemCategory, 0)
		for _, id := range in.ProblemCategories {
			pcs = append(pcs, &models.ProblemCategory{
				ProblemId:  problemBasic.ID,
				CategoryId: uint(id),
				CreatedAt:  models.MyTime(time.Now()),
				UpdatedAt:  models.MyTime(time.Now()),
			})
		}
		err = tx.Create(&pcs).Error
		if err != nil {
			return err
		}
		// 关联测试案例的更新
		// 1、删除已存在的关联关系
		err = tx.Where("problem_identity = ?", in.Identity).Delete(new(models.TestCase)).Error
		if err != nil {
			return err
		}
		// 2、增加新的关联关系
		tcs := make([]*models.TestCase, 0)
		for _, v := range in.TestCases {
			// 举个例子 {"input":"1 2\n","output":"3\n"}
			tcs = append(tcs, &models.TestCase{
				Identity:        helper.GetUUID(),
				ProblemIdentity: in.Identity,
				Input:           v.Input,
				Output:          v.Output,
				CreatedAt:       models.MyTime(time.Now()),
				UpdatedAt:       models.MyTime(time.Now()),
			})
		}
		err = tx.Create(tcs).Error
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Problem Modify Error:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "问题修改成功",
	})
}
