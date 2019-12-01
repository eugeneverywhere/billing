# Billing microservice

Service that listens on rabbit queue for commands, performs operations with accounts and sends results 
into specified rabbit channels for unsuccessful and successful results.

## Interacting format

#### Operation header data

Every input message must contain these fields:

```
{
	"cons_id": 0, // Consumer id
	"op_id": 15,  // Operation id
	"op_code": 3, // Operation code
...
}
```

#### Response messages

Response for each message is reported in output_channel or error_channel.
Response message contains operation header data for tracking.

##### Response example
```
{
	"cons_id": 0,
	"op_id": 5,
	"op_code": 1,
	"res_code": 0, //Result code, 0 corresponds to Ok
	"msg": "Account FTA1 created"
}
```
Common possible error res_codes:

	ErrWrongFormat          = 101
	ErrUnknownOperationCode = 102

#### Create account

Operation code: 1
Example:
```
{
	"cons_id": 0,
	"op_id": 0,
	"op_code": 1, // Operation code 1
	"acc_id": "FTA1" // External Account ID
}
```
Possible error codes:

	ErrAccountAlreadyExists = 103
	ErrEmptyID              = 106
	ErrSpaces               = 107
	ErrIDTooLong            = 109 // ID must be less than 20 symbols

#### Add amount to account

Operation code: 2
```
{
	"cons_id": 0,
	"op_id": 2,
	"op_code": 2,
	"acc_id": "FTA1",
	"amount": -100.20 // Negative amount may be used to write-off funds
}
```

Possible error codes:

	ErrAccountDoesNotExist  = 104
	ErrInsufficient         = 105

#### Transfer amount between accounts

Operation code: 3
```
{
	"cons_id": 0,
	"op_id": 15,
	"op_code": 3,
	"amount": 1000,   // Negative amount here is forbidden
	"src_id": "FTA1", //Source account
	"tgt_id": "FTA2"  //Target account
}
```

Possible error codes:

	ErrAccountAlreadyExists = 103
	ErrAccountDoesNotExist  = 104
	ErrInsufficient         = 105

## Tests
To run unit tests:
```bash
make test
```

## Launch
To setup local environment MySQL and RabbitMQ instances required, 
you may use docker-compose.yml from repo

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