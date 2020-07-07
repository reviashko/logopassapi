Ligthweight API template "logopassapi" for users login/registration/restore password methods

API using PG database which you need to install before using API - database sql script found in db folder.
Database setting stored in config(db json):
    "host":"pg.test.test" --> pg address server
    "port":5432 --> pg port
    "dbname":"test" --> database name
    "user":"test" --> pg user login
    "password":"password" --> pg user password

For Email sending you need to set your email settings in smtp config (smtp.json):
    "host":"smtp.test.test" --> SMTP server address. Request it from yours EMAIL service provider
    "email":"test@test.test" --> your accaunt on SMTP server (yours EMAIL address)
    "password":"password" --> your email accaunt password
    "mock_email":"test@mailinator.com" --> email for mock testing

For token protection API using AES256 encryption algorithm. Please do not use default secrets settings. Crypto config (crypto.json) contain:
    "AES256Key":"mysuperpupersecretkeywith32len00" --> your secret phrase (need 32 character!)
    "SHA256Salt":"megaSalt" --> API store only password hash in binary string, not password. For SHA256 use "salt" to improve security
    "TokenTTL":3600 --> ttl (time to live in seconds) for given token
    "PasswordEmailTTL":60 --> ttl for restore password link. Set in seconds
    "RestorePasswordURL":"http://localhost:3000/auth/changepassword/" --> change password link, which sending via email 

GetTestDataByTokenHandler method is an example of how to send data from API throught token authorization
