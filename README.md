# Documentation for the program
## Description:
In this project to develop a payment gateway that use Golang programming language. For the framework, I am using Gin. For the database, I am using PostgreSQL and Prisma as database frameworks. For the authentication using JWT. 

## Deployment 
- If you never use golang then you have to download golang first
- When you are in terminal can use go command then go to your project folder
- do this command to download library:
```
go mod init
go mod tidy
```
- Download Prisma in your terminal
```
npm install -g prisma
```

then tried to execute the comment below to create a database and application
```
//to generate database to the app
npx prisma generate
//build docker container
docker-compose up --build 
```
then do migrate database using the comment
```
npx prisma migrate dev
```

## API Route, Description
### User
- /user/login: login user and get token
- /user/register: register user
- /user/account/detail/:acct_num: check account detail
- /user/payment/history/:loc_acct: check transaction history by the account

### Transaction
-/transaction/send: using jwt token by adding to Authorization and add token without 'bearer '.
- /transaction/withdraw: using jwt token same use like end endpoint. 
- /transaction/detail: get detail for all transaction.
- /transaction/detail/:trx_id: get detail transaction by id.

## Transaction and Card Type
### Transaction Type
So there are 2 transaction type:
- C = Credit -> If the transaction is C, the balance will be reduced. 
- D = Debit -> If the transaction is D, the balance amount will increase.

### Card Type
There are 3 type for card type:
- C : credit card
- D : debit card
- PL: pay loan has a similarity to a credit card.

### status
- W : write Off (if write off will return cannot do transaction)
- G : good (can do transaction)

# Take home assignment


## Description:
Build 2 Backend services which manages user’s accounts and transactions (send/withdraw). 

In Account Manager service, we have:
- User: Login with Id/Password
- Payment Account: One user can have multiple accounts like credit, debit, loan...
- Payment History: Records of transactions

In Payment Manager service, we have:
- Transaction: Include basic information like amount, timestamp, toAddress, status...
- We have a core transaction process function, that will be executed by `/send` or `/withdraw` API:

```js
function processTransaction(transaction) {
    return new Promise((resolve, reject) => {
        console.log('Transaction processing started for:', transaction);

        // Simulate long running process
        setTimeout(() => {
            // After 30 seconds, we assume the transaction is processed successfully
            console.log('transaction processed for:', transaction);
            resolve(transaction);
        }, 30000); // 30 seconds
    });
}

// Example usage
let transaction = { amount: 100, currency: 'USD' }; // Sample transaction input
processTransaction(transaction)
    .then((processedTransaction) => {
        console.log('transaction processing completed for:', processedTransaction);
    })
    .catch((error) => {
        console.error('transaction processing failed:', error);
    });
```

Features:
- Users need to register/log in and then be able to call APIs.
- APIs for 2 operations send/withdraw. Account statements will be updated after the transaction is successful.
- APIs to retrieve all accounts and transactions per account of the user.
- Write Swagger docs for implemented APIs (Optional)
- Auto Debit/Recurring Payments: Users should be able to set up recurring payments. These payments will automatically be processed at specified intervals. (Optional)

### Tech-stack:
- Recommend using authentication 3rd party: Supertokens, Supabase...
- `NodeJs/Golang` for API server (`Fastify/Gin` framework is the best choices)
- `PostgreSQL/MongoDB` for Database. Recommend using `Prisma` for ORM.
- `Docker` for containerization. Recommend using `docker-compose` for running containers.
 
## Target:
- Good document/README to describe your implementation.
- Make sure app functionality works as expected. Run and test it well.
- Containerized and run the app using Docker.
- Using `docker-compose` or any automation script to run the app with single command is a plus.
- Job schedulers utilization is a plus
