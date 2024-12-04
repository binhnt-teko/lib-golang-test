# Structure
- common:
- crypto:
- + ckd: child key derivation 
- + commitments: commitment function 
- + dlnproof: dlnproof function
- -- NewDLNProof
- + facproof: 
- + modproof: 
- + mta: implement MtA with: 
  - -- AliceInit
  - -- BobMidWC
  - -- AliceEnd
  - -- AliceEndWC
  - -- ProveRangeAlice
  - -- RangeProofAliceFromBytes
  - -- ProveBobWC
  - -- ProveBob
  - -- ProofBobWCFromBytes 
- + pailier: pailier encrypt and decrypt 
- + schnorr: schnorr proof 
- + vss: implement feldman - vss 
- 
# Gen Key 
## ecdsa 
  0. Prepare 
   - prepare for concurrent Paillier and safe prime generation
   - generating the Paillier modulus
   - generate safe primes for ZKPs
  1. Round 1
   - calculate "partial" key share ui
   - compute the vss shares
       - generate Paillier public key E_i, private key and proof
       - generate safe primes for ZKPs
       - compute ntilde, h1, h2
   - generate the dlnproofs for keygen
   - save for this party:
     - + shareID 
     - + Vs
     - + Shamir share 
     - + de-commitments, paillier keys for round 2
   - BROADCAST: commitments, paillier pk + proof
  2. Round 2
    -  verify dln proofs, store r1 message pieces, ensure uniqueness of h1j, h2j => save NTilde_j, h1_j, h2_j, ...
    -  p2p send share ij to Pj
       -  + shares[j]
       -  + facProof
    - BROADCAST: de-commitments of Shamir poly*G:
      -  + deCommitPolyG
      -  + modProof
  3. Round 3
    - calculate xi = sum of all share 
    - Create Vc from vs 
    - vss.Share verify
    - 12-16. compute Xj for each Pj 
    - compute and SAVE the ECDSA public key `y`
    - BROADCAST: paillier proof for Pi
  4. Round 4 
   - r3 messages are assumed to be available and != nil in this function
   - consume unbuffered channels (end the goroutines)


# Sign 
## eddsa 
0. Prepare 
   - PrepareForSigning(), GG18Spec (11) Fig. 14
1. Round 1
   - New HashCommitment
   - Send AliceInit => cA, pi => cA
2. Round 2
   - Get rangeProofAliceJ
   - BobMid => get betas, c1jis, pi1jis
   - Bob_mid_wc => get vs, c2jis, pi2jis
   - create and send messages
3. Round 3
    - get proofBob
    - call mta.AliceEnd  => alphas 
    - Alice_end_wc => us
  
4. Round 4
   - calculate thetaInverse
   - call schnorr.NewZKProof => piGamma 
  
5. Round 5 : Create signature[i]
   - decommit => get bigGammaJ => bigGammaJPoint
   - get proof => verify ZK Proof 
   - Cal R = sum(R of party) => R= R * thetaInverse => rx := R.X() => Cal Si
   - random li, roI => liPoint, bigAi => bigVi
   - call hashCommitment: bigVi, bigVi, bigAi, bigAi
  
6. Round 6
    - call ZKProof (schnorr) => piAi, piV
    - signed message of other parties
7. Round 7
    - verify ZKProof 
    - get VjX, VjY, AjX, AjY  => Vjs, Ajs and pijA, pijV
    - cal UiX, UiY, TiX, TiY => create HashCommitment
    - send  HashCommitment.C 
8. Round 8 
    - Send D commitment  
9.  Round 9
    - Decode using D => UiX, UiY , TiX, TiY 
    - Send si (part signature)
10.  Finalize
    - Get sumS
    - Create signature: 
    - + R = rx
    - + S = sumS
    - Verify signature 