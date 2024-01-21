package extractor

var EXTRACTORS = InitExtractors()

func InitExtractors() []Extractor {
	// Create Extractor with dimensions [204,1]
	extractor204OR := NewExtractor("Extractor204OR", "This Extractor extracts a feature vector of [204,1] it includes 200 chars of the URL, the request type, the request method, the request frame_type and the presence of the referrer header")
	extractor204OR.URL(URL_EXTRACTOR)
	extractor204OR.FrameType(FRAME_TYPE_EXTRACTOR)
	extractor204OR.Method(METHOD_EXTRACTOR)
	extractor204OR.Type(TYPE_EXTRACTOR)
	extractor204OR.RequestHeaders(REQUEST_HEADER_REFERER_EXTRACTOR)
	extractor204OR.Labels(LABEL_EXTRACTOR_OR)

	extractor204EasyPrivacy := NewExtractor("Extractor204EasyPrivacy", "This Extractor extracts a feature vector of [204,1] it includes 200 chars of the URL, the request type, the request method, the request frame_type and the presence of the referrer header")
	extractor204EasyPrivacy.URL(URL_EXTRACTOR)
	extractor204EasyPrivacy.FrameType(FRAME_TYPE_EXTRACTOR)
	extractor204EasyPrivacy.Method(METHOD_EXTRACTOR)
	extractor204EasyPrivacy.Type(TYPE_EXTRACTOR)
	extractor204EasyPrivacy.RequestHeaders(REQUEST_HEADER_REFERER_EXTRACTOR)
	extractor204EasyPrivacy.Labels(LABEL_EXTRACTOR_EASY_PRIVACY)

	extractor204EasyList := NewExtractor("Extractor204EasyList", "This Extractor extracts a feature vector of [204,1] it includes 200 chars of the URL, the request type, the request method, the request frame_type and the presence of the referrer header")
	extractor204EasyList.URL(URL_EXTRACTOR)
	extractor204EasyList.FrameType(FRAME_TYPE_EXTRACTOR)
	extractor204EasyList.Method(METHOD_EXTRACTOR)
	extractor204EasyList.Type(TYPE_EXTRACTOR)
	extractor204EasyList.RequestHeaders(REQUEST_HEADER_REFERER_EXTRACTOR)
	extractor204EasyList.Labels(LABEL_EXTRACTOR_EASY_LIST)
	// Here you init new extractors

	return []Extractor{*extractor204OR, *extractor204EasyPrivacy, *extractor204EasyList}
}
