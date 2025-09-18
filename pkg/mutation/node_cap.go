package mutation

import (
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

type nodeCap struct {
	Logger logrus.FieldLogger
}

var _ nodeMutator = (*nodeCap)(nil)

func (nc *nodeCap) Name() string {
	return "node-cap"
}

func (nc *nodeCap) Mutate(node *corev1.Node) (*corev1.Node, error) {
	nc.Logger.Info("node-cap mutator")

	originalCPUCapacity := node.Status.Capacity.Cpu()
	originalCPUAllocatable := node.Status.Allocatable.Cpu()

	nc.Logger.Info("Original CPU Capacity: ", originalCPUCapacity.String())
	nc.Logger.Info("Original CPU Allocatable: ", originalCPUAllocatable.String())

	reserved := originalCPUCapacity
	reserved.Sub(*originalCPUAllocatable)

	newCPUAllocatable := originalCPUAllocatable.Value() * 2
	newCPUCapacity := newCPUAllocatable + reserved.Value()

	nc.Logger.Info("New CPU Allocatable: ", newCPUAllocatable)
	nc.Logger.Info("New CPU Capacity: ", newCPUCapacity)

	node.Status.Allocatable[corev1.ResourceCPU] = *resource.NewQuantity(newCPUAllocatable, resource.DecimalSI)
	node.Status.Capacity[corev1.ResourceCPU] = *resource.NewQuantity(newCPUCapacity, resource.DecimalSI)

	return node, nil
}
