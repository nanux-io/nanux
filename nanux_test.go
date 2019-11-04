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
	handlers      map[string]THandler
	isRunCalled   bool
	isCloseCalled bool
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
	f.handlers[sub] = tHandler
	return nil
}

func (f *FakeTransporter) HandleError(errHandler ErrorHandler) error { return nil }

/*----------------------------------------------------------------------------*\
	Tests implementation
\*----------------------------------------------------------------------------*/
var _ = Describe("Nanux", func() {
	transporter := &FakeTransporter{
		handlers: make(map[string]THandler),
	}
	nanuxCtx := "Nanux context"
	nanuxInstance := New(transporter, nanuxCtx)

	Context("created with New", func() {
		It("must contain the instance of the listener", func() {
			Expect(nanuxInstance.T).To(Equal(transporter))
		})
	})

	Context("handle action error", func() {
		It("can be add before any other action", func() {
			// Add error manager
			errorHandler := func(error, Request) []byte { return []byte("Error managed") }
			err := nanuxInstance.HandleError(errorHandler)

			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("handle action", func() {
		It("should provide the global context", func() {
			sub := "testRoute"
			handlerFn := func(ctx *interface{}, req Request) ([]byte, error) {
				res := *ctx
				return []byte(res.(string)), nil
			}

			nanuxInstance.Handle(sub, Handler{Fn: handlerFn})
			actualRes, _ := transporter.handlers[sub].Fn(Request{})
			Expect(actualRes).To(Equal([]byte(nanuxCtx)))
		})

		It("should send error if the action not handled successfully by the listener", func() {
			sub := "testError"
			handlerFn := func(ctx *interface{}, req Request) ([]byte, error) {
				return nil, nil
			}

			err := nanuxInstance.Handle(sub, Handler{Fn: handlerFn})
			Expect(err).To(Equal(errors.New("Error occured")))
		})
	})

	// Context("handle action error", func() {
	// 	It("should not be added after other actions", func() {
	// 		// Add error manager
	// 		manageError := func(error) []byte { return []byte("Error managed") }
	// 		err := nanuxInstance.HandleError(manageError)

	// 		Expect(err).To(HaveOccurred())
	// 	})
	// })

	Context("run", func() {
		It("should call the run method of the transporter", func() {
			err := nanuxInstance.Run()
			Expect(transporter.isRunCalled).To(Equal(true))
			Expect(err).To(Equal(errors.New("Run error")))
		})
	})

	Context("close", func() {
		It("should call the close method of the transporter", func() {
			err := nanuxInstance.Close()
			Expect(transporter.isCloseCalled).To(Equal(true))
			Expect(err).To(Equal(errors.New("Close error")))
		})
	})
})
