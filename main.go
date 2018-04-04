package main


import (
	"baseJson/basejson"
)

func main(){
	//str1 := "{\"key\"    : 23423423434234, \"key1\"    : \"value2\", \"key2\" :    { \"inner_key1\" : \"inner_value1\"  } , \"key3\"    : \"value3\",}"
	//str2 := "[{\"key\"    : \"value\", , , , \"key1\"    : \"value2\", \"key2\" : { \"inner_key1\" : \"inner_value1\"  } }]"
	//str3 := "[{\"key\"    : 123123L, , , , \"key1\"    : \"value2\", \"key2\" : { \"inner_key1\" : \"inner_value1\"  } }, {\"key\"    : \"value\", , , , \"key1\"    : \"value2\", \"key2\" : { \"inner_key1\" : \"inner_value1\"  } }]"

	//str4 := "[    \"123123\"   ,    \"sdfsdf123123\"    , \"ge34g34g\"    ]"
	//str5 := "[    123123   ,    2.342342e34    , 234234234    ]"
	//str6 := "[[[\"test\"], [ \"test2\", \"test3\" ]]]"

	//str7 := "[ true, false,, true]"
	//str8 := "{\"key\"    : 23423423434234, \"key1\"    : \"value2\", \"key2\" :    { \"inner_key1\" : \"inner_value1\"  } , \"key3\"    : \"value3\"}}"

	str9 := []byte(`[{
  "person": {
    "id": "d50887ca-a6ce-4e59-b89f-14f0b5d03b03",
    "name": {
      "fullName": "Leonid Bugaev",
      "givenName": "Leonid",
      "familyName": "Bugaev"
    },
    "email": "leonsbox@gmail.com",
    "gender": "male",
    "location": "Saint Petersburg, Saint Petersburg, RU",
    "geo": {
      "city": "Saint Petersburg",
      "state": "Saint Petersburg",
      "country": "Russia",
      "lat": 59.9342802,
      "lng": 30.3350986
    },
    "bio": "Senior engineer at Granify.com",
    "site": "http://flickfaver.com",
    "avatar": "https://d1ts43dypk8bqh.cloudfront.net/v1/avatars/d50887ca-a6ce-4e59-b89f-14f0b5d03b03",
    "employment": {
      "name": "www.latera.ru",
      "title": "Software Engineer",
      "domain": "gmail.com"
    },
    "facebook": {
      "handle": "leonid.bugaev"
    },
    "github": {
      "handle": "buger",
      "id": 14009,
      "avatar": "https://avatars.githubusercontent.com/u/14009?v=3",
      "company": "Granify",
      "blog": "http://leonsbox.com",
      "followers": 95,
      "following": 10
    },
    "twitter": {
      "handle": "flickfaver",
      "id": 77004410,
      "bio": null,
      "followers": 2,
      "following": 1,
      "statuses": 5,
      "favorites": 0,
      "location": "",
      "site": "http://flickfaver.com",
      "avatar": null
    },
    "linkedin": {
      "handle": "in/leonidbugaev"
    },
    "googleplus": {
      "handle": null
    },
    "angellist": {
      "handle": "leonid-bugaev",
      "id": 61541,
      "bio": "Senior engineer at Granify.com",
      "blog": "http://buger.github.com",
      "site": "http://buger.github.com",
      "followers": 41,
      "avatar": "https://d1qb2nb5cznatu.cloudfront.net/users/61541-medium_jpg?1405474390"
    },
    "klout": {
      "handle": null,
      "score": null
    },
    "foursquare": {
      "handle": null
    },
    "aboutme": {
      "handle": "leonid.bugaev",
      "bio": null,
      "avatar": null
    },
    "gravatar": {
      "handle": "buger",
      "urls": [
      ],
      "avatar": "http://1.gravatar.com/avatar/f7c8edd577d13b8930d5522f28123510",
      "avatars": [
        {
          "url": "http://1.gravatar.com/avatar/f7c8edd577d13b8930d5522f28123510",
          "type": "thumbnail"
        }
      ]
    },
    "fuzzy": false
  },
  "company": "hello"
}]`)

	parser := basejson.NewJsonParser(string(str9))
	parser.Parse()
}