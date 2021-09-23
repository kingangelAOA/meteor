package core

import "path/filepath"

const (
	METHOD_GET     = "GET"
	METHOD_POST    = "POST"
	METHOD_PUT     = "PUT"
	METHOD_PATCH   = "PATCH"
	METHOD_DELETE  = "DELETE"
	METHOD_OPTIONS = "OPTIONS"
	METHOD_HEAD    = ""
)
const (
	AUTH_NONE    = "none"
	AUTH_OAUTH_2 = "oauth2"
	AUTH_OAUTH_1 = "oauth1"
	AUTH_BASIC   = "basic"
	AUTH_DIGEST  = "digest"
	AUTH_BEARER  = "bearer"
	AUTH_NTLM    = "ntlm"
	AUTH_HAWK    = "hawk"
	AUTH_AWS_IAM = "iam"
	AUTH_NETRC   = "netrc"
	AUTH_ASAP    = ""
)
const (
	CONTENT_TYPE_JSON            = "application/json"
	CONTENT_TYPE_XML             = "application/xml"
	CONTENT_TYPE_YAML            = "text/yaml"
	CONTENT_TYPE_EDN             = "application/edn"
	CONTENT_TYPE_FORM_URLENCODED = "application/x-www-form-urlencoded"
	CONTENT_TYPE_FORM_DATA       = "multipart/form-data"
	CONTENT_TYPE_FILE            = "application/octet-stream"
	CONTENT_TYPE_GRAPHQL         = "application/graphql"
	CONTENT_TYPE_OTHER           = ""
)

const (
	MULTIPART_TYPE_FILE            = "file"
	MULTIPART_TYPE_TEXT            = "text"
	MULTIPART_TYPE_MULTI_LINE_TEXT = "multi_line_text"
)

const (
	HTTP_COMPONENT = "HTTP_COMPONENT"
	POST_PROCESSOR = "POST_PROCESSOR"
)

var ContentTypeMap = map[string]string{
	".a":       "application/octet-stream",
	".ai":      "application/postscript",
	".aif":     "audio/x-aiff",
	".aifc":    "audio/x-aiff",
	".aiff":    "audio/x-aiff",
	".au":      "audio/basic",
	".avi":     "video/x-msvideo",
	".bat":     "text/plain",
	".bcpio":   "application/x-bcpio",
	".bin":     "application/octet-stream",
	".bmp":     "image/x-ms-bmp",
	".c":       "text/plain",
	".cdf":     "application/x-cdf",
	".cpio":    "application/x-cpio",
	".csh":     "application/x-csh",
	".css":     "text/css",
	".csv":     "text/csv",
	".dll":     "application/octet-stream",
	".doc":     "application/msword",
	".dot":     "application/msword",
	".dvi":     "application/x-dvi",
	".eml":     "message/rfc822",
	".eps":     "application/postscript",
	".etx":     "text/x-setext",
	".exe":     "application/octet-stream",
	".gif":     "image/gif",
	".gtar":    "application/x-gtar",
	".h":       "text/plain",
	".hdf":     "application/x-hdf",
	".htm":     "text/html",
	".html":    "text/html",
	".ico":     "image/vnd.microsoft.icon",
	".ief":     "image/ief",
	".jpe":     "image/jpeg",
	".jpeg":    "image/jpeg",
	".jpg":     "image/jpeg",
	".js":      "application/javascript",
	".ksh":     "text/plain",
	".latex":   "application/x-latex",
	".m1v":     "video/mpeg",
	".m3u":     "application/vnd.apple.mpegurl",
	".m3u8":    "application/vnd.apple.mpegurl",
	".man":     "application/x-troff-man",
	".me":      "application/x-troff-me",
	".mht":     "message/rfc822",
	".mhtml":   "message/rfc822",
	".mif":     "application/x-mif",
	".mov":     "video/quicktime",
	".movie":   "video/x-sgi-movie",
	".mp2":     "audio/mpeg",
	".mp3":     "audio/mpeg",
	".mp4":     "video/mp4",
	".mpa":     "video/mpeg",
	".mpe":     "video/mpeg",
	".mpeg":    "video/mpeg",
	".mpg":     "video/mpeg",
	".ms":      "application/x-troff-ms",
	".nc":      "application/x-netcdf",
	".nws":     "message/rfc822",
	".o":       "application/octet-stream",
	".obj":     "application/octet-stream",
	".oda":     "application/oda",
	".p12":     "application/x-pkcs12",
	".p7c":     "application/pkcs7-mime",
	".pbm":     "image/x-portable-bitmap",
	".pdf":     "application/pdf",
	".pfx":     "application/x-pkcs12",
	".pgm":     "image/x-portable-graymap",
	".pl":      "text/plain",
	".png":     "image/png",
	".pnm":     "image/x-portable-anymap",
	".pot":     "application/vnd.ms-powerpoint",
	".ppa":     "application/vnd.ms-powerpoint",
	".ppm":     "image/x-portable-pixmap",
	".pps":     "application/vnd.ms-powerpoint",
	".ppt":     "application/vnd.ms-powerpoint",
	".ps":      "application/postscript",
	".pwz":     "application/vnd.ms-powerpoint",
	".py":      "text/x-python",
	".pyc":     "application/x-python-code",
	".pyo":     "application/x-python-code",
	".qt":      "video/quicktime",
	".ra":      "audio/x-pn-realaudio",
	".ram":     "application/x-pn-realaudio",
	".ras":     "image/x-cmu-raster",
	".rdf":     "application/xml",
	".rgb":     "image/x-rgb",
	".roff":    "application/x-troff",
	".rtx":     "text/richtext",
	".sgm":     "text/x-sgml",
	".sgml":    "text/x-sgml",
	".sh":      "application/x-sh",
	".shar":    "application/x-shar",
	".snd":     "audio/basic",
	".so":      "application/octet-stream",
	".src":     "application/x-wais-source",
	".sv4cpio": "application/x-sv4cpio",
	".sv4crc":  "application/x-sv4crc",
	".svg":     "image/svg+xml",
	".swf":     "application/x-shockwave-flash",
	".t":       "application/x-troff",
	".tar":     "application/x-tar",
	".tcl":     "application/x-tcl",
	".tex":     "application/x-tex",
	".texi":    "application/x-texinfo",
	".texinfo": "application/x-texinfo",
	".tif":     "image/tiff",
	".tiff":    "image/tiff",
	".tr":      "application/x-troff",
	".tsv":     "text/tab-separated-values",
	".txt":     "text/plain",
	".ustar":   "application/x-ustar",
	".vcf":     "text/x-vcard",
	".wav":     "audio/x-wav",
	".webm":    "video/webm",
	".wiz":     "application/msword",
	".wsdl":    "application/xml",
	".xbm":     "image/x-xbitmap",
	".xlb":     "application/vnd.ms-excel",
	".xls":     "application/excel",
	".xml":     "text/xml",
	".xpdl":    "application/xml",
	".xpm":     "image/x-xpixmap",
	".xsl":     "application/xml",
	".xwd":     "image/x-xwindowdump",
	".zip":     "application/zip",
}

// GetContentTypeByFilename 通过文件名称获取类型
func GetContentTypeByFilename(filename string) (string, error) {
	var suffix = filepath.Ext(filename)
	if v, ok := ContentTypeMap[suffix]; ok {
		return v, nil
	} else {
		return "application/octet-stream", nil
	}
}
