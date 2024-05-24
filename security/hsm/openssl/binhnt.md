1. Generate RSA key pairs

openssl genrsa -out Alice/private.pem
openssl rsa -in Alice/private.pem -pubout -out Alice/public.pem

openssl genrsa -out Bob/private.pem
openssl rsa -in Bob/private.pem -pubout -out Bob/public.pem

2. Exchange public keys
cp Alice/public.pem Bob/public_alice.pem
cp Bob/public.pem Alice/public_bob.pem

1. Generate AES key

AES_KEY=$(openssl rand -hex 32)
AES_IV=$(openssl rand -hex 16)

echo $AES_KEY
9567260a75a667062d046cabef28f855b8df9f2a2df25256fdd2f4fa3a77e092

echo $AES_IV
163f32eb2e45505ac6e7ac377ecd7971

4. Encrypt data with AES key
openssl enc -aes-256-cbc -K $AES_KEY -iv $AES_IV -in Alice/config.yaml -out Alice/config.enc

5. Encrypt AES key with RSA
echo $AES_KEY | openssl pkeyutl -encrypt -pubin -inkey Alice/public_bob.pem  -out Alice/aes_key_bob.enc

1. Sign message
openssl dgst -sha256 -sign Alice/private.pem -out Alice/aes_key_bob.enc.sig Alice/aes_key_bob.enc

1. Send encrypted data, encrypted AES key, signature and iv
- config.enc
- aes_key_bob.enc
- aes_key_bob.enc.sig
- plain text iv

cp Alice/config.enc Bob/config.enc
cp Alice/aes_key_bob.enc Bob/aes_key_alice_bob.enc
cp Alice/aes_key_bob.enc.sig Bob/aes_key_alice_bob.enc.sig

8. Verify signature
openssl dgst -sha256 -verify Bob/public_alice.pem -signature aes_key_alice_bob.enc.sig aes_key_alice_bob.enc

9. Decrypt AES key
AES_KEY=$(openssl pkeyutl -decrypt -inkey Bob/private.pem -in Bob/aes_key_alice_bob.enc)

10. Decrypt data
openssl enc -d -aes-256-cbc -K $AES_KEY -iv $AES_IV -in Bob/config.enc -out Bob/config1.yaml




