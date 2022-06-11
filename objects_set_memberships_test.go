package webpubsub

import (
	"fmt"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	h "github.com/webpubsub/go/v7/tests/helpers"
	"github.com/webpubsub/go/v7/utils"
)

func AssertSetMemberships(t *testing.T, checkQueryParam, testContext bool, withFilter bool, withSort bool) {
	assert := assert.New(t)
	pn := NewWebPubSub(NewDemoConfig())

	incl := []WPSMembershipsInclude{
		WPSMembershipsIncludeCustom,
	}

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	if !checkQueryParam {
		queryParam = nil
	}

	inclStr := EnumArrayToStringArray(incl)

	o := newSetMembershipsBuilder(pn)
	if testContext {
		o = newSetMembershipsBuilderWithContext(pn, backgroundContext)
	}

	spaceID := "id0"
	limit := 90
	start := "Mxmy"
	end := "Nxny"

	o.UUID(spaceID)
	o.Include(incl)
	o.Limit(limit)
	o.Start(start)
	o.End(end)
	o.Count(false)
	o.QueryParam(queryParam)

	id0 := "id0"
	if withFilter {
		o.Filter("name like 'a*'")
	}
	sort := []string{"name", "created:desc"}
	if withSort {
		o.Sort(sort)
	}

	custom := make(map[string]interface{})
	custom["a1"] = "b1"
	custom["c1"] = "d1"

	channel := WPSMembershipsChannel{
		ID: id0,
	}

	in := WPSMembershipsSet{
		Channel: channel,
		Custom:  custom,
	}

	inArr := []WPSMembershipsSet{
		in,
	}

	custom2 := make(map[string]interface{})
	custom2["a2"] = "b2"
	custom2["c2"] = "d2"

	o.Set(inArr)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/objects/%s/uuids/%s/channels", pn.Config.SubscribeKey, spaceID),
		path, []int{})

	body, err := o.opts.buildBody()
	assert.Nil(err)
	expectedBody := "{\"set\":[{\"channel\":{\"id\":\"id0\"},\"custom\":{\"a1\":\"b1\",\"c1\":\"d1\"}}]}"

	assert.Equal(expectedBody, string(body))

	if checkQueryParam {
		u, _ := o.opts.buildQuery()
		assert.Equal("v1", u.Get("q1"))
		assert.Equal("v2", u.Get("q2"))
		assert.Equal(string(utils.JoinChannels(inclStr)), u.Get("include"))
		assert.Equal(strconv.Itoa(limit), u.Get("limit"))
		assert.Equal(start, u.Get("start"))
		assert.Equal(end, u.Get("end"))
		assert.Equal("0", u.Get("count"))
		if withFilter {
			assert.Equal("name like 'a*'", u.Get("filter"))
		}
		if withSort {
			v := &url.Values{}
			SetQueryParamAsCommaSepString(v, sort, "sort")
			assert.Equal(v.Get("sort"), u.Get("sort"))
		}

	}

}

func TestSetMemberships(t *testing.T) {
	AssertSetMemberships(t, true, false, false, false)
}

func TestSetMembershipsContext(t *testing.T) {
	AssertSetMemberships(t, true, true, false, false)
}

func TestSetMembershipsWithFilter(t *testing.T) {
	AssertSetMemberships(t, true, false, true, false)
}

func TestSetMembershipsWithFilterContext(t *testing.T) {
	AssertSetMemberships(t, true, true, true, false)
}

func TestSetMembershipsWithSort(t *testing.T) {
	AssertSetMemberships(t, true, false, false, true)
}

func TestSetMembershipsWithSortContext(t *testing.T) {
	AssertSetMemberships(t, true, true, false, true)
}

func TestSetMembershipsWithFilterWithSort(t *testing.T) {
	AssertSetMemberships(t, true, false, true, true)
}

func TestSetMembershipsWithFilterWithSortContext(t *testing.T) {
	AssertSetMemberships(t, true, true, true, true)
}

func TestSetMembershipsResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewWebPubSub(NewDemoConfig())
	opts := &setMembershipsOpts{
		webpubsub: pn,
	}
	jsonBytes := []byte(`s`)

	_, _, err := newWPSSetMembershipsResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("webpubsub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestSetMembershipsResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewWebPubSub(NewDemoConfig())
	opts := &setMembershipsOpts{
		webpubsub: pn,
	}
	jsonBytes := []byte(`{"status":200,"data":[{"id":"spaceid3","custom":{"a3":"b3","c3":"d3"},"channel":{"id":"spaceid3","name":"spaceid3name","description":"spaceid3desc","custom":{"a":"b"},"created":"2019-08-23T10:34:43.985248Z","updated":"2019-08-23T10:34:43.985248Z","eTag":"Aazjn7vC3oDDYw"},"created":"2019-08-23T10:41:17.156491Z","updated":"2019-08-23T10:41:17.156491Z","eTag":"AamrnoXdpdmzjwE"}],"totalCount":1,"next":"MQ"}`)

	r, _, err := newWPSSetMembershipsResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal(1, r.TotalCount)
	assert.Equal("MQ", r.Next)
	assert.Equal("spaceid3", r.Data[0].ID)
	assert.Equal("spaceid3", r.Data[0].Channel.ID)
	assert.Equal("spaceid3name", r.Data[0].Channel.Name)
	assert.Equal("spaceid3desc", r.Data[0].Channel.Description)
	//assert.Equal("2019-08-23T10:34:43.985248Z", r.Data[0].Channel.Created)
	assert.Equal("2019-08-23T10:34:43.985248Z", r.Data[0].Channel.Updated)
	assert.Equal("Aazjn7vC3oDDYw", r.Data[0].Channel.ETag)
	assert.Equal("b", r.Data[0].Channel.Custom["a"])
	assert.Equal("2019-08-23T10:41:17.156491Z", r.Data[0].Created)
	assert.Equal("2019-08-23T10:41:17.156491Z", r.Data[0].Updated)
	assert.Equal("AamrnoXdpdmzjwE", r.Data[0].ETag)
	assert.Equal("b3", r.Data[0].Custom["a3"])
	assert.Equal("d3", r.Data[0].Custom["c3"])

	assert.Nil(err)
}
