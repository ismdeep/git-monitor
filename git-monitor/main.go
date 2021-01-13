package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func LoadGitPathList(path string) (gitPathList []string) {
	f, err := os.Open(path)
	if err != nil {
		log.Println(err.Error())
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
	cmd := exec.Command("git", "ls-files", "--exclude-standard", "--others")
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

	if string(opBytes) != "" {
		log.Printf("[%v] CHANGED", gitPath)
		return 1
	}

	return 0
}

func main() {
	helpMsg := "Usage: git-monitor <git-path-list-file>"

	if len(os.Args) <= 1 {
		log.Fatalln(helpMsg)
	}

	gitPathList := LoadGitPathList(os.Args[1])

	cnt := 0

	for _, gitPath := range gitPathList {
		cnt += CheckGitChange(gitPath)
	}

	if cnt != 0 {
		log.Println("--------------------------------------------------")
	} else {
		log.Println("NO CHANGE FOUND.")
	}
}
