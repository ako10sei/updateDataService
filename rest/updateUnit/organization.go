package updateUnit

import "visiologyDataUpdate/digital_profile"

type Response struct {
	Count         int                            `json:"count"`
	Next          any                            `json:"next"`
	Previous      any                            `json:"previous"`
	Organizations []digital_profile.Organization `json:"results"`
}
