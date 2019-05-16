# Perceptia API

The Perceptia API is a REST based api. The purpose of this document is to define the general syntax for the API and provide references to the various API specifications that make up the Perceptia API.

## [Contents](#Contents)

* [Overview](#overview)

* [Syntax](#syntax)

* [Making a Request](#making-a-request)

* [API Specifications](#api-specifications)

* [Internal API Specifications](#internal-api-specifications)

## [Overview](#overview)

The Perceptia API (API) is designed to service the Perceptia Application, which is currently being developed as a Web Application. This service, and the APIs it exposes are not intended for use by thrid-parties.

The API is a REST based archetecture, with JSON as the primary format for request and response communication. As the API is still under development, breaking changes may occur. For the current API specs, see the individual OpenApi specification documents located in their respecitive directories (for which a listing can be found: [API Specifications](#api-specifications)).

Note, each services repository may contain an OpenApi Yaml file `*-api.yaml` defining the API the code implements.

## [Syntax](#Syntax)

### [URL Syntax](#url-syntax)

Basic URL: `<scheme>://<host>:<port>/<path>?<queryParameters>`

Example: `https://api.perceptia.info:443/api/v{majorVersion}/{collection}/{collectionSpecific}?{queryParameters}`

Syntax of Perceptia API: `/api/v{majorVersion}/{collection}/{collectionSpecific}?{queryParameters}`

Example: `/api/v1/anyquiz/read/apple`

Meanings:

   `majorVersion:` this is the major version of the api being called, such as "v1", and its impact on the request is dependent on the collection being requested. Its format should be the character "v" followed by an intiger, with no spaces or other characters. Minor version and Patch version values should not be included here (see [Common Elements](#common-elements))

   `collection:` this is the resource being acted on, such as "anyquiz". See [API Specifications](#api-specifications) for possible collection values

   `collectionSpecific:` any part of the path after the collection, whose impact on the request is dependent on the collection being requested

   `queryParameters:` query parameters may be used by the collection resource and/or by the the gateway. See [API Specifications](#api-specifications) for possible parameters used by a collection

### [Misc Syntax](#misc-syntax)

TODO: Describe other common elements of the API's syntax (such as common headers, parameters)

## [Common Elements](#common-elements)

This section lists the common elements used in the API, including query parameters.

### [Query Parameters](#query-parameters)

This subsection lists the query parameters that are common to all API calls. Note, parameters may not be used by all API collections (in which case their inclusion will have no effect on the processing of the request)

Syntax: `/api/v{majorVersion}/{collection}/{collectionSpecific}?{queryParameters}`

Where `queryParameters:` are a set of key value pairs seperated by &

Example: `/api/v1/anyquiz/read/apple?param1=val1&param2=val2`

* [apiVersion](#params-api-version)

#### [Params API Version](#params-api-version)

Parameter: `apiVersion={major.minor.patch}:` where major.minor.patch denotes the minimum required version for the API that this query should be processed by. This parameter can be used to ensure the query won't be run by a version earlier than the one specified

Example: `/api/v1/anyquiz/read/apple?apiVersion=1.0.0`

#### [Params Auth Token](#params-auth-token)

Parameter: `access_token={access token}:` where access token is the access token the client wants to authenticate with. Note, the client should only use this if they are unable to use the Authroization header

Example: `/api/v1/anyquiz/read/apple?access_token=`

## [Making a Request](#making-a-request)

The Perceptia API is available on the public internet. There is currently no access restrictions, although this may change in the future. This section will provide example requests to demonstrate how to make requests to the API.

Production API endpoint: `https://api.perceptia.info:443/api`

Development API endpoint: `https://api.dev.perceptia.info:443/api`

Example request: `/api/v1/anyquiz/read/apple`

Make request using curl: `curl -X GET "https://api.perceptia.info/api/v1/gateway/health`

Example response: `{"name":"Perceptia API Health Report","version":"0.2.0","status":"ready"}`

## [API Specifications](#api-specifications)

The Perceptia API is documented using the [OpenApi standard](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.0.2.md), version 3. Each service of the Perceptia application maintains its own API specification in its respective directory. Additionally, each specification is maintained in the [api](./../api/) directory of this repository.

The [api](./../api/) directory is organized to reflect the major version of the API, such that routes with the {apiVersion} of "v1" are located in the directory [v1](./v1). The version subdirectory is further organized by collection (service or resource), such that the gateway service has its API archived in the [v1/gateway](./v1/gateway) directory. Note, the files in the collection subdirectory are labeled based on the version of the API they document, using the [semver](https://semver.org/spec/v2.0.0.html) format. The `*.yaml` files are the OpenApi specification, and the `*.html` files are the visual form of the specification (Note, opening html file on github will load raw html, instead open html files in browser from local clone of repository).

Note, each collection is versioned independently of each other, and thus versions are specific to a given collection, and each collection handles the processing of different major and minor versions of their respective API specifications.

A note about API versions below 1, such as 0.1.0. These versions are under development, and thus breaking changes may be made. Once the API reaches version 1, no breaking changes will be made, only non-breaking fixes and new features will be added. Version 0 APIs are served under the "v1" url path.

### [Gateway Service API](#gateway-service-api)

[Version 1](./v1/gateway)

* 0.1.0 - [API Specification](./v1/gateway/0.1.0.yaml) | [API Documentation](./v1/gateway/0.1.0.html)

  * 0.1.1 - [API Specification](./v1/gateway/0.1.1.yaml) | [API Documentation](./v1/gateway/0.1.1.html)

* 0.2.0 - [API Specification](./v1/gateway/0.2.0.yaml) | [API Documentation](./v1/gateway/0.2.0.html)

* 0.3.0 - [API Specification](./v1/gateway/0.3.0.yaml) | [API Documentation](./v1/gateway/0.3.0.html)

[Gateway Service API - Current](./../gateway/gateway-service-api.yaml)

### [AnyQuiz Service API](#anyquiz-service-api)

* 0.0.0 - [API Specification](./v1/anyquiz/0.1.0.yaml) | [API Documentation](./v1/anyquiz/0.1.0.html) THIS IS A PLACEHOLDER

## [Internal API Specifications](#internal-api-specifications)

The purpose of this section is to define how the api processes requests internally and how internal services communicate.

### [Authentication](#authentication)

The gateway service handles all authentication for the Perceptia API. For informaiton on how an API client authenticates, see the API spec for the gateway, version 0.3.0 or greater. Once the client has authenticated with the system, an authentication token is generated and returned to the client in an Authroization header. This token should then be supplied in an Authorization header on each subsequent request to authenticate the user without the user having to provide their credentials again.