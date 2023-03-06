# Assignment 1


## Description
This project is a submission to the first assignment in PROG2005. 

In this assignment I have developed a REST web application in Golang that provides the client to retrieve information about universities and countries.
By entering the names of the university and/or the country you can retrieve useful information about the country and the university. 
Information like the languages spoken in the countries, the official websites of the universities, the isocode, the name of the university and country will be retrieved. 

This REST web service relies on: 

- http://universities.hipolabs.com/
- Documentation/Source under: https://github.com/Hipo/university-domains-list/

* https://restcountries.com/

* Documentation/Source under: https://gitlab.com/amatos/rest-countries


## How to use endpoints: 

### uniinfo

This endpoint focuses on returning information about a country that a university resides in.
It retrieves information about the country like what languages are spoken there, the official name of the country, a OpenStreetMap link, and 
a link to the universities official website.

````
Method: GET
Path: uniinfo/{:partial_or_complete_university_name}
````

The resource root path:
* http://localhost:8080/unisearcher/v1/uniinfo/

In order to use this endpoint change the `{partial_or_complete_university_name}` with the partial or complete name of the
university you are searching for. Use space in between the words in order to search with the complete name of the university.

Example request : `http://localhost:8080/unisearcher/v1/uniinfo/stavanger`

Response: 

````
[
    {
        "name": "University of Stavanger",
        "country": "Norway",
        "isocode": "NO",
        "webpages": [
            "http://www.uis.no/"
        ],
        "languages": {
            "nno": "Norwegian Nynorsk",
            "nob": "Norwegian Bokmål",
            "smi": "Sami"
        },
        "maps": "https://www.openstreetmap.org/relation/2978650"
    }
]
````

- Content type: `application/json`


### neighbourunis

This endpoint will retrieve information about the universities in neighbouring countries, 
given that they have the same `{:partial_or_complete_university_name}` component in their institution name (for example: applied science).
This should not include universities from the given country itself.


````
Method: GET
Path: neighbourunis/{:country_name}/{:partial_or_complete_university_name}{?limit={:number}}
````

The resource root path: 
* http://localhost:8080/unisearcher/v1/neighbourunis/


To use neighbourunis your specify a country in the `{:country_name}`  component and
the partial or complete name of a university in the `{:partial_or_complete_university_name}` component.
You could also specifiy the optional component `{?limit={:number}` after the university name component, which will limit the number of universities that are retrieved. 
It is important to first specify which country in the place of the `{:country_name}` , and then which university in the `{:partial_or_complete_university_name}` component.
If you want to use the optional limit component then it has to be written right after the university component. 
An example of how to search could be `/russia/tech?limit=3`, which would retrieve three universities in the neighbouring countries of russia with the "tech" in their name.
- Note : *leave no empty spaces between components in the path when searching, unless the space is used in the `{:country_name}` or `{:partial_or_complete_university_name}` component of the path.*

Example request: `unisearcher/v1/neighbourunis/sweden/science`

Response: 

````

[
    {
        "name": "Central Ostrobothnia University of Applied Sciences",
        "country": "Finland",
        "isocode": "FI",
        "webpages": [
            "http://www.cou.fi/"
        ],
        "languages": {
            "fin": "Finnish",
            "swe": "Swedish"
        },
        "maps": "openstreetmap.org/relation/54224"
    },
    .......
    .....
    ...

````


Example request using optional parameter limit: `unisearcher/v1/neighbourunis/saudi arabia/science and technology?limit=2`
````
[
	{
		"name": "Komar University of Science and Technology",
		"country": "Iraq",
		"isocode": "IQ",
		"webpages": [
			"http://www.komar.edu.iq/"
		],
		"languages": {
			"ara": "Arabic",
			"arc": "Aramaic",
			"ckb": "Sorani"
		},
		"maps": "https://www.openstreetmap.org/relation/304934"
	},
	{
		"name": "Jordan University of Science and Technology",
		"country": "Jordan",
		"isocode": "JO",
		"webpages": [
			"http://www.just.edu.jo/"
		],
		"languages": {
			"ara": "Arabic"
		},
		"maps": "https://www.openstreetmap.org/relation/184818"
	}
]

````

- Content type: `application/json`




### diag

The diagnostic interface shows which services are required for this particular service to function. 
This is reported through the status codes received from the dependent services, along with details about the service's uptime.

The resource root path:
* http://localhost:8080/unisearcher/v1/diag/


````
Method: GET
Path: /unisearcher/v1/diag/
````


Example request: `http://localhost:8080/unisearcher/v1/diag/`

Response:
````
{
    "universitiesapi": "200 OK",
    "countriesapi": "200 OK",
    "version": "v1",
    "uptime": 166.3748636
}

````

- Content type: `application/json`
- Status code: 200 if everything is OK



## Deployment

Render, which is a cloud computing platform that provides a simple way for 
developers to deploy and manage their apps and websites, was used in this assignment to deploy the service.

- URL: `https://assignment1-j59l.onrender.com`

To use the service in the render instance, add the endpoint path after the URL.
- Example (uniinfo) : `https://assignment1-j59l.onrender.com/unisearcher/v1/uniinfo/norway`
- Example (neighbourUnis) : `https://assignment1-j59l.onrender.com/unisearcher/v1/neighbourunis/saudi%20arabia/science%20and%20technology?limit=2`

## Design choices

During the development of this application the main areas of emphasis were maintaining loose coupling, 
high cohesion, robustness and modularity that closely adheres to the conventions of Golang. 
This was achieved by utilizing techniques such as utilizing constants, 
separating handlers, structs and generic functions into different files.
To further improve this service I would implement better ways of giving feedback to the client, and fix other shortcomings 
in the service. 
