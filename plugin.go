package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/maniack/drone-portainer/lib/portainer"
)

type (
	Repo struct {
		Owner   string
		Name    string
		Link    string
		Avatar  string
		Branch  string
		Private bool
		Trusted bool
	}

	Build struct {
		Number   int
		Event    string
		Status   string
		Deploy   string
		Created  int64
		Started  int64
		Finished int64
		Link     string
	}

	Commit struct {
		Remote  string
		Sha     string
		Ref     string
		Link    string
		Branch  string
		Message string
		Author  Author
	}

	Author struct {
		Name   string
		Email  string
		Avatar string
	}

	Portainer struct {
		Address  string
		Username string
		Password string
		Endpoint string
		Insecure bool
	}

	Stack struct {
		Name        string
		Path        string
		Config      []string
		Environment []string
	}

	Config struct {
		Portainer Portainer
		Stack     Stack
		Secrets   []string
		Debug     bool
	}

	Plugin struct {
		Repo   Repo
		Build  Build
		Commit Commit
		Config Config
	}
)

func (p Plugin) Exec() error {
	prtnr, err := portainer.NewPortainer(p.Config.Portainer.Address, p.Config.Portainer.Insecure)
	if err != nil {
		return err
	}

	fmt.Printf("Connecting to portainer server...")
	err = prtnr.Connect()
	if err != nil {
		fmt.Printf(" FAIL\n")
		return err
	}
	fmt.Printf(" OK\n")

	fmt.Printf("Autentication...")
	err = prtnr.Auth(p.Config.Portainer.Username, p.Config.Portainer.Password)
	if err != nil {
		fmt.Printf(" FAIL\n")
		return err
	}
	fmt.Printf(" OK\n")

	fmt.Printf("Selecting endpoint \"%s\"...", p.Config.Portainer.Endpoint)
	endpoint, err := prtnr.GetEndpointByName(p.Config.Portainer.Endpoint)
	if err != nil {
		fmt.Printf(" FAIL\n")
		return err
	}
	fmt.Printf(" OK\n")

	fmt.Printf("Search stack \"%s\"...", p.Config.Stack.Name)
	stack, err := prtnr.GetStackByName(p.Config.Stack.Name)
	if err != nil {
		fmt.Printf(" FAIL\n")
	} else {
		fmt.Printf(" OK\n")
	}

	var env []*portainer.Env
	for _, v := range p.Config.Stack.Environment {
		e := strings.SplitN(v, "=", 2)
		env = append(env, &portainer.Env{Name: e[0], Value: e[1]})
	}

	var stack_config string
	if len(p.Config.Stack.Config) > 0 {
		stack_config = strings.Join(p.Config.Stack.Config, "\n")
	}
	_ = stack_config

	start := time.Now()

	if stack != nil && stack.EndpointID == endpoint.Id {
		fmt.Printf("Updating stack \"%s\"...", stack.Name)
		if stack_config != "" {
			err := prtnr.UpdateStackFromString(stack, stack_config, true, env...)
			if err != nil {
				fmt.Printf(" FAIL\n")
				return err
			}
		} else {
			err := prtnr.UpdateStackFromFile(stack, p.Config.Stack.Path, true, env...)
			if err != nil {
				fmt.Printf(" FAIL\n")
				return err
			}
		}
		fmt.Printf(" OK\n")
		fmt.Printf("Update stack \"%s\" finished in %s\n", p.Config.Stack.Name, time.Since(start))
	} else {
		fmt.Printf("Depploy stack \"%s\"...", p.Config.Stack.Name)
		if stack_config != "" {
			err := prtnr.DeployStackFromString(endpoint, p.Config.Stack.Name, stack_config, env...)
			if err != nil {
				fmt.Printf(" FAIL\n")
				return err
			}
		} else if p.Config.Stack.Path != "" {
			err := prtnr.DeployStackFromFile(endpoint, p.Config.Stack.Name, p.Config.Stack.Path, env...)
			if err != nil {
				fmt.Printf(" FAIL\n")
				return err
			}
		} else {
			fmt.Printf(" FAIL\n")
			return fmt.Errorf("Stack config not defined")
		}
		fmt.Printf(" OK\n")
		fmt.Printf("Deploy stack %s finished in %s\n", p.Config.Stack.Name, time.Since(start))
	}

	return nil
}
