# gURL (go cURL)

gURL is currently in incubation state. Do not trust `master` until it has been tagged with a version.

gURL started as a way to easily test AWS Lambda end-points. The intent is to make token and IAM authentication to AWS, from the command line easier.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

What things you need to install the software and how to install them

- Go ([download](https://golang.org/dl))

### Installing gURL

```
go get github.com/dbyington/gurl
```


```
$ gurl https://api.chucknorris.io/jokes/random
{"category":null,"icon_url":"https:\/\/assets.chucknorris.host\/img\/avatar\/chuck-norris.png","id":"7NDkvIo_RJGxUq183RJNhg","url":"https:\/\/api.chucknorris.io\/jokes\/7NDkvIo_RJGxUq183RJNhg","value":"Chuck Norris has Hitler\u0027s skull hanging off of his key-chain."}
```

## Running the tests

TODO: Explain how to run the automated tests for this system

### Break down into end to end tests

TODO: Explain what these tests test and why

```
Give an example
```



## Built With

* [go](https://golang.org/) - The web framework used

## Contributing

TODO: Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/dbyington/gurl/tags). 

## Authors

* **Don Byington** - *Initial work* - [dbyington](https://github.com/dbyington)

See also the list of [contributors](contributors.md) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments

* Google for sponsoring Go
