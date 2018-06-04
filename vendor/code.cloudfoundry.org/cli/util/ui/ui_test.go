package ui_test

import (
	"errors"
	"time"

	"code.cloudfoundry.org/cli/command/translatableerror/translatableerrorfakes"
	"code.cloudfoundry.org/cli/util/configv3"
	. "code.cloudfoundry.org/cli/util/ui"
	"code.cloudfoundry.org/cli/util/ui/uifakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
)

var _ = Describe("UI", func() {
	var (
		ui         *UI
		fakeConfig *uifakes.FakeConfig
		out        *Buffer
		errBuff    *Buffer
	)

	BeforeEach(func() {
		fakeConfig = new(uifakes.FakeConfig)
		fakeConfig.ColorEnabledReturns(configv3.ColorEnabled)

		var err error
		ui, err = NewUI(fakeConfig)
		Expect(err).NotTo(HaveOccurred())

		out = NewBuffer()
		ui.Out = out
		ui.OutForInteration = out
		errBuff = NewBuffer()
		ui.Err = errBuff
	})

	Describe("DisplayError", func() {
		Context("when passed a TranslatableError", func() {
			var fakeTranslateErr *translatableerrorfakes.FakeTranslatableError

			BeforeEach(func() {
				fakeTranslateErr = new(translatableerrorfakes.FakeTranslatableError)
				fakeTranslateErr.TranslateReturns("I am an error")

				ui.DisplayError(fakeTranslateErr)
			})

			It("displays the error to ui.Err and displays FAILED in bold red to ui.Out", func() {
				Expect(ui.Err).To(Say("I am an error\n"))
				Expect(out).To(Say("\x1b\\[31;1mFAILED\x1b\\[0m\n"))
			})

			Context("when the locale is not set to english", func() {
				It("translates the error text", func() {
					Expect(fakeTranslateErr.TranslateCallCount()).To(Equal(1))
					Expect(fakeTranslateErr.TranslateArgsForCall(0)).NotTo(BeNil())
				})
			})
		})

		Context("when passed a generic error", func() {
			It("displays the error text to ui.Err and displays FAILED in bold red to ui.Out", func() {
				ui.DisplayError(errors.New("I am a BANANA!"))
				Expect(ui.Err).To(Say("I am a BANANA!\n"))
				Expect(out).To(Say("\x1b\\[31;1mFAILED\x1b\\[0m\n"))
			})
		})
	})

	Describe("DisplayHeader", func() {
		It("displays the header colorized and bolded to ui.Out", func() {
			ui.DisplayHeader("some-header")
			Expect(out).To(Say("\x1b\\[1msome-header\x1b\\[0m"))
		})

		Context("when the locale is not set to English", func() {
			BeforeEach(func() {
				fakeConfig.LocaleReturns("fr-FR")

				var err error
				ui, err = NewUI(fakeConfig)
				Expect(err).NotTo(HaveOccurred())

				ui.Out = out
			})

			It("displays the translated header colorized and bolded to ui.Out", func() {
				ui.DisplayHeader("FEATURE FLAGS")
				Expect(out).To(Say("\x1b\\[1mINDICATEURS DE FONCTION\x1b\\[0m"))
			})
		})
	})

	Describe("DisplayNewline", func() {
		It("displays a new line", func() {
			ui.DisplayNewline()
			Expect(out).To(Say("\n"))
		})
	})

	Describe("DisplayOK", func() {
		It("displays 'OK' in green and bold", func() {
			ui.DisplayOK()
			Expect(out).To(Say("\x1b\\[32;1mOK\x1b\\[0m"))
		})
	})

	// Covers the happy paths, additional cases are tested in TranslateText
	Describe("DisplayText", func() {
		It("displays the template with map values substituted in to ui.Out with a newline", func() {
			ui.DisplayText(
				"template with {{.SomeMapValue}}",
				map[string]interface{}{
					"SomeMapValue": "map-value",
				})
			Expect(out).To(Say("template with map-value\n"))
		})

		Context("when the locale is not set to english", func() {
			BeforeEach(func() {
				fakeConfig.LocaleReturns("fr-FR")

				var err error
				ui, err = NewUI(fakeConfig)
				Expect(err).NotTo(HaveOccurred())

				ui.Out = out
			})

			It("displays the translated template with map values substituted in to ui.Out", func() {
				ui.DisplayText(
					"\nTIP: Use '{{.Command}}' to target new org",
					map[string]interface{}{
						"Command": "foo",
					})
				Expect(out).To(Say("\nASTUCE : utilisez 'foo' pour cibler une nouvelle organisation"))
			})
		})
	})

	Describe("DisplayTextWithBold", func() {
		It("displays the template to ui.Out", func() {
			ui.DisplayTextWithBold("some-template")
			Expect(out).To(Say("some-template"))
		})

		Context("when an optional map is passed in", func() {
			It("displays the template with map values bolded and substituted in to ui.Out", func() {
				ui.DisplayTextWithBold(
					"template with {{.SomeMapValue}}",
					map[string]interface{}{
						"SomeMapValue": "map-value",
					})
				Expect(out).To(Say("template with \x1b\\[1mmap-value\x1b\\[0m"))
			})
		})

		Context("when multiple optional maps are passed in", func() {
			It("displays the template with only the first map values bolded and substituted in to ui.Out", func() {
				ui.DisplayTextWithBold(
					"template with {{.SomeMapValue}} and {{.SomeOtherMapValue}}",
					map[string]interface{}{
						"SomeMapValue": "map-value",
					},
					map[string]interface{}{
						"SomeOtherMapValue": "other-map-value",
					})
				Expect(out).To(Say("template with \x1b\\[1mmap-value\x1b\\[0m and <no value>"))
			})
		})

		Context("when the locale is not set to english", func() {
			BeforeEach(func() {
				fakeConfig.LocaleReturns("fr-FR")

				var err error
				ui, err = NewUI(fakeConfig)
				Expect(err).NotTo(HaveOccurred())

				ui.Out = out
			})

			It("displays the translated template with map values bolded and substituted in to ui.Out", func() {
				ui.DisplayTextWithBold(
					"App {{.AppName}} does not exist.",
					map[string]interface{}{
						"AppName": "some-app-name",
					})
				Expect(out).To(Say("L'application \x1b\\[1msome-app-name\x1b\\[0m n'existe pas.\n"))
			})
		})
	})

	Describe("DisplayTextWithFlavor", func() {
		It("displays the template to ui.Out", func() {
			ui.DisplayTextWithFlavor("some-template")
			Expect(out).To(Say("some-template"))
		})

		Context("when an optional map is passed in", func() {
			It("displays the template with map values colorized, bolded, and substituted in to ui.Out", func() {
				ui.DisplayTextWithFlavor(
					"template with {{.SomeMapValue}}",
					map[string]interface{}{
						"SomeMapValue": "map-value",
					})
				Expect(out).To(Say("template with \x1b\\[36;1mmap-value\x1b\\[0m"))
			})
		})

		Context("when multiple optional maps are passed in", func() {
			It("displays the template with only the first map values colorized, bolded, and substituted in to ui.Out", func() {
				ui.DisplayTextWithFlavor(
					"template with {{.SomeMapValue}} and {{.SomeOtherMapValue}}",
					map[string]interface{}{
						"SomeMapValue": "map-value",
					},
					map[string]interface{}{
						"SomeOtherMapValue": "other-map-value",
					})
				Expect(out).To(Say("template with \x1b\\[36;1mmap-value\x1b\\[0m and <no value>"))
			})
		})

		Context("when the locale is not set to english", func() {
			BeforeEach(func() {
				fakeConfig.LocaleReturns("fr-FR")

				var err error
				ui, err = NewUI(fakeConfig)
				Expect(err).NotTo(HaveOccurred())

				ui.Out = out
			})

			It("displays the translated template with map values colorized, bolded and substituted in to ui.Out", func() {
				ui.DisplayTextWithFlavor(
					"App {{.AppName}} does not exist.",
					map[string]interface{}{
						"AppName": "some-app-name",
					})
				Expect(out).To(Say("L'application \x1b\\[36;1msome-app-name\x1b\\[0m n'existe pas.\n"))
			})
		})
	})

	// Covers the happy paths, additional cases are tested in TranslateText
	Describe("DisplayWarning", func() {
		It("displays the warning to ui.Err", func() {
			ui.DisplayWarning(
				"template with {{.SomeMapValue}}",
				map[string]interface{}{
					"SomeMapValue": "map-value",
				})
			Expect(ui.Err).To(Say("template with map-value\n\n"))
		})

		Context("when the locale is not set to english", func() {
			BeforeEach(func() {
				fakeConfig.LocaleReturns("fr-FR")

				var err error
				ui, err = NewUI(fakeConfig)
				Expect(err).NotTo(HaveOccurred())

				ui.Err = NewBuffer()
			})

			It("displays the translated warning to ui.Err", func() {
				ui.DisplayWarning(
					"'{{.VersionShort}}' and '{{.VersionLong}}' are also accepted.",
					map[string]interface{}{
						"VersionShort": "some-value",
						"VersionLong":  "some-other-value",
					})
				Expect(ui.Err).To(Say("'some-value' et 'some-other-value' sont également acceptés.\n"))
			})
		})
	})

	// Covers the happy paths, additional cases are tested in TranslateText
	Describe("DisplayWarnings", func() {
		It("displays the warnings to ui.Err", func() {
			ui.DisplayWarnings([]string{"warning-1", "warning-2"})
			Expect(ui.Err).To(Say("warning-1\n"))
			Expect(ui.Err).To(Say("warning-2\n"))
			Expect(ui.Err).To(Say("\n"))
		})

		Context("when the locale is not set to english", func() {
			BeforeEach(func() {
				fakeConfig.LocaleReturns("fr-FR")

				var err error
				ui, err = NewUI(fakeConfig)
				Expect(err).NotTo(HaveOccurred())

				ui.Err = NewBuffer()
			})

			It("displays the translated warnings to ui.Err", func() {
				ui.DisplayWarnings([]string{"Also delete any mapped routes", "FEATURE FLAGS"})
				Expect(ui.Err).To(Say("Supprimer aussi les routes mappées\n"))
				Expect(ui.Err).To(Say("INDICATEURS DE FONCTION\n"))
				Expect(ui.Err).To(Say("\n"))
			})
		})

		Context("does not display newline when warnings are empty", func() {
			It("does not print out a new line", func() {
				ui.DisplayWarnings(nil)
				Expect(errBuff.Contents()).To(BeEmpty())
			})
		})
	})

	Describe("TranslateText", func() {
		It("returns the template", func() {
			Expect(ui.TranslateText("some-template")).To(Equal("some-template"))
		})

		Context("when an optional map is passed in", func() {
			It("returns the template with map values substituted in", func() {
				expected := ui.TranslateText(
					"template {{.SomeMapValue}}",
					map[string]interface{}{
						"SomeMapValue": "map-value",
					})
				Expect(expected).To(Equal("template map-value"))
			})
		})

		Context("when multiple optional maps are passed in", func() {
			It("returns the template with only the first map values substituted in", func() {
				expected := ui.TranslateText(
					"template with {{.SomeMapValue}} and {{.SomeOtherMapValue}}",
					map[string]interface{}{
						"SomeMapValue": "map-value",
					},
					map[string]interface{}{
						"SomeOtherMapValue": "other-map-value",
					})
				Expect(expected).To(Equal("template with map-value and <no value>"))
			})
		})

		Context("when the locale is not set to english", func() {
			BeforeEach(func() {
				fakeConfig.LocaleReturns("fr-FR")

				var err error
				ui, err = NewUI(fakeConfig)
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns the translated template", func() {
				expected := ui.TranslateText("   View allowable quotas with 'CF_NAME quotas'")
				Expect(expected).To(Equal("   Affichez les quotas pouvant être alloués avec 'CF_NAME quotas'"))
			})
		})
	})

	Describe("UserFriendlyDate", func() {
		It("formats a time into an ISO8601 string", func() {
			Expect(ui.UserFriendlyDate(time.Unix(0, 0))).To(MatchRegexp("\\w{3} [0-3]\\d \\w{3} [0-2]\\d:[0-5]\\d:[0-5]\\d \\w+ \\d{4}"))
		})
	})
})
