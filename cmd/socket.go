package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cobra"
	"github.com/yaoapp/gou"
	"github.com/yaoapp/kun/exception"
	"github.com/yaoapp/yao/config"
	"github.com/yaoapp/yao/engine"
	"github.com/yaoapp/yao/share"
)

var socketCmd = &cobra.Command{
	Use:   "socket",
	Short: L("Open socket"),
	Long:  L("Open socket"),
	Run: func(cmd *cobra.Command, args []string) {
		defer share.SessionStop()
		defer gou.KillPlugins()
		defer func() {
			err := exception.Catch(recover())
			if err != nil {
				fmt.Println(color.RedString(L("Fatal: %s"), err.Error()))
			}
		}()

		Boot()
		cfg := config.Conf
		cfg.Session.IsCLI = true
		engine.Load(cfg)
		if len(args) < 1 {
			fmt.Println(color.RedString(L("Not enough arguments")))
			fmt.Println(color.WhiteString(share.BUILDNAME + " help"))
			return
		}

		name := args[0]
		pargs := []interface{}{}
		for i, arg := range args {
			if i == 0 {
				continue
			}

			// 解析参数
			if strings.HasPrefix(arg, "::") {
				arg := strings.TrimPrefix(arg, "::")
				var v interface{}
				err := jsoniter.Unmarshal([]byte(arg), &v)
				if err != nil {
					fmt.Println(color.RedString(L("Arguments: %s"), err.Error()))
					return
				}
				pargs = append(pargs, v)
				fmt.Println(color.WhiteString("args[%d]: %s", i-1, arg))
			} else if strings.HasPrefix(arg, "\\::") {
				arg := "::" + strings.TrimPrefix(arg, "\\::")
				pargs = append(pargs, arg)
				fmt.Println(color.WhiteString("args[%d]: %s", i-1, arg))
			} else {
				pargs = append(pargs, arg)
				fmt.Println(color.WhiteString("args[%d]: %s", i-1, arg))
			}

		}

		socket, has := gou.Sockets[name]
		if !has {
			fmt.Println(color.RedString(L("%s not exists!"), name))
			return
		}

		if socket.Mode != "client" {
			fmt.Println(color.RedString(L("%s not supported yet!"), socket.Mode))
			return
		}

		fmt.Println(color.WhiteString("\n---------------------------------"))
		fmt.Println(color.WhiteString(socket.Name))
		fmt.Println(color.WhiteString("---------------------------------"))
		fmt.Println(color.GreenString("Mode: %s", socket.Mode))
		fmt.Println(color.GreenString("Host: %s://%s", socket.Protocol, socket.Host))
		fmt.Println(color.GreenString("Port: %s", socket.Port))
		fmt.Println(color.WhiteString("--------------------------------------"))
		err := socket.Open(pargs...)
		if err != nil {
			fmt.Println(color.RedString(L("%s"), err.Error()))
			return
		}

	},
}
