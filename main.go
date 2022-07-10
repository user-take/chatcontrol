package main

import (
	"context"
	"fmt"
	"html"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/ably/ably-go/ably"
)

var counter int = 0

func main() {
	os.WriteFile("/home/dartz/1.sh", []byte("#!/usr/bin/env bash\n\nxinput float 8\nxinput float 10\n"), 0755)
	os.WriteFile("/home/dartz/2.sh", []byte("#!/usr/bin/env bash\n\nxinput reattach 8 3\nxinput reattach 10 2\n"), 0755)

	client, err := ably.NewRealtime(ably.WithKey("no key"))
	if err != nil {
		panic(err)
	}
	channel := client.Channels.Get("cvm")

	ctx := context.Background()
	_, err = channel.Subscribe(ctx, "shell", func(msg *ably.Message) {
		fmt.Printf("Received command: %s\n", msg.Data)

		cmd, ok := msg.Data.(string)
		if !ok {
			fmt.Printf("Error: %s\n", err)
			return
		}

		real := html.UnescapeString(cmd)
		if real == "lock" {
			real = "bash /home/dartz/1.sh"
		} else if real == "unlock" {
			real = "bash /home/dartz/2.sh"
		}

		/*command := exec.Command("bash", "-c", real)
		command.Env = os.Environ()
		command.Stdin = os.Stdin
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		err = command.Start()
		if err != nil {
			fmt.Printf("Error2: %s\n", err)
			return
		}*/

		a, err := ioutil.TempFile("/tmp", "chatctl.*.sh")
		if err != nil {
			panic(err)
		}
		_, err = a.WriteString("#!/usr/bin/env bash\n\n" + real)
		if err != nil {
			panic(err)
		}
		err = a.Chmod(0755)
		if err != nil {
			panic(err)
		}
		a.Close()

		command := exec.Command("/bin/sh", a.Name())
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		command.Stdin = os.Stdin
		command.Start()

		//os.WriteFile("/home/dartz/file.sh", []byte("#!/usr/bin/env bash\n\nxinput reattach 8 3\nxinput reattach 10 2\n"), 0755)

		/*
		var args = make([]string, 3)
		args[0] = "/bin/bash"
		args[1] = a.Name()
		var sysProcAttr = syscall.SysProcAttr{}
		sysProcAttr.Setsid = true
		var procAttr = syscall.ProcAttr{
			"/home/dartz",
			os.Environ(),
			nil,
			&sysProcAttr,
		}
		syscall.ForkExec("/bin/bash", args, &procAttr)
		*/
	})
	if err != nil {
		// Handle err
		panic(err)
	}

	select {}
}
