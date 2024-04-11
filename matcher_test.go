package jmatch

import (
	"os"
	"reflect"
	"testing"
)

type CollectorMatcher struct {
	matches []string
}

func (fm *CollectorMatcher) Match(path string, token Token) {
	fm.matches = append(fm.matches, path)
}

func TestMatcherValid(t *testing.T) {
	testCases := []struct {
		name     string
		expected map[string]Token
	}{
		{name: "testdata/valid/nested.json",
			expected: map[string]Token{
				".name":                    newToken(String, "Chris", 2, 13),
				".age":                     newToken(Number, "23", 3, 12),
				".address.city":            newToken(String, "New York", 5, 15),
				".address.country":         newToken(String, "America", 6, 18),
				".friends.[0].name":        newToken(String, "Emily", 10, 17),
				".friends.[0].hobbies.[0]": newToken(String, "biking", 11, 22),
				".friends.[0].hobbies.[1]": newToken(String, "music", 11, 32),
				".friends.[0].hobbies.[2]": newToken(String, "gaming", 11, 41),
				".friends.[1].name":        newToken(String, "John", 14, 17),
				".friends.[1].hobbies.[0]": newToken(String, "soccer", 15, 22),
				".friends.[1].hobbies.[1]": newToken(String, "gaming", 15, 32),
			},
		},
		{name: "testdata/valid/arrays.json",
			expected: map[string]Token{
				".[0].name": newToken(String, "Chris", 2, 15),
				".[0].age":  newToken(Number, "23", 2, 31),
				".[0].city": newToken(String, "New York", 2, 43),
				".[1].name": newToken(String, "Emily", 3, 15),
				".[1].age":  newToken(Number, "19", 3, 31),
				".[1].city": newToken(String, "Atlanta", 3, 43),
				".[2].name": newToken(String, "Joe", 4, 15),
				".[2].age":  newToken(Number, "32", 4, 29),
				".[2].city": newToken(String, "New York", 4, 41),
				".[3].name": newToken(String, "Kevin", 5, 15),
				".[3].age":  newToken(Number, "19", 5, 31),
				".[3].city": newToken(String, "Atlanta", 5, 43),
				".[4].name": newToken(String, "Michelle", 6, 15),
				".[4].age":  newToken(Number, "27", 6, 34),
				".[4].city": newToken(String, "Los Angeles", 6, 46),
				".[5].name": newToken(String, "Robert", 7, 15),
				".[5].age":  newToken(Number, "45", 7, 32),
				".[5].city": newToken(String, "Manhattan", 7, 44),
				".[6].name": newToken(String, "Sarah", 8, 15),
				".[6].age":  newToken(Number, "31", 8, 31),
				".[6].city": newToken(String, "New York", 8, 43),
			},
		},
		{name: "testdata/valid/date.json",
			expected: map[string]Token{
				".id":      newToken(String, "a98d1377-2270-45fd-8e25-cde720c50bce", 2, 11),
				".message": newToken(String, "Hi Jane 😃 are you busy tonight? Shall we go out for dinner?", 3, 16),
				".date":    newToken(String, "2023-07-24T12:56:15.609Z", 4, 13),
			},
		},
		{name: "testdata/valid/colors.json",
			expected: map[string]Token{
				".[0].calendarId": newToken(String, "e2a5c", 3, 21),
				".[0].color":      newToken(String, "#3997f5", 4, 16),
				".[1].calendarId": newToken(String, "aa027", 7, 21),
				".[1].color":      newToken(String, "#ef5353", 8, 16),
				".[2].calendarId": newToken(String, "5d9a1", 11, 21),
				".[2].color":      newToken(String, "#3fc13f", 12, 16),
			},
		},
		{name: "testdata/valid/geo.json",
			expected: map[string]Token{
				".type":                                          newToken(String, "FeatureCollection", 2, 13),
				".features.[0].type":                             newToken(String, "Feature", 5, 17),
				".features.[0].geometry.type":                    newToken(String, "Point", 8, 19),
				".features.[0].geometry.coordinates.[0]":         newToken(Number, "4.483605784808901", 9, 27),
				".features.[0].geometry.coordinates.[1]":         newToken(Number, "51.907188449679325", 9, 46),
				".features.[1].type":                             newToken(String, "Feature", 13, 17),
				".features.[1].geometry.type":                    newToken(String, "Polygon", 16, 19),
				".features.[1].geometry.coordinates.[0].[0].[0]": newToken(Number, "3.974369110811523", 19, 16),
				".features.[1].geometry.coordinates.[0].[0].[1]": newToken(Number, "51.907355547778565", 19, 36),
				".features.[1].geometry.coordinates.[0].[1].[0]": newToken(Number, "4.173944459020191", 20, 16),
				".features.[1].geometry.coordinates.[0].[1].[1]": newToken(Number, "51.86237166892457", 20, 36),
				".features.[1].geometry.coordinates.[0].[2].[0]": newToken(Number, "4.3808076710679416", 21, 16),
				".features.[1].geometry.coordinates.[0].[2].[1]": newToken(Number, "51.848867725914914", 21, 36),
				".features.[1].geometry.coordinates.[0].[3].[0]": newToken(Number, "4.579822414365026", 22, 16),
				".features.[1].geometry.coordinates.[0].[3].[1]": newToken(Number, "51.874487141880024", 22, 36),
				".features.[1].geometry.coordinates.[0].[4].[0]": newToken(Number, "4.534413416598767", 23, 16),
				".features.[1].geometry.coordinates.[0].[4].[1]": newToken(Number, "51.9495302480326", 23, 36),
				".features.[1].geometry.coordinates.[0].[5].[0]": newToken(Number, "4.365110733567974", 24, 16),
				".features.[1].geometry.coordinates.[0].[5].[1]": newToken(Number, "51.92360787140825", 24, 36),
				".features.[1].geometry.coordinates.[0].[6].[0]": newToken(Number, "4.179550508127079", 25, 16),
				".features.[1].geometry.coordinates.[0].[6].[1]": newToken(Number, "51.97336560819281", 25, 36),
				".features.[1].geometry.coordinates.[0].[7].[0]": newToken(Number, "4.018096293847009", 26, 16),
				".features.[1].geometry.coordinates.[0].[7].[1]": newToken(Number, "52.00236546429852", 26, 36),
				".features.[1].geometry.coordinates.[0].[8].[0]": newToken(Number, "3.9424146309028174", 27, 16),
				".features.[1].geometry.coordinates.[0].[8].[1]": newToken(Number, "51.97681895676649", 27, 36),
				".features.[1].geometry.coordinates.[0].[9].[0]": newToken(Number, "3.974369110811523", 28, 16),
				".features.[1].geometry.coordinates.[0].[9].[1]": newToken(Number, "51.907355547778565", 28, 36),
			}},
		{name: "testdata/valid/youtube.json",
			expected: map[string]Token{
				".kind":                                        newToken(String, "youtube#searchListResponse", 2, 13),
				".etag":                                        newToken(String, "q4ibjmYp1KA3RqMF4jFLl6PBwOg", 3, 13),
				".nextPageToken":                               newToken(String, "CAUQAA", 4, 22),
				".regionCode":                                  newToken(String, "NL", 5, 19),
				".pageInfo.totalResults":                       newToken(Number, "1000000", 6, 34),
				".pageInfo.resultsPerPage":                     newToken(Number, "5", 6, 61),
				".items.[0].kind":                              newToken(String, "youtube#searchResult", 9, 17),
				".items.[0].etag":                              newToken(String, "QCsHBifbaernVCbLv8Cu6rAeaDQ", 10, 17),
				".items.[0].id.kind":                           newToken(String, "youtube#video", 11, 24),
				".items.[0].id.videoId":                        newToken(String, "TvWDY4Mm5GM", 11, 52),
				".items.[0].snippet.publishTime":               newToken(String, "2023-07-24T14:15:01Z", 39, 26),
				".items.[0].snippet.channelId":                 newToken(String, "UCwozCpFp9g9x0wAzuFh0hwQ", 14, 24),
				".items.[0].snippet.title":                     newToken(String, "3 Football Clubs Kylian Mbappe Should Avoid Signing ✍️❌⚽️ #football #mbappe #shorts", 15, 20),
				".items.[0].snippet.description":               newToken(String, "", 16, 26),
				".items.[0].snippet.thumbnails.default.url":    newToken(String, "https://i.ytimg.com/vi/TvWDY4Mm5GM/default.jpg", 19, 22),
				".items.[0].snippet.thumbnails.default.width":  newToken(Number, "120", 20, 24),
				".items.[0].snippet.thumbnails.default.height": newToken(Number, "90", 21, 25),
				".items.[0].snippet.thumbnails.medium.url":     newToken(String, "https://i.ytimg.com/vi/TvWDY4Mm5GM/mqdefault.jpg", 25, 22),
				".items.[0].snippet.thumbnails.medium.height":  newToken(Number, "180", 27, 25),
				".items.[0].snippet.thumbnails.medium.width":   newToken(Number, "320", 26, 24),
				".items.[0].snippet.thumbnails.high.url":       newToken(String, "https://i.ytimg.com/vi/TvWDY4Mm5GM/hqdefault.jpg", 31, 22),
				".items.[0].snippet.thumbnails.high.width":     newToken(Number, "480", 32, 24),
				".items.[0].snippet.thumbnails.high.height":    newToken(Number, "360", 33, 25),
				".items.[0].snippet.channelTitle":              newToken(String, "FC Motivate", 37, 27),
				".items.[0].snippet.liveBroadcastContent":      newToken(String, "none", 38, 35),
				".items.[0].snippet.publishedAt":               newToken(String, "2023-07-24T14:15:01Z", 13, 26),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			json, err := os.ReadFile(tc.name)
			if err != nil {
				t.Fatal(err)
			}

			collector := collectorMatcher{
				collection: make(map[string]Token),
			}

			Match(string(json), &collector)

			if !reflect.DeepEqual(collector.collection, tc.expected) {
				t.Errorf("Expected '%v', got '%v' instead\n", tc.expected, collector.collection)
			}
		})
	}
}
