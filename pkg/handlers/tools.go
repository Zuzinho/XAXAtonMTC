package handlers

import (
	"fmt"
	"net/http"
	"strconv"
)

func HttpJSONErr(w http.ResponseWriter, err error, status int) {
	http.Error(w, fmt.Sprintf("{\"error\": \"%s\"}", err.Error()), status)
}

func UInt32FromMapString(vars map[string]string, key string) (uint32, error) {
	idStr := vars[key]
	if idStr == "" {
		return 0, NoUserIDParamErr
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, InvalidRequestParamsErr
	}

	return uint32(id), nil
}
