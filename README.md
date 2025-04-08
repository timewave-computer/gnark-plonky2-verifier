# Gnark Plonky2 Verifier - Forked by Valence
This repo was forked as part of an exploration effort revolving around Plonkish proving systems.
We are very likely not going to use this as it depends on a lot of deprecated infrastructure 
and complex proof types that were used by succinct in their `succinctx` back in the days.

The only reason why we want to pin this and look into it is because it provides a nice
interface for wrapping Plonky2 proofs (old version 0.2.0, not 1.0.2!) with groth16.

We would love to see something similar for an up-to-date version of Plonky3.

Someone at Succinct [made an effort to work towards this goal](https://github.com/raftedproc/simplified_p3_air_with_sp1_rec/blob/main/src/main.rs), 
but the implementation is not finished.

What they do is essentially wrap a Plonky3 proof as an SP1 VM proof and prove it's verification using their in-house groth16 wrapper.

This is also not idea for what we are trying to achieve.
