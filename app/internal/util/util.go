package util

import (
	"context"
	"fmt"
	"net/http"

	"cloud.google.com/go/firestore"
)

// Context is a context.Context that provides extra utilities for common
// operations.
type Context struct {
	resp   http.ResponseWriter
	req    *http.Request
	client *firestore.Client

	context.Context
}

// NewContext constructs a new Context from an http.ResponseWriter and an
// *http.Request. If an error occurs, NewContext takes care of writing an
// appropriate response to w.
func NewContext(w http.ResponseWriter, r *http.Request) (Context, error) {
	ctx := r.Context()
	client, err := firestore.NewClient(ctx, "test")
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return Context{}, err
	}

	return Context{w, r, client, ctx}, nil
}

// HTTPRequest returns the *http.Request that was used to construct this
// Context.
func (c *Context) HTTPRequest() *http.Request {
	return c.req
}

// HTTPResponseWriter returns the http.ResponseWriter that was used to construct
// this Context.
func (c *Context) HTTPResponseWriter() http.ResponseWriter {
	return c.resp
}

func (c *Context) FirestoreClient() *firestore.Client {
	return c.client
}

// WriteInternalServerError writes an internal server error to
// c.HTTPResponseWriter().
func (c *Context) WriteInternalServerError() {
	http.Error(c.resp, "", http.StatusInternalServerError)
}

// ValidateRequestMethod validates that ctx.HTTPRequest().Method == method, and
// if not, writes an appropriate response to ctx.HTTPResponseWriter().
func ValidateRequestMethod(ctx *Context, method, err string) error {
	m := ctx.HTTPRequest().Method
	if m != method {
		http.Error(ctx.HTTPResponseWriter(), err, http.StatusMethodNotAllowed)
		return fmt.Errorf("unsupported method: %v", m)
	}

	return nil
}
