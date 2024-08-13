# View info 

http://localhost:10001/.well-known/openid-configuration

Casdoor User Field	OIDC UserInfo Field
Id	sub
originBackend	iss
Aud	aud
Name	preferred_username
DisplayName	name
Email	email
Avatar	picture
Location	address
Phone	phone

# build 
go build -o idtoken idtoken/app.go
go build -o userinfo userinfo/app.go
# Run 
OAUTH2_CLIENT_ID=client1 OAUTH2_CLIENT_SECRET=89c557dbfb4494011547ec83277d35a4316583f0
