package command

import (
	"fmt"
	proto "github.com/chremoas/chremoas/proto"
	discord "github.com/chremoas/discord-gateway/proto"
	permsrv "github.com/chremoas/perms-srv/proto"
	"github.com/chremoas/services-common/args"
	common "github.com/chremoas/services-common/command"
	"golang.org/x/net/context"
	"strings"
	"time"
)

type ClientFactory interface {
	NewDiscordGateway() discord.DiscordGatewayService
	NewPermsClient() permsrv.PermissionsService
}

var cmdName = "sig"
var clientFactory ClientFactory
var permissions *common.Permissions

type Command struct {
	//Store anything you need the Help or Exec functions to have access to here
	name    string
	factory ClientFactory
}

func (c *Command) Help(ctx context.Context, req *proto.HelpRequest, rsp *proto.HelpResponse) error {
	rsp.Usage = c.name
	rsp.Description = "Purge old bot related messages"
	return nil
}

func (c *Command) Exec(ctx context.Context, req *proto.ExecRequest, rsp *proto.ExecResponse) error {
	cmd := args.NewArg(cmdName)
	cmd.Add("full", &args.Command{fullPurge, "List all SIGs"})
	cmd.Add("start", &args.Command{startPurgeThread, "Add SIGs"})
	cmd.Add("stop", &args.Command{stopPurgeThread, "Delete SIGs"})
	cmd.Add("keep", &args.Command{keepLength, "Add user to SIG"})
	cmd.Add("frequency", &args.Command{threadFrequency, "Remove user from SIG"})
	// TODO: Add a command for the user to get a list of what SIGs they are members of.
	err := cmd.Exec(ctx, req, rsp)

	//I don't 100% love this, but it'll do for now. -brian
	if err != nil {
		rsp.Result = []byte(common.SendError(err.Error()))
	}
	return nil
}

func getMessages(channelID string) {
	longCtx, _ := context.WithTimeout(context.Background(), 60*time.Minute)
	fmt.Printf("Calling GetMessages: ChannelID=%s\n", channelID)
	discordClient := clientFactory.NewDiscordGateway()
	messages, err := discordClient.GetMessages(longCtx, &discord.GetMessagesRequest{
		ChannelID: channelID,
		Limit:     100,
	})
	if (err != nil) {
		fmt.Printf("Error: %s\n", err.Error())
	}
	fmt.Printf("messages: %+v\n", messages)
}

func fullPurge(ctx context.Context, req *proto.ExecRequest) string {
	if len(req.Args) != 2 {
		return common.SendError("Usage: !purge full")
	}

	canPerform, err := permissions.CanPerform(ctx, req.Sender)
	if err != nil {
		return common.SendFatal(err.Error())
	}

	if !canPerform {
		return common.SendError("User doesn't have permission to this command")
	}

	sender := strings.Split(req.Sender, ":")

	go getMessages(sender[0])
	return "stuff"
}

func startPurgeThread(ctx context.Context, req *proto.ExecRequest) string {
	if len(req.Args) != 2 {
		return common.SendError("Usage: !purge full")
	}

	canPerform, err := permissions.CanPerform(ctx, req.Sender)
	if err != nil {
		return common.SendFatal(err.Error())
	}

	if !canPerform {
		return common.SendError("User doesn't have permission to this command")
	}

	return "stuff"
}

func stopPurgeThread(ctx context.Context, req *proto.ExecRequest) string {
	if len(req.Args) != 2 {
		return common.SendError("Usage: !purge full")
	}

	canPerform, err := permissions.CanPerform(ctx, req.Sender)
	if err != nil {
		return common.SendFatal(err.Error())
	}

	if !canPerform {
		return common.SendError("User doesn't have permission to this command")
	}

	return "stuff"
}

func keepLength(ctx context.Context, req *proto.ExecRequest) string {
	if len(req.Args) != 2 {
		return common.SendError("Usage: !purge full")
	}

	canPerform, err := permissions.CanPerform(ctx, req.Sender)
	if err != nil {
		return common.SendFatal(err.Error())
	}

	if !canPerform {
		return common.SendError("User doesn't have permission to this command")
	}

	return "stuff"
}

func threadFrequency(ctx context.Context, req *proto.ExecRequest) string {
	if len(req.Args) != 2 {
		return common.SendError("Usage: !purge full")
	}

	canPerform, err := permissions.CanPerform(ctx, req.Sender)
	if err != nil {
		return common.SendFatal(err.Error())
	}

	if !canPerform {
		return common.SendError("User doesn't have permission to this command")
	}

	return "stuff"
}

func NewCommand(name string, factory ClientFactory) *Command {
	clientFactory = factory
	permissions = common.NewPermission(clientFactory.NewPermsClient(), []string{"purge_admins"})

	return &Command{name: name, factory: factory}
}
