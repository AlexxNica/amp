package team

import (
	"errors"
	"fmt"
	"text/tabwriter"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/appcelerator/amp/data/accounts"
	"github.com/appcelerator/amp/pkg/time"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc/status"
)

type listTeamOptions struct {
	org   string
	quiet bool
}

// NewTeamListCommand returns a new instance of the team list command.
func NewTeamListCommand(c cli.Interface) *cobra.Command {
	opts := listTeamOptions{}
	cmd := &cobra.Command{
		Use:     "ls [OPTIONS]",
		Short:   "List teams",
		Aliases: []string{"list"},
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listTeam(c, cmd, opts)
		},
	}
	flags := cmd.Flags()
	//flags.StringVar(&opts.org, "org", "", "Organization name")
	flags.BoolVarP(&opts.quiet, "quiet", "q", false, "Only display team names")
	return cmd
}

func listTeam(c cli.Interface, cmd *cobra.Command, opts listTeamOptions) error {
	opts.org = accounts.DefaultOrganization
	//org, err := cli.ReadOrg(c.Server())
	//if !cmd.Flag("org").Changed {
	//	switch {
	//	case err == nil:
	//		opts.org = org
	//		c.Console().Println("organization name:", opts.org)
	//	default:
	//		opts.org = c.Console().GetInput("organization name")
	//	}
	//}
	conn := c.ClientConn()
	client := account.NewAccountClient(conn)
	request := &account.ListTeamsRequest{
		OrganizationName: opts.org,
	}
	reply, err := client.ListTeams(context.Background(), request)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			return errors.New(s.Message())
		}
	}
	if opts.quiet {
		for _, team := range reply.Teams {
			c.Console().Println(team.Name)
		}
		return nil
	}
	w := tabwriter.NewWriter(c.Out(), 0, 0, cli.Padding, ' ', 0)
	fmt.Fprintln(w, "TEAM\tCREATED ON")
	for _, team := range reply.Teams {
		fmt.Fprintf(w, "%s\t%s\n", team.Name, time.ConvertTime(team.CreateDt))
	}
	w.Flush()
	return nil
}
