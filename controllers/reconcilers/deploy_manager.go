package reconcilers

import "time"

type DeployManager struct {
}

func NewDeployManager() *DeployManager {
	return &DeployManager{}
}

func (d *DeployManager) IsDeployApplied() (bool, error) {
	time.Sleep(time.Second * 10)
	return true, nil
}

func (d *DeployManager) IsDeployReady() (bool, error) {
	return true, nil
}

func (d *DeployManager) ApplyDeploy() error {
	time.Sleep(time.Second * 10)
	return nil
}

func (d *DeployManager) CheckDeploy() (bool, error) {
	return true, nil
}
