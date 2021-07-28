# Oengus patron fetcher
A microservice that fetches the patrons from [oengus.io](https://oengus.io/)

## How to run on your computer

### Requirements
- Go 1.15
- Docker

##### Docker (recommended)
This package is distributed on docker hub: https://hub.docker.com/repository/docker/oengusio/patreon-fetcher

Copy `docker-compose.yml` to `docker-compose.override.yml` and ill in the environment variables, then run `docker-compose up --build`

##### With ide (tested with GoLand)
1. Create a build configuration
2. Set the environment variables for the build configuration
3. Press run

## Support

Please join the official [Discord server](https://discord.gg/ZZFS8YT) for questions and support

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate (if tests exist).

## License
[GNU AGPL v3.0](https://choosealicense.com/licenses/agpl-3.0/)
