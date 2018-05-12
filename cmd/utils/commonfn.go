package utils

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

//////////////////// Vars and Consts ///////////////////

const (
	GITHUB_USERNAME_ENV_VAR    = "GITHUB_USERNAME"
	GITHUB_TOKEN_ENV_VAR       = "DEPLOY_GITHUB_TOKEN"
	BINTRAY_USERNAME_ENV_VAR   = "BINTRAY_USERNAME"
	BINTRAY_TOKEN_ENV_VAR      = "BINTRAY_API_KEY"
	BINTRAY_REPO_NAME_ENV_VAR  = "BINTRAY_REPO_NAME"
	DOCKERHUB_USERNAME_ENV_VAR = "DOCKERHUB_USERNAME"
)

////////////////////// Exported fns /////////////////////

func InitConfig(cfgFile string) {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			CheckErr(err.Error())
		}
		viper.SetConfigType("yaml")
		viper.SetConfigFile(filepath.Join(home, ".gotils.yml"))
	}
	if err := viper.ReadInConfig(); err != nil {
		CheckErr(err.Error())
	}
}

func GetDockerConfig(cfgFile string) []byte {
	if _, err := os.Stat(cfgFile); err != nil {
		CheckErr(err.Error())
	}
	b, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		CheckErr(err.Error())
	}
	return b
}

func GetGithubToken() string {
	token := os.Getenv(GITHUB_TOKEN_ENV_VAR)
	if len(token) == 0 {
		CheckErr(fmt.Sprintf("Could not find %s exported in the shell", GITHUB_TOKEN_ENV_VAR))
	}
	return token
}

func GetGithubUser() string {
	token := os.Getenv(GITHUB_USERNAME_ENV_VAR)
	if len(token) == 0 {
		CheckErr(fmt.Sprintf("Could not find %s exported in the shell", GITHUB_USERNAME_ENV_VAR))
	}
	return token
}

func GetBintrayToken() string {
	token := os.Getenv(BINTRAY_TOKEN_ENV_VAR)
	if len(token) == 0 {
		CheckErr(fmt.Sprintf("Could not find %s exported in the shell", BINTRAY_TOKEN_ENV_VAR))
	}
	return token
}

func GetBintrayUser() string {
	token := os.Getenv(BINTRAY_USERNAME_ENV_VAR)
	if len(token) == 0 {
		CheckErr(fmt.Sprintf("Could not find %s exported in the shell", BINTRAY_USERNAME_ENV_VAR))
	}
	return token
}

func GetBintrayRepo() string {
	token := os.Getenv(BINTRAY_REPO_NAME_ENV_VAR)
	if len(token) == 0 {
		CheckErr(fmt.Sprintf("Could not find %s exported in the shell", BINTRAY_REPO_NAME_ENV_VAR))
	}
	return token
}

func GetDockerhubUser() string {
	token := os.Getenv(DOCKERHUB_USERNAME_ENV_VAR)
	if len(token) == 0 {
		CheckErr(fmt.Sprintf("Could not find %s exported in the shell", DOCKERHUB_USERNAME_ENV_VAR))
	}
	return token
}

func GetK8sClientSet() *kubernetes.Clientset {
	home, err := homedir.Dir()
	if err != nil {
		CheckErr(err.Error())
	}
	kubeconfig := flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "absolute path to the kubeconfig file")
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		CheckErr(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		CheckErr(err.Error())
	}

	return clientset
}
