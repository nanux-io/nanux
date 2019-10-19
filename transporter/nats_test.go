package transporter_test

import (
	"errors"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/nats-io/go-nats"

	"github.com/nanux-io/nanux/handler"
	. "github.com/nanux-io/nanux/transporter"
)

var _ = Describe("Nats transporter", func() {
	Context("created with NewNats", func() {
		nt := NewNats(natsUrls, []nats.Option{})
		res := []byte("action called")
		sub := "testSub"

		It("should satisfy Listener interface", func() {
			var i interface{} = &nt
			_, ok := i.(Listener)
			Expect(ok).To(Equal(true))
		})

		It("should allow to add action before listening", func() {
			fn := func(req handler.Request) ([]byte, error) {
				return res, nil
			}

			action := handler.ListenerAction{Fn: fn}

			err := nt.HandleAction(sub, action)
			Expect(err).ShouldNot(HaveOccurred())

		})

		It("should allow action with queued option", func() {
			fn := func(req handler.Request) ([]byte, error) {
				return res, nil
			}

			action := handler.ListenerAction{
				Fn:   fn,
				Opts: []handler.Opt{{Name: NatsOptIsQueued, Value: true}},
			}

			err := nt.HandleAction("queuedsub", action)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should throw error if adding several actions for same subject", func() {
			fn := func(req handler.Request) ([]byte, error) { return nil, nil }

			action := handler.ListenerAction{Fn: fn}

			nt.HandleAction("sameSub", action)
			err := nt.HandleAction("sameSub", action)
			Expect(err).To(BeAssignableToTypeOf(errors.New("")))
		})

		It("should connect to nats server and close connection", func() {
			go func() {
				err := nt.Listen()
				Expect(err).ShouldNot(HaveOccurred())
			}()

			// FIXME: check directly with connection status to know when it is connected
			// Wait for connection to be established
			time.Sleep(100 * time.Millisecond)

			err := nt.Close()
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should respond to subscribed topic", func() {
			// Set nats into listening
			go func() {
				err := nt.Listen()
				Expect(err).ShouldNot(HaveOccurred())
			}()

			// FIXME: check directly with connection status to know when it is connected
			// Wait for connection to be established
			time.Sleep(100 * time.Millisecond)

			msg, err := natsClient.Request(sub, nil, time.Second)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(msg.Data).To(Equal(res))

			// Close nats listening
			nt.Close()
		})

		It("should have error when closing unconnected connection", func() {
			err := nt.Close()

			Expect(err.Error()).To(Equal("Can not close nats connection because the status is not connected"))
		})
	})
})
