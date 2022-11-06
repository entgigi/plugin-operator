package utility

import (
	"fmt"
	"os"

	"github.com/entgigi/plugin-operator.git/common"
)

// GetWatchNamespace returns the Namespace the operator should be watching for changes
func GetWatchNamespace() (string, error) {
	// WatchNamespaceEnvVar is the constant for env variable WATCH_NAMESPACE
	// which specifies the Namespace to watch.
	// An empty value means the operator is running with cluster scope.

	ns, found := os.LookupEnv(common.WatchNamespaceEnvVar)
	if !found {
		return "", fmt.Errorf("%s must be set", common.WatchNamespaceEnvVar)
	}
	return ns, nil
}

func GetOperatorDeploymentType() string {
	operatorType, found := os.LookupEnv(common.OperatorTypeEnvVar)
	if found {
		return operatorType
	} else {
		// default
		return common.OperatorTypeStandard
	}
}
