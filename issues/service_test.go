package issues_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/errskipower/find-issues/issues"
	"github.com/errskipower/find-issues/issues/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Issues Service", func() {
	var (
		service issues.Service
		client  *fakes.HttpClient
	)

	BeforeEach(func() {
		client = &fakes.HttpClient{}
		service = issues.NewService("repo-name", client)
	})

	Describe("Get", func() {
		BeforeEach(func() {
			client.GetCall.Returns.Response = &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`[
					{
						"number": 121,
						"title": "Title of Issue 121"
					},
					{
						"number": 123,
						"title": "Title of Issue 123"
					}
				]`)),
			}
		})

		It("returns the list of open issues", func() {
			issues, err := service.Get("", "")
			Expect(err).NotTo(HaveOccurred())

			Expect(client.GetCall.Receives.Url).To(Equal("https://api.github.com/repos/repo-name/issues"))

			Expect(issues).To(HaveLen(2))
			Expect(issues[0].Number).To(Equal(121))
			Expect(issues[1].Number).To(Equal(123))
		})

		Context("when we pass in a label", func() {
			It("generates the correct rawQuery string", func() {
				_, err := service.Get("help wanted", "")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.GetCall.Receives.Url).To(Equal("https://api.github.com/repos/repo-name/issues?labels=help+wanted"))
			})
		})

		Context("when we pass in a creator", func() {
			It("generates the correct rawQuery string", func() {
				_, err := service.Get("", "genevieve")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.GetCall.Receives.Url).To(Equal("https://api.github.com/repos/repo-name/issues?creator=genevieve"))
			})
		})

		Describe("error cases", func() {
			Context("when the client fails to get the issues endpoint", func() {
				BeforeEach(func() {
					client.GetCall.Returns.Error = errors.New("hai")
				})

				It("returns the error", func() {
					_, err := service.Get("", "")

					Expect(err).To(MatchError("hai"))
				})
			})

			Context("when the response status code is not 200", func() {
				BeforeEach(func() {
					client.GetCall.Returns.Response = &http.Response{
						StatusCode: 999,
						Body:       ioutil.NopCloser(bytes.NewBufferString("[]")),
					}
				})

				It("returns the error", func() {
					_, err := service.Get("", "")

					Expect(err).To(MatchError("invalid path https://api.github.com/repos/repo-name/issues"))
				})
			})

			Context("when the response body is invalid", func() {
				BeforeEach(func() {
					client.GetCall.Returns.Response = &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewBufferString("$$$")),
					}
				})

				It("returns the error", func() {
					_, err := service.Get("", "")

					Expect(err).To(MatchError("invalid character '$' looking for beginning of value"))
				})
			})
		})
	})
})
