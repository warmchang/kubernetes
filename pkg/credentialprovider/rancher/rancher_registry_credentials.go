package rancher_credentials

import (
	"github.com/golang/glog"
	"github.com/rancher/go-rancher/client"
	"k8s.io/kubernetes/pkg/credentialprovider"
	"os"
	"time"
)

// rancher provider fetching credentials from rancher server
type rancherProvider struct {
	client *client.RancherClient
}

type rConfig struct {
	Global configGlobal
}

type configGlobal struct {
	CattleURL       string `gcfg:"cattle-url"`
	CattleAccessKey string `gcfg:"cattle-access-key"`
	CattleSecretKey string `gcfg:"cattle-secret-key"`
}

type registryCredential struct {
	credential *client.RegistryCredential
	serverIP   string
}

func init() {
	client, err := getRancherClient()
	if err != nil {
		glog.Errorf("Failed to get rancher client: %v", err)
	}
	credentialprovider.RegisterCredentialProvider("rancher-registry-creds",
		&credentialprovider.CachingDockerConfigProvider{
			Provider: &rancherProvider{client},
			Lifetime: 30 * time.Minute,
		})
}

// Assuming its always enabled
func (p *rancherProvider) Enabled() bool {
	return p.client != nil
}

// Provide implements DockerConfigProvider.Provide, refreshing Rancher tokens on demand
func (p *rancherProvider) Provide() credentialprovider.DockerConfig {
	cfg := credentialprovider.DockerConfig{}
	for _, cred := range p.getRancherCredentials() {
		entry := credentialprovider.DockerConfigEntry{
			Username: cred.credential.PublicValue,
			Password: cred.credential.SecretValue,
			Email:    cred.credential.Email,
		}
		cfg[cred.serverIP] = entry
	}

	return cfg
}

func (p *rancherProvider) getRancherCredentials() []registryCredential {
	var registryCreds []registryCredential
	credColl, err := p.client.RegistryCredential.List(client.NewListOpts())
	if err != nil {
		glog.Errorf("Failed to pull registry credentials from rancher %v", err)
		return registryCreds
	}
	for _, cred := range credColl.Data {
		registry := &client.Registry{}
		if err = p.client.GetLink(cred.Resource, "registry", registry); err != nil {
			glog.Errorf("Failed to pull registry from rancher %v", err)
			return registryCreds
		}
		registryCred := registryCredential{
			credential: &cred,
			serverIP:   registry.ServerAddress,
		}
		registryCreds = append(registryCreds, registryCred)
	}
	return registryCreds
}

func getRancherClient() (*client.RancherClient, error) {
	url := os.Getenv("CATTLE_URL")
	accessKey := os.Getenv("CATTLE_ACCESS_KEY")
	secretKey := os.Getenv("CATTLE_SECRET_KEY")
	conf := rConfig{
		Global: configGlobal{
			CattleURL:       url,
			CattleAccessKey: accessKey,
			CattleSecretKey: secretKey,
		},
	}

	return client.NewRancherClient(&client.ClientOpts{
		Url:       conf.Global.CattleURL,
		AccessKey: conf.Global.CattleAccessKey,
		SecretKey: conf.Global.CattleSecretKey,
	})
}
