package nanux_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/nanux-io/nanux"
)

/*----------------------------------------------------------------------------*\
	Define fake listener
\*----------------------------------------------------------------------------*/
type FakeTransporter struct {
	tHandlers           map[string]THandler
	errHandler          ErrorHandler
	isRunCalled         bool
	isCloseCalled       bool
	isHandleErrorCalled bool
}

func (f *FakeTransporter) Run() error {
	f.isRunCalled = true
	return errors.New("Run error")
}
func (f *FakeTransporter) Close() error {
	f.isCloseCalled = true
	return errors.New("Close error")
}
func (f *FakeTransporter) Handle(sub string, tHandler THandler) error {
	if sub == "testError" {
		return errors.New("Error occured")
	}
	f.tHandlers[sub] = tHandler
	return nil
}

func (f *FakeTransporter) HandleError(errHandler ErrorHandler) error {
	f.errHandler = errHandler
	f.isHandleErrorCalled = true
	return nil
}

/*----------------------------------------------------------------------------*\
	Tests implementation
\*----------------------------------------------------------------------------*/
var _ = Describe("Nanux", func() {
	var nanuxCtx string
	var transporter *FakeTransporter

	BeforeEach(func() {
		transporter = &FakeTransporter{
			tHandlers: make(map[string]THandler),
		}
		nanuxCtx = "Nanux context"
	})

	It("should create a new instance of Nanux", func() {
		n := New(transporter, nanuxCtx)
		Expect(n).To(BeAssignableToTypeOf(&Nanux{}))
	})

	Context("instance", func() {
		var n *Nanux

		JustBeforeEach(func() {
			n = New(transporter, nanuxCtx)
		})

		It("should contain the instance of the transporter", func() {
			Expect(n.T).To(Equal(transporter))
		})

		It("should add an error handler", func() {
			errorHandler := func(error, Request) []byte { return []byte("Error managed") }
			err := n.HandleError(errorHandler)

			Expect(err).NotTo(HaveOccurred())
			Expect(transporter.isHandleErrorCalled).To(Equal(true))
		})

		Context("when handle is called", func() {
			It("should provide the global context", func() {
				sub := "testRoute"
				handlerFn := func(ctx *interface{}, req Request) ([]byte, error) {
					res := *ctx
					return []byte(res.(string)), nil
				}

				n.Handle(sub, Handler{Fn: handlerFn})
				actualRes, _ := transporter.tHandlers[sub].Fn(Request{})
				Expect(actualRes).To(Equal([]byte(nanuxCtx)))
			})

			It("should send error if the action not handled successfully by the listener", func() {
				sub := "testError"
				handlerFn := func(ctx *interface{}, req Request) ([]byte, error) {
					return nil, nil
				}

				err := n.Handle(sub, Handler{Fn: handlerFn})
				Expect(err).To(Equal(errors.New("Error occured")))
			})

			Context("when middlewares are provided", func() {
				var handler Handler
				route := "mwRoute"
				handlerMsg := "Response from the handler"

				JustBeforeEach(func() {
					handler = Handler{
						Fn: func(ctx *interface{}, req Request) ([]byte, error) {
							return []byte(handlerMsg), nil
						},
					}
				})

				It("should chain the middlewares and forward the req and call the handler", func() {
					mw1 := func(fn HandlerFunc) HandlerFunc {
						return func(ctx *interface{}, req Request) (response []byte, err error) {
							req.M["mw"] = 1
							return fn(ctx, req)
						}
					}
					mw2 := func(fn HandlerFunc) HandlerFunc {
						return func(ctx *interface{}, req Request) (response []byte, err error) {
							req.M["mw"] = 2
							return fn(ctx, req)
						}
					}

					n.Handle(route, handler, mw1, mw2)

					req := &Request{M: make(map[string]interface{})}
					res, err := transporter.tHandlers[route].Fn(*req)

					Expect(err).ToNot(HaveOccurred())
					Expect(res).To(Equal([]byte(handlerMsg)))

					Expect(req.M["mw"].(int)).To(Equal(2))
				})
			})
		})

		Context("run", func() {
			It("should call the run method of the transporter", func() {
				err := n.Run()
				Expect(transporter.isRunCalled).To(Equal(true))
				Expect(err).To(Equal(errors.New("Run error")))
			})
		})

		Context("close", func() {
			It("should call the close method of the transporter", func() {
				err := n.Close()
				Expect(transporter.isCloseCalled).To(Equal(true))
				Expect(err).To(Equal(errors.New("Close error")))
			})
		})
	})
})
