# saarCTF 2020

Services from [saarCTF 2020](https://ctftime.org/event/980/)

## Prerequisites
To build any of the services, first build the base-image
```bash
cd base-image
bash docker-build.sh
```

## Building services
Once the base-image has been built, enter a service directory and use `docker-compose`, e.g.:
```bash
cd saarXiv
docker-compose up --build -d
```

## Running checkers
Every service comes with a `checkers` directory, which contains a python-script named after the service.
Running this script should place a flag in the service and try to retrieve it subsequently.
Caveat: Make sure the `gamelib` is in the `PYTHONPATH`, e.g.:
```bash
PYTHONPATH=../../ python saarxiv.py
```
