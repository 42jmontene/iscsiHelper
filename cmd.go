package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"strconv"
)

func execCmd (cmd string) (out string) {
//	fmt.Printf("%s\n", cmd)
	parts := strings.Fields(cmd)
	command := exec.Command(parts[0],parts[1:len(parts)]...)
//	command.Dir = ""
	output, err := command.Output()
	if err != nil {
		log.Fatalf("Cmd.Run() Failed with %s\n", err)
	}
//	fmt.Printf("Return value :  >%s<  !\n", string(output))
	return string(output)
}

func iscsiGetInfos (infos info, login string) (info) {
	var err error

	output := execCmd(fmt.Sprintf("bash ./scripts/which_home.sh %s", login))
	infos.storageNb, err = strconv.Atoi(output[0:1])
	if err != nil {
		log.Fatal(err)
	}
	return infos
}

