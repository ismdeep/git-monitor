package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func LoadGitPathList(path string) (gitPathList []string) {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer f.Close()

	br := bufio.NewReader(f)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}

		gitPath := strings.TrimSpace(string(a))
		if gitPath != "" {
			gitPathList = append(gitPathList, gitPath)
		}
	}

	return
}

func CheckGitChange(gitPath string) int {
	if err := os.Chdir(gitPath); err != nil {
		return 0
	}
	cmd := exec.Command("git", "status")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return 0
	}
	defer stdout.Close()

	if err := cmd.Start(); err != nil { // 运行命令
		return 0
	}

	opBytes, err := ioutil.ReadAll(stdout)

	if err != nil { // 读取输出结果
		return 0
	}

	content := string(opBytes)

	if strings.Contains(content, "nothing to commit, working tree clean") {
		return 0
	}

	fmt.Printf("[%v] CHANGED\n", gitPath)
	return 1
}

func main() {
	helpMsg := "Usage: git-monitor <git-path-list-file>"

	if len(os.Args) <= 1 {
		fmt.Println(helpMsg)
		return
	}

	gitPathList := LoadGitPathList(os.Args[1])

	cnt := 0

	for _, gitPath := range gitPathList {
		cnt += CheckGitChange(gitPath)
	}

	if cnt != 0 {
		fmt.Println("--------------------------------------------------")
	} else {
		fmt.Println("NO CHANGE FOUND.")
	}
}
