use clap::{Parser, Subcommand};

#[derive(Parser)]
#[command(name = "felapp")]
#[command(about = "Multi-service CLI for FelApp", long_about = None)]
struct Cli {
    #[command(subcommand)]
    command: Commands,
}

#[derive(Subcommand)]
enum Commands {
    Cache {
        #[arg(long, default_value = "0.0.0.0:50051")]
        addr: String,
    },
}

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    let cli = Cli::parse();

    match cli.command {
        Commands::Cache { addr } => {
            cache::run(addr).await.map_err(|e| anyhow::Error::msg(e))?;
        }
    }

    Ok(())
}
