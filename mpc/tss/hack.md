# -- Attack ------- 
# 1. The Forget-And-Forgive Attack
- + The Target: Multi-Party Reshare
- + The Attack: 
- -- The attack aims to prevent parties running the reshare protocol from holding valid shares:
- ---- a corrupted old party will send “bad” secret shares to some parties of the
n′ parties
- ----  “good” secret shares to the others, where “good” means that it will pass the
validation
- ----> broadcast different VSS for each group partes => The system failed to issue signatures 
- + Exploitation: 
- --- fully eliminate the capability to issue a transaction => If the organization has insuﬃcient back-up, the attacker could blackmail the exchange
- + Mitigation: party broadcast "failed VSS" to all parrty 
  
## 2. The Lather, Rinse, Repeat Attack
- + The Target: Two-Party Reshare
- + The Attack: 
- + Exploitation
- + Mitigation
## 3. The Golden Shoe Attack
- + The Target: MtA Share Conversion
- + The Attack
- + -- The Missing Validation
- + -- How the Attack Works
- + Exploitation
- + Mitigation

# ---  Key Extraction --- 
## Attacks on multiplicative-to-additive protocol
-  Attack on absent range proofs
-  Small Paillier attack
-  

# --- Attack dnl proof 
## 1.Bit-Probing Attack
- How to Remediate the Attack.
  
## 2. Attacks when ZK-Proofs are Ommited
- + Bit-Probing for Least Significant Bits
- + Bit-Probing for MSB’s & LSB’s.
- + Guessing P2’s input bit by bit.