use solana_sdk::transaction;
use solana_transaction_status::{UiMessage,UiInstruction,UiParsedInstruction};
use solana_transaction_status::{UiConfirmedBlock,EncodedTransaction};
use crate::models::block::block::whale_trans;
use crate::util::printjson::print_json;
use serde_json::Value;

fn parse_data(
    data:&Value,
    whale_trans:&mut whale_trans
){

}

fn parse_instruction(
    ins:&UiInstruction,
    whale_trans:&mut whale_trans
){
    match ins {
        UiInstruction::Parsed(parsed_ins) => {
            // 匹配 UiParsedInstruction
            match parsed_ins {
                UiParsedInstruction::Parsed(inner) => {
                    // inner 就是 ParsedInstruction
                    let program = &inner.program;
                    let program_id = &inner.program_id;
                    let parsed_value: &Value = &inner.parsed;

                    println!("program: {}, program_id: {}", program, program_id);
                    // println!("parsed data: {}", parsed_value);


                    // parse_data(parsed_value,whale_trans)
                    
                    // TODO: 根据 parsed_value 填充 whale_trans
                }
                UiParsedInstruction::PartiallyDecoded(_partially) => {
                    // 如果是部分解析指令
                    println!("PartiallyDecoded instruction");
                }
            }
        }
        UiInstruction::Compiled(_compiled) => {
            // 如果是编译形式指令
            println!("Compiled instruction, need decoding");
        }
    }
}

fn parse_message(
    mes:&UiMessage,
    whale_trans:&mut whale_trans
){

    match mes {
        UiMessage::Parsed(parsed_mes)=>{
            let account=parsed_mes.account_keys.clone();
            let instructions=parsed_mes.instructions.clone();
            for ins in instructions{
                parse_instruction(&ins,whale_trans);
            }
        }
        _=>{}
    }
}

pub fn parse_trans(
    curblock:&UiConfirmedBlock
){

    let mut whale_trans = whale_trans::default();

    whale_trans.block_hash=curblock.blockhash.clone();
    whale_trans.block_height=curblock.block_height.clone();
    whale_trans.block_time=curblock.block_time.clone();
    print_json(&whale_trans);


    if let Some(trans_lists) = &curblock.transactions {
        for trans in trans_lists{
                if let EncodedTransaction::Json(ui_tx) = &trans.transaction {
                        // let signatures: Vec<String> = &ui_tx.signatures;
                        parse_message(&ui_tx.message,&mut whale_trans);
                        // println!("signatures: {:?}", signatures);
                    } else {
                        // 非 Json 类型交易可以忽略或者处理
                        println!("非 Json 类型交易，不解析签名");
                    }
            }
    }

}