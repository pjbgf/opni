package v2

import "github.com/rancher/opni/pkg/validation"

func (r *OpniReceiver) Validate() error {
	if r == nil {
		return validation.Error("Input is nil")
	}
	if r.Receiver == nil {
		return validation.Error("field receiver is required")
	}
	if err := r.Receiver.Validate(); err != nil {
		return err
	}
	return nil
}
