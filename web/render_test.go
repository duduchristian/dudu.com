package web

import (
	"encoding/json"
	"fmt"
	"github.com/ajg/form"
	"github.com/stretchr/testify/assert"
	"html"
	"reflect"
	"regexp"
	"strings"
	"testing"
)
import v8 "rogchap.com/v8go"

func TestExample(t *testing.T) {
	ctx := v8.NewContext()                                  // creates a new V8 context with a new Isolate aka VM
	ctx.RunScript("const add = (a, b) => a + b", "math.js") // executes a script on the global context
	ctx.RunScript("const result = add(3, 4)", "main.js")    // any functions previously added to the context can be called
	val, _ := ctx.RunScript("result", "value.js")           // return a value in JavaScript back to Go
	fmt.Printf("addition result: %s", val)
}

func TestExample2(t *testing.T) {
	iso := v8.NewIsolate() // create a new VM
	// a template that represents a JS function
	printfn := v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		fmt.Printf("%v", info.Args()) // when the JS function is called this Go callback will execute
		return nil                    // you can return a value back to the JS caller if required
	})
	global := v8.NewObjectTemplate(iso)       // a template that represents a JS Object
	global.Set("print", printfn)              // sets the "print" property of the Object to our function
	ctx := v8.NewContext(iso, global)         // new Context with the global Object set to our object template
	ctx.RunScript("print('foo')", "print.js") // will execute the Go callback with a single argunent 'foo'
}

func TestExample3(t *testing.T) {
	ctx := v8.NewContext()                           // new context with a default VM
	obj := ctx.Global()                              // get the global object from the context
	obj.Set("version", "v1.0.0")                     // set the property "version" on the object
	val, _ := ctx.RunScript("version", "version.js") // global object will have the property set within the JS VM
	fmt.Printf("version: %s", val)

	if obj.Has("version") { // check if a property exists on the object
		obj.Delete("version") // remove the property from the object
	}
}

func TestExample4(t *testing.T) {
	source := "const multiply = (a, b) => a * b"
	iso1 := v8.NewIsolate()                                                         // creates a new JavaScript VM
	ctx1 := v8.NewContext(iso1)                                                     // new context within the VM
	script1, _ := iso1.CompileUnboundScript(source, "math.js", v8.CompileOptions{}) // compile script to get cached data
	script1.Run(ctx1)

	cachedData := script1.CreateCodeCache()

	iso2 := v8.NewIsolate()     // create a new JavaScript VM
	ctx2 := v8.NewContext(iso2) // new context within the VM

	script2, _ := iso2.CompileUnboundScript(source, "math.js", v8.CompileOptions{CachedData: cachedData}) // compile script in new isolate with cached data
	script2.Run(ctx2)
}

const ajaxReq = `
{
 "placementKey":"s4936796149952",
 "deviceType":"PHONE",
 "osType":"ANDROID",
 "osVersion":"7.0",
 "deviceScreenWidth":1680,
 "deviceScreenHeight":1050,
 "appPackageName":"com.opera.interstitial.web",
 "appVersion":"",
 "deviceVendor":"",
 "deviceModel":"",
 "operator":"",
 "connectionType":"UNKNOWN",
 "userConsent":"false",
 "advertisingId":"",
 "operaId":"",
 "availableServices":[
        "GOOGLE_PLAY"
    ],
 "userId":"fcbf0782-a36a-4192-a2ab-c6328a18gfsefe01",
 "renderWidth":500,
 "renderHeight":100,
 "languageCode":"zh",
    "browserLanguage": "zh-CN",
 "adCount":3,
 "ip":"10.7.6.11",
 "timestamp":1581497446,
 "token":"9390a38ee2aaf7d53fef7f711df82f37",
    "viewer": "demo-test.adx.opera.com",
 "supportedCreativeTypes":
    [
     "BIG_CARD",
     "SMALL_CARD",
     "VIDEO_16x9",
     "NATIVE_BANNER_6x5",
     "DISPLAY_HTML_300x250",
                    "DISPLAY_HTML_320x480",
     "VIDEO_9x16", 
                    "VAST_3_XML"
    ]
}
`

type PageAdReq struct {
	AdvertisingId          string   `json:"advertisingId" form:"aid"`
	AvailableServices      []string `json:"availableServices" form:"avs"`
	CountryCode            string   `json:"countryCode" form:"cc"`
	XCountry               string   `json:"xCountry,omitempty" form:"xcc"`
	XOperator              string   `json:"xOperator,omitempty" form:"xop"`
	ChannelId              string   `json:"channelId" from:"cid"`
	ConnectionType         string   `json:"connectionType" form:"conn"`
	SupportedCreativeTypes []string `json:"supportedCreativeTypes" form:"ct"`
	CityCode               string   `json:"cityCode" form:"ctc"`
	DeviceModel            string   `json:"deviceModel" form:"dm"`
	DeviceType             string   `json:"deviceType" form:"dt"`
	DeviceVendor           string   `json:"deviceVendor" form:"dv"`
	FloorPriceInLi         int      `json:"floorPriceInLi" form:"fpil"`
	RenderHeight           int      `json:"renderHeight" form:"h"`
	HashedAndroidId        string   `json:"hashedAndroidId" form:"haid"`
	HashedImei             string   `json:"hashedImei" form:"himei"`
	Latitude               float64  `json:"latitude" form:"lat"`
	LanguageCode           string   `json:"languageCode" form:"lc"`
	Longitude              float64  `json:"longitude" form:"lng"`
	OperaId                string   `json:"operaId" form:"oid"`
	Operator               string   `json:"operator" form:"opr"`
	OSType                 string   `json:"osType" form:"ost"`
	OSVersion              string   `json:"osVersion" form:"osv"`
	PlacementKey           string   `json:"placementKey" form:"pk"`
	PositionTimestamp      int64    `json:"positionTimestamp" form:"pts"`
	AppPackageName         string   `json:"appPackageName" form:"pkg"`
	AppVersion             string   `json:"appVersion" form:"pkgv"`
	DeviceScreenHeight     int      `json:"deviceScreenHeight" form:"sh"`
	DeviceScreenWidth      int      `json:"deviceScreenWidth" form:"sw"`
	Token                  string   `json:"token" form:"tk"`
	Timestamp              int64    `json:"timestamp" form:"ts"`
	UserConsent            string   `json:"userConsent" form:"uc"`
	UserId                 string   `json:"userId" form:"uid"`
	RenderWidth            int      `json:"renderWidth" form:"w"`

	BrowserLanguage string `json:"browserLanguage" form:"bl"`
	DocumentCharset string `json:"documentCharset" form:"cst"`
	Debug           bool   `json:"debug" form:"debug"`
	Referrer        string `json:"referrer" form:"rf"`
	ScrollHeight    int    `json:"scrollHeight" form:"sch"`
	ScrollLeft      int    `json:"scrollLeft" form:"scl"`
	ScrollTop       int    `json:"scrollTop" form:"sct"`
	ScrollWidth     int    `json:"scrollWidth" form:"scw"`
	Title           string `json:"-" form:"title"`
	Timezone        string `json:"timezone" form:"tz"`
	URL             string `json:"url" form:"url"`
	ViewportHeight  int    `json:"viewportHeight" form:"vph"`
	ViewportWidth   int    `json:"viewportWidth" form:"vpw"`
	Viewer          string `json:"viewer" form:"vr"`

	Pubcid string `json:"pubcid" form:"pubcid"`
	Tdid   string `json:"tdid" form:"tdid"`

	AdCount   int    `json:"adCount" form:"-"`
	PageTitle string `json:"pageTitle" form:"-"`
}

func (req PageAdReq) FilterXSS() (shouldFilter bool) {
	return filterCore(reflect.ValueOf(req))
}

func filterCore(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		str := value.String()
		if str != "" && html.EscapeString(str) != str {
			return true
		}
	case reflect.Pointer:
		if !value.IsNil() {
			return filterCore(value.Elem())
		}
	case reflect.Struct:
		for i := 0; i < value.NumField(); i++ {
			if filterCore(value.Field(i)) {
				return true
			}
		}
	case reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			if filterCore(value.Index(i)) {
				return true
			}
		}
	}
	return false
}

func TestAjaxReqToPageAdReq(t *testing.T) {
	var req PageAdReq
	err := json.Unmarshal([]byte(ajaxReq), &req)
	assert.NoError(t, err)
	q, e := form.EncodeToValues(&req)
	assert.NoError(t, e)
	var re = regexp.MustCompile(`\.[0-9]=`)
	s := strings.Replace(re.ReplaceAllString(q.Encode(), `=`), "&", "\n", -1)
	t.Log(s)
}

func TestXSSFilter(t *testing.T) {
	var req PageAdReq
	err := json.Unmarshal([]byte(ajaxReq), &req)
	assert.NoError(t, err)
	assert.Equal(t, false, req.FilterXSS(), "should not be filtered")
	req.URL = `</ins><svg onload=alert(1)>`
	assert.Equal(t, true, req.FilterXSS(), "should be filtered")
	req.URL = ""
	req.AvailableServices = append(req.AvailableServices, `</ins><svg onload=alert(1)>`)
	assert.Equal(t, true, req.FilterXSS(), "should be filtered")
}

// BenchmarkXSSFilter 714.4 ns/op on M1 Pro
func BenchmarkXSSFilter(b *testing.B) {
	var req PageAdReq
	json.Unmarshal([]byte(ajaxReq), &req)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req.FilterXSS()
	}
}
