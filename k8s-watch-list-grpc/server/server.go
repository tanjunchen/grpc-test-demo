package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"grpc-test-demo/k8s-watch-list-grpc/proto"
	"grpc-test-demo/k8s-watch-list-grpc/service"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

// Kubeconfig returns kube config path.
func Kubeconfig() string {
	return filepath.Join(Home(), ".kube", "config")
}

// Home returns home path.
func Home() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

var K8SClient *kubernetes.Clientset

func init() {
	clusterConfig, err := clientcmd.BuildConfigFromFlags("", Kubeconfig())
	if err != nil {
		fmt.Println(err)
	}
	k8s, err := kubernetes.NewForConfig(clusterConfig)
	K8SClient = k8s
}

func main() {

	address := ":9999"

	listen, err := net.Listen("tcp", address)

	if err != nil {
		log.Printf("failed to listen: %v", err)
		os.Exit(1)
	}

	fmt.Println("listening: ", address)

	var opts []grpc.ServerOption

	//  gRPC server instances
	s := grpc.NewServer(opts...)

	k8s.RegisterServiceServiceServer(s, service.K8sServiceService{})

	reflection.Register(s)

	stopCh := signals.SetupSignalHandler()

	// 开启 watch event
	go watch(stopCh)
	// wait for SIGTERM or SIGINT

	// It is prepared for multiple instances of pod, for the watch event and LeaderElection mechanism.
	// TODO LeaderElection?

	// start HTTP server
	// run server in background
	go func() {
		if err := s.Serve(listen); err != http.ErrServerClosed {
			log.Printf("HTTP server crashed %v", err)
		}
	}()

	// wait for SIGTERM or SIGINT
	<-stopCh

	if err := s.Stop; err != nil {
		log.Printf("HTTP server graceful shutdown failed %v", err)
	} else {
		log.Printf("HTTP server stopped")
	}
}

func watch(stopCh <-chan struct{}) {
	options := func(options *v1.ListOptions) {
		options.LabelSelector = "test.io"
	}

	sharedOptions := []informers.SharedInformerOption{
		informers.WithNamespace(v1.NamespaceAll),
		informers.WithTweakListOptions(options),
	}
	// 同步周期是 30 秒
	informerFactory := informers.NewSharedInformerFactoryWithOptions(K8SClient, time.Second*30, sharedOptions...)

	serviceInformer := informerFactory.Core().V1().Services()
	informer := serviceInformer.Informer()

	defer runtime.HandleCrash()

	// 启动 informer，list & watch
	go informerFactory.Start(stopCh)

	// 从 apiserver 同步资源，即 list
	if !cache.WaitForCacheSync(stopCh, informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
		return
	}

	// 使用自定义 handler
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: add,
		// 此处省略 workqueue 的使用
		UpdateFunc: update,
		DeleteFunc: delete,
	})
	<-stopCh
}

func add(obj interface{}) {
	fmt.Println("========add=========")
	svc := obj.(*corev1.Service)
	serviceInfo := &k8s.SyncServiceInfo{
		Name:              svc.Name,
		ResourceVersion:   svc.ResourceVersion,
		Operation:         0,
		CreationTimeStamp: svc.CreationTimestamp.String(),
		Labels:            svc.Labels,
		Selector:          svc.Spec.Selector,
	}
	var syncServiceInfo []*k8s.SyncServiceInfo
	syncServiceInfo = append(syncServiceInfo, serviceInfo)
	response := k8s.SyncServiceResponse{}
	response.SyncServiceInfo = syncServiceInfo
	response.Namespace = svc.Namespace

	// rpc 推送出去 连接对方的 client 建立连接吗？
	service.MsgChan <- response

	fmt.Println(response)

	fmt.Println("response:", response)
	fmt.Println("========add=========")
}

func update(old interface{}, new interface{}) {
	fmt.Println("========update=========")
}

func delete(object interface{}) {
	fmt.Println("------delete-------")
	fmt.Println(object)
	fmt.Println("-------delete------")
}
