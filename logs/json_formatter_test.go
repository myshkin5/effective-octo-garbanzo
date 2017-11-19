package logs_test

import (
	"encoding/json"
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"

	"github.com/myshkin5/effective-octo-garbanzo/logs"
)

var _ = Describe("Health", func() {
	It("Error not lost", func() {
		formatter := &logs.JSONFormatter{}

		b, err := formatter.Format(logrus.WithField("error", errors.New("wild walrus")))
		Expect(err).NotTo(HaveOccurred())

		entry := make(map[string]interface{})
		err = json.Unmarshal(b, &entry)
		Expect(err).NotTo(HaveOccurred())

		Expect(entry["error"]).To(Equal("wild walrus"))
	})

	It("Error not lost on field not named error", func() {
		formatter := &logs.JSONFormatter{}

		b, err := formatter.Format(logrus.WithField("omg", errors.New("wild walrus")))
		Expect(err).NotTo(HaveOccurred())

		entry := make(map[string]interface{})
		err = json.Unmarshal(b, &entry)
		Expect(err).NotTo(HaveOccurred())

		Expect(entry["omg"]).To(Equal("wild walrus"))
	})

	It("Field clash with time", func() {
		formatter := &logs.JSONFormatter{}

		b, err := formatter.Format(logrus.WithField("time", "right now!"))
		Expect(err).NotTo(HaveOccurred())

		entry := make(map[string]interface{})
		err = json.Unmarshal(b, &entry)
		Expect(err).NotTo(HaveOccurred())

		Expect(entry["fields.time"]).To(Equal("right now!"))

		Expect(entry["time"]).To(Equal("0001-01-01T00:00:00Z"))
	})

	It("Field clash with msg", func() {
		formatter := &logs.JSONFormatter{}

		b, err := formatter.Format(logrus.WithField("msg", "something"))
		Expect(err).NotTo(HaveOccurred())

		entry := make(map[string]interface{})
		err = json.Unmarshal(b, &entry)
		Expect(err).NotTo(HaveOccurred())

		Expect(entry["fields.msg"]).To(Equal("something"))
	})

	It("Field clash with level", func() {
		formatter := &logs.JSONFormatter{}

		b, err := formatter.Format(logrus.WithField("level", "something"))
		Expect(err).NotTo(HaveOccurred())

		entry := make(map[string]interface{})
		err = json.Unmarshal(b, &entry)
		Expect(err).NotTo(HaveOccurred())

		Expect(entry["fields.level"]).To(Equal("something"))
	})

	It("JSON entry ends with new line", func() {
		formatter := &logs.JSONFormatter{}

		b, err := formatter.Format(logrus.WithField("level", "something"))
		Expect(err).NotTo(HaveOccurred())

		Expect(b[len(b)-1]).To(BeEquivalentTo('\n'))
	})

	It("JSON message key", func() {
		formatter := &logs.JSONFormatter{
			FieldMap: logs.FieldMap{
				logrus.FieldKeyMsg: "message",
			},
		}

		b, err := formatter.Format(&logrus.Entry{Message: "oh hai"})
		Expect(err).NotTo(HaveOccurred())

		s := string(b)
		Expect(s).To(ContainSubstring("message"))
		Expect(s).To(ContainSubstring("oh hai"))
	})

	It("JSON level key", func() {
		formatter := &logs.JSONFormatter{
			FieldMap: logs.FieldMap{
				logrus.FieldKeyLevel: "somelevel",
			},
		}

		b, err := formatter.Format(logrus.WithField("level", "something"))
		Expect(err).NotTo(HaveOccurred())

		s := string(b)
		Expect(s).To(ContainSubstring("somelevel"))
	})

	It("JSON time key", func() {
		formatter := &logs.JSONFormatter{
			FieldMap: logs.FieldMap{
				logrus.FieldKeyTime: "timeywimey",
			},
		}

		b, err := formatter.Format(logrus.WithField("level", "something"))
		Expect(err).NotTo(HaveOccurred())

		s := string(b)
		Expect(s).To(ContainSubstring("timeywimey"))
	})

	It("JSON disable timestamp", func() {
		formatter := &logs.JSONFormatter{
			DisableTimestamp: true,
		}

		b, err := formatter.Format(logrus.WithField("level", "something"))
		Expect(err).NotTo(HaveOccurred())

		s := string(b)
		Expect(s).NotTo(ContainSubstring(logrus.FieldKeyTime))
	})

	It("JSON enable timestamp", func() {
		formatter := &logs.JSONFormatter{}

		b, err := formatter.Format(logrus.WithField("level", "something"))
		Expect(err).NotTo(HaveOccurred())

		s := string(b)
		Expect(s).To(ContainSubstring(logrus.FieldKeyTime))
	})

	It("JSON level map", func() {
		formatter := &logs.JSONFormatter{
			LevelMap: logs.LevelMap{
				logrus.InfoLevel: "INFO",
			},
		}

		b, err := formatter.Format(&logrus.Entry{Message: "oh hai", Level: logrus.InfoLevel})
		Expect(err).NotTo(HaveOccurred())

		s := string(b)
		Expect(s).To(ContainSubstring("\"level\":\"INFO\""))
	})
})
