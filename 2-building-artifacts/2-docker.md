# Running Go programs in Docker containers

Before looking at building Docker images, let's use it to compile go code

## Compile using golang image
Using golang image to build binaries can be useful for :
- Continuous integration
- Cross building binary
- Local development

### Continuous integration
When using a CI server like Jenkins, it can become cumbersome to have to install all the go versions required by your different projects.
And even more if you use your CI server to build projects in other languages, each with its own version (java8, java7, node, ruby, python, ...).
For this use case, it can be easier to build the projects in docker containers. An example script could be :
```bash
 $ docker run --rm -v $PWD:/go/src/github.com/repo/project -w /go/src/github.com/repo/project golang:1.8.1 go build
```
This will run the `go build` command with go version 1.8.1, sharing your current directory with the container. Obviously, you can run `go fmt/vet/test/...` commands the same way

### Cross building binary
Usually, cross building a binary requires to install go from sources. Using the golang image, you can install the binary distribution on your computer for day-to-day use and rely on the docker image to build binaries for different targets.
As seen on previous section, cross compiling binaries is only a matter of environment variables. So, an example command could be :
```bash
 $ docker run --rm -v $PWD:/go/src/github.com/repo/project -w /go/src/github.com/repo/project -e GOOS=windows -e GOARCH=amd64 golang:1.8.1 go build -o dist/binary-windows.exe
 $ docker run --rm -v $PWD:/go/src/github.com/repo/project -w /go/src/github.com/repo/project -e GOOS=linux -e GOARCH=amd64 golang:1.8.1 go build -o dist/binary-linux-x64
 $ docker run --rm -v $PWD:/go/src/github.com/repo/project -w /go/src/github.com/repo/project -e GOOS=darwin -e GOARCH=amd64 golang:1.8.1 go build -o dist/binary-darwin
```

### Local development
Another use case would be to build a custom image (based on golang) including the version of go of your choice and all the tools your need to lint / check / test your code. Then you can distribute this image to all the developers in your team so you can easily manage environment development for all your team from a simple Dockerfile !

## Build docker images
When comes time to distribute your awesome project, you might want to provide a docker image instead of a binary. It can ease the installation instruction (particularly if your binary needs to be deployed alongside other files like a webserver serving static files) or to run it in a kubernetes cluster.
Let's see the options we have to achieve this.

### FROM alpine
Using the base image from which to build your own is a matter of taste. Nethertheless, alpine is massively adopted, even by docker who publishes an alpine flavor for all the official images. This is because the alpine base image is a lot smaller than any other linux distribution.
Assuming you've already built your binary in your current directory, an example Dockerfile would be :
```
FROM alpine:latest
ADD mybinary .
CMD ./mybinary
```

### FROM scratch
`scratch` is a special keyword to tell docker to build an image from "nothing". So the resulting image will only contain the files and directories you explicitly add. As go binaries don't require anything to run, you can be tempted to build your image from scratch instead of alpine, this will produce an image as small as your binary. There are just a few tricks you should be aware of :
- To have a true autonomous binary that don't rely on any linux libray, your must build your binary with `CGO_ENABLED=0` flag and ` -installsuffix cgo` option
- If your application needs to make HTTPS calls, you need to add an up-to-date `/etc/ssl/certs/ca-certificates.crt` file in your container so your application can validate https certs.

Also, as `scratch` is not really a distribution, don't expect to find any command like `tar, unzip, curl, apt-get`... So scratch is really useful to package your binary but can be trcky to use if your container must embed more stuffs.


### Multistage build
Until docker 17.05-ce, if you wanted to build your binary and tiny images you had to :
- Use a `docker run` command as explained earlier to generate a binary in your current directory
- and then run a `docker build` based from alpine or scratch and copy the binary obtained in the previous step into it.

Now, docker support multi-stage builds, which means you can use a single Dockerfile for those 2 steps, eg :
```
FROM golang:latest as builder
ADD . /go/src/github.com/repo/project
RUN go build -o /tmp/binary .

FROM alpine:latest
COPY --from=builder /tmp/binary /usr/local/binary
CMD /usr/local/binary
```
Running `docker build` will build the binary in a temp image, then build your final image and copy the binary from the temp image into the final one !

The only drawback of the multistage approach is you will run a single `docker build` command for the whole process. If your build process produces additional files such as test coverage reports, you won't be able to get them back. Indeed, `docker build` don't support volume mounts. So you can `ADD` files from your current directory into the image but you can't extract generated artifacts from the temp images (Obviously, you can still use `docker cp` from the final image if you wanted to extract your binary or something else).

### Best practices
Docker provides a lot of [best pratcices](https://docs.docker.com/engine/userguide/eng-image/dockerfile_best-practices/). One particularly cool to use if your binary provides multiple commands is the use of the entrypoint.
With a Dockerfile like
```
FROM
...
ENTRYPOINT [/usr/local/binary]
CMD [run]
```
the user can run your binary with a single `docker run imagename`. But if your binary supports other commands like `help`, it can be simply invoked with `docker run imagename help`.  


## Congratulations

You're now able to build Docker images that not only run your Go programs correctly,
but also they are as light as possible, and easy to maintain!

Next, let's talk about how to better understand the behavior of running programs in
the [dynamic analysis section](../3-dynamic-analysis/README.md).
