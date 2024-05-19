
# Add key to hsm 
softhsm2-util --init-token --slot 1 --label "configkey" --so-pin 5462 --pin 8764329

/opt/hyperledger/softhsm2/bin/softhsm2-util --init-token  --free --token orderer --label orderer --so-pin 123456 --pin 654321

/opt/hyperledger/softhsm2/bin/softhsm2-util --show-slots | grep orderer |  wc -l 


The token has been initialized and is reassigned to slot 2058063310

softhsm2-util --pin 8764329 --import test/hsm/key/api.pem --token configkey --label configkey --id A1 --no-public-key

Found slot 2058063310 with matching token label.


softhsm2-util --init-token --slot 2 --label "configkey2" --so-pin 5462 --pin 8764329
The token has been initialized and is reassigned to slot 777297901


softhsm2-util --pin 8764329 --import test/hsm/key/session.key --token configkey2 --label configkey2 --id A2 --no-public-key



softhsm2-util --pin 8764329  --token configkey2  

Found slot 777297901 with matching token label.



softhsm2-util --init-token --slot 3 --label "configkey3" --so-pin 5462 --pin 8764329
The token has been initialized and is reassigned to slot 424394087

softhsm2-util --pin 8764329 --import test/hsm/key/api.pem --token configkey3  --label configkey3 --id A3 --no-public-key

# Add private key to hsm  

softhsm2-util --init-token --free  --token test_private --label "test_private" --so-pin 5462 --pin 8764329

softhsm2-util --pin 8764329 --import test/hsm/cert/priv_sk --token test_private  --label test_private --id A3 --no-public-key


softhsm2-util --init-token --free  --token test_import --label "test_import" --so-pin 5462 --pin 8764329

# Show slots m
softhsm2-util --show-slots


softhsm2-util --init-token --label "configkey" --token "configkey"  --so-pin 5462 --pin 8764329

# Delete token 

for i in $(./softhsm2-util --show-slots | grep Serial | cut -d ":"  -f 2 |  awk '{$1=$1};1' ) ; do  ./softhsm2-util  --delete-token  --serial $i  ;  done  



