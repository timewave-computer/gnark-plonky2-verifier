use std::str::FromStr;

use ark_bls12_381::Bls12_381;
use ark_bls12_381::Fq;
use ark_bls12_381::Fq2;
use ark_bls12_381::Fr;
use ark_bls12_381::G1Affine;
use ark_bls12_381::G1Projective;
use ark_bls12_381::G2Affine;
use ark_bls12_381::G2Projective;
use ark_groth16::Proof;
use ark_serialize::CanonicalDeserialize;
use ark_serialize::CanonicalSerialize;
use serde::Deserialize;
use serde::Deserializer;
use serde::Serialize;
use serde::Serializer;

pub fn str_to_fq(s: &str) -> anyhow::Result<Fq> {
    Ok(Fq::from_str({
        let decimal_str = hex::decode(s)?;
        &std::str::from_utf8(&decimal_str)?.to_string()
    })
    .map_err(|_| anyhow::anyhow!("str to fq conversion failed"))?)
}
pub fn str_to_fr(s: &str) -> anyhow::Result<Fr> {
    Ok(Fr::from_str({
        let decimal_str = hex::decode(s)?.clone();
        &std::str::from_utf8(&decimal_str)?.to_string()
    })
    .map_err(|_| anyhow::anyhow!("str to fr conversion failed"))?)
}

#[derive(Serialize, Deserialize, Clone, PartialEq, Debug)]
pub struct Groth16ProofSerializable {
    pub pi_a: [String; 2],
    pub pi_b: [[String; 2]; 2],
    pub pi_c: [String; 2],
    pub public_inputs: Vec<String>,
}

impl Groth16ProofSerializable {
    pub fn to_proof_groth16_bls12_381(&self) -> anyhow::Result<Proof<Bls12_381>> {
        let a_g1 = G1Affine::from(G1Projective::new(
            str_to_fq(&self.pi_a[0])?,
            str_to_fq(&self.pi_a[1])?,
            Default::default(),
        ));
        let b_g2 = G2Affine::from(G2Projective::new(
            // x
            Fq2::new(str_to_fq(&self.pi_b[0][0])?, str_to_fq(&self.pi_b[0][1])?),
            // y
            Fq2::new(str_to_fq(&self.pi_b[1][0])?, str_to_fq(&self.pi_b[1][1])?),
            Default::default(),
        ));

        let c_g1 = G1Affine::from(G1Projective::new(
            str_to_fq(&self.pi_c[0])?,
            str_to_fq(&self.pi_c[1])?,
            Default::default(),
        ));

        Ok(Proof::<Bls12_381> {
            a: a_g1,
            b: b_g2,
            c: c_g1,
        })
    }
    pub fn to_proof_with_public_inputs_groth16_bls12_381(
        &self,
    ) -> anyhow::Result<Groth16BLS12381ProofData> {
        let proof = self.to_proof_groth16_bls12_381()?;
        let mut public_inputs: Vec<Fr> = Vec::new();
        for pi in self.public_inputs.iter() {
            let v = str_to_fr(pi)?;
            public_inputs.push(v);
        }
        Ok(Groth16BLS12381ProofData {
            proof: proof,
            public_inputs: public_inputs,
        })
    }
    pub fn from_proof_with_public_inputs_groth16_bls12_381(
        proof: &Groth16BLS12381ProofData,
    ) -> Self {
        let a_g1_projective = G1Projective::from(proof.proof.a);
        let b_g2_projective = G2Projective::from(proof.proof.b);
        let c_g1_projective = G1Projective::from(proof.proof.c);
        let public_inputs = proof
            .public_inputs
            .iter()
            .map(|x| x.to_string())
            .collect::<Vec<_>>();
        Self {
            pi_a: [a_g1_projective.x.to_string(), a_g1_projective.y.to_string()],
            pi_b: [
                [
                    b_g2_projective.x.c0.to_string(),
                    b_g2_projective.x.c1.to_string(),
                ],
                [
                    b_g2_projective.y.c0.to_string(),
                    b_g2_projective.y.c1.to_string(),
                ],
            ],
            pi_c: [c_g1_projective.x.to_string(), c_g1_projective.y.to_string()],
            public_inputs: public_inputs,
        }
    }
}

#[derive(Clone, Debug, PartialEq, CanonicalSerialize, CanonicalDeserialize)]
pub struct Groth16BLS12381ProofData {
    pub proof: Proof<Bls12_381>,
    pub public_inputs: Vec<Fr>,
}

impl Serialize for Groth16BLS12381ProofData {
    fn serialize<S>(&self, serializer: S) -> Result<S::Ok, S::Error>
    where
        S: Serializer,
    {
        let raw = Groth16ProofSerializable::from_proof_with_public_inputs_groth16_bls12_381(&self);

        raw.serialize(serializer)
    }
}

impl<'de> Deserialize<'de> for Groth16BLS12381ProofData {
    fn deserialize<D>(deserializer: D) -> Result<Self, D::Error>
    where
        D: Deserializer<'de>,
    {
        use serde::de::Error;
        let raw = Groth16ProofSerializable::deserialize(deserializer)?;

        if let Ok(proof) = raw.to_proof_with_public_inputs_groth16_bls12_381() {
            Ok(proof)
        } else {
            Err(Error::custom("invalid Groth16BLS12381ProofData JSON"))
        }
    }
}
