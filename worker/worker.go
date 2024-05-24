package worker

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"

	"github.com/GopherJ/doge-covenant/serialize"
	gl "github.com/cf/gnark-plonky2-verifier/goldilocks"
	"github.com/cf/gnark-plonky2-verifier/types"
	"github.com/cf/gnark-plonky2-verifier/variables"
	"github.com/cf/gnark-plonky2-verifier/verifier"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
	"github.com/consensys/gnark/backend/groth16"
	groth16_bls12381 "github.com/consensys/gnark/backend/groth16/bls12-381"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

type CRVerifierCircuit struct {
	PublicInputs            []frontend.Variable               `gnark:",public"`
	Proof                   variables.Proof                   `gnark:"-"`
	VerifierOnlyCircuitData variables.VerifierOnlyCircuitData `gnark:"-"`

	OriginalPublicInputs []gl.Variable `gnark:"_"`

	// This is configuration for the circuit, it is a constant not a variable
	CommonCircuitData types.CommonCircuitData
}

func (c *CRVerifierCircuit) Define(api frontend.API) error {
	verifierChip := verifier.NewVerifierChip(api, c.CommonCircuitData)
	if len(c.PublicInputs) != 2 {
		panic("invalid public inputs, should contain 2 BLS12_381 elements")
	}
	if len(c.OriginalPublicInputs) != 512 {
		panic("invalid original public inputs, should contain 512 goldilocks elements")
	}

	two := big.NewInt(2)

	blockStateHashAcc := frontend.Variable(0)
	sighashAcc := frontend.Variable(0)
	for i := 255; i >= 0; i-- {
		blockStateHashAcc = api.MulAcc(c.OriginalPublicInputs[i].Limb, blockStateHashAcc, two)
	}
	for i := 511; i >= 256; i-- {
		sighashAcc = api.MulAcc(c.OriginalPublicInputs[i].Limb, sighashAcc, two)
	}

	api.AssertIsEqual(c.PublicInputs[0], blockStateHashAcc)
	api.AssertIsEqual(c.PublicInputs[1], sighashAcc)

	verifierChip.Verify(c.Proof, c.OriginalPublicInputs, c.VerifierOnlyCircuitData)

	return nil
}

func initKeyStorePath(id string) {
	_, err := os.Stat(KEY_STORE_PATH + id)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(KEY_STORE_PATH + id, os.ModePerm)
		}
	}
}

func GenerateProof(common_circuit_data string, proof_with_public_inputs string, verifier_only_circuit_data string, id string) (string, string) {
	initKeyStorePath(id)

	commonCircuitData := types.ReadCommonCircuitDataRaw(common_circuit_data)

	verifierOnlyCircuitDataRaw := types.ReadVerifierOnlyCircuitDataRaw(verifier_only_circuit_data)
	fmt.Println("circuit digest: ", verifierOnlyCircuitDataRaw.CircuitDigest)
	verifierOnlyCircuitData := variables.DeserializeVerifierOnlyCircuitData(verifierOnlyCircuitDataRaw)

	rawProofWithPis := types.ReadProofWithPublicInputsRaw(proof_with_public_inputs)
	proofWithPis := variables.DeserializeProofWithPublicInputs(rawProofWithPis)

	two := big.NewInt(2)

	blockStateHashAcc := big.NewInt(0)
	sighashAcc := big.NewInt(0)
	for i := 255; i >= 0; i-- {
		blockStateHashAcc = new(big.Int).Mul(blockStateHashAcc, two)
		blockStateHashAcc = new(big.Int).Add(blockStateHashAcc, new(big.Int).SetUint64(rawProofWithPis.PublicInputs[i]))
	}
	for i := 511; i >= 256; i-- {
		sighashAcc = new(big.Int).Mul(sighashAcc, two)
		sighashAcc = new(big.Int).Add(sighashAcc, new(big.Int).SetUint64(rawProofWithPis.PublicInputs[i]))
	}
	blockStateHash := frontend.Variable(blockStateHashAcc)
	sighash := frontend.Variable(sighashAcc)

	circuit := CRVerifierCircuit{
		PublicInputs:            make([]frontend.Variable, 2),
		Proof:                   proofWithPis.Proof,
		OriginalPublicInputs:    proofWithPis.PublicInputs,
		VerifierOnlyCircuitData: verifierOnlyCircuitData,
		CommonCircuitData:       commonCircuitData,
	}

	assignment := CRVerifierCircuit{
		PublicInputs:            []frontend.Variable{blockStateHash, sighash},
		Proof:                   circuit.Proof,
		OriginalPublicInputs:    circuit.OriginalPublicInputs,
		VerifierOnlyCircuitData: circuit.VerifierOnlyCircuitData,
		CommonCircuitData:       commonCircuitData,
	}

	// NewWitness() must be called before Compile() to avoid gnark panicking.
	// ref: https://github.com/Consensys/gnark/issues/1038
	witness, err := frontend.NewWitness(&assignment, ecc.BLS12_381.ScalarField())
	if err != nil {
		panic(err)
	}

	cs, pk, vk, err := Setup(&circuit, id)
	if err != nil {
		panic(err)
	}

	proof, err := groth16.Prove(cs, pk, witness)
	if err != nil {
		panic(err)
	}

	publicWitness, err := witness.Public()
	if err != nil {
		panic(err)
	}

	err = groth16.Verify(proof, vk, publicWitness)
	if err != nil {
		panic(err)
	}

	blsProof := proof.(*groth16_bls12381.Proof)
	blsWitness := publicWitness.Vector().(fr.Vector)

	proof_city, err := serialize.ToJsonCityProof(blsProof, blsWitness)
	if err != nil {
		panic(err)
	}

	proof_bytes, err := json.Marshal(&proof_city)
	if err != nil {
		panic(err)
	}

	var g16VerifyingKey = G16VerifyingKey{
		VK: vk,
	}

	vk_bytes, err := json.Marshal(g16VerifyingKey)
	if err != nil {
		panic(err)
	}

	return string(proof_bytes), string(vk_bytes)

}

func VerifyProof(proofString string, vkString string) string {
	g16ProofWithPublicInputs := NewG16ProofWithPublicInputs()

	if err := json.Unmarshal([]byte(proofString), g16ProofWithPublicInputs); err != nil {
		panic(err)
	}

	g16VerifyingKey := NewG16VerifyingKey()

	if err := json.Unmarshal([]byte(vkString), g16VerifyingKey); err != nil {
		panic(err)
	}

	if err := groth16.Verify(g16ProofWithPublicInputs.Proof, g16VerifyingKey.VK, g16ProofWithPublicInputs.PublicInputs); err != nil {
		return "false"
	}
	return "true"
}

func Setup(circuit *CRVerifierCircuit, id string) (constraint.ConstraintSystem, groth16.ProvingKey, groth16.VerifyingKey, error) {
	if _, err := os.Stat(KEY_STORE_PATH + id + VK_PATH); err == nil {
		ccs, err := ReadCircuit(ecc.BLS12_381, KEY_STORE_PATH+id+CIRCUIT_PATH)
		if err != nil {
			return nil, nil, nil, err
		}
		vk, err := ReadVerifyingKey(ecc.BLS12_381, KEY_STORE_PATH+id+VK_PATH)
		if err != nil {
			return nil, nil, nil, err
		}
		pk, err := ReadProvingKey(ecc.BLS12_381, KEY_STORE_PATH+id+PK_PATH)
		if err != nil {
			return nil, nil, nil, err
		}
		return ccs, pk, vk, nil
	}

	ccs, err := frontend.Compile(ecc.BLS12_381.ScalarField(), r1cs.NewBuilder, circuit)
	if err != nil {
		return nil, nil, nil, err
	}

	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		return nil, nil, nil, err
	}

	if err := WriteCircuit(ccs, KEY_STORE_PATH+id+CIRCUIT_PATH); err != nil {
		return nil, nil, nil, err
	}

	if err := WriteVerifyingKey(vk, KEY_STORE_PATH+id+VK_PATH); err != nil {
		return nil, nil, nil, err
	}

	if err := WriteProvingKey(pk, KEY_STORE_PATH+id+PK_PATH); err != nil {
		return nil, nil, nil, err
	}

	return ccs, pk, vk, nil
}
