var paramTypes 		= null	// 参数类型字符串拼接，例如"uint256,string"
var paramVals 		= null	// 参数值的字符串拼接，例如"2345675643,,Hello!%"   值使用`,,`拼接，用于区分参数中包含数组的情况

var args = process.argv
var len = args.length

for ( var i = 0; i < args.length; i++ )
{
    if ( args[i].indexOf( "-types" ) >= 0 && i + 1 < len ) {
        paramTypes = args[i+1]
    } else if ( args[i].indexOf( "-vals" ) >= 0 && i + 1 < len ) {
        paramVals = args[i+1]
    }
}

if ( paramTypes == null || paramVals == null ) {
    process.exit()
}


var Web3 = require("web3");
var fs = require("fs");
var web3 = new Web3();

var typeList = paramTypes.split( "," )
var valList  = paramVals.split( ",," )


if ( typeList.length != valList.length ) {
	process.exit()
}

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

var result = web3.eth.abi.encodeParameters( typeList, valList )
console.log( result.substring( 2, result.length ) )
