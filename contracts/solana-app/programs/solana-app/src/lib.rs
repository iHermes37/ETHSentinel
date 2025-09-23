use anchor_lang::prelude::*;

declare_id!("CRXua4uvt4hw6jTcpUZ6VgRymBhM3WZb9GUadz9jLf57");

#[program]
pub mod solana_app {
    use super::*;

    pub fn initialize(ctx: Context<Initialize>) -> Result<()> {
         msg!("Hello Solana with Anchor!");
        Ok(())
    }
}

#[derive(Accounts)]
pub struct Initialize {}
