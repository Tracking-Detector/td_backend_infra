package extractor

import (
	"errors"
)

func GetTypes() []string {
	return []string{
		"xmlhttprequest",
		"image",
		"font",
		"script",
		"stylesheet",
		"ping",
		"sub_frame",
		"other",
		"main_frame",
		"csp_report",
		"object",
		"media"}
}

func GetFrameTypes() []string {
	return []string{"outermost_frame", "fenced_frame", "sub_frame"}
}

func GetMethods() []string {
	return []string{"GET",
		"POST",
		"OPTIONS",
		"HEAD",
		"PUT",
		"DELETE",
		"SEARCH",
		"PATCH"}
}

func URL_EXTRACTOR(s string) ([]int, error) {
	if s == "" {
		return nil, errors.New("Url is not set")
	}
	encoded := make([]int, 200)
	count := 199
	for i := len(s) - 1; i >= 0; i-- {
		c := []rune(s)[i]
		encoded[count] = (int(c) % 89) + 1
		if count == 0 {
			break
		}
		count--
	}
	return encoded, nil
}

func FRAME_TYPE_EXTRACTOR(s string) ([]int, error) {
	if s == "" {
		return nil, errors.New("Frame_type is not set")
	}
	for i, val := range GetFrameTypes() {
		if val == s {
			return []int{i + 1}, nil
		}
	}
	return nil, errors.New("Unknown frame_type encounter")
}

func METHOD_EXTRACTOR(s string) ([]int, error) {
	if s == "" {
		return nil, errors.New("Method is not set")
	}
	for i, val := range GetMethods() {
		if val == s {
			return []int{i + 1}, nil
		}
	}
	return nil, errors.New("Unknown method encountered")
}

func TYPE_EXTRACTOR(s string) ([]int, error) {
	if s == "" {
		return nil, errors.New("Type is not set")
	}
	for i, val := range GetTypes() {
		if val == s {
			return []int{i + 1}, nil
		}
	}
	return nil, errors.New("Unknown type encountered")
}

func TRACKER_EXTRACTOR(tracker bool) ([]int, error) {
	if tracker {
		return []int{1}, nil
	}
	return []int{0}, nil
}

func REQUEST_HEADER_REFERER_EXTRACTOR(headers []map[string]string) ([]int, error) {
	if len(headers) == 0 {
		return nil, errors.New("Headers are not set")
	}
	for _, header := range headers {
		if val, exists := header["name"]; exists && val == "Referer" {
			return []int{1}, nil
		}
	}
	return []int{0}, nil
}
