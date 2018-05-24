package main

import (
	"fmt"
	proto "github.com/chremoas/chremoas/proto"
	discord "github.com/chremoas/discord-gateway/proto"
	permsrv "github.com/chremoas/perms-srv/proto"
	"github.com/chremoas/services-common/config"
	"github.com/chremoas/purge-cmd/command"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
)

var Version = "SET ME YOU KNOB"
var service micro.Service
var name = "purge"

func main() {
	service = config.NewService(Version, "cmd", name, initialize)

	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}

// This function is a callback from the config.NewService function.  Read those docs
func initialize(config *config.Configuration) error {
	clientFactory := clientFactory{
		discordGateway: config.LookupService("gateway", "discord"),
		permsSrv:       config.LookupService("srv", "perms"),
		client:         service.Client()}

	proto.RegisterCommandHandler(service.Server(),
		command.NewCommand(name,
			&clientFactory,
		),
	)

	return nil
}

type clientFactory struct {
	discordGateway string
	permsSrv       string
	client         client.Client
}

func (c clientFactory) NewDiscordGateway() discord.DiscordGatewayService {
	return discord.NewDiscordGatewayService(c.discordGateway, c.client)
}

func (c clientFactory) NewPermsClient() permsrv.PermissionsService {
	return permsrv.NewPermissionsService(c.permsSrv, c.client)
}
