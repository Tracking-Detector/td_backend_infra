var FRAME_TYPES = ["outermost_frame", "fenced_frame", "sub_frame"]
var METHODS = ["GET", "POST", "OPTIONS", "HEAD", "PUT", "DELETE", "SEARCH", "PATCH"]
var TYPES = ["xmlhttprequest", "image", "font", "script", "stylesheet", "ping", "sub_frame", "other", "main_frame", "csp_report", "object", "media"]

function extractUrl(s) {
    if (s === "") {
        throw new Error("Url is not set");
    }

    var encoded = new Array(200)
    for (var i = 0; i < 200;i++) {
        encoded[i] = 0
    }
    var count = 199;

    for (var i = s.length - 1; i >= 0; i--) {
        var c = s.charCodeAt(i);
        encoded[count] = (c % 89) + 1;

        if (count === 0) {
            break;
        }

        count--;
    }

    return encoded;
}

function extractFrameTypes(s) {
    if (s === "") {
        throw new Error("Frame_type is not set");
    }

    for (var i = 0; i < FRAME_TYPES.length; i++) {
        if (FRAME_TYPES[i] === s) {
            return [i + 1];
        }
    }

    throw new Error("Unknown frame_type encountered");
}

function methodExtractor(s) {
    if (s === "") {
        throw new Error("Method is not set");
    }

    for (var i = 0; i < METHODS.length; i++) {
        if (METHODS[i] === s) {
            return [i + 1];
        }
    }

    throw new Error("Unknown Method encountered");
}

function typeExtractor(s) {
    if (s === "") {
        throw new Error("Type is not set");
    }

    for (var i = 0; i < TYPES.length; i++) {
        if (TYPES[i] === s) {
            return [i + 1];
        }
    }

    throw new Error("Unknown Type encountered");
}

function requestHeaderRefererExtractor(headers) {
    if (headers.length === 0) {
        throw new Error("Headers are not set");
    }

    for (var i = 0; i < headers.length; i++) {
        if (headers[i].name === "Referer") {
            return [1];
        }
    }

    return [0];
}

function extract(request) {
    encodedUrl = extractUrl(request.url)
    encodedFrametype = extractFrameTypes(request.frameType)
    encodedMethod = methodExtractor(request.method)
    encodedType = typeExtractor(request.type)
    encodedRequestHeader = requestHeaderRefererExtractor(request.requestHeaders)
    encodedTracker = request.tracker ? [1] : [0]
    return encodedUrl.concat(encodedFrametype, encodedMethod, encodedType, encodedRequestHeader, encodedTracker)
}
