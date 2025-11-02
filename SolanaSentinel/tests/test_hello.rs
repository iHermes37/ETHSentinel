use SolanaSentinel::modules::{scanner::scanchain, track_whales::{trackwhale}};

#[test]
fn test_hello() {
    let msg = "hello test";
    println!("{}", msg); // 输出 hello test
    assert_eq!(msg, "hello test");
}


#[tokio::test(flavor = "multi_thread", worker_threads = 4)]
async fn test_scanner() {
    scanchain::scan_main().await;
}


#[tokio::test(flavor = "current_thread")]
async fn test_track_whale() {
    if let Err(e) = trackwhale::track_main().await {
        eprintln!("track_main 运行失败: {:?}", e);
    }
}



