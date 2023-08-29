# Blog system using Echo Framework

## How to Run?

1. Clone project
2. Create new database
3. Create `.env` file (example contents are in `.env.example`) and customize the values.
4. Open terminal and then run :

```sh
   go mod tidy
```

5. Running app use command :

```sh
   go run main.go
```

6. Recommended flow

-  Post new user : `\api\v1\users`
   ex. :
   {
   "Username": "Lana",
   "Email": "lana@mail.com",
   "Password": "123abc"
   }
-  Login to get the token :`\api\v1\login`
-  In Authorization/Auth, select Bearer Token and enter the token
-  Enjoy to try other API
