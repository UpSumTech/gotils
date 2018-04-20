package utils

import (
	"flag"
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
