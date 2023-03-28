// Package metering provides objects to interact with metering API
package metering

import (
	"net/http"
	"time"
)

// UsageEventReq a type to represent the usage metering event request
type usageEvent struct {
	Dimension          string    `json:"dimension"`
	Quantity           float32   `json:"quantity"`
	EffectiveStartTime time.Time `json:"effectiveStartTime"`

	ResourceURI string `json:"resourceUri"` // unique identifier of the resource against which usage is emitted.
	PlanID      string `json:"planId"`      // id of the plan purchased for the offer
}

// UsageEventRes a type to represent the usage metering event response
type usageEventResponse struct {
	*http.Response     `json:"-"`
	UsageEventId       string                `json:"usageEventId"`       // unique identifier associated with the usage event in Microsoft records
	Status             string                `json:"status"`             // this is the only value in case of single usage event
	MessageTime        time.Time             `json:"messageTime"`        // time in UTC this event was accepted
	ResourceId         string                `json:"resourceId"`         // unique identifier of the resource against which usage is emitted. For SaaS it's the subscriptionId.
	Quantity           float32               `json:"quantity"`           // amount of emitted units as recorded by Microsoft
	Dimension          string                `json:"dimension"`          // custom dimension identifier
	EffectiveStartTime time.Time             `json:"effectiveStartTime"` // time in UTC when the usage event occurred, as sent by the ISV
	PlanId             string                `json:"planId"`             // id of the plan purchased for the offer
	Error              usageEventErrorDetail `json:"error"`
}

// UsageEventErrorDetail represents a detail error mensage.
type usageEventErrorDetail struct {
	Message string                  `json:"message"`
	Target  string                  `json:"target"`
	Code    string                  `json:"code"`
	Details []usageEventErrorDetail `json:"details,omitempty"`
}

// UsageEventBatchReq a type to represent the usage metering batch events request
type usageEventBatch struct {
	Events []usageEvent `json:"request"` // batch events.
}

// UsageEventRes a type to represent the usage metering event response
type usageEventBatchResponse struct {
	*http.Response `json:"-"`
	Count          int                  `json:"count"`  // number of records in the response
	Result         []usageEventResponse `json:"result"` // result
}
