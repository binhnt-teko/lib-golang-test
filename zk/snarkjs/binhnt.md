# Install circom 

- Install dependencies:
  
  curl --proto '=https' --tlsv1.2 https://sh.rustup.rs -sSf | sh

- Installing circom:
git clone https://github.com/iden3/circom.git
cargo build --release
cargo install --path circom

- Installing snarkjs
npm install -g snarkjs

# Compile circuit 
circom circuits/multiplier2.circom --r1cs --wasm --sym --c 
circom circuits/multiplier2.circom --r1cs --wasm --sym --c

# Computing the witness 
- with WebAssembly
cd multiplier2_js
node generate_witness.js multiplier2.wasm input.json witness.wtns
- with C++
cd multiplier2_cpp 
make 

# Proving circuits
- Powers of Tau: which is independent of the circuit
 +   new "powers of tau" ceremony 
snarkjs powersoftau new bn128 12 pot12_0000.ptau -v
 + contribute to the ceremony
  snarkjs powersoftau contribute pot12_0000.ptau pot12_0001.ptau --name="First contribution" -v

+ Phase 2: which depends on the circuit
snarkjs powersoftau prepare phase2 pot12_0001.ptau pot12_final.ptau -v

- generate a .zkey file that will contain the proving and verification keys together with all phase 2 contributions:

snarkjs groth16 setup multiplier2.r1cs pot12_final.ptau multiplier2_0000.zkey

- Contribute to the phase 2 of the ceremony:

snarkjs zkey contribute multiplier2_0000.zkey multiplier2_0001.zkey --name="1st Contributor Name" -v

- Export the verification key: 112233

snarkjs zkey export verificationkey multiplier2_0001.zkey verification_key.json

## Generating a Proof
snarkjs groth16 prove multiplier2_0001.zkey multiplier2_js/witness.wtns proof.json public.json

proof.json: it contains the proof.
public.json: it contains the values of the public inputs and outputs.


## Verifying a Proof
snarkjs groth16 verify verification_key.json public.json proof.json

## Verifying from a Smart Contract
- First, we need to generate the Solidity code using the command:

snarkjs zkey export solidityverifier multiplier2_0001.zkey verifier.sol

- call
snarkjs generatecall
