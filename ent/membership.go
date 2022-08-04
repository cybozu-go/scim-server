package ent

type Membership struct {
	Value   string `json:"value"`
	Ref     string `json:"$ref"`
	Type    string `json:"type"`
	Display string `json:"display"`
}
