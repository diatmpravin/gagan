## Gagan

Gagan is a Golang package that provides a REST client for the Cloud Foundry REST API.

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

### Organization

* List All Organizations

API endpoint:

    /listallorganizations

cURL

    curl "http://localhost:8080/listallorganizations" -d '{"AccessToken":"TOKEN"}' -X POST -H "Content-Type: application/json"

### Space

* List All Spaces
 
API endpoint:

    /listallspaces

cURL

    curl "http://localhost:8080/listallspaces" -d '{"AccessToken":"TOKEN","Organization":{"Name":"ORGANIZATION-NAME","Guid":"ORGANIZATION-GUID"}}' -X POST -H "Content-Type: application/json"

### App

* List All Apps

API endpoint:

    /listallapps

cURL

    curl "http://localhost:8080/listallapps" -d '{"AccessToken":"TOKEN","Organization":{"Name":"ORGANIZATION-NAME","Guid":"ORGANIZATION-GUID"},"Space":{"Name":"SPACE-NAME","Guid":"SPACE-GUID","Applications":null,"ServiceInstances":null}}' -X POST -H "Content-Type: application/json"

* Creating An App
* Get App Summary
* Stoping An App
* Starting An App
* Delete a Praticular App
* Get The Instance Information
* Retrieve Particular App Usage Event

### Service

* Creating Service Instance
* Creating Service Binding
* Delete Particular Service Binding
* Delete Particular Service

