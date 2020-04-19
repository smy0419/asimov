var funcName 	= null		// 调用的函数名称
var paramTypes 	= ""		// 参数类型的字符串拼接形式，例如"string,uint256"
var paramVals 	= ""		// 参数值的字符串拼接形式，例如"aaa,123"

var args = process.argv
var len = args.length

for ( var i = 0; i < args.length; i++ )
{
	if ( args[i].indexOf( "-func" ) >= 0 && i + 1 < len  ) {
		funcName = args[i+1]
    } else if ( args[i].indexOf( "-types" ) >= 0 && i + 1 < len ) {
		paramTypes = args[i+1]
	} else if ( args[i].indexOf( "-vals" ) >= 0 && i + 1 < len ) {
		paramVals = args[i+1]
	}
}

if ( funcName == null ) {
	console.log( "funcName is null" )
	process.exit()
}

typeArr = paramTypes.split( "," )
valArr = paramVals.split( ",," )

if ( typeArr.length != valArr.length ) {
	console.log( "types and vals are not match" )
	process.exit()
}

var Web3 = require("web3");
var fs = require("fs");
var web3 = new Web3();

var inputs = new Array()
for (var i = 0; i < typeArr.length; i++ ) {
    inputs.push( {
		type: typeArr[i],
		name: "unknown",
	} )
}

var funcSigJsonObj 	= {
	name : funcName,
	type : "function",
	inputs : inputs,
}

/*
for ( var i = 0; i < valArr.length; i++ ) {
	if ( valArr[i].indexOf(",") >= 0 ) {
		valArr[i] = valArr[i].split( "," )
	}
}*/

var valList = valArr
for ( var i = 0; i < valList.length; i++ ) {
    if ( valList[i].length < 2 ||
        (valList[i][0] != "[" && valList[i][valList[i].length-1] != "]") ) {
        continue
    }

    var val = ""
    if (valList[i].length > 2 ) {
        val = valList[i].substring( 1, valList[i].length - 1 )
    }

    valList[i] = val.split( "," )
}

var result = null

if ( paramTypes.length > 0 ) {
    result = web3.eth.abi.encodeFunctionCall( funcSigJsonObj, valArr )
} else {
	result = web3.eth.abi.encodeFunctionSignature( funcSigJsonObj )
}

console.log( result.substring( 2, result.length ) )
