# Host a Hack

A Hassle free hosting solution for your hackathon projects! [Project link](github.com/chebro/host-a-hack) (https://github.com/chebro/host-a-hack).

### To run

```bash
docker build -t hostahack:latest .
go run .
```

## Inspiration

One of the most challenging tasks for a begginer participating in a hackathon is to figure out how to host their hackathon idea. Moving from `localhost:8080` to `mydomain.com` is not an easy task, figuring out how to work with virtual machinces, finding a domain provider and understading the nitty gritties of DNS management. It's all a huge hassle that takes up a significant amount of time that you can otherwise dedicate to developing your product.

So we decided to make it easy for out fellow hackers by creating an all in one platform that can do all of this for you.

Our motto is simple: if you can host it on `localhost`, let us take care of the rest.


## What it does?

Our platform provides:
- a shell for any user to jump right in.
- upload their projects and install the dependencies on our server, just like you would in your local machine.

And everything is setup and the user hits the run command, we take care of the rest: 
- a private subdomain for your project is generated on the fly 
- any open port gets immediately picked up by our reverse proxy and you can see your project on the web!

Our platform also provides supporting guides for new users to familiarize themselves with hosting their projects on the top cloud platforms.

## How we built it

Our backend was purely written in golang as we wanted to explore a new language for this hackathon. We use docker containers to host the projects submitted by our users. And all the port forwarding from the containers to the hosts is handled by Nginx.

Our idea was to make it extremely easy for a new user to jump in, so we create a fresh container for our users to hack away right as they visit our website. A project upload button dumps the project files directly into their freshly generated containers.

## Challenges we ran into

- The biggest challenge was to figure out a way to make the whole process instantaneous for our users, so we decided to allocate a container pool which is maintained by our server side go program that ensures that there are enough containers available at all times.
- We faced a huge hurdle while sketching up the solution for generating a weblink for each container, since our reverse proxy has to handle webservers running on every container.
- Not having dedicated front-end engineers in our team made it difficult for us to make an even better and polished user experience.

## Accomplishments that we're proud of

- Within the small amount of time we were able to build an entire project in a fairly new language, and were able to host and run dynamic websites within minutes just as out motto promises.
- Even though we didn't have many front-end devs on our we team, we are proud of how our website ended up looking!

## What's next for Host a Hack

There's a lot of future scope for the project as our supported tech stacks is quite limited. We can definitely make progress in improving our website UI and UX department. We can also add aditional guides for new comers to get used to our platform and how they can leverage it to host multiple projects.

## Built With

golang, nginx, docker, html, css, javascript
