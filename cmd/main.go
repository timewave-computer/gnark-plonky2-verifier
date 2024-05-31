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
	"encoding/hex"
	"fmt"
	"math/big"

	// "os"

	"github.com/cf/gnark-plonky2-verifier/worker"
	curve "github.com/consensys/gnark-crypto/ecc/bls12-381"
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
  // correct := "{\"pi_a\":\"8109066c4e4aef4f4f14a884a21994206a2e6d40a7afb94b150da01adcd23b909e8e8441a66f4397ecc794cda5dd9180\",\"pi_b_a0\":\"e656e320f770b3451755401206541faec8f33f7776b03d6707ed60a1af9a5d2cc8b29a4476c8886559be6bee5c8e0c02\",\"pi_b_a1\":\"ecdfd81c27a959b5d063464cf6250aaa91fa7a9d8177f76229908faaf4720e9fba1d9d5e1700f2f49a3374567ca5a311\",\"pi_c\":\"fbd932e0149814c5feafafe329f83df743896f3b18c41d4c07de825b621b10187a79573672acd2b39dad92cf72c0350f\",\"public_input_0\":\"caec89bedcab3c6707981a6fa86d27387b2c0a732feee3717aaddd6728877c5b\",\"public_input_1\":\"63ae8af9e3a21dcaa0d7e25af7314bc69e17ca08ee86ac6365fa508c7bd80a00\"}"
  // wrong := "{\"pi_a\":\"a21cd8a05bfa970df753e328a9d38739a429932a4cb4a5978da06a2e1afcee8b70d94a88a5dede1c143cefcced207b01\",\"pi_b_a0\":\"8fbb072475d03897109c41c3f4f8851c95af6871dfdc13116aae3784bb2bb9231dd2db2e48675ea464449eb32ad24b08\",\"pi_b_a1\":\"8e6f494b9785cb86fd1e2645f6e0ebe0d8e97b7c42c8bb5ae83f85713aacc5a742a7af979361491c5b6d1be53de5d38c\",\"pi_c\":\"d859e9884edf1f5a537895a5cd17a68cab46185a9ef0472060a4d98ed230d5d29d68d1132b9004cae7682763b6729988\",\"public_input_0\":\"66cf50f0fbac9928f8bcda00df07535d4cbb61b29b656d1a714fc1cca1f18728\",\"public_input_1\":\"0a0f684d7f385d854611316de64bcf9544d7a7284a5493e94d9f518df1defb00\"}"
  // wrong2 := "{\"pi_a\":\"632209f3ace26ba838157904e10f4648d8e894ed4d4b70a5a5272575b50732cba1c6394489b54a4af69d4b7a5525a606\",\"pi_b_a0\":\"f11a9ea0fba01b37605cf5a30d77a65956b0e577efd1a24ee3a40de03f66013e513b94e5c2ba002ae7905e5a5b73740b\",\"pi_b_a1\":\"1b000aab7973420232fc001c24c0fab5fa3b38d5658ec0fdf505822897c12df485773052938c861ff45b90686cae6597\",\"pi_c\":\"6e58fe2692dec5f31b6a1bce251ebcfd5a56dcbe69199ede65e36105beecdfabf89c1ecdb056584bd72122485d75ef03\",\"public_input_0\":\"0363ee02c3b769234a6a00def67f1f57de9094d83c4cc233e2fe6116fe0b0b02\",\"public_input_1\":\"6f42e3aa8f3b50f5d702ee7212ba311db6860059dac81e093db3e90e3d4b5500\"}"
  var pi_c = new(curve.G1Affine)
  b, _ := hex.DecodeString("01f6d8edb442c362d0bf0c44df0d0a6d80caf9017956abde0d897d1ceebb63e398688c99959586e135ec2c3675327d01056d06d924ff546b313483d6210b5fcef5212ae3b88b276330f5588ecc0f2e755ad2509f3f344ef4c698d7670e38126e")
  pi_c.Unmarshal(b)
  fmt.Println(pi_c.X.BigInt(big.NewInt(0)))
  fmt.Println(pi_c.Y.BigInt(big.NewInt(0)))
  wrong3 := "{\"pi_a\":\"e60d2bfef6abe2cb42b20f658480728a8f6409e8a93f69cdb11367eda0ce98df78d33f2213761620c5e8c51683215e0e\",\"pi_b_a0\":\"dd4cb222f9669cbcc95e56de14d9cda538a535c0c5c4a085a11adcc00af8b1a795ebb869dec6cb1a37a3f5faac8d1704\",\"pi_b_a1\":\"712b229d7819f9cbe181c759e2a2891e165d760068044c253f8607a23444382639d41d7a53503698029d9c617248d18c\",\"pi_c\":\"017d3275362cec35e1869595998c6898e363bbee1c7d890ddeab567901f9ca806d0a0ddf440cbfd062c342b4edd8f601\",\"public_input_0\":\"caec89bedcab3c6707981a6fa86d27387b2c0a732feee3717aaddd6728877c5b\",\"public_input_1\":\"f408bbd8d7f1719d4d251d477d73444f2fdc7da7e0ed54adcab6f463db97f600\"}"
	if worker.VerifyProof(
	  wrong3,
		"{\"alpha_g1\":\"b14314e4591e346d10a5e3ff6f27593e02d6c911bcfda1cdb554b29fa57863d6448d623f2dd9c3e3f6037840b522d48a\",\"k\":[\"4c7e92a74ff2275b2c23f2d8a07e265e8d4ace7748d37cfdab5139dd8e22ed729c06800675aa1e198ad2f2e07370338a\",\"d768918f786556e92955f09a82b3987cf138d978096f8ba1d7d309cb230b97afa01ae7e52cec6d4154bc82fb38b5418b\",\"c0847c7b309db151b70b294c904ca62ddd39aa59fdf20b2fd02903d1f3a8b08bb6eec58bc6fdfcf87d37441d3ae6ea8f\"],\"beta_g2\":\"c0c9949c6859905000a83aebe0aad9b550d672c9c3849a7ce5cad295939c11c96daaf36db518ff802ebb4b36e37155156aa989ee7392f2b64aceed795188b47df2dbbf3863e56bd59b2f0bea2c8fe03777d9c28d55ac2e1ccf4c4618f5383e06\",\"delta_g2\":\"2fdae7da1e4a4d87532e44ee3ef62eaa80e5990ed959f97e20c5b7e00d1080e11991e77d0f38c0e925c51a8db4ceda19085a90ec39cb7fd747e8becb6ae6fac36ebf56694349ec7513a2af85d2241ab7ec6d8f7d42de14067efa2160d3cb7105\",\"gamma_g2\":\"9388044478c3b8ddcb64bc53f1fd04647d8805b159f0333feff9a1d4b7c0d969dcec8f82d61b18cfe83b9a6175d17203b394331b26f61899d73efe55d5b5a2de21d44cdb0fe2829bba8a195aa8700981cdb45bb357f278903a047cbd37a63285\"}") == "false" {
  fmt.Println("failed to verify")
  }


	// path := "/tmp/plonky2_proof/0"
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
