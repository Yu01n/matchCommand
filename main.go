/*
@Title : 匹配linux/Windows操作系统命令，并返回相应选项的帮助信息
@Author : yuan Y
@File : main.go
@Software: GoLand
*/
package main

import (
	"C"
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

// Path 路径对象 存放操作系统命令的json文件的路径
type Path struct {
	Path []string `json:"path"`
}

// Cmd 命令对象，包含了linux的部分命令
type Cmd struct {
	Os_id    string     `json:"os_id" yaml:"os_id"`       // 操作系统命令id
	Os       string     `json:"os" yaml:"os"`             // 操作系统名称
	Cmd_info []Cmd_info `json:"cmd_info" yaml:"cmd_info"` // 操作系统命令详细信息
}

// Cmd_info 操作系统命令详细信息
type Cmd_info struct {
	Cmd_category_id string     `json:"cmd_category_id" yaml:"cmd_category_id"` // 一级类别id
	Cmd_category    string     `json:"cmd_category" yaml:"cmd_category"`       // 一级类别名称
	Cmd_list        []Cmd_list `json:"cmd_list" yaml:"cmd_list"`               // 二级类别列表
}

// Cmd_list 操作系统命令列表
type Cmd_list struct {
	Cmd_sub_category_id string         `json:"cmd_sub_category_id" yaml:"cmd_sub_category_id"` // 二级类别id
	Cmd_sub_category    string         `json:"cmd_sub_category" yaml:"cmd_sub_category"`       // 二级类别名称
	Cmd_sub_list        []Cmd_sub_list `json:"cmd_sub_list" yaml:"cmd_sub_list"`               // 命令列表
}

// Cmd_sub_list 操作系统每条命令详细信息
type Cmd_sub_list struct {
	Cmd_id         string           `json:"cmd_id" yaml:"cmd_id"`                 // 命令id
	Cmd_name       string           `json:"cmd_name" yaml:"cmd_name"`             // 命令名称
	Cmd_desc       string           `json:"cmd_desc" yaml:"cmd_desc"`             // 命令描述
	Cmd_parameters []Cmd_parameters `json:"cmd_parameters" yaml:"cmd_parameters"` // 命令参数列表
	Cmd_examples   string           `json:"cmd_examples" yaml:"cmd_examples"`     // 命令实例
	Cmd_os         string           `json:"cmd_os" yaml:"cmd_os"`                 // 命令所属操作系统名称
	Cmd_type       string           `json:"cmd_type" yaml:"cmd_type"`             // 命令参数的类型设置
}

// Cmd_parameters 参数列表的详细信息
type Cmd_parameters struct {
	Cmd_parameter_key   string `json:"cmd_parameter_key" yaml:"cmd_parameter_key"`     // 参数名称
	Cmd_parameter_value string `json:"cmd_parameter_value" yaml:"cmd_parameter_value"` // 参数描述
}

// 声明存放所有信息及所有参数的cmdMatchSlice切片
var cmdMatchSlice []cmdMatch

var Linux_cmd_Parameters map[string]mapCmdInfo
var Windows_cmd_Parameters map[string]mapCmdInfo

//对json每一层的值是否有效进行检查，并获取value值
func (cmd *Cmd) validCheck() {
	if cmd.Os == "linux" {
		Linux_cmd_Parameters = make(map[string]mapCmdInfo)
	} else if cmd.Os == "Windows" {
		Windows_cmd_Parameters = make(map[string]mapCmdInfo)
	} else {
		return
	}

	for _, item := range cmd.Cmd_info {
		item.validCheck(cmd.Os)
	}
}
func (cmdInfo *Cmd_info) validCheck(cmd_os string) {
	for _, item := range cmdInfo.Cmd_list {
		item.validCheck(cmd_os)
	}
}
func (cmdList *Cmd_list) validCheck(cmd_os string) {
	for _, item := range cmdList.Cmd_sub_list {
		item.validCheck(cmd_os)
	}
}
func (cmdSubList *Cmd_sub_list) validCheck(cmd_os string) {
	constructparameter(cmd_os, cmdSubList.Cmd_name, cmdSubList.Cmd_desc, cmdSubList.Cmd_parameters)
}


type mapCmdInfo struct {
	Cmd_desc string
	Cmd_os string
	Cmd_parameters []Cmd_parameters
}

func constructparameter(cmd_os string, cmd_name string, cmd_desc string, cmd_params []Cmd_parameters){
	if cmd_os == "linux" {
		Linux_cmd_Parameters[cmd_name] = mapCmdInfo{cmd_desc, cmd_os,cmd_params}
	}else if cmd_os == "Windows" {
		Windows_cmd_Parameters[cmd_name] = mapCmdInfo{cmd_desc, cmd_os, cmd_params}
	} else {
		return
	}

}

// 将需要获取的全部参数内容存在cmdAllMatch结构体中
type cmdMatch struct {
	Cmd_line           string           // 存放传入的命令行
	Cmd_os           string           // 存放获取的获取Cmd_os
	Cmd_name           string           // 存放获取的获取Cmd_name
	Cmd_desc           string           // 存放获取的获取Cmd_desc
	Cmd_parameters []Cmd_parameters // 存放获取到的所有参数或选项
}

// 存放对多条命令同时查询时所有参数的最终输出结果
type cmdMatchslice struct {
	CmdMatchSlice []cmdMatch `json:"Cmd_info"`
}

// 存放获取的需求参数结果
type parameterslice struct {
	Cmd_current_parameters []Cmd_parameters
	Cmd_all_parameters     []Cmd_parameters
}


func main() {}

//读取json文件内容
func file_get_contents(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

// 保存每个json文件的修改时间
var jsonFileModTime map[string]time.Time

// 保存每个json文件Unmarshal后的Cmd结构体
var jsonCmd map[string]*Cmd

// 设置只解析一次json文件
func getJsonOnce(item string) {
	fileStat, err :=os.Stat(item)
	if err != nil{
		return
	}
	itemModTime := fileStat.ModTime() // 获取json文件修改时间

	var needUnmarshal = false

	// 初始化
	if jsonFileModTime == nil {
		jsonFileModTime = make(map[string]time.Time)
	}
	if jsonCmd == nil {
		jsonCmd = make(map[string]*Cmd)
	}

	if _, ok := jsonFileModTime[item]; !ok {
		jsonFileModTime[item] = itemModTime
		needUnmarshal = true
	}

	// 判断json文件是否被修改
	if(jsonFileModTime[item] != itemModTime) {
		needUnmarshal = true
	}

	// Unmarshal
	if(needUnmarshal) {
		var cmd Cmd

		// 调用方法获取json文件的内容
		content, err := file_get_contents(item)
		err = json.Unmarshal([]byte(content), &cmd)
		if err != nil {
			return
		}

		cmd.validCheck()

		// 保存Unmarshal后的cmd结构
		jsonCmd[item] = &cmd
	}
}

// 获取json文件路径
func getFilePathList(filepath string) []string {
	var p Path
	var content []byte
	// 调用方法获取路径json文件的内容
	var fileName string = filepath + "/configs/path.json"
	content, err := file_get_contents(fileName)
	if err != nil {
		return OPENERROR
	}
	err = json.Unmarshal([]byte(content), &p)
	if err != nil {
		return GETERROR
	}
	return p.Path
}

var OPENERROR = []string{"open file error"}
var GETERROR = []string{"get file content error"}

// 获取命令json文件中的所有内容
func getPathContent(filepath string) []string {
	// 获取json文件路径
	path := getFilePathList(filepath)
	return path
}


// man方法：返回单个命令行所有的信息以及全部参数信息
//export man
func man(getfilepath, process_name, os *C.char) *C.char {
//func man(getfilepath, process_name, os string) string {
	// timeout control
	timeout := 1000 * time.Millisecond
	ctx, cancle := context.WithTimeout(context.Background(), timeout)
	defer cancle()

	resChan := make(chan string, 1)

	if C.GoString(process_name) == "" {
		return C.CString("")
	}

	//go func(getfilepath, process_name, os string, resChan chan string) {
	go func(getfilepath, process_name, os *C.char, resChan chan string) {
		// recover panic
		defer func() {
			if re := recover(); re != nil {
				resChan <- ""
				return
			}
		}()

		//cmdOs := os
		//cmdLine := process_name
		//filepath := getfilepath
		cmdOs := C.GoString(os)
		cmdLine := C.GoString(process_name)
		filepath := C.GoString(getfilepath)

		var cmdInfo string
		// 根据操作系统的类型不同执行不同的函数
		if cmdOs == "" {
			cmdInfo = cmdAllParamenters(filepath, cmdLine)
		} else if strings.EqualFold(cmdOs, "linux") {
			cmdInfo = linuxCmdAllCommand(filepath, cmdLine)
		} else if strings.EqualFold(cmdOs, "windows") {
			cmdInfo = windowsCmdAllCommand(filepath, cmdLine)
		}
		resChan <- cmdInfo
	}(getfilepath, process_name, os, resChan)

	select {
	case <- ctx.Done():
		return C.CString("")
		//return ""
	case result := <-resChan:
		return C.CString(result)
		//return result
	}

	return C.CString("")
	//return ""

}

// man_list方法：返回多个命令行所有的信息以及全部参数信息
//export man_list
func man_list(getfilepath, process_name, os *C.char) *C.char {
//func man_list(getfilepath, process_name, os string) string {

	// timeout control
	timeout := 10000 * time.Millisecond
	ctx, cancle := context.WithTimeout(context.Background(), timeout)
	defer cancle()

	resChan := make(chan string, 1)

	if C.GoString(process_name) == "" {
		return C.CString("")
	}

	//go func(getfilepath, process_name, os string, resChan chan string) {
	go func(getfilepath, process_name, os *C.char, resChan chan string) {
		// recover panic
		defer func() {
			if re := recover(); re != nil {
				resChan <- ""
				return
			}
		}()
		//cmdOs := os
		//cmdLine := process_name
		//filepath := getfilepath
		cmdOs := C.GoString(os)
		cmdLine := C.GoString(process_name)
		filepath := C.GoString(getfilepath)

		var cmdInfo string
		// 根据操作系统的类型不同执行不同的函数
		if cmdOs == "" {
			cmdInfo = allParamenters(filepath, cmdLine)
		} else if strings.EqualFold(cmdOs, "linux") {
			cmdInfo = linuxAllCommand(filepath, cmdLine)
		} else if strings.EqualFold(cmdOs, "windows") {
			cmdInfo = windowsAllCommand(filepath, cmdLine)
		}
		resChan <- cmdInfo
		}(getfilepath, process_name, os, resChan)

		select {
		case <-ctx.Done():
			return C.CString("")
			//return ""
		case result := <-resChan:
			return C.CString(result)
			//return result
	}
	return C.CString("")
	//return ""
}

// explain_cmd方法：返回单个命令行所有的信息以及当前参数信息
//export explain_cmd
func explain_cmd(getfilepath, cmd_string, os *C.char) *C.char {
//func explain_cmd(getfilepath, cmd_string, os string) string {

	// timeout control
	timeout := 1000 * time.Millisecond
	ctx, cancle := context.WithTimeout(context.Background(), timeout)
	defer cancle()
	resChan := make(chan string, 1)

	if C.GoString(cmd_string) == "" {
		return C.CString("")
	}

	//go func(getfilepath, cmd_string, os string, resChan chan string) {
	go func(getfilepath, cmd_string, os *C.char, resChan chan string) {
		// recover panic
		defer func() {
			if re := recover(); re != nil {
				resChan <- ""
				return
			}
		}()
		//cmdOs := os
		//cmdLine := cmd_string
		//filepath := getfilepath
		cmdOs := C.GoString(os)
		cmdLine := C.GoString(cmd_string)
		filepath := C.GoString(getfilepath)


		var cmdInfo string
		if cmdOs == "" {
			cmdInfo = currentCmdParamenters(filepath, cmdLine)
		} else if strings.EqualFold(cmdOs, "linux") {
			cmdInfo = linuxCurrentCmdCommand(filepath, cmdLine)
		} else if strings.EqualFold(cmdOs, "windows") {
			cmdInfo = windowsCurrentCmdCommand(filepath, cmdLine)
		}
			resChan <- cmdInfo
		}(getfilepath, cmd_string, os, resChan)

	select {
	case <-ctx.Done():
		return C.CString("")
		//return ""
	case result := <-resChan:
		return C.CString(result)
		//return result
	}
	return C.CString("")
	//return ""
}

// explain_cmd_list方法：返回多个命令行所有的信息以及当前参数信息
//export explain_cmd_list
func explain_cmd_list(getfilepath, cmd_string, os *C.char) *C.char {
//func explain_cmd_list(getfilepath, cmd_string, os string) string {

	// timeout control
	timeout := 10000 * time.Millisecond
	ctx, cancle := context.WithTimeout(context.Background(), timeout)
	defer cancle()
	resChan := make(chan string, 1)

	if C.GoString(cmd_string) == "" {
		return C.CString("")
	}

	//go func(getfilepath, cmd_string, os string, resChan chan string) {
	go func(getfilepath, cmd_string, os *C.char, resChan chan string) {
		// recover panic
		defer func() {
			if re := recover(); re != nil {
				resChan <- ""
				return
			}
		}()
		//cmdOs := os
		//cmdLine := cmd_string
		//filepath := getfilepath
		cmdOs := C.GoString(os)
		cmdLine := C.GoString(cmd_string)
		filepath := C.GoString(getfilepath)

		if cmdLine == ""{
			return
		}
		var cmdInfo string
		if cmdOs == "" {
			cmdInfo = currentParamenters(filepath, cmdLine)
		} else if strings.EqualFold(cmdOs, "linux") {
			cmdInfo = linuxCurrentCommand(filepath, cmdLine)
		} else if strings.EqualFold(cmdOs, "windows") {
			cmdInfo = windowsCurrentCommand(filepath, cmdLine)
		}
			resChan <- cmdInfo
	}(getfilepath, cmd_string, os, resChan)

	select {
	case <-ctx.Done():
		return C.CString("")
		//return ""
	case result := <-resChan:
		return C.CString(result)
		//return result
	}
	return C.CString("")
	//return ""
}

// 设置按照 ':' 和 '=' 分割
func Split(r rune) bool {
	return r == ':' || r == '='
}

// 对单个命令行进行匹配，并获取匹配到的内容(linux)
func getcmdslicecontent(context []string, item string, cmdCurrentslice *cmdMatchslice) {
	var parameter parameterslice
	//判断命令是否匹配成功
	for CmdName, Parameters := range Linux_cmd_Parameters{
		//keys = append(keys, Cmd_name)
		if context[0] != CmdName {
			continue
		}
		for m := 0; m < len(context) -1; m++ {
			if len(context[m+1]) == 1{
				continue
			}
			if context[m+1][0] != '-' || context[m+1][1] == '-' {
				// 判断选项的值是否和参数列表中的key值匹配，匹配成功获取参数的key -- value值
				for i :=0; i< len(Parameters.Cmd_parameters);i++ {
					if context[m+1] == Parameters.Cmd_parameters[i].Cmd_parameter_key {
						parameter.Cmd_current_parameters = append(parameter.Cmd_current_parameters, Cmd_parameters{Cmd_parameter_key: Parameters.Cmd_parameters[i].Cmd_parameter_key, Cmd_parameter_value: Parameters.Cmd_parameters[i].Cmd_parameter_value})
					}
				}
			} else if context[m+1][1] != '-' { // 选项或参数存在'-'时执行
				str := context[m+1][1:]
				var a int
				// 将选项的每个字母分开，存放在tmpArr中
				var tmpArr []string
				for a = 0; a < len(str); a++ {
					// 给每个字母前添加 '-' 用于做后续的匹配
					tmpArr = append(tmpArr, "-" + string(str[a]))
					// 判断选项的值是否和参数列表中的key值匹配，匹配成功获取参数的key -- value值
					for i :=0; i< len(Parameters.Cmd_parameters);i++ {
						if tmpArr[a] == Parameters.Cmd_parameters[i].Cmd_parameter_key {
							parameter.Cmd_current_parameters = append(parameter.Cmd_current_parameters, Cmd_parameters{Cmd_parameter_key: Parameters.Cmd_parameters[i].Cmd_parameter_key, Cmd_parameter_value: Parameters.Cmd_parameters[i].Cmd_parameter_value})
						}
					}
				}
			}
		}
		cmdCurrentslice.CmdMatchSlice = append(cmdCurrentslice.CmdMatchSlice, cmdMatch{Cmd_line: item, Cmd_name: CmdName,Cmd_os: Parameters.Cmd_os,  Cmd_desc: Parameters.Cmd_desc, Cmd_parameters: parameter.Cmd_current_parameters})
	}
	parameter.Cmd_current_parameters = nil
}

// 对单个命令行进行匹配，并获取匹配到的内容(windows)
func getwindowscmdslicecontent(context []string,item string, cmdCurrentslice *cmdMatchslice) {
	var parameter parameterslice
	for CmdName, Parameters := range Windows_cmd_Parameters {
		//判断命令是否匹配成功
		if !strings.EqualFold(context[0], CmdName) {
			continue
		}
			for l := 0; l < len(Parameters.Cmd_parameters); l++ {
				for m := 0; m < len(context)-1; m++ {
					// 分割出包含 ':' 和 '=' 的参数
					if strings.Contains(context[m+1], ":") || strings.Contains(context[m+1], "=") {
						var para = strings.FieldsFunc(context[m+1], Split)
						// 判断选项的值是否和参数列表中的key值匹配，匹配成功获取参数的key -- value值
						if strings.EqualFold(para[0], Parameters.Cmd_parameters[l].Cmd_parameter_key) {
							parameter.Cmd_current_parameters = append(parameter.Cmd_current_parameters, Cmd_parameters{Cmd_parameter_key: Parameters.Cmd_parameters[l].Cmd_parameter_key, Cmd_parameter_value: Parameters.Cmd_parameters[l].Cmd_parameter_value})
						}
					} else if strings.EqualFold(context[m+1], Parameters.Cmd_parameters[l].Cmd_parameter_key) { // 判断参数是否能匹配上，匹配成功获取参数的key -- value值
						parameter.Cmd_current_parameters = append(parameter.Cmd_current_parameters, Cmd_parameters{Cmd_parameter_key: Parameters.Cmd_parameters[l].Cmd_parameter_key, Cmd_parameter_value: Parameters.Cmd_parameters[l].Cmd_parameter_value})
					}
				}
			}

		cmdCurrentslice.CmdMatchSlice = append(cmdCurrentslice.CmdMatchSlice, cmdMatch{Cmd_line: item, Cmd_name: CmdName, Cmd_os: Parameters.Cmd_os,  Cmd_desc: Parameters.Cmd_desc, Cmd_parameters: parameter.Cmd_current_parameters})
	}
	parameter.Cmd_current_parameters = nil
}

// 获取linux下的CurrentSlice的值
func getCurrentSlice(context []string, cmdLine string, currentslice *[]cmdMatch) {
	var parameter parameterslice
	for CmdName, Parameters := range Linux_cmd_Parameters {
		//判断命令是否匹配成功
		if context[0] != CmdName {
			continue
		}
		for m := 0; m < len(context)-1; m++ {
			if len(context[m+1]) == 1{
				continue
			}
			if context[m+1][0] != '-' || context[m+1][1] == '-' {
				// 判断选项的值是否和参数列表中的key值匹配，匹配成功获取参数的key -- value值
				for i := 0; i < len(Parameters.Cmd_parameters); i++ {
					if context[m+1] == Parameters.Cmd_parameters[i].Cmd_parameter_key {
						parameter.Cmd_current_parameters = append(parameter.Cmd_current_parameters, Cmd_parameters{Cmd_parameter_key: Parameters.Cmd_parameters[i].Cmd_parameter_key, Cmd_parameter_value: Parameters.Cmd_parameters[i].Cmd_parameter_value})
					}
				}
			} else if context[m+1][1] != '-' { // 选项或参数存在'-'时执行
				str := context[m+1][1:]
				var a int
				// 将选项的每个字母分开，存放在tmpArr切片中
				var tmpArr []string
				for a = 0; a < len(str); a++ {
					// 给每个字母前添加 '-' 用于做后续的匹配
					tmpArr = append(tmpArr, "-"+string(str[a]))
					// 判断选项的值是否和参数列表中的key值匹配，匹配成功获取参数的key -- value值
					for i := 0; i < len(Parameters.Cmd_parameters); i++ {
						if tmpArr[a] == Parameters.Cmd_parameters[i].Cmd_parameter_key {
							parameter.Cmd_current_parameters = append(parameter.Cmd_current_parameters, Cmd_parameters{Cmd_parameter_key: Parameters.Cmd_parameters[i].Cmd_parameter_key, Cmd_parameter_value: Parameters.Cmd_parameters[i].Cmd_parameter_value})
						}
					}
				}
			}
		}
		*currentslice = append(*currentslice, cmdMatch{Cmd_line: cmdLine, Cmd_name: CmdName, Cmd_os: Parameters.Cmd_os,  Cmd_desc: Parameters.Cmd_desc, Cmd_parameters: parameter.Cmd_current_parameters})
	}
	parameter.Cmd_current_parameters = nil


}

// 获取Windows下的CurrentSlice的值
func getWindowsCurrentSlice(context []string, cmdLine string, currentslice *[]cmdMatch) {
	var parameter parameterslice
	for CmdName, Parameters := range Windows_cmd_Parameters {
	//判断命令是否匹配成功
		if !strings.EqualFold(context[0], CmdName) {
			continue
		}
		for m := 0; m < len(context)-1; m++ {
			if len(context[m+1]) == 1{
				continue
			}
			if context[m+1][0] != '-' || context[m+1][1] == '-' {
				// 判断选项的值是否和参数列表中的key值匹配，匹配成功获取参数的key -- value值
				for i :=0; i< len(Parameters.Cmd_parameters);i++ {
					if context[m+1] == Parameters.Cmd_parameters[i].Cmd_parameter_key {
						parameter.Cmd_current_parameters = append(parameter.Cmd_current_parameters, Cmd_parameters{Cmd_parameter_key: Parameters.Cmd_parameters[i].Cmd_parameter_key, Cmd_parameter_value: Parameters.Cmd_parameters[i].Cmd_parameter_value})
					}
				}
			} else if context[m+1][1] != '-' { // 选项或参数存在'-'时执行
				str := context[m+1][1:]
				var a int
				// 将选项的每个字母分开，存放在tmpArr中
				var tmpArr []string
				for a = 0; a < len(str); a++ {
					// 给每个字母前添加 '-' 用于做后续的匹配
					tmpArr = append(tmpArr, "-" + string(str[a]))
					// 判断选项的值是否和参数列表中的key值匹配，匹配成功获取参数的key -- value值
					for i :=0; i< len(Parameters.Cmd_parameters);i++ {
						if tmpArr[a] == Parameters.Cmd_parameters[i].Cmd_parameter_key {
							parameter.Cmd_current_parameters = append(parameter.Cmd_current_parameters, Cmd_parameters{Cmd_parameter_key: Parameters.Cmd_parameters[i].Cmd_parameter_key, Cmd_parameter_value: Parameters.Cmd_parameters[i].Cmd_parameter_value})
						}
					}
				}
			}
		}
	*currentslice = append(*currentslice, cmdMatch{Cmd_line: cmdLine, Cmd_name: CmdName, Cmd_os: Parameters.Cmd_os,  Cmd_desc: Parameters.Cmd_desc, Cmd_parameters: parameter.Cmd_current_parameters})
	}
	parameter.Cmd_current_parameters = nil
}

// 获取Linux的allSlice的值
func getlLinuxAllSlice(cmdLine string, context []string, allslice *[]cmdMatch) {
	var parameter parameterslice
	for CmdName, Parameters := range Linux_cmd_Parameters {
		//判断命令是否匹配成功
		if context[0] != CmdName {
			continue
		}
		for l := 0; l < len(Parameters.Cmd_parameters); l++ {
			parameter.Cmd_all_parameters = append(parameter.Cmd_all_parameters, Cmd_parameters{Cmd_parameter_key: Parameters.Cmd_parameters[l].Cmd_parameter_key, Cmd_parameter_value: Parameters.Cmd_parameters[l].Cmd_parameter_value})
		}
		*allslice = append(*allslice, cmdMatch{Cmd_line: cmdLine, Cmd_os: Parameters.Cmd_os, Cmd_name: CmdName, Cmd_desc: Parameters.Cmd_desc, Cmd_parameters: parameter.Cmd_all_parameters})
	}
	parameter.Cmd_all_parameters = nil
}

// 获取Windows下的allSlice的值
func getlWindowsAllSlice(cmdLine string, context []string, allslice *[]cmdMatch) {
	var parameter parameterslice
	for CmdName, Parameters := range Windows_cmd_Parameters {
		//判断命令是否匹配成功
		if !strings.EqualFold(context[0], CmdName) {
			continue
		}
		for l := 0; l < len(Parameters.Cmd_parameters); l++ {
			parameter.Cmd_all_parameters = append(parameter.Cmd_all_parameters, Cmd_parameters{Cmd_parameter_key: Parameters.Cmd_parameters[l].Cmd_parameter_key, Cmd_parameter_value: Parameters.Cmd_parameters[l].Cmd_parameter_value})
		}
		*allslice = append(*allslice, cmdMatch{Cmd_line: cmdLine,  Cmd_name: CmdName, Cmd_os: Parameters.Cmd_os, Cmd_desc: Parameters.Cmd_desc, Cmd_parameters: parameter.Cmd_all_parameters})
	}
	parameter.Cmd_all_parameters = nil


}

// 不区分大小写比较
func ContainsI(a string, b string) bool{
	return strings.Contains(
		strings.ToLower(a),
		strings.ToLower(b),
	)
}

// 去掉Linux中带有路径命令行的路径
func cmdLinuxPrefix(context []string) []string{
	temp_linux := []string{"/bin", "/sbin", "/lib", "/opt", "/home", "/data", "/hids", "/usr"}
	for _, item_temp := range temp_linux{
		if ContainsI(context[0],item_temp){
			var temp = strings.Split(context[0],"/")
			context[0] = temp[len(temp)-1]
		}
	}
	return context
}

// 去掉Windows中带有路径命令行的路径
func cmdWindowsPrefix(context []string) []string{
	if ContainsI(context[0], ".exe") {
		tmp_context := strings.ToLower(context[0])
		context[0] = strings.TrimSuffix(tmp_context, ".exe")
	}
	temp_windows := []string{"\\java", "\\root","\\AppData","\\Local","\\System32","\\TEMP","\\Windows","\\app","\\application", "\\batch","\\Program Files (x86)","\\software","\\Program Files"}
	for _, item_temp := range temp_windows{
		if ContainsI(context[0], item_temp) {
			var temp = strings.Split(context[0], "\\")
			context[0] = temp[len(temp)-1]
		}
	}
	return context
}

// 获取linux下的 cmdMatchslice的值
func getCmdMathchSlice(cmdLine string, cmdMatchslice *cmdMatchslice) {
	var parameter parameterslice
	// 将多行命令按照有分号的地方分割
	contexts := strings.Split(cmdLine, ";")
	for _, item := range contexts {
		var context = strings.Fields(item)
		// 去掉带有路径命令行的路径
		cmdLinuxPrefix(context)
		for CmdName, Parameters := range Linux_cmd_Parameters {
			//判断命令是否匹配成功
			if context[0] != CmdName {
				continue
			}
			for l := 0; l < len(Parameters.Cmd_parameters); l++ {
				parameter.Cmd_all_parameters = append(parameter.Cmd_all_parameters, Cmd_parameters{Cmd_parameter_key: Parameters.Cmd_parameters[l].Cmd_parameter_key, Cmd_parameter_value: Parameters.Cmd_parameters[l].Cmd_parameter_value})
			}
			cmdMatchslice.CmdMatchSlice = append(cmdMatchslice.CmdMatchSlice, cmdMatch{Cmd_line: item, Cmd_name: CmdName, Cmd_os: Parameters.Cmd_os,  Cmd_desc: Parameters.Cmd_desc, Cmd_parameters: parameter.Cmd_all_parameters})
		}
		parameter.Cmd_all_parameters = nil
	}
}

// 获取windows下的CmdMatchslice的值
func getWindowsCmdMathchSlice(cmdLine string, cmdMatchslice *cmdMatchslice) {
	var parameter parameterslice
	// 将多行命令按照有分号的地方分割
	contexts := strings.Split(cmdLine, ";")
	for _, item := range contexts {
		var context = strings.Fields(item)
		// 去掉带有路径命令行的路径
		cmdWindowsPrefix(context)
		for CmdName, Parameters := range Windows_cmd_Parameters {
			//判断命令是否匹配成功
			if !strings.EqualFold(context[0], CmdName) {
				continue
			}
			for l := 0; l < len(Parameters.Cmd_parameters); l++ {
				parameter.Cmd_all_parameters = append(parameter.Cmd_all_parameters, Cmd_parameters{Cmd_parameter_key: Parameters.Cmd_parameters[l].Cmd_parameter_key, Cmd_parameter_value: Parameters.Cmd_parameters[l].Cmd_parameter_value})
			}
			cmdMatchslice.CmdMatchSlice = append(cmdMatchslice.CmdMatchSlice, cmdMatch{Cmd_line: item, Cmd_name: CmdName, Cmd_os: Parameters.Cmd_os,  Cmd_desc: Parameters.Cmd_desc, Cmd_parameters: parameter.Cmd_all_parameters})
		}
		parameter.Cmd_all_parameters = nil


	}
}

// 对多个命令行进行分割匹配，并获取值(linux)
func getCmdCurrentSlice( cmdLine string, cmdCurrentslice *cmdMatchslice) {
	var contexts []string
	contexts = strings.Split(cmdLine, ";")
	for _, item := range contexts {
		var context = strings.Fields(item)
		// 去掉带有路径命令行的路径
		cmdLinuxPrefix(context)
		getcmdslicecontent(context, item, cmdCurrentslice)
	}
}

// 对多个命令行进行分割匹配，并获取值(Windows)
func getWindowsCmdCurrentSlice(cmdLine string, cmdCurrentslice *cmdMatchslice) {
	var contexts []string
	contexts = strings.Split(cmdLine, ";")
	for _, item := range contexts {
		var context = strings.Fields(item)
		// 去掉带有路径命令行的路径
		cmdWindowsPrefix(context)
		getwindowscmdslicecontent(context,item, cmdCurrentslice)
	}
}

// 获取单条命令当前所有参数的信息，包括从linux和windows的json中分别获取
func currentCmdParamenters(filepath, commandline string) string {
	// 存放对单条命令查询时当前参数的最终输出结果
	var currentslice []cmdMatch
	cmdLine := commandline
	var result string
	for _, item := range getPathContent(filepath) {
		getJsonOnce(item)
		var context = strings.Fields(cmdLine)
		if ContainsI(item, "linux_cmd.json") {
			// 去掉带有路径命令行的路径
			cmdLinuxPrefix(context)
			getCurrentSlice(context, cmdLine, &currentslice)
		} else if ContainsI(item, "Windows_cmd.json") {
			// 去掉带有路径命令行的路径
			cmdWindowsPrefix(context)
			getWindowsCurrentSlice(context, cmdLine, &currentslice)
		}
	}
	res, err := json.Marshal(currentslice)
	if err != nil {
		return "JSON ERR:" + err.Error()
	}
	if currentslice != nil {
		result = string(res[1 : len(res)-1])
	}
	return result
}

// 获取多条命令当前所有参数的信息，包括从linux和windows的json中分别获取
func currentParamenters(filepath, commandline string) string {
	cmdLine := commandline
	var cmdCurrentslice cmdMatchslice
	var result string
	// 将传入的命令行拆分为一个或多个连续空格的每个实例
	for _, item := range getPathContent(filepath) {
		getJsonOnce(item)
		if ContainsI(item, "linux_cmd.json") {
			getCmdCurrentSlice(cmdLine, &cmdCurrentslice)
		} else if ContainsI(item, "Windows_cmd.json") {
			getWindowsCmdCurrentSlice(cmdLine, &cmdCurrentslice)
		}
	}
	res, err := json.Marshal(cmdCurrentslice)
	if err != nil {
		return "JSON ERR:" + err.Error()
	}
	if cmdCurrentslice.CmdMatchSlice != nil {
		result = string(res)
	}
	return result

}

// 获取多条命令所有参数的信息，包括从linux和windows的json中分别获取
func allParamenters(filepath, commandline string) string {
	cmdLine := commandline
	var cmdMatchslice cmdMatchslice
	var result string
	for _, item := range getPathContent(filepath) {
		getJsonOnce(item)
		if ContainsI(item, "linux_cmd.json"){
			getCmdMathchSlice(cmdLine, &cmdMatchslice)
		}else if ContainsI(item, "Windows_cmd.json"){
			getWindowsCmdMathchSlice(cmdLine, &cmdMatchslice)
		}
	}
	res, err := json.Marshal(cmdMatchslice)
	if err != nil {
		return "JSON ERR:" + err.Error()
	}
	if cmdMatchslice.CmdMatchSlice != nil {
		result = string(res)
	}
	return result
}

// 获取单条命令所有参数的信息，包括从linux和windows的json中分别获取
func cmdAllParamenters(filepath, commandline string) string {
	// 存放对单条命令查询时所有参数的最终输出结果
	var allslice []cmdMatch
	cmdLine := commandline
	var result string
	for _, item := range getPathContent(filepath) {
		getJsonOnce(item)
		var context = strings.Fields(cmdLine)
		if ContainsI(item, "linux_cmd.json"){
			// 去掉带有路径命令行的路径
			cmdLinuxPrefix(context)
			getlLinuxAllSlice(cmdLine,context , &allslice)
		}else if ContainsI(item, "Windows_cmd.json"){
			// 去掉带有路径命令行的路径
			cmdWindowsPrefix(context)
			getlWindowsAllSlice(cmdLine, context, &allslice)
		}
	}
	res, err := json.Marshal(allslice)
	if err != nil {
		return "JSON ERR:" + err.Error()
	}
	if allslice != nil {
		result = string(res[1 : len(res)-1])
	}
	return result
}

// linux下多条命令的所有参数的查询
func linuxAllCommand(filepath, commandline string) string {

	cmdLine := commandline
	var cmdMatchslice cmdMatchslice
	var result string
	for _, item := range getPathContent(filepath) {
		getJsonOnce(item)

		if !ContainsI(item, "linux_cmd.json") {
			continue
		}
		getCmdMathchSlice(cmdLine, &cmdMatchslice)
		res, err := json.Marshal(cmdMatchslice)
		if err != nil {
			return "JSON ERR:" + err.Error()
		}

		if cmdMatchslice.CmdMatchSlice != nil {
			result = string(res)
		} else {
			result = ""
		}
	}
	return result
}

// linux下单条命令的所有参数的查询
func linuxCmdAllCommand(filepath, commandline string) string {
	// 存放对单条命令查询时所有参数的最终输出结果
	var allslice []cmdMatch
	var result string
	cmdLine := commandline
	for _, item := range getPathContent(filepath) {
		getJsonOnce(item)
		if !ContainsI(item, "linux_cmd.json") {
			continue
		}
		var context = strings.Fields(cmdLine)
		cmdLinuxPrefix(context)
		getlLinuxAllSlice(cmdLine, context, &allslice)
	}
	res, err := json.Marshal(allslice)
	if err != nil {
		return "JSON ERR:" + err.Error()
	}

	if allslice != nil {
		result = string(res[1 : len(res)-1])
	} else {
		result = ""
	}
	return result
}

// windows下多条命令所有参数的查询
func windowsAllCommand(filepath, commandline string) string {
	cmdLine := commandline
	var result string
	for _, item := range getPathContent(filepath) {
		getJsonOnce(item)
		if !ContainsI(item, "Windows_cmd.json") {
			continue
		}

		var cmdMatchslice cmdMatchslice
		getWindowsCmdMathchSlice(cmdLine, &cmdMatchslice)
		res, err := json.Marshal(cmdMatchslice)
		if err != nil {
			return "JSON ERR:" + err.Error()
		}

		if cmdMatchslice.CmdMatchSlice != nil {
			result = string(res)
		} else {
			result = ""
		}
	}
	return result
}

// windows下单条命令全部参数的查询
func windowsCmdAllCommand(filepath, commandline string) string {
	// 存放对单条命令查询时所有参数的最终输出结果
	var allslice []cmdMatch
	cmdLine := commandline
	var result string
	for _, item := range getPathContent(filepath) {
		getJsonOnce(item)
		if !ContainsI(item, "Windows_cmd.json") {
			continue
		}

		var context = strings.Fields(cmdLine)
		// 去掉带有路径命令行的路径
		cmdWindowsPrefix(context)
		getlWindowsAllSlice(cmdLine,context, &allslice)
	}
	res, err := json.Marshal(allslice)
	if err != nil {
		return "JSON ERR:" + err.Error()
	}

	if allslice != nil {
		return string(res[1 : len(res)-1])
	} else {
		result = ""
	}
	return result
}

// linux下多条命令当前参数的查询
func linuxCurrentCommand(filepath, commandline string) string {
	var cmdCurrentslice cmdMatchslice
	cmdLine := commandline
	var result string
	for _, item := range getPathContent(filepath) {
		getJsonOnce(item)
		if !ContainsI(item, "linux_cmd.json") {
			continue
		}

		getCmdCurrentSlice(cmdLine, &cmdCurrentslice)
		res, err := json.Marshal(cmdCurrentslice)
		if err != nil {
			return "JSON ERR:" + err.Error()
		}

		if cmdCurrentslice.CmdMatchSlice != nil {
			result = string(res)
		} else {
			result = ""
		}
	}
	return result
}

// linux下单条命令当前参数的查询
func linuxCurrentCmdCommand(filepath, commandline string) string {
	// 存放对单条命令查询时当前参数的最终输出结果
	var currentslice []cmdMatch
	cmdLine := commandline
	var result string
	for _, item := range getPathContent(filepath) {
		getJsonOnce(item)
		if !ContainsI(item, "linux_cmd.json") {
			continue
		}

		var context = strings.Fields(cmdLine)
		// 去掉带有路径命令行的路径
		cmdLinuxPrefix(context)
		getCurrentSlice(context, cmdLine, &currentslice)
		res, err := json.Marshal(currentslice)
		if err != nil {
			return "JSON ERR:" + err.Error()
		}
		if currentslice != nil {
			result = string(res[1 : len(res)-1])
		} else {
			result = ""
		}
	}
	return result
}

// windows下多条命令当前参数的查询
func windowsCurrentCommand(filepath, commandline string) string {
	var cmdCurrentslice cmdMatchslice
	cmdLine := commandline
	var result string
	for _, item := range getPathContent(filepath) {
		// 获取Windows_cmd.json的内容
		getJsonOnce(item)
		if !ContainsI(item, "Windows_cmd.json") {
			continue
		}
		getWindowsCmdCurrentSlice(cmdLine, &cmdCurrentslice)
		res, err := json.Marshal(cmdCurrentslice)
		if err != nil {
			return "JSON ERR:" + err.Error()
		}

		if cmdCurrentslice.CmdMatchSlice != nil {
			result = string(res)
		} else {
			result = ""
		}
	}
	return result
}

// windows 下单条命令当前参数的查询
func windowsCurrentCmdCommand(filepath, commandline string) string {
	// 存放对单条命令查询时当前参数的最终输出结果
	var currentslice []cmdMatch
	cmdLine := commandline
	var result string
	for _, item := range getPathContent(filepath) {
		getJsonOnce(item)
		// 获取Windows_cmd.json的内容
		if !ContainsI(item, "Windows_cmd.json") {
			continue
		}
		var context = strings.Fields(cmdLine)
		// 去掉带有路径命令行的路径
		cmdWindowsPrefix(context)
		getWindowsCurrentSlice(context, cmdLine, &currentslice)
		res, err := json.Marshal(currentslice)
		if err != nil {
			return "JSON ERR:" + err.Error()
		}
		if currentslice != nil {
			result = string(res[1 : len(res)-1])
		} else {
			result = ""
		}
	}
	return result
}
