package engine

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/yaoapp/gou"
	"github.com/yaoapp/kun/exception"
	"github.com/yaoapp/kun/log"
	"github.com/yaoapp/yao/api"
	"github.com/yaoapp/yao/app"
	"github.com/yaoapp/yao/cert"
	"github.com/yaoapp/yao/chart"
	"github.com/yaoapp/yao/config"
	"github.com/yaoapp/yao/flow"
	"github.com/yaoapp/yao/importer"
	"github.com/yaoapp/yao/model"
	"github.com/yaoapp/yao/page"
	"github.com/yaoapp/yao/plugin"
	"github.com/yaoapp/yao/query"
	"github.com/yaoapp/yao/script"
	"github.com/yaoapp/yao/server"
	"github.com/yaoapp/yao/share"
	"github.com/yaoapp/yao/store"
	"github.com/yaoapp/yao/table"
	"github.com/yaoapp/yao/workflow"
)

// Load 根据配置加载 API, FLow, Model, Plugin
func Load(cfg config.Config) (err error) {
	defer func() { err = exception.Catch(recover()) }()

	// 加载应用信息
	// 第一步: 加载应用信息
	app.Load(cfg)

	// 加密密钥函数
	gou.LoadCrypt(fmt.Sprintf(`{"key":"%s"}`, cfg.DB.AESKey), "AES")
	gou.LoadCrypt(`{}`, "PASSWORD")

	// Load Certs
	err = cert.Load(cfg)
	if err != nil {
		log.Debug(err.Error())
	}

	// 第二步: 建立数据库 & 会话连接
	share.DBConnect(cfg.DB) // 创建数据库连接
	// share.SessionConnect(cfg.Session) // 创建会话服务器链接

	// 加载应用引擎
	if os.Getenv("YAO_DEV") != "" {
		LoadEngine(filepath.Join(os.Getenv("YAO_DEV"), "/yao"))
	} else {
		LoadEngine()
	}

	// 第三步: 加载数据分析引擎
	query.Load(cfg) // 加载数据分析引擎

	// 第四步: 加载共享库 & JS 处理器
	err = share.Load(cfg) // 加载共享库 lib
	if err != nil {
		log.Debug(err.Error())
	}

	err = script.Load(cfg) // 加载JS处理器 script
	if err != nil {
		log.Debug(err.Error())
	}

	// 第五步: 加载数据模型等
	err = model.Load(cfg) // 加载数据模型 model
	if err != nil {
		log.Debug(err.Error())
	}

	err = flow.Load(cfg) // 加载业务逻辑 Flow
	if err != nil {
		log.Debug(err.Error())
	}

	err = store.Load(cfg) // Load stores
	if err != nil {
		log.Debug(err.Error())
	}

	err = plugin.Load(cfg) // 加载业务插件 plugin
	if err != nil {
		log.Debug(err.Error())
	}

	err = table.Load(cfg) // 加载数据表格 table
	if err != nil {
		log.Debug(err.Error())
	}

	err = chart.Load(cfg) // 加载分析图表 chart
	if err != nil {
		log.Debug(err.Error())
	}

	page.Load(cfg)     // 加载页面 page 忽略错误
	importer.Load(cfg) // 加载数据导入 imports
	workflow.Load(cfg) // 加载工作流  workflow

	err = api.Load(cfg) // 加载业务接口 API
	if err != nil {
		log.Debug(err.Error())
	}

	server.Load(cfg) // 加载服务
	return nil
}

// Reload 根据配置重新加载 API, FLow, Model, Plugin
func Reload(cfg config.Config) {
	gou.APIs = map[string]*gou.API{}
	gou.Models = map[string]*gou.Model{}
	gou.Flows = map[string]*gou.Flow{}
	gou.Plugins = map[string]*gou.Plugin{}
	Load(cfg)
}

// LoadEngine 加载引擎的 API, Flow, Model 配置
func LoadEngine(from ...string) {
	var scripts []share.Script
	if len(from) > 0 {
		scripts = share.GetFilesFS(from[0], ".json")
	} else {
		scripts = share.GetFilesBin("yao", ".json")
	}

	if scripts == nil {
		exception.New("读取文件失败", 500, from).Throw()
	}

	if len(scripts) == 0 {
		exception.New("读取文件失败, 未找到任何可执行脚本", 500, from).Throw()
	}

	// 加载 API, Flow, Models, Table, Chart, Screens
	for _, script := range scripts {
		switch script.Type {
		case "models":
			gou.LoadModel(string(script.Content), "xiang."+script.Name)
			break
		case "flows":
			gou.LoadFlow(string(script.Content), "xiang."+script.Name)
			break
		case "apis":
			gou.LoadAPI(string(script.Content), "xiang."+script.Name)
			break
		}
	}

	// 加载数据应用
	for _, script := range scripts {
		switch script.Type {
		case "tables":
			table.LoadTable(string(script.Content), "xiang."+script.Name)
			break
		}
	}
}
