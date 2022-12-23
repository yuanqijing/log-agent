package elector

import (
	"context"
	"github.com/go-logr/logr"
	"github.com/yuanqijing/log-agent/pkg/kube"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/klog/v2"
	"sync"

	"time"
)

type Elector struct {
	logger logr.Logger

	// leader default is false
	leader      bool
	leaderMutex sync.Mutex

	// leaseLockName is the lease lock resource name
	leaseLockName string
	// leaseLockNamespace is the lease lock resource namespace
	leaseLockNamespace string

	// kubeconfig path to kubeconfig file
	kubeconfig string

	// Pod ID
	podId string
}

func NewElector(cfg *Config) *Elector {
	err := cfg.Validate()

	if err != nil {
		panic(err)
	}

	podID := kube.GetPodID()

	e := &Elector{
		logger:             cfg.Logger,
		leaseLockName:      cfg.LeaseLockName,
		leaseLockNamespace: cfg.LeaseLockNamespace,
		kubeconfig:         cfg.KubeConfig,
		podId:              podID,
		leader:             false,
	}

	return e
}

func (e *Elector) GetLeader() bool {
	e.leaderMutex.Lock()
	defer e.leaderMutex.Unlock()
	return e.leader
}

func (e *Elector) SetLeader(leader bool) {
	e.leaderMutex.Lock()
	defer e.leaderMutex.Unlock()
	e.leader = leader
}

func (e *Elector) Run(stopCh <-chan struct{}) {
	e.logger.Info("Elector is running...")

	go e.run(stopCh)

	<-stopCh
	e.logger.Info("Elector is shutting down...")
}

func (e *Elector) run(stopCh <-chan struct{}) {
	if e.leaseLockName == "" {
		klog.Fatal("unable to get lease lock resource name (missing lease-lock-name flag).")
	}

	if e.leaseLockNamespace == "" {
		klog.Fatal("unable to get lease lock resource namespace (missing lease-lock-namespace flag).")
	}

	// leader election uses the Kubernetes API by writing to a
	// lock object, which can be a LeaseLock object (preferred),
	// a ConfigMap, or an Endpoints (deprecated) object.
	// Conflicting writes are detected and each client handles those actions
	// independently.
	cfg, err := buildConfig(e.kubeconfig)
	if err != nil {
		klog.Fatal(err)
	}
	client := clientset.NewForConfigOrDie(cfg)

	runWithContext := func(ctx context.Context) {
		// complete your controller loop here
		klog.Info("starting Controller loop...")
		for {
			select {
			case <-ctx.Done():
				klog.Info("stopping Controller loop...")
				return
			default:
				time.Sleep(5 * time.Second)
			}
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// signal handler.
	go func() {
		<-stopCh
		cancel()
	}()

	// we use the Lease lock type since edits to Leases are less common
	// and fewer objects in the cluster watch "all Leases".
	lock := &resourcelock.LeaseLock{
		LeaseMeta: metav1.ObjectMeta{
			Name:      e.leaseLockName,
			Namespace: e.leaseLockNamespace,
		},
		Client: client.CoordinationV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: e.podId,
		},
	}

	// start the leader election code loop
	leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
		Lock: lock,
		// IMPORTANT: you MUST ensure that any code you have that
		// is protected by the lease must terminate **before**
		// you call cancel. Otherwise, you could have a background
		// loop still running and another process could
		// get elected before your background loop finished, violating
		// the stated goal of the lease.
		ReleaseOnCancel: true,
		LeaseDuration:   30 * time.Second,
		RenewDeadline:   15 * time.Second,
		RetryPeriod:     5 * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {
				// we're notified when we start - this is where you would
				// usually put your code
				klog.Info("I am the leader")
				e.SetLeader(true)
				runWithContext(ctx)
			},
			OnStoppedLeading: func() {
				e.SetLeader(false)
				klog.Fatalf("leader election lost")
				return // exit
			},
			OnNewLeader: func(identity string) {
				// we're notified when new leader elected
				e.SetLeader(false)
				if identity == e.podId {
					// I just got the lock
					return
				}
				klog.Infof("new leader elected: %s", identity)
			},
		},
	})
}

func buildConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
		return cfg, nil
	}

	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
