package util

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"github.com/missdeer/golib/fsutil"

	"github.com/missdeer/hannah/config"
)

// ExternalPlay play song by external player
func ExternalPlay(song string) error {
	b, err := fsutil.FileExists(config.Player)
	if err == nil && b {
		cmd := exec.Command(config.Player, song)

		// create a pipe for the output of the script
		cmdReader, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
			return err
		}

		scanner := bufio.NewScanner(cmdReader)
		go func() {
			for scanner.Scan() {
				fmt.Println(scanner.Text())
			}
		}()

		err = cmd.Start()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
			return err
		}

		err = cmd.Wait()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
			return err
		}
		return nil
	}
	return err
}
