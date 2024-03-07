package gobpmn_reflection

import (
	"strings"

	"github.com/deemount/gobpmnReflection/internals/utils"
)

// Count ...
type Count struct {
	Process          int
	Participant      int
	Message          int
	StartEvent       int
	EndEvent         int
	BusinessRuleTask int
	ManualTask       int
	ReceiveTask      int
	ScriptTask       int
	SendTask         int
	ServiceTask      int
	Task             int
	UserTask         int
	Flow             int
	Shape            int
	Edge             int
}

/*
 * @Base
 */

// countPool ...
func (r *Reflection) countPool(field, builderField string) {
	if strings.ToLower(field) == "pool" {
		if strings.Contains(builderField, "Process") {
			r.Process++
		}
		if strings.Contains(builderField, "ID") {
			r.Participant++
			r.Shape++
		}
	}
}

// countMessage ...
func (r *Reflection) countMessage(field, builderField string) {
	if strings.ToLower(field) == "message" {
		if strings.Contains(builderField, "Message") {
			r.Message++
			r.Edge++
		}
	}
}

/*
 * @Processes
 */

// countProcess ...
func (r *Reflection) countProcess(builderField string) {
	if strings.Contains(builderField, "Process") {
		r.Process++
	}
}

/*
 * @Elements
 */

func (r *Reflection) countElements(builderField string) {

	if utils.After(builderField, "From") == "" {

		switch true {
		case strings.Contains(builderField, "StartEvent"):
			r.StartEvent++
		case strings.Contains(builderField, "EndEvent"):
			r.EndEvent++
		case strings.Contains(builderField, "BusinessRuleTask"):
			r.BusinessRuleTask++
		case strings.Contains(builderField, "ManualTask"):
			r.ManualTask++
		case strings.Contains(builderField, "ReceiveTask"):
			r.ReceiveTask++
		case strings.Contains(builderField, "ScriptTask"):
			r.ScriptTask++
		case strings.Contains(builderField, "SendTask"):
			r.SendTask++
		case strings.Contains(builderField, "ServiceTask"):
			r.ServiceTask++
		case strings.Contains(builderField, "Task"):
			r.Task++
		case strings.Contains(builderField, "UserTask"):
			r.UserTask++
		}

		r.Shape++

	}
}
