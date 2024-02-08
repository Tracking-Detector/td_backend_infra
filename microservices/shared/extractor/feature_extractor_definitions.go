package extractor

var EXTRACTORS = InitExtractors()

func InitExtractors() []Extractor {
	// Create Extractor with dimensions [204,1]
	extractor204 := NewExtractor("GoExtractor204",
		`This Extractor extracts a feature vector of [204,1].
	This vectorspace is build by the last 200 encoded chars of the request url,
	the request type, the request method, the request frame_type and
	the presence of the referrer header`, []int{204, 1})
	extractor204.URL(URL_EXTRACTOR)
	extractor204.FrameType(FRAME_TYPE_EXTRACTOR)
	extractor204.Method(METHOD_EXTRACTOR)
	extractor204.Type(TYPE_EXTRACTOR)
	extractor204.RequestHeaders(REQUEST_HEADER_REFERER_EXTRACTOR)
	extractor204.Tracker(TRACKER_EXTRACTOR)

	return []Extractor{*extractor204}
}
