package validation

import (
	"errors"
	"os/exec"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

// Validator is a container for mutation
type Validator struct {
	Logger *logrus.Entry
}

// NewValidator returns an initialised instance of Validator
func NewValidator(logger *logrus.Entry) *Validator {
	return &Validator{Logger: logger}
}

// podValidators is an interface used to group functions mutating pods
type podValidator interface {
	Validate(*corev1.Pod) (validation, error)
	Name() string
}

type validation struct {
	Valid  bool
	Reason string
}

// ValidatePod returns true if a pod is valid
func (v *Validator) ValidatePod(pod *corev1.Pod) (validation, error) {
/*	var podName string
	if pod.Name != "" {
		podName = pod.Name
	} else {
		if pod.ObjectMeta.GenerateName != "" {
			podName = pod.ObjectMeta.GenerateName
		}
	}
	//log := logrus.WithField("pod_name", podName)
  //log.Print("delete me")*/
	log := logrus.WithField("image_name", pod.Spec.Containers[0].Image)
  	log.Print("test sboms")

  	imageName := "registry:" + pod.Spec.Containers[0].Image
  	logrus.Print(imageName)
  	outsbom, err:= exec.Command("syft", imageName, "-o json").Output()
  	logrus.Print(string(outsbom))
  	
	err1 := errors.New(string(outsbom))
	if err == nil {
		return validation{Valid: false, Reason: err1.Error()}, err1
	}
	
	// list of all validations to be applied to the pod
	validations := []podValidator{
		nameValidator{v.Logger},
	}

	// apply all validations
	for _, v := range validations {
		var err error
		vp, err := v.Validate(pod)
		if err != nil {
			return validation{Valid: false, Reason: err.Error()}, err
		}
		if !vp.Valid {
			return validation{Valid: false, Reason: vp.Reason}, err
		}
	}

	return validation{Valid: true, Reason: "valid pod"}, nil
}
