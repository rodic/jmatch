package jmatch

import (
	"os"
	"reflect"
	"testing"

	z "github.com/rodic/jmatch/tokenizer"
)

type CollectorMatcher struct {
	matches map[string]z.Token
}

func (fm *CollectorMatcher) Match(path string, token z.Token) {
	fm.matches[path] = token
}

func TestMatcherValid(t *testing.T) {
	testCases := []struct {
		name     string
		expected map[string]z.Token
	}{
		{name: "testdata/valid/empty.json",
			expected: map[string]z.Token{}},
		{name: "testdata/valid/nested.json",
			expected: map[string]z.Token{
				".name":                  z.NewStringToken("Chris", 2, 13),
				".age":                   z.NewNumberToken("23", 3, 12),
				".address.city":          z.NewStringToken("New York", 5, 15),
				".address.country":       z.NewStringToken("America", 6, 18),
				".friends[0].name":       z.NewStringToken("Emily", 10, 17),
				".friends[0].hobbies[0]": z.NewStringToken("biking", 11, 22),
				".friends[0].hobbies[1]": z.NewStringToken("music", 11, 32),
				".friends[0].hobbies[2]": z.NewStringToken("gaming", 11, 41),
				".friends[1].name":       z.NewStringToken("John", 14, 17),
				".friends[1].hobbies[0]": z.NewStringToken("soccer", 15, 22),
				".friends[1].hobbies[1]": z.NewStringToken("gaming", 15, 32),
			},
		},
		{name: "testdata/valid/arrays.json",
			expected: map[string]z.Token{
				".[0].name": z.NewStringToken("Chris", 2, 15),
				".[0].age":  z.NewNumberToken("23", 2, 31),
				".[0].city": z.NewStringToken("New York", 2, 43),
				".[1].name": z.NewStringToken("Emily", 3, 15),
				".[1].age":  z.NewNumberToken("19", 3, 31),
				".[1].city": z.NewStringToken("Atlanta", 3, 43),
				".[2].name": z.NewStringToken("Joe", 4, 15),
				".[2].age":  z.NewNumberToken("32", 4, 29),
				".[2].city": z.NewStringToken("New York", 4, 41),
				".[3].name": z.NewStringToken("Kevin", 5, 15),
				".[3].age":  z.NewNumberToken("19", 5, 31),
				".[3].city": z.NewStringToken("Atlanta", 5, 43),
				".[4].name": z.NewStringToken("Michelle", 6, 15),
				".[4].age":  z.NewNumberToken("27", 6, 34),
				".[4].city": z.NewStringToken("Los Angeles", 6, 46),
				".[5].name": z.NewStringToken("Robert", 7, 15),
				".[5].age":  z.NewNumberToken("45", 7, 32),
				".[5].city": z.NewStringToken("Manhattan", 7, 44),
				".[6].name": z.NewStringToken("Sarah", 8, 15),
				".[6].age":  z.NewNumberToken("31", 8, 31),
				".[6].city": z.NewStringToken("New York", 8, 43),
			},
		},
		{name: "testdata/valid/date.json",
			expected: map[string]z.Token{
				".id":      z.NewStringToken("a98d1377-2270-45fd-8e25-cde720c50bce", 2, 11),
				".message": z.NewStringToken("Hi Jane 😃 are you busy tonight? Shall we go out for dinner?", 3, 16),
				".date":    z.NewStringToken("2023-07-24T12:56:15.609Z", 4, 13),
			},
		},
		{name: "testdata/valid/colors.json",
			expected: map[string]z.Token{
				".[0].calendarId": z.NewStringToken("e2a5c", 3, 21),
				".[0].color":      z.NewStringToken("#3997f5", 4, 16),
				".[1].calendarId": z.NewStringToken("aa027", 7, 21),
				".[1].color":      z.NewStringToken("#ef5353", 8, 16),
				".[2].calendarId": z.NewStringToken("5d9a1", 11, 21),
				".[2].color":      z.NewStringToken("#3fc13f", 12, 16),
			},
		},
		{name: "testdata/valid/geo.json",
			expected: map[string]z.Token{
				".type":                                      z.NewStringToken("FeatureCollection", 2, 13),
				".features[0].type":                          z.NewStringToken("Feature", 5, 17),
				".features[0].geometry.type":                 z.NewStringToken("Point", 8, 19),
				".features[0].geometry.coordinates[0]":       z.NewNumberToken("4.483605784808901", 9, 27),
				".features[0].geometry.coordinates[1]":       z.NewNumberToken("51.907188449679325", 9, 46),
				".features[1].type":                          z.NewStringToken("Feature", 13, 17),
				".features[1].geometry.type":                 z.NewStringToken("Polygon", 16, 19),
				".features[1].geometry.coordinates[0][0][0]": z.NewNumberToken("3.974369110811523", 19, 16),
				".features[1].geometry.coordinates[0][0][1]": z.NewNumberToken("51.907355547778565", 19, 36),
				".features[1].geometry.coordinates[0][1][0]": z.NewNumberToken("4.173944459020191", 20, 16),
				".features[1].geometry.coordinates[0][1][1]": z.NewNumberToken("51.86237166892457", 20, 36),
				".features[1].geometry.coordinates[0][2][0]": z.NewNumberToken("4.3808076710679416", 21, 16),
				".features[1].geometry.coordinates[0][2][1]": z.NewNumberToken("51.848867725914914", 21, 36),
				".features[1].geometry.coordinates[0][3][0]": z.NewNumberToken("4.579822414365026", 22, 16),
				".features[1].geometry.coordinates[0][3][1]": z.NewNumberToken("51.874487141880024", 22, 36),
				".features[1].geometry.coordinates[0][4][0]": z.NewNumberToken("4.534413416598767", 23, 16),
				".features[1].geometry.coordinates[0][4][1]": z.NewNumberToken("51.9495302480326", 23, 36),
				".features[1].geometry.coordinates[0][5][0]": z.NewNumberToken("4.365110733567974", 24, 16),
				".features[1].geometry.coordinates[0][5][1]": z.NewNumberToken("51.92360787140825", 24, 36),
				".features[1].geometry.coordinates[0][6][0]": z.NewNumberToken("4.179550508127079", 25, 16),
				".features[1].geometry.coordinates[0][6][1]": z.NewNumberToken("51.97336560819281", 25, 36),
				".features[1].geometry.coordinates[0][7][0]": z.NewNumberToken("4.018096293847009", 26, 16),
				".features[1].geometry.coordinates[0][7][1]": z.NewNumberToken("52.00236546429852", 26, 36),
				".features[1].geometry.coordinates[0][8][0]": z.NewNumberToken("3.9424146309028174", 27, 16),
				".features[1].geometry.coordinates[0][8][1]": z.NewNumberToken("51.97681895676649", 27, 36),
				".features[1].geometry.coordinates[0][9][0]": z.NewNumberToken("3.974369110811523", 28, 16),
				".features[1].geometry.coordinates[0][9][1]": z.NewNumberToken("51.907355547778565", 28, 36),
			}},
		{name: "testdata/valid/youtube.json",
			expected: map[string]z.Token{
				".kind":                                       z.NewStringToken("youtube#searchListResponse", 2, 13),
				".etag":                                       z.NewStringToken("q4ibjmYp1KA3RqMF4jFLl6PBwOg", 3, 13),
				".nextPageToken":                              z.NewStringToken("CAUQAA", 4, 22),
				".regionCode":                                 z.NewStringToken("NL", 5, 19),
				".pageInfo.totalResults":                      z.NewNumberToken("1000000", 6, 34),
				".pageInfo.resultsPerPage":                    z.NewNumberToken("5", 6, 61),
				".items[0].kind":                              z.NewStringToken("youtube#searchResult", 9, 17),
				".items[0].etag":                              z.NewStringToken("QCsHBifbaernVCbLv8Cu6rAeaDQ", 10, 17),
				".items[0].id.kind":                           z.NewStringToken("youtube#video", 11, 24),
				".items[0].id.videoId":                        z.NewStringToken("TvWDY4Mm5GM", 11, 52),
				".items[0].snippet.publishTime":               z.NewStringToken("2023-07-24T14:15:01Z", 39, 26),
				".items[0].snippet.channelId":                 z.NewStringToken("UCwozCpFp9g9x0wAzuFh0hwQ", 14, 24),
				".items[0].snippet.title":                     z.NewStringToken("3 Football Clubs Kylian Mbappe Should Avoid Signing ✍️❌⚽️ #football #mbappe #shorts", 15, 20),
				".items[0].snippet.description":               z.NewStringToken("", 16, 26),
				".items[0].snippet.thumbnails.default.url":    z.NewStringToken("https://i.ytimg.com/vi/TvWDY4Mm5GM/default.jpg", 19, 22),
				".items[0].snippet.thumbnails.default.width":  z.NewNumberToken("120", 20, 24),
				".items[0].snippet.thumbnails.default.height": z.NewNumberToken("90", 21, 25),
				".items[0].snippet.thumbnails.medium.url":     z.NewStringToken("https://i.ytimg.com/vi/TvWDY4Mm5GM/mqdefault.jpg", 25, 22),
				".items[0].snippet.thumbnails.medium.height":  z.NewNumberToken("180", 27, 25),
				".items[0].snippet.thumbnails.medium.width":   z.NewNumberToken("320", 26, 24),
				".items[0].snippet.thumbnails.high.url":       z.NewStringToken("https://i.ytimg.com/vi/TvWDY4Mm5GM/hqdefault.jpg", 31, 22),
				".items[0].snippet.thumbnails.high.width":     z.NewNumberToken("480", 32, 24),
				".items[0].snippet.thumbnails.high.height":    z.NewNumberToken("360", 33, 25),
				".items[0].snippet.channelTitle":              z.NewStringToken("FC Motivate", 37, 27),
				".items[0].snippet.liveBroadcastContent":      z.NewStringToken("none", 38, 35),
				".items[0].snippet.publishedAt":               z.NewStringToken("2023-07-24T14:15:01Z", 13, 26),
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
				matches: make(map[string]z.Token),
			}

			Match(string(json), &collector)

			if !reflect.DeepEqual(collector.matches, tc.expected) {
				t.Errorf("Expected '%v', got '%v' instead\n", tc.expected, collector.matches)
			}
		})
	}
}
