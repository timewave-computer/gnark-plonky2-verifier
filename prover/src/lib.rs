use plonky2::field::goldilocks_field::GoldilocksField;
use plonky2::plonk::circuit_data::CircuitData;
use plonky2::plonk::config::PoseidonGoldilocksConfig;
use plonky2::plonk::proof::ProofWithPublicInputs;
use plonky2x::backend::wrapper::wrap::WrappedCircuit;
use plonky2x::frontend::builder::CircuitBuilder as WrapperBuilder;
use plonky2x::prelude::DefaultParameters;

use crate::parameters::Groth16WrapperParameters;

pub mod parameters;

pub mod fr;
pub mod logger;
pub mod plonky2_config;
pub mod poseidon_bls12_381;
pub mod poseidon_bls12_381_constants;
pub mod serde;

pub type F = GoldilocksField;
pub type C = PoseidonGoldilocksConfig;
pub const D: usize = 2;

pub fn wrap_plonky2_proof(
    circuit_data: CircuitData<F, C, D>,
    proof: &ProofWithPublicInputs<F, C, D>,
) -> anyhow::Result<String> {
    let wrapper_builder = WrapperBuilder::<DefaultParameters, D>::new();
    let mut circuit = wrapper_builder.build();
    circuit.data = circuit_data;
    let wrapped_circuit =
        WrappedCircuit::<DefaultParameters, Groth16WrapperParameters, D>::build(circuit);
    let wrapped_proof = wrapped_circuit.prove(&proof)?;
    let g16_proof = gnark_plonky2_verifier_ffi::generate_groth16_proof(
        &serde_json::to_string(&wrapped_proof.common_data)?,
        &serde_json::to_string(&wrapped_proof.proof)?,
        &serde_json::to_string(&wrapped_proof.verifier_data)?,
    );
    Ok(g16_proof)
}

#[cfg(test)]
mod tests {
    use plonky2::field::types::Field;
    use plonky2::iop::witness::PartialWitness;
    use plonky2::iop::witness::WitnessWrite;
    use plonky2::plonk::circuit_builder::CircuitBuilder;
    use plonky2::plonk::circuit_data::CircuitConfig;

    use super::*;

    #[test]
    fn test_prover() -> anyhow::Result<()> {
        logger::setup_logger();
        // config defines number of wires of gates, FRI strategies etc.
        let config = CircuitConfig::standard_recursion_config();

        // We use GoldilocksField as circuit arithmetization
        type F = GoldilocksField;

        // We use Poseidon hash on GoldilocksField as FRI hasher
        type C = PoseidonGoldilocksConfig;

        // We use the degree D extension Field when soundness is required.
        const D: usize = 2;

        tracing::info!("prove that the prover knows x such that x^2 - 2x + 1 = 0");
        let mut builder = CircuitBuilder::<F, D>::new(config.clone());

        let x_t = builder.add_virtual_target();
        let minus_x_t = builder.neg(x_t);
        let minus_2x_t = builder.mul_const(F::from_canonical_u64(2), minus_x_t);
        let x2_t = builder.exp_u64(x_t, 2);
        let one_t = builder.one();
        let zero_t = builder.zero();
        let poly_t = builder.add_many(&[x2_t, minus_2x_t, one_t]);
        builder.connect(poly_t, zero_t); // x^2 - 2x + 1 = 0
        builder.register_public_input(x_t);
        builder.register_public_input(x_t);
        builder.register_public_input(x_t);
        builder.register_public_input(x_t);
        builder.register_public_input(x_t);
        builder.register_public_input(x_t);
        builder.register_public_input(x_t);
        builder.register_public_input(x_t);

        tracing::info!("compiling circuits...");
        let data = builder.build::<C>();
        let mut pw = PartialWitness::<F>::new();
        tracing::info!("setting witness...");
        pw.set_target(x_t, GoldilocksField(1)); // set x = 1

        tracing::info!("proving...");
        let proof = data.prove(pw)?;
        tracing::info!("verifying...");
        data.verify(proof.clone())?;
        tracing::info!("done!");

        tracing::info!("compiling wrapping circuits...");
        let g16_proof = wrap_plonky2_proof(data, &proof)?;
        tracing::info!("proving...");
        tracing::info!("saving...");
        tracing::info!("done!");

        println!("{}", g16_proof);

        Ok(())
    }

    #[test]
    fn test_serde() {
       let proof = r#"{"CommitmentPok":"0afccd689faefc15cdfea39be55b9df65023517f66e79bbdaaa6439d54adbc70ab0aa4cf1c7ca1579601b610fedec52c178d81e8efb2c36d1d733f71aa36587f50e6965227f40e87fa7e3da6a69e11e587f8a4b20744b14e67a0031685318f8a","Commitments":"00e5a7f40bd986df68ff103023172b37e37ce6d4811d938513db024cdc6fb232e0cd4f6923851004f37f953fd2808a9c04cf818d21f58215597d680129ef958f9c3f4543d5ac7a8d0c72b88bd123223482616da95a081badad282f160d071371","pi_a":["15936e02581ccba2a1e1f0f27d63646b4d2eba0533bcebe524622e0ccae94c9c6b82128c042208b040257ce4b3be510f","0b4e43f99a5bbd59170aa208f0911bfaa233c32232386435ba541f58ad340906cce1a3b82ee7c30dec1fb6afeff32aa3"],"pi_b":[["147f8aa633a8485d805720b631a69d923b3ba0d6a7458e0c610d78c2ffcf9e42d5948529b6531a20a662e79868e681ac","14d02521951fcf49fc20e6318da1e7d38efc153e3d99289f722a17171aba10bd36c301e87ba12cf1c82f3bb9b0af39e1"],["173c3c851748ac45e34251d7e6dd0fb8b6b048882ec983959b957f52f3f185916b9034e6dd2579bacded3230cfb73e0d","12fcbfcfff3a540f63bb10de67758db3ac126999624d997507a61d2ad207adc46dd003a9c350e6cb15b4e366da5b706f"]],"pi_c": ["14ffbd3a74497118ea4548a88cc4a69c3dfbab9e8ff807556079198f3f16aa5b292c02f7c214d0ccda1502e320b8ba52","0e0641ce74fe3e2decd86e900c0ce67c8031a4e63b8e32d654806ec8629bfb1bcb270e5a5d6ca285e9313931b3667aa0"],"public_inputs":["0c700243a92fe885c6688844ec22b3db16e09b9494509b74f88fb1864991e044","0cc66451cf0cf4382a118d3241c6d919ca6c891ea487fe23c3ac6e37fe6fcc13"]}"#;

        serde_json::from_str::<serde::Groth16BLS12381ProofData>(proof).unwrap();
    }
}
