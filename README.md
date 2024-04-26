# go-boilerplate
### What is this repository for? ###

go-boilerplate with DDD architecture which will be base for the projects going forward.

### How do I get set up? ###

#### Prerequisites ####
Golang with version 1.19
or Docker with the latest version

#### There are 2 ways of setting up the service ####
* Using Dockerfile, 
<br />
1. Update your bitbucket username & password with your credentials
<br />
2. Run ``docker build -t go-boilerplate .`` -> this command builds the docker image required for running our service.
<br />
3. Run ``docker run -p 8080:8080 go-boilerplate`` -> this command starts the go-boilerplate inside a docker container.
* Using Go,
<br />
1. Run ``export GOPRIVATE="bitbucket.org/kodnest"``
<br />
2. Run ``git config --global url."git@bitbucket.org”.insteadOf "https://bitbucket.org" ``
<br />
3. Above steps are required as we are using our own common go-packages
<br />
4. Finally, run ``go run main.go``

#### Dependencies ####
Our own go-common-libraries, https://bitbucket.org/kodnest/go-common-libraries/src/master/

#### For BackEnd developers I would suggest the 2nd step to set up the service. ####