package extendedjob

import (
	"time"

	"go.uber.org/zap"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"

	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type ExtendedJobController struct {
	client   client.Client
	log      *zap.SugaredLogger
	cache    cache.Cache
	waitFunc waitFunc
	runner   Runner
}

type waitFunc func(func(), time.Duration, <-chan struct{})

// Add creates a new ExtendedJob controller and adds it to the Manager
func Add(log *zap.SugaredLogger, mgr manager.Manager) error {
	c := NewExtendedJobController(log, mgr, wait.Until, NewRunner(log, mgr))
	return mgr.Add(c)
}

func NewExtendedJobController(
	log *zap.SugaredLogger,
	mgr manager.Manager,
	waitFunc waitFunc,
	runner Runner,
) *ExtendedJobController {
	return &ExtendedJobController{
		log:      log,
		client:   mgr.GetClient(),
		cache:    mgr.GetCache(),
		waitFunc: waitFunc,
		runner:   runner,
	}
}

func (ejc *ExtendedJobController) Start(stopCh <-chan struct{}) error {

	defer utilruntime.HandleCrash()
	ejc.log.Infof("Starting CronJob Manager")
	// Check things every 10 second.
	go ejc.waitFunc(ejc.wakeUp, 10*time.Second, stopCh)
	<-stopCh
	ejc.log.Infof("Shutting down CronJob Manager")
	return nil
}

func (ejc *ExtendedJobController) wakeUp() {
	ejc.log.Debugf("extendedjob controller wakeup")
	ejc.runner.Run()
}
