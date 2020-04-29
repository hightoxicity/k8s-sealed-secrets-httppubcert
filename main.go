package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	api "k8s.io/kubernetes/pkg/apis/core"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeconfig      = flag.String("kubeconfig", "", "Absolute path to the kubeconfig file")
	sealedsecretsns = flag.String("sealedsecretsns", "kube-system", "sealed-secrets controller namespace")
	secretprefix    = flag.String("secretprefix", "sealed-secrets-key", "sealed-secrets secret prefix in the sealedsecretsns")
	verbose         = flag.Bool("verbose", false, "Turn on verbosity")
	listenaddress   = flag.String("listenaddress", ":8080", "Webserver address")
	certpath        = flag.String("certpath", "/cert", "Exposed cert path on webserver")

	config  *rest.Config
	crtMu   sync.RWMutex
	lastCrt string
)

func GetClientset() (cs *kubernetes.Clientset, retErr error) {

	var err error

	if *kubeconfig == "" {
		config, err = rest.InClusterConfig()
		if err != nil {
			retErr = errors.New(err.Error())
		}
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			retErr = errors.New(fmt.Sprintf("Error with config file %s: %s", *kubeconfig, err))
		}
	}

	cs, err = kubernetes.NewForConfig(config)
	if err != nil {
		retErr = errors.New(fmt.Sprintf("Bad config file: %s", err))
	}

	return cs, retErr
}

func main() {
	flag.Parse()
	clientset, err := GetClientset()

	if err != nil {
		log.Fatalf("%s", err)
	} else {
		if *verbose == true {
			log.Println("Clientset properly retrieved!")
		}
	}

	fieldSelector := fields.Set{api.SecretTypeField: string(v1.SecretTypeTLS)}.AsSelector()

	watchlist := cache.NewListWatchFromClient(
		clientset.CoreV1().RESTClient(),
		"secrets",
		*sealedsecretsns,
		fieldSelector,
	)

	_, controller := cache.NewInformer(
		watchlist,
		&v1.Secret{},
		0, //Duration is int64
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				if strings.HasPrefix(obj.(*v1.Secret).ObjectMeta.Name, *secretprefix) {
					secret := obj.(*v1.Secret)
					if *verbose {
						log.Printf("%s matches secret prefix\n", secret.ObjectMeta.Name)
					}
					for labName, labValue := range secret.ObjectMeta.Labels {
						if "sealedsecrets.bitnami.com/sealed-secrets-key" == labName && "active" == labValue {
							log.Printf("I got the fish: %s\n", secret.ObjectMeta.Name)
							crtMu.Lock()
							lastCrt = string(secret.Data["tls.crt"])
							crtMu.Unlock()
						}
					}
				}
			},
		},
	)
	stop := make(chan struct{})
	defer close(stop)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	http.HandleFunc(*certpath, func(w http.ResponseWriter, r *http.Request) {

		if *verbose {
			log.Printf("Received a request from %s", r.RemoteAddr)
		}

		var b bytes.Buffer
		crtMu.Lock()
		b.WriteString(lastCrt)
		crtMu.Unlock()
		w.Write(b.Bytes())
	})

	go controller.Run(stop)
	go http.ListenAndServe(*listenaddress, nil)
	<-sigs
}
