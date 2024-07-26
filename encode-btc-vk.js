const vk = {
  alpha_g1:
    '3f537fd07503940b785177602fa1942fdd301d970fbfc2904bc5096e1aecba92345815782479f9b653ea2b3c16cb3519',
  beta_g2:
    'db4c2501fcb96816784d01d71b2291f9b3a6834da6fe258dc949721ed582642fbdf8b5b1d38b791c6c59f92233773703408eaa790c19fb052fc047ce3fe8f01984b80460eae38fdfca0407dbd448ad4bbdfed7e1678b8fd14b3f02d9decaa18e',
  gamma_g2:
    '98ab1a3dc22437194b15cb0e9aacbad648f2e3c3a282bb2e4331273957095d57d66063a1fe5fad7b4028c315e337950f0bd32774584689e5b4762503d7915083f57317b5bb5b801870d4c24d51d88b8e669d98a5e792245c76f2a593aa34ac03',
  delta_g2:
    '2277e8f519fc59614347a2545f7da3361aed1e7a16350f09ee36931749eb4ab03d8a308c9e6e030e7a96867b41848c016f9c6721b13ce2f2ef264cfc8235a52bef9c7c9fe72864365e10eddc4ad5ba4b3b5631c15cb8ecd944843ef744516c0c',
  k: [
    '5dc5f406b2746714af04af01483568c81c99de1c4492bb09d50dc481eadce66d10d1f7b10a94a17fad93ba345482b700',
    '4a01651246a221290309677321cef71d3a8d287fdb81f9f5e4f73ebc992766f1cbcf4e587b32ae85f0ff9b9a028a3210',
    '98260f07b570257778f4f97ef8c765ae34ec4a045eae99881d7b24e1bbe81fe195ae9c63a2b377c555350f12b3cc5106',
  ],
}


const combinedHex = vk.alpha_g1 + vk.k.join('') + vk.beta_g2 + vk.delta_g2 + vk.gamma_g2

for (let i = 0; i < combinedHex.length; i+=160) {
  console.log(combinedHex.slice(i, i + 160))
  console.log()
}
