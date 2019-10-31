package controllers

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"
)

type CommandController struct {
	WorkDir string `json:"work_dir"`
	Branch  string `json:"branch"`
}

// 主函数 命令行接收参数 第一个参数为当前目录下的项目名 第二个参数为分支名
// 默认和所有项目同级目录 运行更新所有项目 如果从某个项目路径运行 则只更新本项目
func main() {
	var Git CommandController
	var argsLen int = len(os.Args)
	if 1 == argsLen {
		RunPath, _ := os.Getwd()
		FilePath := Git.getCurrentPath()
		RunPath = RunPath + "\\"
		if RunPath != FilePath {
			Git.WorkDir = RunPath
			Git.updateProject()
		} else {
			Git.AllProject()
		}
	} else if 2 == argsLen {
		Git.WorkDir = Git.getCurrentPath() + os.Args[1]
		Git.updateProject()
	} else if 3 == argsLen {
		Git.WorkDir = Git.getCurrentPath() + os.Args[1]
		Git.Branch = os.Args[2]
		Git.updateProject()
	}
}

// 获取分支列表
func (G *CommandController) GetBranchList(path string) []string {
	G.WorkDir = path
	allBranch := G.gitBranch([]string{"-a"})
	for i := 0; i < len(allBranch); i++ {
		allBranch[i] = strings.Replace(allBranch[i], "  ", "", 1)
		allBranch[i] = strings.Replace(allBranch[i], "remotes/origin/HEAD -> origin/", "", 1)
		allBranch[i] = strings.Replace(allBranch[i], "remotes/origin/", "", 1)
		allBranch[i] = strings.Replace(allBranch[i], "* ", "", 1)
	}
	return G.RemoveRep(allBranch)
}

// 排序去重
func (G *CommandController) RemoveRep(slc []string) []string {
	var temp []string
	// 排序
	sort.Strings(slc)
	for i := 0; i < len(slc)-1; i++ {
		flag := strings.Compare(slc[i], slc[i+1])
		if -1 == flag {
			temp = append(temp, slc[i])
		}
	}
	temp = append(temp, slc[len(slc)-1])
	return temp
}

func (G *CommandController) GetNowBranch(path string) string {
	branch := G.gitSymbolic([]string{"--short", "-q", "HEAD"})
	if branch == nil{
		return "branch error"
	}
	return branch[0]
}

func (G *CommandController) CheckoutBranch(path string, branch string) []string {
	G.WorkDir = path
	result := G.gitCheckout([]string{branch})
	return result
}

// 更新所有Git项目
func (G *CommandController) AllProject() {
	var rootPath string = G.getCurrentPath()
	var dir []string = G.getDir()
	for _, d := range dir {
		G.WorkDir = rootPath + d
		G.updateProject()
	}
}

// 更新Git项目
func (G *CommandController) updateProject() map[string][]string {
	// updateInfo[G.WorkDir][0] = "*****************************************************************\n"
	// updateInfo[G.WorkDir][1] = "                  " + G.WorkDir + "\n"
	// updateInfo[G.WorkDir][2] = "*****************************************************************\n"
	var NowBranch string
	if "" == G.Branch {
		NowBranch = G.IoutilBranch()
	} else {
		NowBranch = G.Branch
	}
	if "0" == NowBranch {
		message := make(map[string][]string, 1)
		message["message"] = []string{"This is not a GIT project\n"}
		return message
	}
	result := make(map[string][]string, 7)
	result["a_resetInfo"] = G.gitReset([]string{"--hard"})
	result["b_checkoutMaster"] = G.gitCheckout([]string{"master"})
	result["c_fetch"] = G.gitFetch([]string{})
	result["d_rebaseMaster"] = G.gitRebase([]string{})
	if "master" != NowBranch {
		result["e_clean"] = G.gitClean([]string{"-df"})
		result["f_checkoutNowBranch"] = G.gitCheckout([]string{NowBranch})
		result["g_rebaseNowBranch"] = G.gitRebase([]string{"origin/" + NowBranch})
	}
	return result
}

// 获取当前完整路径
func (G *CommandController) getCurrentPath() string {
	s, _ := exec.LookPath(os.Args[0])
	i := strings.LastIndex(s, "\\")
	path := string(s[0 : i+1])
	return path
}

// 获取当前目录下所有非隐藏文件夹
func (G *CommandController) getDir() []string {
	var dir []string
	files, _ := ioutil.ReadDir(G.getCurrentPath())
	for _, f := range files {
		if f.IsDir() {
			str := f.Name()
			matched, err := regexp.MatchString("^\\.\\S*", str)
			if err == nil && !matched {
				dir = append(dir, f.Name())
			}
		}
	}
	return dir
}

// 获取Git项目当前分支
func (G *CommandController) IoutilBranch() string {
	var name string = G.WorkDir + "/.git/HEAD"
	if contents, err := ioutil.ReadFile(name); err == nil {
		// 因为contents是[]byte类型，直接转换成string类型后会多一行空格,需要使用strings.Replace替换换行符
		result := strings.Replace(string(contents), "\n", "", 1)
		result = strings.Replace(result, "ref: refs/heads/", "", 1)
		return result
	}
	return "0"
}

// 重置项目
func (G *CommandController) gitReset(command []string) []string {
	commandName := "git"
	params := append([]string{"reset"}, command...)
	return G.execCommand(commandName, params)
}

// 获得最新代码
func (G *CommandController) gitFetch(command []string) []string {
	commandName := "git"
	params := append([]string{"fetch"}, command...)
	return G.execCommand(commandName, params)
}

// 更新本地分支到最新分支
func (G *CommandController) gitRebase(command []string) []string {
	commandName := "git"
	params := append([]string{"rebase"}, command...)
	return G.execCommand(commandName, params)
}

// 切换分支
func (G *CommandController) gitCheckout(command []string) []string {
	commandName := "git"
	params := append([]string{"checkout"}, command...)
	return G.execCommand(commandName, params)
}

// 清除多余文件
func (G *CommandController) gitClean(command []string) []string {
	commandName := "git"
	params := append([]string{"clean"}, command...)
	return G.execCommand(commandName, params)
}

// 操作分支
func (G *CommandController) gitBranch(command []string) []string {
	commandName := "git"
	params := append([]string{"branch"}, command...)
	return G.execCommand(commandName, params)
}

// 操作分支
func (G *CommandController) gitSymbolic(command []string) []string {
	commandName := "git"
	params := append([]string{"symbolic-ref"}, command...)
	return G.execCommand(commandName, params)
}

// CD
func (G *CommandController) cd(command []string) []string {
	commandName := "cd"
	return G.execCommand(commandName, command)
}

// 调用命令
func (G *CommandController) execCommand(commandName string, params []string) []string {
	var commandArray []string
	cmd := exec.Command(commandName, params...)
	cmd.Dir = G.WorkDir

	// 显示运行的命令
	// fmt.Println(cmd.Args)

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Println(err)
		return commandArray
	}
	cmd.Start()

	reader := bufio.NewReader(stdout)
	// 实时循环读取输出流中的一行内容
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 || line == "" {
			break
		}
		commandArray = append(commandArray, strings.Replace(line, "\n", "", -1))
	}
	return commandArray
}
