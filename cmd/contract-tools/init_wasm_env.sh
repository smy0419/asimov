## 安装rust编译工具链，如果慢可以搜索其他安装地址
curl https://sh.rustup.rs -sSf | sh     				
echo "install rustup finish"

## 设置环境变量
source $HOME/.cargo/env  
echo "set env finish"												

## 安装指定版本 工具链
rustup install nightly-2018-11-12
echo "install nightly finish"

rustup target add wasm32-unknown-unknown
echo "add target wasm32 finish"

## 安装合约打包工具，固定版本0.6.0，其他版本合约格式可能会变(0.7.0不能有deploy方法)
cargo install pwasm-utils-cli --version 0.6.0 --bin wasm-build
echo "install wasm-build finish"

## 设置Rust源
export RUSTUP_DIST_SERVER=https://mirrors.ustc.edu.cn/rust-static
export RUSTUP_UPDATE_ROOT=https://mirrors.ustc.edu.cn/rust-static/rustup

## 创建cargo配置文件
echo "cargo config"
cat>~/.cargo/config<<EOF
[registry]
index = "https://mirrors.ustc.edu.cn/crates.io-index/"
[source.crates-io]
registry = "https://github.com/rust-lang/crates.io-index"
replace-with = 'ustc'
[source.ustc]
registry = "https://mirrors.ustc.edu.cn/crates.io-index/"
EOF

echo "init confif finish"
