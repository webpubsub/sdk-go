package webpubsub

// WPSMPNSData is the struct used for the MPNS paylod
type WPSMPNSData struct {
	Title       string `json:"title"`
	Type        string `json:"type"`
	Count       int    `json:"count"`
	BackTitle   string `json:"back_title"`
	BackContent string `json:"back_content"`
	Custom      map[string]interface{}
}

// WPSFCMData is the struct used for the FCM paylod
type WPSFCMData struct {
	Data   WPSFCMDataFields `json:"data"`
	Custom map[string]interface{}
}

// WPSFCMDataFields is the helper struct used for the FCM paylod
type WPSFCMDataFields struct {
	Summary interface{} `json:"summary"`
	Custom  map[string]interface{}
}

// WPSAPSData is the helper struct used for the APNS paylod
type WPSAPSData struct {
	Alert    interface{} `json:"alert"`
	Badge    int         `json:"badge"`
	Sound    string      `json:"sound"`
	Title    string      `json:"title"`
	Subtitle string      `json:"subtitle"`
	Body     string      `json:"body"`
	Custom   map[string]interface{}
}

// WPSAPNSData is the struct used for the APNS paylod
type WPSAPNSData struct {
	APS    WPSAPSData `json:"aps"`
	Custom map[string]interface{}
}

// WPSAPNS2Data is the struct used for the APNS2 paylod
type WPSAPNS2Data struct {
	CollapseID string          `json:"collapseId"`
	Expiration string          `json:"expiration"`
	Targets    []WPSPushTarget `json:"targets"`
	Version    string          `json:"version"`
}

// WPSPushTarget is the helper struct used for the APNS2 paylod
type WPSPushTarget struct {
	Topic          string             `json:"topic"`
	ExcludeDevices []string           `json:"exclude_devices"`
	Environment    WPSPushEnvironment `json:"environment"`
}

type publishPushHelperBuilder struct {
	opts *publishPushHelperOpts
}

func newPublishPushHelperBuilder(webpubsub *WebPubSub) *publishPushHelperBuilder {
	builder := publishPushHelperBuilder{
		opts: &publishPushHelperOpts{
			webpubsub: webpubsub,
		},
	}

	return &builder
}

func newPublishPushHelperBuilderWithContext(webpubsub *WebPubSub,
	context Context) *publishPushHelperBuilder {
	builder := publishPushHelperBuilder{
		opts: &publishPushHelperOpts{
			webpubsub: webpubsub,
			ctx:       context,
		},
	}

	return &builder
}

// BuildPayload builds the push payload and returns an map of interface
func (b *publishPushHelperBuilder) BuildPayload() map[string]interface{} {
	response := make(map[string]interface{})
	apns := b.opts.buildAPNSPayload()
	if apns != nil {
		response["pn_apns"] = apns
		apns2 := b.opts.PushAPNS2Data
		if apns2 != nil {
			response["pn_push"] = apns2
		}
	}

	mpns := b.opts.buildMPNSPayload()
	if mpns != nil {
		response["pn_mpns"] = mpns
	}

	fcm := b.opts.buildFCMPayload()
	if fcm != nil {
		response["pn_gcm"] = fcm
	}

	if b.opts.CommonPayload != nil {
		for key, value := range b.opts.CommonPayload {
			response[key] = value
		}
	}

	return response
}

func (o *publishPushHelperOpts) buildAPNSPayload() map[string]interface{} {
	apns := make(map[string]interface{})
	if o.PushAPNSData != nil {
		aps := &o.PushAPNSData.APS
		if aps != nil {
			apsData := make(map[string]interface{})
			if aps.Alert != nil {
				apsData["alert"] = aps.Alert
			} else if aps.Subtitle != "" || aps.Body != "" || aps.Title != "" {
				alert := make(map[string]interface{})
				if aps.Subtitle != "" {
					alert["subtitle"] = aps.Subtitle
				}
				if aps.Title != "" {
					alert["title"] = aps.Title
				}
				if aps.Body != "" {
					alert["body"] = aps.Body
				}
				apsData["alert"] = alert
			}

			apsData["badge"] = aps.Badge

			if aps.Sound != "" {
				apsData["sound"] = aps.Sound
			}

			if aps.Custom != nil {
				for key, value := range aps.Custom {
					apsData[key] = value
				}
			}
			apns["aps"] = apsData
		}
		custom := o.PushAPNSData.Custom
		if custom != nil {
			for key, value := range custom {
				apns[key] = value
			}
		}
	}

	return apns
}

func (o *publishPushHelperOpts) buildMPNSPayload() map[string]interface{} {
	mpns := make(map[string]interface{})
	if o.PushMPNSData != nil {
		if o.PushMPNSData.Title != "" {
			mpns["title"] = o.PushMPNSData.Title
		}
		if o.PushMPNSData.Type != "" {
			mpns["type"] = o.PushMPNSData.Type
		}
		if o.PushMPNSData.BackTitle != "" {
			mpns["back_title"] = o.PushMPNSData.BackTitle
		}
		if o.PushMPNSData.BackContent != "" {
			mpns["back_content"] = o.PushMPNSData.BackContent
		}
		mpns["count"] = o.PushMPNSData.Count

		custom := o.PushMPNSData.Custom
		if custom != nil {
			for key, value := range custom {
				mpns[key] = value
			}
		}
	}

	return mpns
}

func (o *publishPushHelperOpts) buildFCMPayload() map[string]interface{} {
	fcm := make(map[string]interface{})
	if o.PushFCMData != nil {
		data := &o.PushFCMData.Data
		if data != nil {
			fcmData := make(map[string]interface{})
			if data.Summary != nil {
				fcmData["summary"] = data.Summary
			}

			custom := data.Custom
			if custom != nil {
				for key, value := range custom {
					fcmData[key] = value
				}
			}
			fcm["data"] = fcmData
		}

		custom := o.PushFCMData.Custom
		if custom != nil {
			for key, value := range custom {
				fcm[key] = value
			}
		}
	}

	return fcm
}

// SetAPNSPayload sets the APNS payload
func (b *publishPushHelperBuilder) SetAPNSPayload(pnAPNSData WPSAPNSData, pnAPNS2Data []WPSAPNS2Data) *publishPushHelperBuilder {
	b.opts.PushAPNSData = &pnAPNSData
	b.opts.PushAPNS2Data = pnAPNS2Data

	return b
}

// SetMPNSPayload sets the MPNS payload
func (b *publishPushHelperBuilder) SetMPNSPayload(pnMPNSData WPSMPNSData) *publishPushHelperBuilder {
	b.opts.PushMPNSData = &pnMPNSData

	return b
}

// SetCommonPayload sets the common payload
func (b *publishPushHelperBuilder) SetCommonPayload(commonPayload map[string]interface{}) *publishPushHelperBuilder {
	b.opts.CommonPayload = commonPayload

	return b
}

// SetFCMPayload sets the FCM payload
func (b *publishPushHelperBuilder) SetFCMPayload(pnFCMData WPSFCMData) *publishPushHelperBuilder {
	b.opts.PushFCMData = &pnFCMData

	return b
}

type publishPushHelperOpts struct {
	webpubsub *WebPubSub

	PushAPNS2Data  []WPSAPNS2Data
	PushAPNSData   *WPSAPNSData
	PushMPNSData   *WPSMPNSData
	PushFCMData    *WPSFCMData
	CommonPayload  map[string]interface{}
	PushCustomData map[string]interface{}

	ctx Context
}

func (o *publishPushHelperOpts) context() Context {
	return o.ctx
}
