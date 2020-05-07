package util

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

// WriteStatusError writes err to c.HTTPResponseWriter().
func (c *Context) WriteStatusError(err StatusError) {
	http.Error(c.resp, err.HTTPResponseBody(), err.HTTPStatusCode())
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

// StatusError is implemented by error types which correspond to a particular
// HTTP status code.
type StatusError interface {
	error

	// HTTPStatusCode returns the HTTP status code for this error.
	HTTPStatusCode() int
	// HTTPResponseBody returns a string which will be used as the body of the
	// HTTP response.
	HTTPResponseBody() string
}

// InternalServerError wraps an error and implements StatusError. Its
// implementation of HTTPStatusCode returns http.StatusInternalServerError, and
// its implementation of HTTPResponseBody returns the empty string.
type InternalServerError struct{ error }

func (e InternalServerError) HTTPStatusCode() int {
	return http.StatusInternalServerError
}

func (e InternalServerError) HTTPResponseBody() string {
	// We don't want to leak any potentially sensitive data that might be
	// contained in the error.
	return ""
}

// BadRequestError wraps an error and implements StatusError. Its implementation
// of HTTPStatusCode returns http.StatusBadRequest, and its implementation of
// HTTPResponseBody returns the result of calling the Error method.
type BadRequestError struct{ error }

// NewBadRequestError wraps err in a BadRequestError.
func NewBadRequestError(err error) BadRequestError {
	return BadRequestError{err}
}

func (e BadRequestError) HTTPStatusCode() int {
	return http.StatusBadRequest
}

func (e BadRequestError) HTTPResponseBody() string {
	return e.Error()
}

var (
	notFoundError = BadRequestError{errors.New("not found")}
)

// FirestoreToStatusError converts an error returned from the
// "cloud.google.com/go/firestore" package to a StatusError.
func FirestoreToStatusError(err error) StatusError {
	if status.Code(err) == codes.NotFound {
		return notFoundError
	}

	return InternalServerError{err}
}

// JSONToStatusError converts an error returned from the "encoding/json" package
// to a StatusError. It assumes that all error types defined in the
// "encoding/json" package and io.EOF are bad request errors and all others are
// internal server errors.
func JSONToStatusError(err error) StatusError {
	switch err := err.(type) {
	case *json.MarshalerError, *json.SyntaxError, *json.UnmarshalFieldError,
		*json.UnmarshalTypeError, *json.UnsupportedTypeError, *json.UnsupportedValueError:
		return NewBadRequestError(err)
	default:
		if err == io.EOF {
			return NewBadRequestError(err)
		}
		return InternalServerError{err}
	}
}
