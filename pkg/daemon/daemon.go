package daemon

import (
	"fmt"
	k8sClient "github.com/Mellanox/ib-kubernetes/pkg/k8s-client"
	"github.com/Mellanox/ib-kubernetes/pkg/utils"
	"github.com/Mellanox/ib-kubernetes/pkg/watcher"
	resEvenHandler "github.com/Mellanox/ib-kubernetes/pkg/watcher/resouce-event-handler"
	"reflect"
	"strconv"
	"strings"

	"github.com/golang/glog"
)

type Daemon interface {
	// ReadConfig read the configurations of the daemon
	ReadConfig() error
	// ValidateConfig validate the configurations of the daemon
	ValidateConfig() error
	// Init initialize the daemon with the needed component
	Init() error
	// Run run listener for k8s pod events.
	Run()
}

type daemonConfig struct {
	SubnetClientPlugin            string `json:"sub_net_manager_plugin"`
	GuidRangeStart                string `json:"guid_range_start"`
	GuidRangeEnd                  string `json:"guid_range_end"`
	SubnetManagerSecretConfigName string `json:"sub_net_manager_secret_config_name"`
	PeriodicUpdate                int    `json:"periodic_update"`
}

type daemon struct {
	watcher watcher.Watcher
}

const daemonNamespace = "kube-system"

// NewDaemon initializes the need components including k8s client, subnet manager client plugins, and guid pool.
// It returns error in case of failure.
func NewDaemon() (Daemon, error) {
	glog.Info("daemon NewDaemon():")
	podEventHandler := resEvenHandler.NewPodEventHandler()
	client, err := k8sClient.NewK8sClient()

	if err != nil {
		glog.Error(err)
		return nil, err
	}

	podWatcher := watcher.NewWatcher(podEventHandler, client)
	return &daemon{watcher: podWatcher}, nil
}

func (d *daemon) ReadConfig() error {
	glog.Info("daemon ReadConfig():")
	conf := &daemonConfig{}
	elements := reflect.ValueOf(conf).Elem()
	for i := 0; i < elements.NumField(); i++ {
		field := elements.Type().Field(i)
		envVar := strings.ToUpper(field.Tag.Get("json"))
		value, err := utils.LoadEnvVar(envVar)
		if err != nil {
			glog.Warningf("daemon ReadConfig: didn't env var for filed %s", envVar)
			continue
		}
		// string filed
		if field.Type == reflect.TypeOf("") {
			elements.FieldByName(field.Name).SetString(value)

		} else {
			// int filed
			i, err := strconv.Atoi(value)
			if err != nil {
				err := fmt.Errorf("daemon ReadConfig(): failed to parse int value of %s, with error: %v", envVar, err)
				glog.Error(err)
				return err
			}
			elements.FieldByName(field.Name).SetInt(int64(i))
		}
	}

	glog.V(3).Infof("daemon ReadConfig(): daemon config %+v", conf)
	return nil
}

func (d *daemon) ValidateConfig() error {
	glog.Info("daemon ValidateConfig():")

	return nil
}

func (d *daemon) Init() error {
	glog.Info("daemon Init():")
	return nil
}

func (d *daemon) Run() {
	glog.Info("daemon Run():")
	d.watcher.Run()
}
