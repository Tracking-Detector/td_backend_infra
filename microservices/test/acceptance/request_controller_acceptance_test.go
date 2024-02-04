package acceptance

import (
	"context"
	"fmt"
	"net/http"
	"tds/shared/configs"
	"tds/shared/controller"
	"tds/shared/models"
	"tds/shared/repository"
	"tds/shared/service"
	"tds/test/testsupport"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

func TestRequestControllerAcceptance(t *testing.T) {
	suite.Run(t, &RequestControllerAcceptanceTest{})
}

type RequestControllerAcceptanceTest struct {
	suite.Suite
	requestRepo       models.RequestRepository
	requestController *controller.RequestController
	requestService    *service.RequestService
	ctx               context.Context
}

func (suite *RequestControllerAcceptanceTest) SetupTest() {
	suite.ctx = context.Background()
	suite.requestRepo = repository.NewMongoRequestRepository(configs.GetDatabase(configs.ConnectDB(suite.ctx)))
	suite.requestService = service.NewRequestService(suite.requestRepo)
	suite.requestController = controller.NewRequestController(suite.requestService)
	go func() {
		fmt.Println("Starting server...")
		suite.requestController.Start()

	}()
	time.Sleep(5 * time.Second)

	suite.requestRepo.DeleteAll(suite.ctx)
}

func (suite *RequestControllerAcceptanceTest) TearDownTest() {
	suite.requestController.Stop()
}

func (suite *RequestControllerAcceptanceTest) TestHealth_Success() {
	// given

	// when
	resp, err := testsupport.Get("http://localhost:8081/requests/health")

	// then
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	fmt.Println(resp.Body, "{\"message\":\"System is running correct.\",\"status\":200}")
}

func (suite *RequestControllerAcceptanceTest) TestCreateRequest_Success() {
	// given
	request := `{
        "documentId": "CFB05D1A2E1B7E6B44813CCBB3ED7638",
        "documentLifecycle": "active",
        "frameId": 0,
        "frameType": "outermost_frame",
        "initiator": "https://www.sportsnet.ca",
        "method": "GET",
        "parentFrameId": -1,
        "requestId": "979320",
        "tabId": 19381,
        "timeStamp": 1660340377939.461,
        "type": "image",
        "url": "https://dpm.demdex.net/ibs:dpid=411&dpuuid=YvYavQAF2aAqsAA0&d_uuid=78854490201948989640033682043832183304",
        "requestHeaders": [
            {
                "name": "sec-ch-ua",
                "value": "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"101\", \"Google Chrome\";v=\"101\""
            },
            {
                "name": "sec-ch-ua-mobile",
                "value": "?0"
            },
            {
                "name": "User-Agent",
                "value": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.67 Safari/537.36"
            },
            {
                "name": "sec-ch-ua-platform",
                "value": "\"Windows\""
            },
            {
                "name": "Accept",
                "value": "image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8"
            },
            {
                "name": "Sec-Fetch-Site",
                "value": "cross-site"
            },
            {
                "name": "Sec-Fetch-Mode",
                "value": "no-cors"
            },
            {
                "name": "Sec-Fetch-Dest",
                "value": "image"
            },
            {
                "name": "Referer",
                "value": "https://www.sportsnet.ca/"
            },
            {
                "name": "Accept-Encoding",
                "value": "gzip, deflate, br"
            },
            {
                "name": "Accept-Language",
                "value": "en-US,en;q=0.9"
            },
            {
                "name": "Cookie",
                "value": "demdex=78854490201948989640033682043832183304; dpm=78854490201948989640033682043832183304; DST=; dextp=359-1-1660295910237|477-1-1660295911238|771-1-1660295912202|903-1-1660295913168|1957-1-1660295914176|30646-1-1660295915161|57282-1-1660295916177|129099-1-1660296474248|21-1-1660296879877|269-1-1660296881216|282-1-1660296883280|3-1-1660296884304|375-1-1660296885171|358-1-1660296887163|843-1-1660296889178|540-1-1660296890174|832-1-1660296892175|1083-1-1660296893177|1085-1-1660296894173|1086-1-1660296895166|1087-1-1660296896175|1088-1-1660296897173|1175-1-1660296898172|6835-1-1660296899173|19913-1-1660296900171|83349-1-1660296901168|144230-1-1660297642942|144231-1-1660297643175|144232-1-1660297644222|144233-1-1660297645176|144234-1-1660297646186|144235-1-1660297647176|144236-1-1660297648179|144237-1-1660297649174|20-1-1660298040164|3462-1-1660298042171|70027-1-1660298043216|152416-1-1660298052393|60-1-1660298546233|22052-1-1660298548162|30064-1-1660298549211|73426-1-1660298550223|121998-1-1660298551455|199624-1-1660298552271|420-1-1660298717173|1121-1-1660298719167|28645-1-1660298721218|575-1-1660298723188|53196-1-1660298724168|87898-1-1660298729163|208568-1-1660298730226|175765-1-1660298748173|963840-1-1660298749181|796-1-1660298788880|601-1-1660298924182|348447-1-1660298938248|127444-1-1660298939164|470-1-1660299158731|1123-1-1660299160387|139200-1-1660299162277|3047-1-1660299214215|22054-1-1660299215165|22069-1-1660299216164|49276-1-1660299218268|66013-1-1660299219175|81309-1-1660299220187|13870-1-1660299725215|80742-1-1660299727176|275754-1-1660299730358|481-1-1660299787310|19566-1-1660299790197|23728-1-1660299791163|30432-1-1660299792187|66757-1-1660299794181|134096-1-1660299795173|147592-1-1660299804223|461447-1-1660299805165|144228-1-1660299852226|144229-1-1660299853217|67587-1-1660300211165|75557-1-1660300213164|408820-1-1660300218165|72352-1-1660300423416|285689-1-1660301434219|12105-1-1660301647279|30862-1-1660301787174|178522-1-1660301790172|992-1-1660301872329|2299-1-1660301874840|96678-1-1660302242178|640-1-1660302618232|782-1-1660302665162|444422-1-1660303572174|822-1-1660303943258|79908-1-1660303962167|2340-1-1660303963167|161033-1-1660303975180|16292-1-1660304245312|47438-1-1660304246163|57289-1-1660304247169|58342-1-1660305505179|1586-1-1660306188186|82530-1-1660309708213|58051-1-1660309792167|96420-1-1660309793177|466-1-1660313817167|445-1-1660313856378|530-1-1660313859309|1127-1-1660313865198|1342-1-1660313869245|13485-1-1660313871190|75884-1-1660313883164|143525-1-1660314082221|390122-1-1660317069182|38117-1-1660319516165|19360-1-1660321530184|134084-1-1660321547175|139423-1-1660322341168|1265-1-1660322435165|61283-1-1660324353163|411-1-1660325686185|87880-1-1660326432202"
            }
        ],
        "response": {
            "documentId": "CFB05D1A2E1B7E6B44813CCBB3ED7638",
            "documentLifecycle": "active",
            "frameId": 0,
            "frameType": "outermost_frame",
            "fromCache": false,
            "initiator": "https://www.sportsnet.ca",
            "ip": "54.154.38.9",
            "method": "GET",
            "parentFrameId": -1,
            "requestId": "979320",
            "responseHeaders": [
                {
                    "name": "Cache-Control",
                    "value": "no-cache,no-store,must-revalidate,max-age=0,proxy-revalidate,no-transform,private"
                },
                {
                    "name": "DCS",
                    "value": "dcs-prod-irl1-1-v038-0bef0d017.edge-irl1.demdex.com 3 ms"
                },
                {
                    "name": "Expires",
                    "value": "Thu, 01 Jan 1970 00:00:00 UTC"
                },
                {
                    "name": "P3P",
                    "value": "policyref=\"/w3c/p3p.xml\", CP=\"NOI NID CURa ADMa DEVa PSAa PSDa OUR SAMa BUS PUR COM NAV INT\""
                },
                {
                    "name": "Pragma",
                    "value": "no-cache"
                },
                {
                    "name": "set-cookie",
                    "value": "dpm=78854490201948989640033682043832183304; Max-Age=15552000; Expires=Wed, 08 Feb 2023 21:39:39 GMT; Path=/; Domain=.dpm.demdex.net; Secure; SameSite=None"
                },
                {
                    "name": "set-cookie",
                    "value": "demdex=78854490201948989640033682043832183304; Max-Age=15552000; Expires=Wed, 08 Feb 2023 21:39:39 GMT; Path=/; Domain=.demdex.net; Secure; SameSite=None"
                },
                {
                    "name": "Strict-Transport-Security",
                    "value": "max-age=31536000; includeSubDomains"
                },
                {
                    "name": "X-Content-Type-Options",
                    "value": "nosniff"
                },
                {
                    "name": "X-TID",
                    "value": "GxuhoSwmRBI="
                },
                {
                    "name": "Content-Length",
                    "value": "0"
                },
                {
                    "name": "Connection",
                    "value": "keep-alive"
                }
            ],
            "statusCode": 200,
            "statusLine": "HTTP/1.1 200 OK",
            "tabId": 19381,
            "timeStamp": 1660340377967.0042,
            "type": "image",
            "url": "https://dpm.demdex.net/ibs:dpid=411&dpuuid=YvYavQAF2aAqsAA0&d_uuid=78854490201948989640033682043832183304"
        },
        "success": true,
        "labels": [
            {
                "isLabeled": false,
                "blocklist": "EasyList"
            },
            {
                "isLabeled": true,
                "rule": [
                    [
                        "csp",
                        null
                    ],
                    [
                        "filter",
                        null
                    ],
                    [
                        "hostname",
                        "demdex.net"
                    ],
                    [
                        "mask",
                        335609855
                    ],
                    [
                        "domains",
                        null
                    ],
                    [
                        "denyallow",
                        null
                    ],
                    [
                        "redirect",
                        null
                    ],
                    [
                        "rawLine",
                        null
                    ],
                    [
                        "id",
                        3052061144
                    ],
                    [
                        "regex",
                        null
                    ]
                ],
                "blocklist": "EasyPrivacy"
            }
        ]
    }`

	// when
	resp, err := testsupport.Post("http://localhost:8081/requests", request, "application/json")

	// then
	suite.NoError(err)
	suite.Equal(http.StatusCreated, resp.StatusCode)
	count, err := suite.requestRepo.Count(suite.ctx)
	suite.NoError(err)
	suite.Equal(int64(1), count)
}
