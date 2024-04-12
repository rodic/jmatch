package jmatch

import (
	"os"
	"reflect"
	"testing"

	token "github.com/rodic/jmatch/token"
)

type CollectorMatcher struct {
	matches map[string]token.Token
}

func (fm *CollectorMatcher) Match(path string, token token.Token) {
	fm.matches[path] = token
}

func TestMatcherValid(t *testing.T) {
	testCases := []struct {
		name     string
		expected map[string]token.Token
	}{
		{name: "testdata/valid/nested.json",
			expected: map[string]token.Token{
				".name":                  token.NewStringToken("Chris", 2, 13),
				".age":                   token.NewNumberToken("23", 3, 12),
				".address.city":          token.NewStringToken("New York", 5, 15),
				".address.country":       token.NewStringToken("America", 6, 18),
				".friends[0].name":       token.NewStringToken("Emily", 10, 17),
				".friends[0].hobbies[0]": token.NewStringToken("biking", 11, 22),
				".friends[0].hobbies[1]": token.NewStringToken("music", 11, 32),
				".friends[0].hobbies[2]": token.NewStringToken("gaming", 11, 41),
				".friends[1].name":       token.NewStringToken("John", 14, 17),
				".friends[1].hobbies[0]": token.NewStringToken("soccer", 15, 22),
				".friends[1].hobbies[1]": token.NewStringToken("gaming", 15, 32),
			},
		},
		{name: "testdata/valid/arrays.json",
			expected: map[string]token.Token{
				".[0].name": token.NewStringToken("Chris", 2, 15),
				".[0].age":  token.NewNumberToken("23", 2, 31),
				".[0].city": token.NewStringToken("New York", 2, 43),
				".[1].name": token.NewStringToken("Emily", 3, 15),
				".[1].age":  token.NewNumberToken("19", 3, 31),
				".[1].city": token.NewStringToken("Atlanta", 3, 43),
				".[2].name": token.NewStringToken("Joe", 4, 15),
				".[2].age":  token.NewNumberToken("32", 4, 29),
				".[2].city": token.NewStringToken("New York", 4, 41),
				".[3].name": token.NewStringToken("Kevin", 5, 15),
				".[3].age":  token.NewNumberToken("19", 5, 31),
				".[3].city": token.NewStringToken("Atlanta", 5, 43),
				".[4].name": token.NewStringToken("Michelle", 6, 15),
				".[4].age":  token.NewNumberToken("27", 6, 34),
				".[4].city": token.NewStringToken("Los Angeles", 6, 46),
				".[5].name": token.NewStringToken("Robert", 7, 15),
				".[5].age":  token.NewNumberToken("45", 7, 32),
				".[5].city": token.NewStringToken("Manhattan", 7, 44),
				".[6].name": token.NewStringToken("Sarah", 8, 15),
				".[6].age":  token.NewNumberToken("31", 8, 31),
				".[6].city": token.NewStringToken("New York", 8, 43),
			},
		},
		{name: "testdata/valid/date.json",
			expected: map[string]token.Token{
				".id":      token.NewStringToken("a98d1377-2270-45fd-8e25-cde720c50bce", 2, 11),
				".message": token.NewStringToken("Hi Jane 😃 are you busy tonight? Shall we go out for dinner?", 3, 16),
				".date":    token.NewStringToken("2023-07-24T12:56:15.609Z", 4, 13),
			},
		},
		{name: "testdata/valid/colors.json",
			expected: map[string]token.Token{
				".[0].calendarId": token.NewStringToken("e2a5c", 3, 21),
				".[0].color":      token.NewStringToken("#3997f5", 4, 16),
				".[1].calendarId": token.NewStringToken("aa027", 7, 21),
				".[1].color":      token.NewStringToken("#ef5353", 8, 16),
				".[2].calendarId": token.NewStringToken("5d9a1", 11, 21),
				".[2].color":      token.NewStringToken("#3fc13f", 12, 16),
			},
		},
		{name: "testdata/valid/geo.json",
			expected: map[string]token.Token{
				".type":                                      token.NewStringToken("FeatureCollection", 2, 13),
				".features[0].type":                          token.NewStringToken("Feature", 5, 17),
				".features[0].geometry.type":                 token.NewStringToken("Point", 8, 19),
				".features[0].geometry.coordinates[0]":       token.NewNumberToken("4.483605784808901", 9, 27),
				".features[0].geometry.coordinates[1]":       token.NewNumberToken("51.907188449679325", 9, 46),
				".features[1].type":                          token.NewStringToken("Feature", 13, 17),
				".features[1].geometry.type":                 token.NewStringToken("Polygon", 16, 19),
				".features[1].geometry.coordinates[0][0][0]": token.NewNumberToken("3.974369110811523", 19, 16),
				".features[1].geometry.coordinates[0][0][1]": token.NewNumberToken("51.907355547778565", 19, 36),
				".features[1].geometry.coordinates[0][1][0]": token.NewNumberToken("4.173944459020191", 20, 16),
				".features[1].geometry.coordinates[0][1][1]": token.NewNumberToken("51.86237166892457", 20, 36),
				".features[1].geometry.coordinates[0][2][0]": token.NewNumberToken("4.3808076710679416", 21, 16),
				".features[1].geometry.coordinates[0][2][1]": token.NewNumberToken("51.848867725914914", 21, 36),
				".features[1].geometry.coordinates[0][3][0]": token.NewNumberToken("4.579822414365026", 22, 16),
				".features[1].geometry.coordinates[0][3][1]": token.NewNumberToken("51.874487141880024", 22, 36),
				".features[1].geometry.coordinates[0][4][0]": token.NewNumberToken("4.534413416598767", 23, 16),
				".features[1].geometry.coordinates[0][4][1]": token.NewNumberToken("51.9495302480326", 23, 36),
				".features[1].geometry.coordinates[0][5][0]": token.NewNumberToken("4.365110733567974", 24, 16),
				".features[1].geometry.coordinates[0][5][1]": token.NewNumberToken("51.92360787140825", 24, 36),
				".features[1].geometry.coordinates[0][6][0]": token.NewNumberToken("4.179550508127079", 25, 16),
				".features[1].geometry.coordinates[0][6][1]": token.NewNumberToken("51.97336560819281", 25, 36),
				".features[1].geometry.coordinates[0][7][0]": token.NewNumberToken("4.018096293847009", 26, 16),
				".features[1].geometry.coordinates[0][7][1]": token.NewNumberToken("52.00236546429852", 26, 36),
				".features[1].geometry.coordinates[0][8][0]": token.NewNumberToken("3.9424146309028174", 27, 16),
				".features[1].geometry.coordinates[0][8][1]": token.NewNumberToken("51.97681895676649", 27, 36),
				".features[1].geometry.coordinates[0][9][0]": token.NewNumberToken("3.974369110811523", 28, 16),
				".features[1].geometry.coordinates[0][9][1]": token.NewNumberToken("51.907355547778565", 28, 36),
			}},
		{name: "testdata/valid/youtube.json",
			expected: map[string]token.Token{
				".kind":                                       token.NewStringToken("youtube#searchListResponse", 2, 13),
				".etag":                                       token.NewStringToken("q4ibjmYp1KA3RqMF4jFLl6PBwOg", 3, 13),
				".nextPageToken":                              token.NewStringToken("CAUQAA", 4, 22),
				".regionCode":                                 token.NewStringToken("NL", 5, 19),
				".pageInfo.totalResults":                      token.NewNumberToken("1000000", 6, 34),
				".pageInfo.resultsPerPage":                    token.NewNumberToken("5", 6, 61),
				".items[0].kind":                              token.NewStringToken("youtube#searchResult", 9, 17),
				".items[0].etag":                              token.NewStringToken("QCsHBifbaernVCbLv8Cu6rAeaDQ", 10, 17),
				".items[0].id.kind":                           token.NewStringToken("youtube#video", 11, 24),
				".items[0].id.videoId":                        token.NewStringToken("TvWDY4Mm5GM", 11, 52),
				".items[0].snippet.publishTime":               token.NewStringToken("2023-07-24T14:15:01Z", 39, 26),
				".items[0].snippet.channelId":                 token.NewStringToken("UCwozCpFp9g9x0wAzuFh0hwQ", 14, 24),
				".items[0].snippet.title":                     token.NewStringToken("3 Football Clubs Kylian Mbappe Should Avoid Signing ✍️❌⚽️ #football #mbappe #shorts", 15, 20),
				".items[0].snippet.description":               token.NewStringToken("", 16, 26),
				".items[0].snippet.thumbnails.default.url":    token.NewStringToken("https://i.ytimg.com/vi/TvWDY4Mm5GM/default.jpg", 19, 22),
				".items[0].snippet.thumbnails.default.width":  token.NewNumberToken("120", 20, 24),
				".items[0].snippet.thumbnails.default.height": token.NewNumberToken("90", 21, 25),
				".items[0].snippet.thumbnails.medium.url":     token.NewStringToken("https://i.ytimg.com/vi/TvWDY4Mm5GM/mqdefault.jpg", 25, 22),
				".items[0].snippet.thumbnails.medium.height":  token.NewNumberToken("180", 27, 25),
				".items[0].snippet.thumbnails.medium.width":   token.NewNumberToken("320", 26, 24),
				".items[0].snippet.thumbnails.high.url":       token.NewStringToken("https://i.ytimg.com/vi/TvWDY4Mm5GM/hqdefault.jpg", 31, 22),
				".items[0].snippet.thumbnails.high.width":     token.NewNumberToken("480", 32, 24),
				".items[0].snippet.thumbnails.high.height":    token.NewNumberToken("360", 33, 25),
				".items[0].snippet.channelTitle":              token.NewStringToken("FC Motivate", 37, 27),
				".items[0].snippet.liveBroadcastContent":      token.NewStringToken("none", 38, 35),
				".items[0].snippet.publishedAt":               token.NewStringToken("2023-07-24T14:15:01Z", 13, 26),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			json, err := os.ReadFile(tc.name)
			if err != nil {
				t.Fatal(err)
			}

			collector := CollectorMatcher{
				matches: make(map[string]token.Token),
			}

			Match(string(json), &collector)

			if !reflect.DeepEqual(collector.matches, tc.expected) {
				t.Errorf("Expected '%v', got '%v' instead\n", tc.expected, collector.matches)
			}
		})
	}
}
