package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	_ "github.com/joho/godotenv/autoload"
)

var version string // build number set at compile-time

func main() {
	app := cli.NewApp()
	app.Name = "drone-portainer"
	app.Usage = "drone-portainer usage"
	app.Action = run
	app.Version = version
	app.Flags = []cli.Flag{

		//
		// repo args
		//

		cli.StringFlag{
			Name:   "repo.fullname",
			Usage:  "repository full name",
			EnvVar: "DRONE_REPO",
		},
		cli.StringFlag{
			Name:   "repo.owner",
			Usage:  "repository owner",
			EnvVar: "DRONE_REPO_OWNER",
		},
		cli.StringFlag{
			Name:   "repo.name",
			Usage:  "repository name",
			EnvVar: "DRONE_REPO_NAME",
		},
		cli.StringFlag{
			Name:   "repo.link",
			Usage:  "repository link",
			EnvVar: "DRONE_REPO_LINK",
		},
		cli.StringFlag{
			Name:   "repo.avatar",
			Usage:  "repository avatar",
			EnvVar: "DRONE_REPO_AVATAR",
		},
		cli.StringFlag{
			Name:   "repo.branch",
			Usage:  "repository default branch",
			EnvVar: "DRONE_REPO_BRANCH",
		},
		cli.BoolFlag{
			Name:   "repo.private",
			Usage:  "repository is private",
			EnvVar: "DRONE_REPO_PRIVATE",
		},
		cli.BoolFlag{
			Name:   "repo.trusted",
			Usage:  "repository is trusted",
			EnvVar: "DRONE_REPO_TRUSTED",
		},

		//
		// commit args
		//

		cli.StringFlag{
			Name:   "remote.url",
			Usage:  "git remote url",
			EnvVar: "DRONE_REMOTE_URL",
		},
		cli.StringFlag{
			Name:   "commit.sha",
			Usage:  "git commit sha",
			EnvVar: "DRONE_COMMIT_SHA",
		},
		cli.StringFlag{
			Name:   "commit.ref",
			Value:  "refs/heads/master",
			Usage:  "git commit ref",
			EnvVar: "DRONE_COMMIT_REF",
		},
		cli.StringFlag{
			Name:   "commit.branch",
			Value:  "master",
			Usage:  "git commit branch",
			EnvVar: "DRONE_COMMIT_BRANCH",
		},
		cli.StringFlag{
			Name:   "commit.message",
			Usage:  "git commit message",
			EnvVar: "DRONE_COMMIT_MESSAGE",
		},
		cli.StringFlag{
			Name:   "commit.link",
			Usage:  "git commit link",
			EnvVar: "DRONE_COMMIT_LINK",
		},
		cli.StringFlag{
			Name:   "commit.author.name",
			Usage:  "git author name",
			EnvVar: "DRONE_COMMIT_AUTHOR",
		},
		cli.StringFlag{
			Name:   "commit.author.email",
			Usage:  "git author email",
			EnvVar: "DRONE_COMMIT_AUTHOR_EMAIL",
		},
		cli.StringFlag{
			Name:   "commit.author.avatar",
			Usage:  "git author avatar",
			EnvVar: "DRONE_COMMIT_AUTHOR_AVATAR",
		},

		//
		// build args
		//

		cli.StringFlag{
			Name:   "build.event",
			Value:  "push",
			Usage:  "build event",
			EnvVar: "DRONE_BUILD_EVENT",
		},
		cli.IntFlag{
			Name:   "build.number",
			Usage:  "build number",
			EnvVar: "DRONE_BUILD_NUMBER",
		},
		cli.IntFlag{
			Name:   "build.created",
			Usage:  "build created",
			EnvVar: "DRONE_BUILD_CREATED",
		},
		cli.IntFlag{
			Name:   "build.started",
			Usage:  "build started",
			EnvVar: "DRONE_BUILD_STARTED",
		},
		cli.IntFlag{
			Name:   "build.finished",
			Usage:  "build finished",
			EnvVar: "DRONE_BUILD_FINISHED",
		},
		cli.StringFlag{
			Name:   "build.status",
			Usage:  "build status",
			Value:  "success",
			EnvVar: "DRONE_BUILD_STATUS",
		},
		cli.StringFlag{
			Name:   "build.link",
			Usage:  "build link",
			EnvVar: "DRONE_BUILD_LINK",
		},
		cli.StringFlag{
			Name:   "build.deploy",
			Usage:  "build deployment target",
			EnvVar: "DRONE_DEPLOY_TO",
		},
		cli.BoolFlag{
			Name:   "yaml.verified",
			Usage:  "build yaml is verified",
			EnvVar: "DRONE_YAML_VERIFIED",
		},
		cli.BoolFlag{
			Name:   "yaml.signed",
			Usage:  "build yaml is signed",
			EnvVar: "DRONE_YAML_SIGNED",
		},

		//
		// prev build args
		//

		cli.IntFlag{
			Name:   "prev.build.number",
			Usage:  "previous build number",
			EnvVar: "DRONE_PREV_BUILD_NUMBER",
		},
		cli.StringFlag{
			Name:   "prev.build.status",
			Usage:  "previous build status",
			EnvVar: "DRONE_PREV_BUILD_STATUS",
		},
		cli.StringFlag{
			Name:   "prev.commit.sha",
			Usage:  "previous build sha",
			EnvVar: "DRONE_PREV_COMMIT_SHA",
		},

		//
		// plugin-specific parameters
		//

		cli.BoolFlag{
			Name:   "debug",
			Usage:  "debug mode",
			EnvVar: "PLUGIN_DEBUG",
		},
		cli.StringSliceFlag{
			Name:   "secrets",
			Usage:  "plugin secret",
			EnvVar: "PLUGIN_SECRETS",
		},
		cli.StringFlag{
			Name:   "portainer.address",
			Usage:  "portainer server address",
			EnvVar: "PLUGIN_PORTAINER_ADDRESS,PLUGIN_PORTAINER,PLUGIN_ADDRESS,PORTAINER_ADDRESS",
		},
		cli.BoolFlag{
			Name:   "portainer.insecure",
			Usage:  "portainer insecure connection",
			EnvVar: "PLUGIN_PORTAINER_INSECURE,PLUGIN_INSECURE,PORTAINER_INSECURE",
		},
		cli.StringFlag{
			Name:   "portainer.endpoint",
			Usage:  "portainer endpoint name",
			EnvVar: "PLUGIN_PORTAINER_ENDPOINT,PLUGIN_ENDPOINT,PORTAINER_ENDPOINT",
			Value:  "local",
		},
		cli.StringFlag{
			Name:   "stack.name",
			Usage:  "stack name",
			EnvVar: "PLUGIN_STACK_NAME,PLUGIN_STACK,STACK_NAME",
			Value:  "stack",
		},
		cli.StringFlag{
			Name:   "stack.file",
			Usage:  "stack file path",
			EnvVar: "PLUGIN_STACK_FILE,PLUGIN_FILE,STACK_FILE",
			Value:  "docker-compose.yml",
		},
		cli.StringSliceFlag{
			Name:   "stack.config",
			Usage:  "stack config",
			EnvVar: "PLUGIN_STACK_CONFIG,PLUGIN_CONFIG,STACK_CONFIG",
		},
		cli.StringSliceFlag{
			Name:   "stack.environment",
			Usage:  "stack environment",
			EnvVar: "PLUGIN_STACK_ENVIRONMENT,PLUGIN_STACK_ENV,PLUGIN_ENVIRONMENT,PLUGIN_ENV,STACK_ENVIRONMENT,STACK_ENV",
		},
		cli.StringFlag{
			Name:   "portainer.username",
			Usage:  "portainer server username",
			EnvVar: "PLUGIN_PORTAINER_USERNAME,PLUGIN_USERNAME,PORTAINER_USERNAME",
		},
		cli.StringFlag{
			Name:   "portainer.password",
			Usage:  "portainer server password",
			EnvVar: "PLUGIN_PORTAINER_PASSWORD,PLUGIN_PASSWORD,PORTAINER_PASSWORD",
		},
	}

	app.Run(os.Args)
}

func run(c *cli.Context) {

	plugin := Plugin{
		Repo: Repo{
			Owner:   c.String("repo.owner"),
			Name:    c.String("repo.name"),
			Link:    c.String("repo.link"),
			Avatar:  c.String("repo.avatar"),
			Branch:  c.String("repo.branch"),
			Private: c.Bool("repo.private"),
			Trusted: c.Bool("repo.trusted"),
		},
		Build: Build{
			Number:   c.Int("build.number"),
			Event:    c.String("build.event"),
			Status:   c.String("build.status"),
			Deploy:   c.String("build.deploy"),
			Created:  int64(c.Int("build.created")),
			Started:  int64(c.Int("build.started")),
			Finished: int64(c.Int("build.finished")),
			Link:     c.String("build.link"),
		},
		Commit: Commit{
			Remote:  c.String("remote.url"),
			Sha:     c.String("commit.sha"),
			Ref:     c.String("commit.sha"),
			Link:    c.String("commit.link"),
			Branch:  c.String("commit.branch"),
			Message: c.String("commit.message"),
			Author: Author{
				Name:   c.String("commit.author.name"),
				Email:  c.String("commit.author.email"),
				Avatar: c.String("commit.author.avatar"),
			},
		},
		Config: Config{
			Portainer: Portainer{
				Address:  c.String("portainer.address"),
				Username: c.String("portainer.username"),
				Password: c.String("portainer.password"),
				Endpoint: c.String("portainer.endpoint"),
				Insecure: c.Bool("portainer.insecure"),
			},
			Stack: Stack{
				Name:        c.String("stack.name"),
				Path:        c.String("stack.file"),
				Config:      c.StringSlice("stack.config"),
				Environment: c.StringSlice("stack.environment"),
			},
			Secrets: c.StringSlice("secrets"),
			Debug:   c.Bool("debug"),
		},
	}

	if err := plugin.Exec(); err != nil {
		fmt.Printf("Exited with error: %v\n", err)
		os.Exit(1)
	}
}
