# Samples for validating

该项目由go语言编写，用于验证工作流互操作

## 项目结构

- cards/

    该业务网络参与者的凭证， 角色要求见`run.sh`

- model/

    利用zeebe建模工具所构建的模型，用以部署到zeebe工作流引擎上

- worker/roles/seller

    参与该次业务的seller角色

- worker/roles/user

    参与该次业务的user角色

## 如何运行

1. 完成`hyperledger composer`的安装以及业务网络的部署

1. 完成zeebe工作流引擎的安装和部署

1. (可选，建议执行)安装`zeebe-monitor`

1. 生成`run.sh`文件中所需要的两个角色，并将对应的card保存到`card/`目录下

1. `./run.sh`，文件最后两行提供区块链rest接口以及事件监听接口，请确保端口未被占用

1. 在zeebe引擎中添加一个`user`类型的工作流实例

1. `export GOMODULE=on`

1. `cd worker/roles/seller && go run main.go`

1. `cd worker/roles/user && go run main.go`

完成以上步骤后，就可以看到我们的两个worker处理分配到的task并更新区块链上数据的状态