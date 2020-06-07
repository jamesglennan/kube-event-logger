package main

import (
	"encoding/json"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func signalHandler() (stopCh <-chan struct{}) {
	stop := make(chan struct{})
	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		close(stop)
		log.Info("Shutting Down")
		<-sigChan
		os.Exit(1) // second signal, exit immediately.
	}()

	return stop
}

func getKubeConfig(path string) *rest.Config {
	log.Debugf("Kubeconfig path: %s:", path)
	clientConfigLoadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	clientConfigLoadingRules.DefaultClientConfig = &clientcmd.DefaultClientConfig
	clientConfigLoadingRules.ExplicitPath = path
	clientConfig := clientcmd.NewInteractiveDeferredLoadingClientConfig(clientConfigLoadingRules, &clientcmd.ConfigOverrides{}, os.Stdin)
	cfg, err := clientConfig.ClientConfig()

	if err != nil {
		log.Panic(err)
	}
	return cfg
}

func printEvent(obj interface{}) {
	event := obj.(*corev1.Event)

	jsByteArray, err := json.Marshal(*event)
	if err != nil {
		log.Error("Error marshalling event to json.")
		return
	}

	log.Info(string(jsByteArray))

	return
}

func main() {
	stopCh := signalHandler()
	// add flags
	var kubeConfigPath string
	var debug bool
	flag.StringVar(&kubeConfigPath, "kubeconfig", "", "Path to kubeconfig, defaults to use .kube in home direcotry or in-cluster config if run in a container")
	flag.BoolVar(&debug, "debug", false, "Set debug logs")
	flag.Parse()

	log.SetFormatter(&log.JSONFormatter{})
	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	cfg := getKubeConfig(kubeConfigPath)
	kubeClient := kubernetes.NewForConfigOrDie(cfg)

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Minute)
	eventInformer := kubeInformerFactory.Core().V1().Events().Informer()

	eventInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: printEvent,
	})
	eventInformer.Run(stopCh)

}
