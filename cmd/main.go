package main

/*
#include <stdlib.h> // Include C standard library, if necessary
#include <string.h>
typedef struct {
    char* proof;
    char* vk;
} Groth16ProofWithVK;
*/
import "C"
import (
	"fmt"
	// "os"

	"github.com/cf/gnark-plonky2-verifier/worker"
)

type Groth16ProofWithVK struct {
	Proof string
	Vk    string
}

//export GenerateGroth16Proof
func GenerateGroth16Proof(common_circuit_data *C.char, proof_with_public_inputs *C.char, verifier_only_circuit_data *C.char, keystore_path *C.char) *C.Groth16ProofWithVK {
	proof_str, vk_str := worker.GenerateProof(C.GoString(common_circuit_data), C.GoString(proof_with_public_inputs), C.GoString(verifier_only_circuit_data), C.GoString(keystore_path))

	cProofWithVk := (*C.Groth16ProofWithVK)(C.malloc(C.sizeof_Groth16ProofWithVK))
	cProofWithVk.proof = C.CString(proof_str)
	cProofWithVk.vk = C.CString(vk_str)
	return cProofWithVk
}

//export VerifyGroth16Proof
func VerifyGroth16Proof(proofString *C.char, vkString *C.char) *C.char {
	return C.CString(worker.VerifyProof(C.GoString(proofString), C.GoString(vkString)))
}

//export Initialize
func Initialize(keyPath *C.char) {
	worker.Initialize(C.GoString(keyPath))
}

func main() {
	fmt.Println(worker.VerifyProof(
		"{\"pi_a\":\"8109066c4e4aef4f4f14a884a21994206a2e6d40a7afb94b150da01adcd23b909e8e8441a66f4397ecc794cda5dd9180\",\"pi_b_a0\":\"e656e320f770b3451755401206541faec8f33f7776b03d6707ed60a1af9a5d2cc8b29a4476c8886559be6bee5c8e0c02\",\"pi_b_a1\":\"ecdfd81c27a959b5d063464cf6250aaa91fa7a9d8177f76229908faaf4720e9fba1d9d5e1700f2f49a3374567ca5a311\",\"pi_c\":\"fbd932e0149814c5feafafe329f83df743896f3b18c41d4c07de825b621b10187a79573672acd2b39dad92cf72c0350f\",\"public_input_0\":\"caec89bedcab3c6707981a6fa86d27387b2c0a732feee3717aaddd6728877c5b\",\"public_input_1\":\"63ae8af9e3a21dcaa0d7e25af7314bc69e17ca08ee86ac6365fa508c7bd80a00\"}",
		"{\"alpha_g1\":\"b14314e4591e346d10a5e3ff6f27593e02d6c911bcfda1cdb554b29fa57863d6448d623f2dd9c3e3f6037840b522d48a\",\"k\":[\"4c7e92a74ff2275b2c23f2d8a07e265e8d4ace7748d37cfdab5139dd8e22ed729c06800675aa1e198ad2f2e07370338a\",\"d768918f786556e92955f09a82b3987cf138d978096f8ba1d7d309cb230b97afa01ae7e52cec6d4154bc82fb38b5418b\",\"c0847c7b309db151b70b294c904ca62ddd39aa59fdf20b2fd02903d1f3a8b08bb6eec58bc6fdfcf87d37441d3ae6ea8f\"],\"beta_g2\":\"c0c9949c6859905000a83aebe0aad9b550d672c9c3849a7ce5cad295939c11c96daaf36db518ff802ebb4b36e37155156aa989ee7392f2b64aceed795188b47df2dbbf3863e56bd59b2f0bea2c8fe03777d9c28d55ac2e1ccf4c4618f5383e06\",\"delta_g2\":\"2fdae7da1e4a4d87532e44ee3ef62eaa80e5990ed959f97e20c5b7e00d1080e11991e77d0f38c0e925c51a8db4ceda19085a90ec39cb7fd747e8becb6ae6fac36ebf56694349ec7513a2af85d2241ab7ec6d8f7d42de14067efa2160d3cb7105\",\"gamma_g2\":\"9388044478c3b8ddcb64bc53f1fd04647d8805b159f0333feff9a1d4b7c0d969dcec8f82d61b18cfe83b9a6175d17203b394331b26f61899d73efe55d5b5a2de21d44cdb0fe2829bba8a195aa8700981cdb45bb357f278903a047cbd37a63285\"}"))

	// path := "/tmp/plonky2-proofs/0"
  //
	// common_circuit_data, _ := os.ReadFile(path + "/common_circuit_data.json")
  //
	// proof_with_public_inputs, _ := os.ReadFile(path + "/proof_with_public_inputs.json")
  //
	// verifier_only_circuit_data, _ := os.ReadFile(path + "/verifier_only_circuit_data.json")
  //
  // worker.Initialize("/tmp/groth16-keystore/0")
	// proof_city, vk_city := worker.GenerateProof(string(common_circuit_data), string(proof_with_public_inputs), string(verifier_only_circuit_data), "/tmp/groth16-keystore/0/")
	// fmt.Println("proof city", proof_city)
	// fmt.Println("vk city", vk_city)
}
