package task

import (
	"os"

	"github.com/dcos/dcos-cli/api"
	"github.com/dcos/dcos-core-cli/pkg/pluginutil"
	"github.com/dcos/dcos-core-cli/pkg/mesos"
	"github.com/mesos/mesos-go/api/v1/lib/httpcli/httpagent"
	"github.com/spf13/cobra"
)

func newCmdTaskExec(ctx api.Context) *cobra.Command {
	var interactive, tty bool
	var user string

	cmd := &cobra.Command{
		Use:   "exec <task> <cmd> [<args>...]",
		Short: "Launch a process (<cmd>) inside of a container for a task (<task>).",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			task, err := findTask(ctx, args[0])
			if err != nil {
				return err
			}

			httpClient, err := mesosHTTPClient(ctx, task.AgentID.Value)
			if err != nil {
				return err
			}

			containerID := task.Statuses[0].ContainerStatus.ContainerID

			taskIO, err := mesos.NewTaskIO(*containerID, mesos.TaskIOOpts{
				Stdin:       ctx.Input(),
				Stdout:      ctx.Out(),
				Stderr:      ctx.ErrOut(),
				Interactive: interactive,
				TTY:         tty,
				User:        user,
				Sender:      httpagent.NewSender(httpClient.Send),
				Logger:      pluginutil.Logger(),
			})

			if err != nil {
				return err
			}

			exitCode, err := taskIO.Exec(args[1], args[2:]...)
			if err != nil {
				return err
			}
			os.Exit(exitCode)
			return nil
		},
	}

	cmd.Flags().BoolVar(&interactive, "interactive", false, "Attach a STDIN stream to the remote command for an interactive session")
	cmd.Flags().BoolVar(&tty, "tty", false, "Print completed tasks")
	cmd.Flags().StringVar(&user, "user", "", "Run as the given user")
	return cmd
}
