package html_handler

import (
	"net/http"

)


func (h *HtmlHandlers) RegisterHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "web/html/register.html")
}
