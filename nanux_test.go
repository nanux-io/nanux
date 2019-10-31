package nanux_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/nanux-io/nanux"
	. "github.com/nanux-io/nanux/handler"
)

/*----------------------------------------------------------------------------*\
	Define fake listener
\*----------------------------------------------------------------------------*/
type FakeListener struct {
	actions        map[string]ActionListener
	isListenCalled bool
	isCloseCalled  bool
}

func (f *FakeListener) Listen() error {
	f.isListenCalled = true
	return errors.New("Listen error")
}
func (f *FakeListener) Close() error {
	f.isCloseCalled = true
	return errors.New("Close error")
}
func (f *FakeListener) HandleAction(sub string, action ActionListener) error {
	if sub == "testError" {
		return errors.New("Error occured")
	}
	f.actions[sub] = action
	return nil
}

func (f *FakeListener) HandleError(errHandler ErrorHandler) error { return nil }

/*----------------------------------------------------------------------------*\
	Tests implementation
\*----------------------------------------------------------------------------*/
var _ = Describe("Nanux", func() {
	listener := &FakeListener{
		actions: make(map[string]ActionListener),
	}
	nanuxCtx := "Nanux context"
	nanuxInstance := New(listener, nanuxCtx)

	Context("created with New", func() {
		It("must contain the instance of the listener", func() {
			Expect(nanuxInstance.L).To(Equal(listener))
		})
	})

	Context("handle action error", func() {
		It("can be add before any other action", func() {
			// Add error manager
			errorHandler := func(error) []byte { return []byte("Error managed") }
			err := nanuxInstance.HandleError(errorHandler)

			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("handle action", func() {
		It("should provide the global context", func() {
			sub := "testRoute"
			actionFunc := func(ctx *interface{}, req Request) ([]byte, error) {
				res := *ctx
				return []byte(res.(string)), nil
			}

			nanuxInstance.Handle(sub, Action{Fn: actionFunc})
			actualRes, _ := listener.actions[sub].Fn(Request{})
			Expect(actualRes).To(Equal([]byte(nanuxCtx)))
		})

		It("should send error if the action not handled successfully by the listener", func() {
			sub := "testError"
			actionFunc := func(ctx *interface{}, req Request) ([]byte, error) {
				return nil, nil
			}

			err := nanuxInstance.Handle(sub, Action{Fn: actionFunc})
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

	Context("listen", func() {
		It("should call the listen method of the listener", func() {
			err := nanuxInstance.Listen()
			Expect(listener.isListenCalled).To(Equal(true))
			Expect(err).To(Equal(errors.New("Listen error")))
		})
	})

	Context("close", func() {
		It("should call the close method of the listener", func() {
			err := nanuxInstance.Close()
			Expect(listener.isCloseCalled).To(Equal(true))
			Expect(err).To(Equal(errors.New("Close error")))
		})
	})
})
