package web

import (
	"encoding/json"
	"fmt"
	"github.com/ajg/form"
	"github.com/stretchr/testify/assert"
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
            "advertisingId":"",
            "availableServices":null,
            "countryCode":"za",
            "channelId":"",
            "connectionType":"unknown",
            "supportedCreativeTypes":[
                "BIG_CARD",
                "DISPLAY_HTML_180x150",
                "DISPLAY_HTML_250x250",
                "DISPLAY_HTML_300x100",
                "DISPLAY_HTML_300x250",
                "DISPLAY_HTML_300x50",
                "DISPLAY_HTML_300x600",
                "DISPLAY_HTML_320x100",
                "DISPLAY_HTML_320x140",
                "DISPLAY_HTML_320x480",
                "DISPLAY_HTML_320x50",
                "DISPLAY_HTML_336x280",
                "DISPLAY_HTML_360x375",
                "DISPLAY_HTML_468x60",
                "DISPLAY_HTML_728x90",
                "DISPLAY_HTML_970x90",
                "NATIVE_BANNER_2x1",
                "NATIVE_BANNER_3x1",
                "NATIVE_BANNER_4x1",
                "NATIVE_BANNER_5x1",
                "NATIVE_BANNER_6x1",
                "NATIVE_BANNER_6x5",
                "NATIVE_SMALL_BANNER",
                "JS_TAG",
                "JS_TAG_LIST",
                "SMALL_CARD",
                "VIDEO_16x9",
                "VIDEO_9x16",
                "VAST_3_URL",
                "VAST_3_XML"
            ],
            "cityCode":"",
            "deviceModel":"PPA-LX2",
            "deviceType":"PHONE",
            "deviceVendor":"HUAWEI",
            "floorPriceInLi":0,
            "renderHeight":0,
            "hashedAndroidId":"",
            "hashedImei":"",
            "latitude":0,
            "languageCode":"en",
            "longitude":0,
            "operaId":"0c95fe3ccac1d82d41fa8e15f496e190b8b87ec8ca2b8adb16ef0d4a09b0e218",
            "operator":"65501",
            "osType":"ANDROID",
            "osVersion":"10",
            "placementKey":"s4592841290048",
            "positionTimestamp":0,
            "appPackageName":"com.opera.mini.native",
            "appVersion":"67.0.2254.64762",
            "deviceScreenHeight":0,
            "deviceScreenWidth":0,
            "token":"6aa76fdd3c8d8d273076a83fcba8322a",
            "timestamp":1677747796,
            "userConsent":"true",
            "userId":"ed27a29c-089a-41fc-a126-cd812d3ab943",
            "renderWidth":328,
            "browserLanguage":"en-US",
            "documentCharset":"UTF-8",
            "debug":false,
            "referrer":"",
            "scrollHeight":3746,
            "scrollLeft":0,
            "scrollTop":0,
            "scrollWidth":360,
            "timezone":"GMT+0200",
            "url":"https://cdn-af.feednews.com/news/detail/b2e79092529e0a7cfa61e6f52fad21a5?features=2114057\u0026country=za\u0026uid=37c9a154ff0e4f4699f88d2fa3b6ca2107088822\u0026like_count=1\u0026client=mini\u0026language=en",
            "viewportHeight":679,
            "viewportWidth":360,
            "viewer":"cdn-af.feednews.com",
            "userAgent":"Mozilla/5.0 (Linux; U; Android 10; PPA-LX2 Build/HUAWEIPPA-LX2; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/88.0.4324.93 Mobile Safari/537.36 OPR/67.0.2254.64762",
            "pubcid":"aa6c938c-f8bd-4387-9fa2-c76de4d4f607",
            "tdid":"",
            "adCount":0,
            "pageTitle":""
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
		if str != "" && strings.ContainsAny(str, "<>") {
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
