package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/Smartling/api-sdk-go"
	"github.com/stretchr/testify/assert"
)

func (suite *MainSuite) TestProjectsList() {
	suite.Mock.Handler = func(
		writer http.ResponseWriter,
		request *http.Request,
	) {
		list := smartling.ProjectsList{
			TotalCount: 2,
			Items: []smartling.Project{
				{
					ProjectID:      "01234ab",
					ProjectName:    "Rick and Morty",
					SourceLocaleID: "en-US",
				},
				{
					ProjectID:      "defaaaa",
					ProjectName:    "Adventure Time",
					SourceLocaleID: "en-GB",
				},
			},
		}

		err := writeSmartlingReply(writer, codeSuccess, list)
		if err != nil {
			panic(err)
		}
	}

	suite.assertStdout(
		[]string{
			"01234ab  Rick and Morty  en-US",
			"defaaaa  Adventure Time  en-GB",
		},
		"projects", "list",
	)

	suite.assertStdout(
		[]string{
			"01234ab",
			"defaaaa",
		},
		"projects", "list", "--short",
	)
}

func (suite *MainSuite) TestProjectsInfo() {
	suite.Mock.Handler = func(
		writer http.ResponseWriter,
		request *http.Request,
	) {
		assert.True(
			suite.T(),
			strings.HasSuffix(request.URL.Path, "/01234ab"),
		)

		list := smartling.Project{
			ProjectID:               "01234ab",
			ProjectName:             "Rick and Morty",
			AccountUID:              "xxyyzz",
			SourceLocaleID:          "en-US",
			SourceLocaleDescription: "English (United States)",
			Archived:                false,
		}

		err := writeSmartlingReply(writer, codeSuccess, list)
		if err != nil {
			panic(err)
		}
	}

	suite.assertStdout(
		[]string{
			"ID       01234ab",
			"ACCOUNT  xxyyzz",
			"NAME     Rick and Morty",
			"LOCALE   en-US: English (United States)",
			"STATUS   active",
		},
		"projects", "info", "-p", "01234ab",
	)
}

func (suite *MainSuite) TestProjectsLocales() {
	suite.Mock.Handler = func(
		writer http.ResponseWriter,
		request *http.Request,
	) {
		assert.True(
			suite.T(),
			strings.HasSuffix(request.URL.Path, "/01234ab"),
		)

		list := smartling.ProjectDetails{
			Project: smartling.Project{
				SourceLocaleID: "en-US",
			},
			TargetLocales: []smartling.Locale{
				{
					LocaleID:    "zh-CN",
					Description: "Chinese",
					Enabled:     true,
				},
				{
					LocaleID:    "nl-NL",
					Description: "Dutch",
					Enabled:     false,
				},
			},
		}

		err := writeSmartlingReply(writer, codeSuccess, list)
		if err != nil {
			panic(err)
		}
	}

	suite.assertStdout(
		[]string{
			"zh-CN  Chinese  true",
			"nl-NL  Dutch    false",
		},
		"projects", "locales", "-p", "01234ab",
	)

	suite.assertStdout(
		[]string{
			"zh-CN",
			"nl-NL",
		},
		"projects", "locales", "-p", "01234ab", "--short",
	)

	suite.assertStdout(
		[]string{
			"en-US",
		},
		"projects", "locales", "-p", "01234ab", "--source",
	)

	suite.assertStdout(
		[]string{
			"X",
			"Y",
		},
		"projects", "locales", "-p", "01234ab", "--format",
		`{{if eq .LocaleID "zh-CN"}}X{{else}}Y{{end}}\n`,
	)
}

func (suite *MainSuite) TestFilesList() {
	suite.Mock.Handler = func(
		writer http.ResponseWriter,
		request *http.Request,
	) {
		assert.True(
			suite.T(),
			strings.Contains(request.URL.Path, "/01234ab/"),
		)

		list := smartling.FilesList{
			TotalCount: 2,
			Items: []smartling.File{
				{
					FileURI:      "/Rick/portal-gun.java",
					LastUploaded: utc("2016-09-16T16:06:16Z"),
					FileType:     "javaProperties",
				},
				{
					FileURI:      "/Morty/stupidness.txt",
					LastUploaded: utc("1989-01-09T05:00:00Z"),
					FileType:     "plaintext",
				},
			},
		}

		err := writeSmartlingReply(writer, codeSuccess, list)
		if err != nil {
			panic(err)
		}
	}

	suite.assertStdout(
		[]string{
			"/Rick/portal-gun.java  2016-09-16T16:06:16Z  javaProperties",
			"/Morty/stupidness.txt  1989-01-09T05:00:00Z  plaintext",
		},
		"files", "list", "-p", "01234ab",
	)

	suite.assertStdout(
		[]string{
			"/Rick/portal-gun.java",
			"/Morty/stupidness.txt",
		},
		"files", "list", "-p", "01234ab", "--short",
	)

	suite.assertStdout(
		[]string{
			"/Rick/portal-gun.java",
		},
		"files", "list", "-p", "01234ab", "--short", "**.java",
	)

	suite.assertStdout(
		[]string{
			"/Rick/portal-gun.java javaProperties",
		},
		"files", "list", "-p", "01234ab", "**.java", "--format",
		"{{.FileURI}} {{.FileType}}\n",
	)
}

func (suite *MainSuite) TestFilesPull() {
	suite.Mock.Handler = func(
		writer http.ResponseWriter,
		request *http.Request,
	) {
		assert.True(
			suite.T(),
			strings.Contains(request.URL.Path, "/01234ab/"),
		)

		var reply interface{}

		switch {
		case strings.HasSuffix(request.URL.Path, "/file"):
			writer.WriteHeader(http.StatusOK)

			switch request.URL.Query().Get("fileUri") {
			case "/Rick/portal-gun.java":
				switch {
				case strings.Contains(request.URL.Path, "/de-DE/"):
					io.WriteString(writer, "Rick:de-DE\n")
				default:
					io.WriteString(writer, "Rick:original\n")
				}

			case "/Morty/stupidness.txt":
				switch {
				case strings.Contains(request.URL.Path, "/es/"):
					io.WriteString(writer, "Morty:es\n")
				default:
					io.WriteString(writer, "Morty:original\n")
				}
			}

			return

		case strings.HasSuffix(request.URL.Path, "/status"):
			switch request.URL.Query().Get("fileUri") {
			case "/Rick/portal-gun.java":
				reply = smartling.FileStatus{
					TotalStringCount: 12,
					TotalWordCount:   120,
					Items: []smartling.FileStatusTranslation{
						{
							LocaleID:             "de-DE",
							CompletedStringCount: 10,
							CompletedWordCount:   100,
						},
					},
				}

			case "/Morty/stupidness.txt":
				reply = smartling.FileStatus{
					TotalStringCount: 2,
					TotalWordCount:   12,
					Items: []smartling.FileStatusTranslation{
						{
							LocaleID:             "es",
							CompletedStringCount: 1,
							CompletedWordCount:   10,
						},
					},
				}
			}

		case strings.HasSuffix(request.URL.Path, "/list"):
			reply = smartling.FilesList{
				TotalCount: 2,
				Items: []smartling.File{
					{
						FileURI:      "/Rick/portal-gun.java",
						LastUploaded: utc("2016-09-16T16:06:16Z"),
						FileType:     "javaProperties",
					},
					{
						FileURI:      "/Morty/stupidness.txt",
						LastUploaded: utc("1989-01-09T05:00:00Z"),
						FileType:     "plaintext",
					},
				},
			}
		}

		err := writeSmartlingReply(writer, codeSuccess, reply)
		if err != nil {
			panic(err)
		}
	}

	assertFileEquals := func(path string, contents string) {
		output, err := ioutil.ReadFile(path)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), string(output), contents)
	}

	defer func() {
		err := os.RemoveAll("_test")
		assert.NoError(suite.T(), err)
	}()

	suite.assertStdout(
		[]string{
			"downloaded _test/Morty/stupidness_es.txt 50%",
			"downloaded _test/Rick/portal-gun_de-DE.java 83%",
		},
		"files", "pull", "-p", "01234ab", "-d", "_test",
	)

	assertFileEquals("_test/Morty/stupidness_es.txt", "Morty:es\n")
	assertFileEquals("_test/Rick/portal-gun_de-DE.java", "Rick:de-DE\n")

	suite.assertStdout(
		[]string{
			"downloaded _test/x/Morty/stupidness_es.txt 50%",
			"downloaded _test/x/Rick/portal-gun_de-DE.java 83%",
		},
		"files", "pull", "-p", "01234ab", "-d", "_test/x",
	)

	assertFileEquals("_test/x/Morty/stupidness_es.txt", "Morty:es\n")
	assertFileEquals("_test/x/Rick/portal-gun_de-DE.java", "Rick:de-DE\n")

	suite.assertStdout(
		[]string{
			"downloaded _test/Morty/stupidness 50%",
			"downloaded _test/Rick/portal-gun 83%",
		},
		"files", "pull", "-p", "01234ab", "-d", "_test", "--format",
		"{{name .FileURI}}",
	)

	assertFileEquals("_test/Morty/stupidness", "Morty:es\n")
	assertFileEquals("_test/Rick/portal-gun", "Rick:de-DE\n")

	suite.assertStdout(
		[]string{
			"downloaded _test/Rick/portal-gun_de-DE.java 83%",
		},
		"files", "pull", "-p", "01234ab", "-d", "_test", "--progress", "80%",
	)

	suite.assertStdout(
		[]string{
			"downloaded _test/Rick/portal-gun.java",
			"downloaded _test/Morty/stupidness.txt",
		},
		"files", "pull", "-p", "01234ab", "-d", "_test", "--source",
	)

	assertFileEquals("_test/Morty/stupidness.txt", "Morty:original\n")
	assertFileEquals("_test/Rick/portal-gun.java", "Rick:original\n")
}
