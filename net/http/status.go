package http

type Status int

const (
	StatusOK                    = Status(200)
	StatusCreated               = 201
	StatusSeeOther              = 303
	StatusBadRequest            = 400
	StatusUnauthorized          = 401
	StatusForbidden             = 403
	StatusNotFound              = 404
	StatusMethodNotAllowed      = 405
	StatusRequestTimeout        = 408
	StatusConflict              = 409
	StatusRequestEntityTooLarge = 413
	StatusInternalServerError   = 500
	StatusServiceUnavailable    = 503
)

var Status2String = [...]string{
	0:                           "200",
	StatusOK:                    "200",
	StatusCreated:               "201",
	StatusSeeOther:              "303",
	StatusBadRequest:            "400",
	StatusUnauthorized:          "401",
	StatusForbidden:             "403",
	StatusNotFound:              "404",
	StatusMethodNotAllowed:      "405",
	StatusRequestTimeout:        "408",
	StatusConflict:              "409",
	StatusRequestEntityTooLarge: "413",
	StatusInternalServerError:   "500",
	StatusServiceUnavailable:    "503",
}

var Status2Reason = [...]string{
	0:                           "OK",
	StatusOK:                    "OK",
	StatusCreated:               "Created",
	StatusSeeOther:              "See Other",
	StatusBadRequest:            "Bad Request",
	StatusUnauthorized:          "Unauthorized",
	StatusForbidden:             "Forbidden",
	StatusNotFound:              "Not Found",
	StatusMethodNotAllowed:      "Method Not Allowed",
	StatusRequestTimeout:        "Request Timeout",
	StatusConflict:              "Conflict",
	StatusRequestEntityTooLarge: "Request Entity Too Large",
	StatusInternalServerError:   "Internal Server Error",
	StatusServiceUnavailable:    "Service Unavailable",
}

var StatusLines [VersionCount][]string

func (s Status) String() string {
	return Status2String[s]
}

func init() {
	for i := VersionNone; i <= Version11; i++ {
		var row []string

		for j := 0; j < len(Status2String); j++ {
			if len(Status2String[j]) == 0 {
				row = append(row, "")
			} else {
				row = append(row, Version2String[i]+" "+Status2String[j]+" "+Status2Reason[j]+"\r\n")
			}
		}

		StatusLines[i] = row
	}
}
