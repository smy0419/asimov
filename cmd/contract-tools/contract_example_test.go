package contract_tools

import (
	"fmt"
	"github.com/AsimovNetwork/asimov/common"
	"github.com/AsimovNetwork/asimov/vm/fvm"
	"math/big"
	"os/exec"
	"strings"
	"testing"
)

func TestEncodeMiningMemebers( t *testing.T ) {
	argTypes := "address[]"

	members := make( []string, 0 )
	//members = append( members, "0x66e3054b411051da5492aec7a823b00cb3add772d7" )
	//members = append( members, "0x663cf8b865f2f7e522ff639059e0a437c849ee5ab0" )
	//members = append( members, "0x6694c93303ac6605901e08b033d130e001bdeb9157" )
	//members = append( members, "0x66549a20042cfbf58b7416591a293c008a9ca3dd5b" )
	members = append( members, "0x66a8bd7839720bc7f49beab72fea5e85fcab5eecb8" )

	argVals := strings.Join( members, "," )

	fmt.Printf( "%s\n", argVals )

	cmd := exec.Command("node", "contract-tools/encode_parameter.js",
		"-types", argTypes, "-vals", fmt.Sprintf( "[%s]", argVals ) )

	fmt.Printf( "%v\n", cmd )
	ret, err := Output()

	if err != nil {
		fmt.Errorf( "%s\n", err.Error() )
		return
	}

	fmt.Printf( "mining members %s\n", ret )
}

//// test demo1
//// init env: run init_env.sh that located in dir `contract-tools`
//// make test as follow steps
//// 		step 1 : create a contract project dir in testdata
//// 		step 2 : copy contract file to project dir
//// 		step 3 : write test func like function TestHelloWorldContract or function TestSimpleOrganizationContract
//func TestHelloWorldContract( t *testing.T ) {
//	strAddrHex  := "662250f9452ac336daaeee722615619d2ba1422793"
//	container, err := NewContractContainer( strAddrHex, nil,false )
//
//	if err != nil {
//		t.Error( err.Error() )
//		return
//	}
//
//	// compile
//	action := newCompileEvmAction( "helloworld", "helloworld.sol", "HelloWorld" )
//	err = container.AddAction( action )
//	if err != nil {
//		t.Error( err.Error() )
//		return
//	}
//
//	// deploy
//	action = newDeployAction( "HelloWorld" )
//	action.(*DeployAction).addParam( "string", "HelloWorld" )
//	err = container.AddAction( action )
//	if err != nil {
//		t.Error( err.Error() )
//		return
//	}
//
//	// call speak()
//	action = newCallAction("HelloWorld", "speak", nil )
//	err = container.AddAction( action )
//	if err != nil {
//		t.Error( err.Error() )
//		return
//	}
//
//	// call speak1()
//	action = newCallAction("HelloWorld", "speak1", nil )
//	err = container.AddAction( action )
//	if err != nil {
//		t.Error( err.Error() )
//		return
//	}
//
//	_, err = container.Execute()
//	if err != nil {
//		t.Error( err.Error() )
//	}
//}


func addContractToTemplateWarehouse(
	container *ContractContainer,
	projectName, srcFile, mainClass string,
	categoryId uint16,
	templateName string,
	gas int64 ) error {

	// deploy evm system contract
	action := newDeployEvmSystemContractAction()
	err := container.AddAction( action )
	if err != nil {
		return err
	}

	_, err = container.Execute()
	if err != nil {
		return err
	}

	// compile target contract
	action = newCompileEvmAction( projectName, srcFile, mainClass )
	err = container.AddAction( action )
	if err != nil {
		return err
	}

	// execute and return last action result( CompileEvmAction will be return byte code )
	byteCode, err := container.Execute()
	if err != nil {
		return err
	}

	// call TemplateWarehouse.create to deploy target
	category := categoryId
	tempName := templateName
	codeHex :=  "0x" + common.Bytes2Hex( byteCode )
	cost := big2Byte4ArrText( big.NewInt( gas ) )

	action = newCallAction( common.TemplateWarehouse.String(), "create", nil )
	action.(*CallAction).addParam( "uint16", fmt.Sprintf("%d", category ) )
	action.(*CallAction).addParam( "string", tempName )
	action.(*CallAction).addParam( "bytes", codeHex )
	action.(*CallAction).addParam( "bytes4[]", fmt.Sprintf("[%s]", cost) )
	err = container.AddAction( action )
	if err != nil {
		return err
	}

	// call TemplateWarehouse.getTemplate
	action = newCallAction( common.TemplateWarehouse.String(), "getTemplate", nil )
	action.(*CallAction).addParam( "uint16", fmt.Sprintf("%d", category ) )
	action.(*CallAction).addParam( "string", tempName )
	err = container.AddAction( action )
	if err != nil {
		return err
	}

	// execute and return last action result( CallAction will be return contract method result )
	ret, err := container.Execute()
	if err != nil {
		return err
	}

	// decode template info
	var code []byte
	cTime := new(big.Int)
	returnData := &[]interface{}{new(string), &code, &cTime, new(uint8), new(uint8), new(uint8), new(uint8)}
	abi, err := container.ctx.getContractAbi( common.TemplateWarehouse.String() )
	if err != nil {
		return err
	}

	err = fvm.UnPackFunctionResult( abi, returnData,"getTemplate", ret )
	if err != nil {
		return err
	}

	container.ctx.setContractCodeByte( mainClass, code )

	name 		:= *((*returnData)[0]).(*string)			// template namne
	codeHex 	= common.Bytes2Hex( code )					// template code hex
	createTime 	:= cTime.Int64()							// template create time
	approveCount 	:= *((*returnData)[3]).(*uint8)			// template approve count
	rejectCount 	:= *((*returnData)[4]).(*uint8)			// template reject count
	reviewers 	:= *((*returnData)[5]).(*uint8)				// template reviewers
	status 		:= *((*returnData)[6]).(*uint8)				// template status

	fmt.Printf("name %v\n", name )
	fmt.Printf("codeHex %v\n", codeHex )
	fmt.Printf("createTime %v\n", createTime )
	fmt.Printf("approveCount %v\n", approveCount )
	fmt.Printf("rejectCount %v\n", rejectCount )
	fmt.Printf("reviewers %v\n", reviewers )
	fmt.Printf("status %v\n", status )

	return nil
}

//// test demo2
//// run init-contract-library.sh to download contract library that imported by simple_organization.sol
//func TestSimpleOrganizationContract( t *testing.T ) {
//	strAddrHex  := "662250f9452ac336daaeee722615619d2ba1422793"
//
//	// contract container
//	container, err := NewContractContainer( strAddrHex, nil, false )
//	if err != nil {
//		t.Error( err.Error() )
//		return
//	}
//
//	simpleOrgProjectName 	:= "simple_organization"
//	simpleOrgSrcFile		:= "simple_organization.sol"
//	simpleOrgMainClass 		:= "Association"
//
//	categoryId := uint16(1)
//	templateName := "TempOrg"
//	// add contract to template warehouse
//	err = addContractToTemplateWarehouse (
//		container,
//		simpleOrgProjectName,
//		simpleOrgSrcFile,
//		simpleOrgMainClass,
//		categoryId,
//		templateName,
//		3000)
//
//	if err != nil {
//		t.Error( err.Error() )
//		return
//	}
//
//	// get code from template warehouse, then deploy
//	action := newDeployAction( simpleOrgMainClass )
//	action.(*DeployAction).addParam( "string", "org1" )
//	action.(*DeployAction).addParam( "string", "ass1" )
//	action.(*DeployAction).addParam( "address[]", "[0x662250f9452ac336daaeee722615619d2ba1422793,0x662250f9452ac336daaeee722615619d2ba1422794]" )
//
//	err = container.AddAction( action )
//	if err != nil {
//		t.Log( err.Error() )
//		return
//	}
//	_, err = container.Execute()
//	if err != nil {
//		t.Error( err.Error() )
//		return
//	}
//
//	contractRes := container.ctx.getContractRes( simpleOrgMainClass )
//	if contractRes == nil {
//		t.Error( "contract res is not exist" )
//		return
//	}
//
//	officialAddr := common.BytesToAddress( container.chain.GetOfficialAddr() )
//
//	// init template, the caller must be official address
//	action = newCallAction( simpleOrgMainClass, "initTemplate", &officialAddr )
//	action.(*CallAction).addParam( "uint16", fmt.Sprintf("%d", categoryId) )
//	action.(*CallAction).addParam("string", templateName )
//	err = container.AddAction( action )
//	if err != nil {
//		t.Error( err.Error() )
//		return
//	}
//
//	action = newCallAction( simpleOrgMainClass, "maintainRegister", nil )
//	err = container.AddAction( action )
//	if err != nil {
//		t.Error( err.Error() )
//		return
//	}
//
//	action = newCallAction( simpleOrgMainClass, "getTemplateInfo", nil )
//	err = container.AddAction( action )
//	if err != nil {
//		t.Error( err.Error() )
//		return
//	}
//
//	action = newCallAction( simpleOrgMainClass, "createAsset",  nil )
//	action.(*CallAction).addParam( "string", "a" )
//	action.(*CallAction).addParam( "string", "b" )
//	action.(*CallAction).addParam( "string", "c" )
//	action.(*CallAction).addParam( "uint32", fmt.Sprintf("%d", 1 ) )
//	action.(*CallAction).addParam( "uint32", fmt.Sprintf("%d", 1 ) )
//	action.(*CallAction).addParam( "uint256", fmt.Sprintf("%d", 1 ) )
//
//	err = container.AddAction( action )
//	if err != nil {
//		t.Error( err.Error() )
//		return
//	}
//
//	_, err = container.Execute()
//	if err != nil {
//		t.Error( err.Error() )
//	}
//}


func big2Byte4ArrText( val *big.Int ) string {
	var text string
	byteArr := val.Bytes()
	for _, b := range byteArr {
		if len(text) > 0 {
			text = text + ","
		}
		text = text + fmt.Sprintf( "0x%08x", b )
	}

	return text
}

//func TestRustHelloWorldContract( t *testing.T ) {
//	strAddrHex  := "662250f9452ac336daaeee722615619d2ba1422793"
//
//	// contract container
//	container, err := NewContractContainer( strAddrHex,nil, true )
//	if err != nil {
//		t.Error( err.Error() )
//		return
//	}
//
//	projectName := "rust_helloworld"
//	targetName 	:= "hello_world"
//	buildScript := "build.sh"
//
//	action := newCompileWasmAction( projectName, targetName, buildScript )
//	err = container.AddAction( action )
//
//	if err != nil {
//		t.Error( err.Error() )
//		return
//	}
//
//	action = newDeployAction( targetName )
//	err = container.AddAction( action )
//
//	if err != nil {
//		t.Error( err.Error() )
//		return
//	}
//
//	action = newCallAction( targetName, "display", nil )
//	err = container.AddAction( action )
//
//	if err != nil {
//		t.Error( err.Error() )
//		return
//	}
//
//	_, err = container.Execute()
//	if err != nil {
//		t.Error( err.Error() )
//	}
//}

//func TestRustTokenContract( t *testing.T ) {
//	strAddrHex  := "662250f9452ac336daaeee722615619d2ba1422793"
//
//	// contract container
//	container, err := NewContractContainer( strAddrHex, nil, true )
//	if err != nil {
//		t.Error( err.Error() )
//		return
//	}
//
//	projectName := "rust_token"
//	targetName 	:= "token_contract"
//	buildScript := "build.sh"
//
//	action := newCompileWasmAction( projectName, targetName, buildScript )
//	err = container.AddAction( action )
//
//	if err != nil {
//		t.Error( err.Error() )
//		return
//	}
//
//	action = newDeployAction( targetName )
//	action.(*DeployAction).addParam( "uint256", "17" )
//	err = container.AddAction( action )
//
//	if err != nil {
//		t.Error( err.Error() )
//		return
//	}
//
//	action = newCallAction( targetName, "totalSupply", nil )
//	err = container.AddAction( action )
//
//	if err != nil {
//		t.Error( err.Error() )
//		return
//	}
//
//	action = newCallAction( targetName, "balanceOf", nil )
//	action.(*CallAction).addParam("address", "0x"+strAddrHex )
//	err = container.AddAction( action )
//
//	if err != nil {
//		t.Error( err.Error() )
//		return
//	}
//
//	action = newCallAction( targetName, "transfer",nil )
//	action.(*CallAction).addParam( "address", "0x"+strAddrHex )
//	action.(*CallAction).addParam( "uint256", fmt.Sprintf( "%d", 10 ) )
//	err = container.AddAction( action )
//
//	if err != nil {
//		t.Error( err.Error() )
//		return
//	}
//
//	// test 21 bytes address
//	action = newCallAction( targetName, "getSender",nil )
//	err = container.AddAction( action )
//
//	if err != nil {
//		t.Error( err.Error() )
//		return
//	}
//
//	_, err = container.Execute()
//	if err != nil {
//		t.Error( err.Error() )
//	}
//
//	fmt.Printf( "%v\n", common.Hex2Bytes( strAddrHex ) )
//}
//

//func TestRustCustomInstruction( t *testing.T ) {
//	strAddrHex := "662250f9452ac336daaeee722615619d2ba1422793"
//
//	container, err := NewContractContainer( strAddrHex, protos.NewAssets(1,1,1 ), true )
//	if err != nil {
//		t.Error( err.Error() )
//		return
//	}
//
//	action := newDeployEvmSystemContractAction()
//	err = container.AddAction( action )
//	if err != nil {
//		t.Error( err.Error() )
//		return
//	}
//
//	_, err = container.Execute()
//	if err != nil {
//		t.Error( err.Error() )
//		return
//	}
//
//	projectName := "rust_custom_instruction"
//	targetName 	:= "custom_instruction"
//	buildScript := "build.sh"
//
//	action = newCompileWasmAction( projectName, targetName, buildScript )
//	err = container.AddAction( action )
//
//	if err != nil {
//		t.Error( err.Error() )
//		return
//	}
//
//	action = newDeployAction( targetName )
//	err =  container.AddAction( action )
//
//	if err != nil {
//		t.Error( err.Error() )
//		return
//	}
//
//	action = newCallAction( targetName, "test", nil )
//	err = container.AddAction( action )
//
//	if err != nil {
//		t.Error( err.Error() )
//		return
//	}
//
//	action = newCallAction( targetName, "testGetAsset",  nil )
//	err = container.AddAction( action )
//
//	if err != nil {
//		t.Error( err.Error() )
//		return
//	}
//
//	// we should deploy a organization contract, and register organization before creating asset
//	// if not register organization, it will return `no authority to xxx asset`
//	action = newCallAction( targetName, "testCreateAsset", nil )
//	err = container.AddAction( action )
//
//	if err != nil {
//		t.Error( err.Error() )
//		return
//	}
//
//	action = newCallAction( targetName, "testMintAsset", nil )
//	err = container.AddAction( action )
//
//	if err != nil {
//		t.Error( err.Error() )
//		return
//	}
//
//	action = newCallAction( targetName, "testDeployContract", nil )
//	action.(*CallAction).addParam( "uint32", fmt.Sprintf("%d", 1 ) )
//	action.(*CallAction).addParam( "string", "org1" )
//	action.(*CallAction).addParam( "string", "667879" )
//
//	err = container.AddAction( action )
//
//	if err != nil {
//		t.Error( err.Error() )
//		return
//	}
//
//	action = newCallAction( targetName, "testTransfer",nil )
//	action.(*CallAction).addParam( "address", "0x"+strAddrHex )
//	action.(*CallAction).addParam( "uint256", "100" )
//	action.(*CallAction).addParam( "uint256", "0" )
//	err = container.AddAction( action )
//
//	if err != nil {
//		t.Error( err.Error() )
//		return
//	}
//
//	_, err = container.Execute()
//	if err != nil {
//		t.Error( err.Error() )
//	}
//}
