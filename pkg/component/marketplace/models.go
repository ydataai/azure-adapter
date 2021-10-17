package marketplace

import (
	"time"

	"github.com/Azure/go-autorest/autorest"
)

// UsageEventReq a type to represent the usage metering event request
type UsageEventReq struct {
	ResourceId string    `json:"resourceId"` // unique identifier of the resource against which usage is emitted.
	Quantity   float32   `json:"quantity"`   // how many units were consumed for the date and hour specified in effectiveStartTime, must be greater than 0, can be integer or float value
	Dimension  string    `json:"dimension"`  // custom dimension identifier
	StartTime  time.Time `json:"startTime"`  // time in UTC when the usage event occurred, from now and until 24 hours back
	Plan       string    `json:"plan"`       // id of the plan purchased for the offer
}

// UsageEventRes a type to represent the usage metering event response
type UsageEventRes struct {
	autorest.Response  `json:"-"`
	UsageEventId       string    `json:"usageEventId"`       // unique identifier associated with the usage event in Microsoft records
	Status             string    `json:"status"`             // this is the only value in case of single usage event
	MessageTime        time.Time `json:"messageTime"`        // time in UTC this event was accepted
	ResourceId         string    `json:"resourceId"`         // unique identifier of the resource against which usage is emitted. For SaaS it's the subscriptionId.
	Quantity           float32   `json:"quantity"`           // amount of emitted units as recorded by Microsoft
	Dimension          string    `json:"dimension"`          // custom dimension identifier
	EffectiveStartTime time.Time `json:"effectiveStartTime"` // time in UTC when the usage event occurred, as sent by the ISV
	PlanId             string    `json:"planId"`             // id of the plan purchased for the offer
}

// UsageEventBatchReq a type to represent the usage metering batch events request
type UsageEventBatchReq struct {
	Request []UsageEventReq `json:"request"` // batch events.
}

// UsageEventRes a type to represent the usage metering event response
type UsageEventBatchRes struct {
	autorest.Response `json:"-"`
	Count             int             `json:"count"`  // number of records in the response
	Result            []UsageEventRes `json:"result"` // result
}
