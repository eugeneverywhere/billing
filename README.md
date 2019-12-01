# Billing microservice

Service that listens on rabbit queue for commands, performs operations with accounts and sends results 
into specified rabbit channels for unsuccessful and successful results.

## Available operations
#### Create account

```json
Here goes your json object definition
```

#### Create account

```json
Here goes your json object definition
```

#### Add amount to account

```json
Here goes your json object definition
```

#### Transfer amount between accounts

```json
Here goes your json object definition
```


## Work flow

### Tests

## Local launch

To setup local environment MySQL and RabbitMQ instances required, you may use docker-compose.yml from repo


### Requirements
    
### Launch 

1. Clone the repository.
2. Install dependencies, create the config file.
3. For first time make database migrations
4. Launch the project.

```bash
git clone https://github.com/eugeneverywhere/billing.git && cd billing
make setup && make config
make db:migrate
make run
```