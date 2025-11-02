// https://solana.com/zh/docs/rpc/websocket/accountsubscribe


// {
//   "jsonrpc": "2.0",
//   "id": 1,
//   "method": "accountSubscribe",
//   "params": [
//     "CM78CPUeXjn8o3yroDHxUtKsZZgoy4GPkPPXfouKNH12",
//     {
//       "encoding": "jsonParsed",
//       "commitment": "finalized"
//     }
//   ]
// }



// {
//   "jsonrpc": "2.0",
//   "method": "accountNotification",
//   "params": {
//     "result": {
//       "context": {
//         "slot": 5199307
//       },
//       "value": {
//         "data": {
//           "program": "nonce",
//           "parsed": {
//             "type": "initialized",
//             "info": {
//               "authority": "Bbqg1M4YVVfbhEzwA9SpC9FhsaG83YMTYoR4a8oTDLX",
//               "blockhash": "LUaQTmM7WbMRiATdMMHaRGakPtCkc2GHtH57STKXs6k",
//               "feeCalculator": {
//                 "lamportsPerSignature": 5000
//               }
//             }
//           }
//         },
//         "executable": false,
//         "lamports": 33594,
//         "owner": "11111111111111111111111111111111",
//         "rentEpoch": 635,
//         "space": 80
//       }
//     },
//     "subscription": 23784
//   }
// }