package main

import (
	"fmt"
	"go.i3wm.org/i3/v4"
	"os"
	"os/exec"
	"log"
)

func RunCommand(command string) {
	_, err := i3.RunCommand(command)
	if err != nil && !i3.IsUnsuccessful(err) {
		log.Fatalf("i3.RunCommand() failed with %s\n", err)
	}
}

func main() {
	if len(os.Args) == 1 {
		log.Fatalf("Arguments required\n")
	}
	tree, err := i3.GetTree()
	if err != nil {
		log.Fatalf("i3.GetTree failed with %s\n", err)
	}

	var nodes []*i3.Node
	focused := tree.Root.FindFocused(func(n *i3.Node) bool {
		nodes = append(nodes, n)
		return n.Focused
	})
	if focused == nil {
		log.Fatalf("Could not locate focused node\n")
	}

	focusedParent := nodes[len(nodes)-2]

	if focusedParent.Type == i3.FloatingCon {
		log.Fatalf("Floating windows not supported\n")
	}

	oldLayout := focusedParent.Layout
	if oldLayout == i3.SplitH || oldLayout == i3.SplitV {
		RunCommand("split v, layout stacking")
	} else {
		log.Fatalf("Layout %s not supported\n", oldLayout)
	}

	cmd := exec.Command(os.Args[1], os.Args[2:]...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}

	RunCommand(fmt.Sprintf("layout %s", oldLayout))
}
