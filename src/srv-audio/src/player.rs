use std::fs::File;
use std::io::BufReader;
use std::sync::mpsc as std_mpsc;
use std::thread;

use tokio::sync::mpsc as tokio_mpsc;

pub type Remote = tokio_mpsc::UnboundedSender<Command>;

#[derive(Clone, Debug, PartialEq)]
pub enum Command {
    Start,
    Stop,
}

pub async fn setup() -> Remote {
    // Start player in new thread
    let (stx, srx) = std_mpsc::channel::<Command>();
    thread::spawn(move || {
        let (_stream, stream_handle) = rodio::OutputStream::try_default().unwrap();
        let mut sink = rodio::Sink::try_new(&stream_handle).unwrap();
        while let Ok(command) = srx.recv() {
            if command == Command::Start {
                let file = File::open("default.mp3").unwrap();
                let source = rodio::Decoder::new_looped(BufReader::new(file)).unwrap();
                sink = rodio::Sink::try_new(&stream_handle).unwrap();
                sink.append(source);
                sink.play();
            } else {
                sink.stop();
            }
        }
    });

    // Start conversion loop from tokio channel to std channel
    // This is required as we need to keep _stream in scope. Since that is a raw pointer, it can't be mixed
    // with async/await code.
    let (ttx, mut trx) = tokio_mpsc::unbounded_channel::<Command>();
    tokio::spawn(async move {
        while let Some(command) = trx.recv().await {
            stx.send(command).unwrap();
        }
    });
    return ttx;
}
