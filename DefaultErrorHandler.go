package SimpleRouter

import (
	"net/http"
)

type DefaultHttpErrorHandler struct {
}

func (d DefaultHttpErrorHandler) MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}
