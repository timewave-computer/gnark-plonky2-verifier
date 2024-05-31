let proof_hex = "e60d2bfef6abe2cb42b20f658480728a8f6409e8a93f69cdb11367eda0ce98df78d33f2213761620c5e8c51683215e0edd4cb222f9669cbcc95e56de14d9cda538a535c0c5c4a085a11adcc00af8b1a795ebb869dec6cb1a37a3f5faac8d1704712b229d7819f9cbe181c759e2a2891e165d760068044c253f8607a23444382639d41d7a53503698029d9c617248d18c017d3275362cec35e1869595998c6898e363bbee1c7d890ddeab567901f9ca806d0a0ddf440cbfd062c342b4edd8f601caec89bedcab3c6707981a6fa86d27387b2c0a732feee3717aaddd6728877c5bf408bbd8d7f1719d4d251d477d73444f2fdc7da7e0ed54adcab6f463db97f600"
let vk_hex = "b14314e4591e346d10a5e3ff6f27593e02d6c911bcfda1cdb554b29fa57863d6448d623f2dd9c3e3f6037840b522d48a4c7e92a74ff2275b2c23f2d8a07e265e8d4ace7748d37cfdab5139dd8e22ed729c06800675aa1e198ad2f2e07370338ad768918f786556e92955f09a82b3987cf138d978096f8ba1d7d309cb230b97afa01ae7e52cec6d4154bc82fb38b5418bc0847c7b309db151b70b294c904ca62ddd39aa59fdf20b2fd02903d1f3a8b08bb6eec58bc6fdfcf87d37441d3ae6ea8fc0c9949c6859905000a83aebe0aad9b550d672c9c3849a7ce5cad295939c11c96daaf36db518ff802ebb4b36e37155156aa989ee7392f2b64aceed795188b47df2dbbf3863e56bd59b2f0bea2c8fe03777d9c28d55ac2e1ccf4c4618f5383e062fdae7da1e4a4d87532e44ee3ef62eaa80e5990ed959f97e20c5b7e00d1080e11991e77d0f38c0e925c51a8db4ceda19085a90ec39cb7fd747e8becb6ae6fac36ebf56694349ec7513a2af85d2241ab7ec6d8f7d42de14067efa2160d3cb71059388044478c3b8ddcb64bc53f1fd04647d8805b159f0333feff9a1d4b7c0d969dcec8f82d61b18cfe83b9a6175d17203b394331b26f61899d73efe55d5b5a2de21d44cdb0fe2829bba8a195aa8700981cdb45bb357f278903a047cbd37a63285"

console.log(JSON.stringify({
  "pi_a": proof_hex.slice(0, 96),
  "pi_b_a0": proof_hex.slice(96, 96 * 2),
  "pi_b_a1": proof_hex.slice(96 * 2, 96 * 3),
  "pi_c": proof_hex.slice(96 * 3, 96 * 4),
  "public_input_0": proof_hex.slice(96 * 4, 96 * 4 + 64),
  "public_input_1": proof_hex.slice(96 * 4 + 64),
}).replace(/"/g, '\\"'))

console.log(JSON.stringify({
  "alpha_g1": vk_hex.slice(0, 96),
  "k": [
    vk_hex.slice(96, 96 * 2),
    vk_hex.slice(96 * 2, 96 * 3),
    vk_hex.slice(96 * 3, 96 * 4),
  ],
  "beta_g2": vk_hex.slice(96 * 4, 96 * 6),
  "delta_g2": vk_hex.slice(96 * 6, 96 * 8),
  "gamma_g2": vk_hex.slice(96 * 8, 96 * 10),
}).replace(/"/g, '\\"'))


function reverseHexString(hexStr) {
    let reversed = "";
    for (let i = 0; i < hexStr.length; i += 2) {
        reversed = hexStr[i] + hexStr[i + 1] + reversed;
    }
    return reversed;
}

console.log(reverseHexString("017d3275362cec35e1869595998c6898e363bbee1c7d890ddeab567901f9ca806d0a0ddf440cbfd062c342b4edd8f601"))
