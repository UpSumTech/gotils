package utils

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
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
	token := os.Getenv("DEPLOY_GITHUB_TOKEN")
	if len(token) == 0 {
		CheckErr("Could not find GITHUB_TOKEN exported in the shell")
	}
	return token
}

func GetGithubUser() string {
	token := os.Getenv("GITHUB_USERNAME")
	if len(token) == 0 {
		CheckErr("Could not find GITHUB_USERNAME exported in the shell")
	}
	return token
}

func GetDockerhubUser() string {
	token := os.Getenv("DOCKERHUB_USERNAME")
	if len(token) == 0 {
		CheckErr("Could not find DOCKERHUB_USERNAME exported in the shell")
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
