<div align="center">
<img height="150px" src="https://raw.githubusercontent.com/wisdom-oss/brand/main/svg/standalone_color.svg">
<h1>Water Rights</h1>
<h3>service-water-rights</h3>
<p>🚰 reading parsed and crawled water right information</p>
<img src="https://img.shields.io/github/go-mod/go-version/wisdom-oss/service-water-rights?style=for-the-badge"
alt="Go Lang Version"/>
<a href="openapi.yaml">
<img src="https://img.shields.io/badge/Schema%20Version-3.0.0-6BA539?style=for-the-badge&logo=OpenAPI%20Initiative" alt="Open
API Schema Version"/></a>
</div>

## About
This service allows users to read the water rights that have been parsed using
[nlwkn-rs](https://github.com/wisdom-oss/nlwkn-rs) and its predecessor.
Currently, the service focuses on the predecessor's format and data types.
However, this may be changed in the future.

## Using the service
The service may be accessed using the [api documentation](openapi.yaml).
The service is present on the demonstration system.
However, the service is not included in every standard deployment, due to the
data not being delivered by the service and a required manual crawling process.