## Gagan

Gagan is a Golang package that provides a REST client for the Cloud Foundry REST API.

## Sample application

We are developing a open source dashboard [Aakash](https://github.com/diatmpravin/aakash) for cloudfoudry by using this REST API.

## Development tools

This package was created using Golang and Redis.

## Changes

Change the existing Target and AuthorizationEndpoint 

    configuration/configuration.go

### Example:

    func GetDefaultConfig() (c *Configuration) {
        c = new(Configuration)
        c.Target = "https://api.run.pivotal.io"
        c.ApiVersion = "2"
        c.AuthorizationEndpoint = "https://login.run.pivotal.io"
        return
    }

## Cloud Foundry API

### Session

* Create Session
 
API endpoint:

    /session/new

cURL

    curl "http://localhost:8080/session/new" -d '{"email":"pravinmishra_88@yahoo.com","password":"cf@rest12"}' -X POST -H "Content-Type: application/json"

* Delete Session

API endpoint:

    /session/delete

cURL

    curl "http://localhost:8080/session/delete?sessionid=8" -X DELETE -H "Content-Type: application/json"



### Organization

* List All Organizations

API endpoint:

    /listallorganizations

cURL

    curl "http://localhost:8080/listallorganizations" -d '{"sessionid":20}' -X POST -H "Content-Type: application/json"

### Space

* List All Spaces
 
API endpoint:

    /listallspaces

cURL

    curl "http://localhost:8080/listallspaces" -d '{"sessionid":20,"timestamp":"2016-01-16T16:49:16.2528982+05:30","Organization":{"Name":"diatm1","Guid":"a2dd3cd3-d85b-4fe4-8c53-8fb57da4aafb"}}' -X POST -H "Content-Type: application/json"

### App

* List All Apps

API endpoint:

    /listallapps

cURL

    curl "http://localhost:8080/listallapps" -d '{"sessionid":20,"accesstoken":"","timestamp":"2016-01-16T16:49:16.2528982+05:30","Organization":{"Name":"diatm1","Guid":"a2dd3cd3-d85b-4fe4-8c53-8fb57da4aafb"},"Space":{"Name":"development","Guid":"25cd438e-e0b3-476c-8998-261b9f8e807c","Applications":null,"ServiceInstances":null}}' -X POST -H "Content-Type: application/json"

* Creating An App
* Get App Summary

API endpoint:

    /getappsummary?appname=APPNAME

cURL

    curl "http://localhost:8080/getappsummary?appname=APPNAME" -d '{"sessionid":20,"accesstoken":"","timestamp":"2016-01-16T16:49:16.2528982+05:30","Organization":{"Name":"diatm1","Guid":"a2dd3cd3-d85b-4fe4-8c53-8fb57da4aafb"},"Space":{"Name":"development","Guid":"25cd438e-e0b3-476c-8998-261b9f8e807c","Applications":null,"ServiceInstances":null}}' -X POST -H "Content-Type: application/json"

* Stoping An App

API endpoint:

    /stopinganapp?appname=APPNAME

cURL

    curl "http://localhost:8080/stopinganapp?appname=APPNAME" -d '{"sessionid":20,"accesstoken":"","timestamp":"2016-01-16T16:49:16.2528982+05:30","Organization":{"Name":"diatm1","Guid":"a2dd3cd3-d85b-4fe4-8c53-8fb57da4aafb"},"Space":{"Name":"development","Guid":"25cd438e-e0b3-476c-8998-261b9f8e807c","Applications":null,"ServiceInstances":null}}' -X POST -H "Content-Type: application/json"

* Starting An App

API endpoint:

    /startinganapp?appname=APPNAME

cURL

    curl "http://localhost:8080/startinganapp?appname=APPNAME" -d '{"sessionid":20,"accesstoken":"","timestamp":"2016-01-16T16:49:16.2528982+05:30","Organization":{"Name":"diatm1","Guid":"a2dd3cd3-d85b-4fe4-8c53-8fb57da4aafb"},"Space":{"Name":"development","Guid":"25cd438e-e0b3-476c-8998-261b9f8e807c","Applications":null,"ServiceInstances":null}}' -X POST -H "Content-Type: application/json"

* Delete a Praticular App

API endpoint:

    /deleteapraticularapp?appname=APPNAME

cURL

    curl "http://localhost:8080/deleteapraticularapp?appname=APPNAME" -d '{"sessionid":20,"accesstoken":"","timestamp":"2016-01-16T16:49:16.2528982+05:30","Organization":{"Name":"diatm1","Guid":"a2dd3cd3-d85b-4fe4-8c53-8fb57da4aafb"},"Space":{"Name":"development","Guid":"25cd438e-e0b3-476c-8998-261b9f8e807c","Applications":null,"ServiceInstances":null}}' -X POST -H "Content-Type: application/json"

* Get The Instance Information

API endpoint:

    /gettheinstanceinformation?appname=APPNAME

cURL

    curl "http://localhost:8080/gettheinstanceinformation?appname=APPNAME" -d '{"sessionid":20,"accesstoken":"","timestamp":"2016-01-16T16:49:16.2528982+05:30","Organization":{"Name":"diatm1","Guid":"a2dd3cd3-d85b-4fe4-8c53-8fb57da4aafb"},"Space":{"Name":"development","Guid":"25cd438e-e0b3-476c-8998-261b9f8e807c","Applications":null,"ServiceInstances":null}}' -X POST -H "Content-Type: application/json"

* Retrieve Particular App Usage Event

### Service

* Creating Service Instance
* Creating Service Binding
* Delete Particular Service Binding
* Delete Particular Service

## Contributing

We encourage you to contribute to Gagan! Please check out the [Contributing to Gagan guide]() for guidelines about how to proceed.

## License

Gagan is released under the [MIT License](http://www.opensource.org/licenses/MIT).