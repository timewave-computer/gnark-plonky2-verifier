use plonky2::field::goldilocks_field::GoldilocksField;
use plonky2::plonk::circuit_data::CircuitData;
use plonky2::plonk::config::PoseidonGoldilocksConfig;
use plonky2::plonk::proof::ProofWithPublicInputs;
use plonky2x::backend::wrapper::wrap::WrappedCircuit;
use plonky2x::frontend::builder::CircuitBuilder as WrapperBuilder;
use plonky2x::prelude::DefaultParameters;
use serde_json::to_string as json;

use crate::parameters::Groth16WrapperParameters;

pub mod parameters;

pub mod fr;
pub mod logger;
pub mod plonky2_config;
pub mod poseidon_bls12_381;
pub mod poseidon_bls12_381_constants;

pub type F = GoldilocksField;
pub type C = PoseidonGoldilocksConfig;
pub const D: usize = 2;

pub fn wrap_plonky2_proof(
    circuit_data: CircuitData<F, C, D>,
    proof: &ProofWithPublicInputs<F, C, D>,
    save_wrapped_data_path: Option<&str>,
    id: &str,
) -> anyhow::Result<(String, String)> {
    circuit_data.verify(proof.clone())?;
    let wrapper_builder = WrapperBuilder::<DefaultParameters, D>::new();
    let mut circuit = wrapper_builder.build();
    circuit.data = circuit_data;
    let wrapped_circuit =
        WrappedCircuit::<DefaultParameters, Groth16WrapperParameters, D>::build(circuit);
    let wrapped_proof = wrapped_circuit.prove(&proof)?;
    if let Some(save_wrapped_data_path) = save_wrapped_data_path {
        wrapped_proof.save(save_wrapped_data_path)?;
    }
    Ok(gnark_plonky2_verifier_ffi::generate_groth16_proof(
        &json(&wrapped_proof.common_data)?,
        &json(&wrapped_proof.proof)?,
        &json(&wrapped_proof.verifier_data)?,
        id
    ))
}

pub fn verify_groth16_proof(
    proof_with_public_inputs: &str,
    vk: &str,
) -> anyhow::Result<bool>{

    let result = gnark_plonky2_verifier_ffi::verify_groth16_proof(&proof_with_public_inputs, &vk);
    Ok(result == "true")
}

pub fn initialize(key_path: &str) -> anyhow::Result<()>{
    gnark_plonky2_verifier_ffi::initialize(key_path);
    Ok(())
}

#[cfg(test)]
mod tests {
    use plonky2::field::types::Field;
    use plonky2::iop::witness::PartialWitness;
    use plonky2::iop::witness::WitnessWrite;
    use plonky2::plonk::circuit_builder::CircuitBuilder;
    use plonky2::plonk::circuit_data::CircuitConfig;

    use super::*;

    fn test_prover(save_wrapped_data_path: Option<&str>, id: &str) -> anyhow::Result<()> {
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
        for i in 0..224 {
            builder.register_public_input(if i & 1 == 1 { one_t } else { zero_t });
        }
        for _i in 224..256 {
            builder.register_public_input(zero_t);
        }
        for i in 256..480 {
            builder.register_public_input(if i & 1 == 1 { one_t } else { zero_t });
        }
        for _i in 480..512 {
            builder.register_public_input(zero_t);
        }

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
        // tracing::info!("original public inputs: {:?}", proof.public_inputs);

        tracing::info!("compiling wrapping circuits...");
        let (g16_proof, g16_vk) = wrap_plonky2_proof(data, &proof, save_wrapped_data_path, id)?;
        tracing::info!("done!");

        println!("proof {}", g16_proof);
        println!("vk {}", g16_vk);

        println!("verify {:?}", verify_groth16_proof(&g16_proof, &g16_vk));

        Ok(())
    }

    #[test]
    fn test_setup_once() {
        // First, generate the key files
        test_prover(Some("../testdata/0"), "/tmp/groth16/").unwrap();
        
        // Now initialize with the generated key files
        initialize("/tmp/groth16/").unwrap();
        
        // Run the test again to verify it works
        test_prover(Some("../testdata/0"), "/tmp/groth16/").unwrap();
        test_prover(Some("../testdata/0"), "/tmp/groth16/").unwrap();
        test_prover(Some("../testdata/0"), "/tmp/groth16/").unwrap();
    }

    #[test]
    fn test_generate_wrapped_data() {
        // Create a temporary directory for the wrapped data
        let temp_dir = std::env::temp_dir().join("gnark_plonky2_test");
        std::fs::create_dir_all(&temp_dir).unwrap();
        
        // Generate the wrapped data for your circuit
        test_prover(Some(temp_dir.to_str().unwrap()), "testdata/my_circuit").unwrap();
        
        // Verify that the wrapped data was generated
        assert!(temp_dir.join("testdata/my_circuit/common_circuit_data.json").exists());
        assert!(temp_dir.join("testdata/my_circuit/proof_with_public_inputs.json").exists());
        assert!(temp_dir.join("testdata/my_circuit/verifier_only_circuit_data.json").exists());
        
        // Clean up
        std::fs::remove_dir_all(&temp_dir).unwrap();
    }
}
