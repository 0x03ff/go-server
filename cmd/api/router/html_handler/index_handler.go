// router/html_handler/router.go
package html_handler

import (
	"net/http"

)


func (h *HtmlHandlers) IndexHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "web/html/index.html")
}