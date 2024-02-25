package factory

type AnonymMap struct {
	Field int
	Value []string
}

// hashTable ...
// Note:
// @self is a pseudo and means "this hash is set to this element"
// Ruleset for a straight process without any decision:
// A Collaboration has one hash value: self
// A Participant has two hash values: id, process
// A Pool has one hash value: id
// A StartEvent has two hash values: self, next (flow: from)
// A Flow has three hash values: self, previous, next (task ...)
// A Task, Event has three hash values: self, previous (flow: from), next (flow. from)
// A EndEvent has two hash values: self, previous (flow: from)
// A Gateway has three hash values: self, previous (flow: from), next (flow: from)
// A SubProcess has two hash values: self, next (flow: from)
// A Process has one hash value: self
// A Diagram has one hash value: self
